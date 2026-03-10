package handler_auth

import (
	"net/http"

	coreconfig "gitlab.com/alienspaces/playbymail/core/config"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
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

			// In development mode, query the actual account user record by email bypassing
			// the need for actual session tokens and using real account data.
			accountUserRecs, err := mm.GetManyAccountUserRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: account_record.FieldAccountUserEmail, Val: bypassEmail},
				},
				Limit: 1,
			})
			if err != nil {
				l.Warn("development mode: failed to get account by email >%s< >%v<", bypassEmail, err)
				return server.AuthenData{}, err
			}

			if len(accountUserRecs) == 0 {
				l.Warn("development mode: no account found for email >%s<", bypassEmail)
				return server.AuthenData{}, nil
			}

			accountUserRec := accountUserRecs[0]
			l.Info("development mode: found account user ID >%s< for email >%s<", accountUserRec.ID, bypassEmail)

			// Get account contact name if available
			accountName := ""
			contactRecs, err := mm.GetManyAccountUserContactRecs(&coresql.Options{
				Params: []coresql.Param{
					{Col: account_record.FieldAccountUserContactAccountUserID, Val: accountUserRec.ID},
				},
				Limit: 1,
				OrderBy: []coresql.OrderBy{
					{Col: account_record.FieldAccountUserContactCreatedAt, Direction: coresql.OrderDirectionASC},
				},
			})
			accountUserContactID := ""
			if err == nil && len(contactRecs) > 0 {
				accountName = nullstring.ToString(contactRecs[0].Name)
				accountUserContactID = contactRecs[0].ID
			}

			return server.AuthenData{
				Type: server.AuthenticatedTypeToken,
				AccountUser: server.AuthenticatedAccountUser{
					ID:                   accountUserRec.ID,
					AccountID:            accountUserRec.AccountID,
					AccountUserContactID: accountUserContactID,
					Name:                 accountName,
					Email:                accountUserRec.Email,
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
	accountUserRec, err := mm.VerifyAccountUserSessionToken(token)
	if err != nil {
		l.Warn("failed to verify account session token >%v<", err)
		return server.AuthenData{}, err
	}

	if accountUserRec == nil {
		l.Warn("no account found for session token")
		return server.AuthenData{}, nil
	}

	l.Info("verified session token for account user ID >%s< Email >%s<", accountUserRec.ID, accountUserRec.Email)

	// Get account contact name and ID if available
	accountUserContactRec, err := mm.GetAccountUserContactRecByAccountUserID(accountUserRec.ID, nil)
	if err != nil {
		l.Warn("failed to get account user contact for account user ID >%s< >%v<", accountUserRec.ID, err)
		return server.AuthenData{}, err
	}

	if accountUserContactRec == nil {
		l.Warn("no account user contact found for account user ID >%s<", accountUserRec.ID)
		return server.AuthenData{}, nil
	}

	l.Info("found account user contact ID >%s< for account user ID >%s<", accountUserContactRec.ID, accountUserRec.ID)

	// Get active account subscriptions
	l.Info("getting account subscriptions for account ID >%s< account user ID >%s<", accountUserRec.AccountID, accountUserRec.ID)

	accountSubscriptions, err := mm.GetManyAccountSubscriptionRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountSubscriptionAccountUserID, Val: accountUserRec.ID},
			{Col: account_record.FieldAccountSubscriptionStatus, Val: account_record.AccountSubscriptionStatusActive},
		},
	})
	if err != nil {
		l.Warn("failed to get account subscriptions by account ID >%v<", err)
		return server.AuthenData{}, coreerror.NewInternalError("failed to get account subscriptions: %v", err)
	}

	l.Info("found >%d< active account subscriptions for account user ID >%s<", len(accountSubscriptions), accountUserRec.ID)

	for _, sub := range accountSubscriptions {
		l.Info("account subscription ID >%s< Type >%s< Status >%s< AccountUserID >%s<",
			sub.ID, sub.SubscriptionType, sub.Status, sub.AccountUserID)
	}

	// All accounts must have at least some account subscriptions
	if len(accountSubscriptions) == 0 {
		l.Warn("account user ID >%s< has no active account subscriptions", accountUserRec.ID)
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
		Permissions: permissions,
		AccountUser: server.AuthenticatedAccountUser{
			ID:                   accountUserRec.ID,
			AccountID:            accountUserRec.AccountID,
			AccountUserContactID: accountUserContactRec.ID,
			Name:                 nullstring.ToString(accountUserContactRec.Name),
			Email:                accountUserRec.Email,
		},
	}

	l.Info("authenticated account: ID=%s Email=%s Name=%s Permissions=%v Subscriptions=%d",
		authenData.AccountUser.ID, authenData.AccountUser.Email, authenData.AccountUser.Name, permissions, len(accountSubscriptions))

	return authenData, nil
}
