package handler_auth

import (
	"net/http"

	coreconfig "gitlab.com/alienspaces/playbymail/core/config"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

const (
	PermissionGameDesign     server.AuthorizedPermission = "game_design"
	PermissionGameManagement server.AuthorizedPermission = "game_management"
	PermissionGamePlaying    server.AuthorizedPermission = "game_playing"
)

// Map account subscription types to permissions
var subscriptionPermissionMap = map[string]server.AuthorizedPermission{
	account_record.AccountSubscriptionTypeBasicGameDesigner:        PermissionGameDesign,
	account_record.AccountSubscriptionTypeProfessionalGameDesigner: PermissionGameDesign,
	account_record.AccountSubscriptionTypeBasicManager:             PermissionGameManagement,
	account_record.AccountSubscriptionTypeProfessionalManager:      PermissionGameManagement,
	account_record.AccountSubscriptionTypeBasicPlayer:              PermissionGamePlaying,
	account_record.AccountSubscriptionTypeProfessionalPlayer:       PermissionGamePlaying,
}

// authenticateRequestTokenFunc authenticates a request based on a session token. Returning anything
// other than an AuthenData{} with a valid Typewill result in a 401 Unauthorized response.
func AuthenticateRequestTokenFunc(cfg config.Config, l logger.Logger, m domainer.Domainer, r *http.Request) (server.AuthenData, error) {

	l.Info("authenticateRequestTokenFunc called")
	l.Info("request Authorization header >%s<", r.Header.Get("Authorization"))

	mm := m.(*domain.Domain)

	// Check for development mode authentication bypass
	if cfg.AppEnv == coreconfig.AppEnvDevelop {

		if bypassEmail := r.Header.Get("X-Bypass-Authentication"); bypassEmail != "" {
			l.Info("development mode: using bypass authentication for email >%s<", bypassEmail)

			// In development mode, query the actual account record by email bypassing
			// the need for actual session tokens and using real account data.
			accountRecs, err := mm.GetManyAccountRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: account_record.FieldAccountUserEmail, Val: bypassEmail},
				},
				Limit: 1,
			})
			if err != nil {
				l.Warn("development mode: failed to get account by email >%s< >%v<", bypassEmail, err)
				return server.AuthenData{}, err
			}

			if len(accountRecs) == 0 {
				l.Warn("development mode: no account found for email >%s<", bypassEmail)
				return server.AuthenData{}, nil
			}

			accountRec := accountRecs[0]
			l.Info("development mode: found account ID >%s< for email >%s<", accountRec.ID, bypassEmail)

			// Get account contact name if available
			accountName := ""
			contactRecs, err := mm.GetManyAccountUserContactRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: account_record.FieldAccountUserContactAccountUserID, Val: accountRec.ID},
				},
				Limit: 1,
				OrderBy: []coresql.OrderBy{
					{Col: account_record.FieldAccountUserContactCreatedAt, Direction: coresql.OrderDirectionASC},
				},
			})
			if err == nil && len(contactRecs) > 0 {
				accountName = contactRecs[0].Name
			}

			return server.AuthenData{
				Type:    server.AuthenticatedTypeToken,
				RLSType: server.RLSTypeRestricted,
				AccountUser: server.AuthenticatedAccountUser{
					ID:        accountRec.ID,
					AccountID: accountRec.AccountID,
					Name:      accountName,
					Email:     accountRec.Email,
				},
			}, nil
		}
	}

	// Session token for authentication.
	var token string

	// First, try Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && len(authHeader) >= 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
		l.Debug("token extracted from Authorization header")
	}

	// Fallback: check query parameter (used by iframe-based previews that cannot set headers)
	if token == "" {
		token = r.URL.Query().Get("token")
		if token != "" {
			l.Debug("token extracted from query parameter")
		}
	}

	if token == "" {
		l.Warn("no token found in Authorization header or query parameters")
		return server.AuthenData{}, nil
	}

	// Verify the session token
	tokenPreview := token
	if len(token) > 20 {
		tokenPreview = token[:20]
	}
	l.Info("verifying session token >%s<", tokenPreview)
	accountRec, err := mm.VerifyAccountSessionToken(token)
	if err != nil {
		l.Warn("failed to verify account session token >%v<", err)
		return server.AuthenData{}, err
	}

	if accountRec == nil {
		l.Warn("no account found for session token")
		return server.AuthenData{}, nil
	}

	l.Info("verified session token for account ID >%s< Email >%s<", accountRec.ID, accountRec.Email)

	// Get account contact name if available
	accountName := ""
	contactRecs, err := mm.GetManyAccountUserContactRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountUserContactAccountUserID, Val: accountRec.ID},
		},
		OrderBy: []coresql.OrderBy{
			{Col: account_record.FieldAccountUserContactCreatedAt, Direction: coresql.OrderDirectionASC},
		},
		Limit: 1,
	})
	if err == nil && len(contactRecs) > 0 {
		accountName = contactRecs[0].Name
	}

	// Get active account subscriptions
	l.Info("getting account subscriptions for account ID >%s<", accountRec.ID)
	accountSubscriptions, err := mm.GetManyAccountSubscriptionRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountSubscriptionAccountID, Val: accountRec.AccountID},
			{Col: account_record.FieldAccountSubscriptionStatus, Val: account_record.AccountSubscriptionStatusActive},
		},
	})
	if err != nil {
		l.Warn("failed to get account subscriptions >%v<", err)
		return server.AuthenData{}, coreerror.NewInternalError("failed to get account subscriptions: %v", err)
	}

	l.Info("found >%d< active account subscriptions for account ID >%s<", len(accountSubscriptions), accountRec.ID)
	for _, sub := range accountSubscriptions {
		l.Info("account subscription ID >%s< Type >%s< Status >%s< AccountID >%s<",
			sub.ID, sub.SubscriptionType, sub.Status, sub.AccountID)
	}

	// All accounts must have at least some account subscriptions
	if len(accountSubscriptions) == 0 {
		l.Warn("account ID >%s< has no active account subscriptions", accountRec.ID)
		return server.AuthenData{}, coreerror.NewUnauthorizedError()
	}

	// Build permissions set from account subscriptions only
	permissionSet := make(map[server.AuthorizedPermission]bool)

	// Add account subscriptions permissions
	for _, sub := range accountSubscriptions {
		if perm, ok := subscriptionPermissionMap[sub.SubscriptionType]; ok {
			l.Info("mapping subscription type >%s< to permission >%s<", sub.SubscriptionType, perm)
			permissionSet[perm] = true
		} else {
			l.Warn("subscription type >%s< not found in permission map", sub.SubscriptionType)
		}
	}

	// Convert set to slice
	permissions := make([]server.AuthorizedPermission, 0, len(permissionSet))
	for perm := range permissionSet {
		permissions = append(permissions, perm)
	}

	authenData := server.AuthenData{
		Type:        server.AuthenticatedTypeToken,
		RLSType:     server.RLSTypeRestricted,
		Permissions: permissions,
		AccountUser: server.AuthenticatedAccountUser{
			ID:        accountRec.ID,
			AccountID: accountRec.AccountID,
			Name:      accountName,
			Email:     accountRec.Email,
		},
	}

	l.Info("authenticated account: ID=%s Email=%s Name=%s Permissions=%v Subscriptions=%d",
		authenData.AccountUser.ID, authenData.AccountUser.Email, authenData.AccountUser.Name, permissions, len(accountSubscriptions))

	return authenData, nil
}
