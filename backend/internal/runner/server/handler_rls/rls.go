package handler_rls

import (
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

const (
	rlsIdentifierAccountID          = "account_id"
	rlsIdentifierGameID             = "game_id"
	rlsIdentifierGameSubscriptionID = "game_subscription_id"
)

// rlsFunc determines what game resources the authenticated user has access to
func HandlerRLSFunc(l logger.Logger, m domainer.Domainer, authedReq server.AuthenData) (server.RLS, error) {

	l.Info("(playbymail) rlsFunc called for account: ID=%s Email=%s",
		authedReq.Account.ID, authedReq.Account.Email)

	mm := m.(*domain.Domain)

	// Get all games the user has access to through subscriptions
	gameSubscriptions, err := mm.GameSubscriptionRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: game_record.FieldGameSubscriptionAccountID,
				Val: authedReq.Account.ID,
			},
		},
	})
	if err != nil {
		l.Warn("(playbymail) failed to get game subscriptions >%v<", err)
		return server.RLS{}, err
	}

	// Extract game IDs from subscriptions (for subscription-based access)
	subscribedGameIDs := make([]string, 0, len(gameSubscriptions))
	for _, sub := range gameSubscriptions {
		subscribedGameIDs = append(subscribedGameIDs, sub.GameID)
	}

	// Get all game subscription IDs for the user
	gameSubscriptionIDs := make([]string, 0, len(gameSubscriptions))
	for _, sub := range gameSubscriptions {
		gameSubscriptionIDs = append(gameSubscriptionIDs, sub.ID)
	}

	// Create RLS identifiers map
	// account_id will automatically filter games owned by this account
	identifiers := map[string][]string{
		rlsIdentifierAccountID: {authedReq.Account.ID},
	}

	// Add game IDs for subscription-based access
	if len(subscribedGameIDs) > 0 {
		identifiers[rlsIdentifierGameID] = subscribedGameIDs
	}

	// Add game subscription IDs
	if len(gameSubscriptionIDs) > 0 {
		identifiers[rlsIdentifierGameSubscriptionID] = gameSubscriptionIDs
	}

	l.Info("(playbymail) RLS applied: account_id=%s subscription_games=%d game_ids=%v subscription_ids=%v",
		authedReq.Account.ID, len(subscribedGameIDs), subscribedGameIDs, gameSubscriptionIDs)

	return server.RLS{
		Identifiers: identifiers,
	}, nil
}

// HandlerRLSAccountIdentifierValue returns the account identifier value for the current authenticated user
func HandlerRLSAccountIdentifierValue(l logger.Logger, r *http.Request) (string, error) {
	rlsData := server.GetRequestRLSData(l, r)
	if rlsData == nil {
		return "", nil
	}

	return rlsData.Identifiers[rlsIdentifierAccountID][0], nil
}

// HandlerRLSGameIdentifierValues returns the list of game identifier the current authenticated user has access to
func HandlerRLSGameIdentifierValues(l logger.Logger, r *http.Request) ([]string, error) {
	rlsData := server.GetRequestRLSData(l, r)
	if rlsData == nil {
		return nil, nil
	}

	return rlsData.Identifiers[rlsIdentifierGameID], nil
}

// HandlerRLSGameSubscriptionIdentifierValues returns the list of game subscription identifier the current authenticated user has access to
func HandlerRLSGameSubscriptionIdentifierValues(l logger.Logger, r *http.Request) ([]string, error) {
	rlsData := server.GetRequestRLSData(l, r)
	if rlsData == nil {
		return nil, nil
	}

	return rlsData.Identifiers[rlsIdentifierGameSubscriptionID], nil
}
