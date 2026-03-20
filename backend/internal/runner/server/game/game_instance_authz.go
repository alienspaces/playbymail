package game

import (
	"net/http"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// requireManagerSubscription verifies the authenticated account has an active manager
// subscription for the given game. Returns both the auth data and the manager subscription record.
// Authentication is already guaranteed by the token middleware before any handler runs.
func requireManagerSubscription(l logger.Logger, r *http.Request, mm *domain.Domain, gameID string) (*server.AuthenData, *game_record.GameSubscription, error) {
	authenData := server.GetRequestAuthenData(l, r)

	managerSubRec, err := mm.GetGameSubscriptionRecByAccountUserAndGame(
		authenData.AccountUser.ID,
		gameID,
		game_record.GameSubscriptionTypeManager,
	)
	if err != nil {
		l.Warn("failed to find manager subscription for account_user >%s< and game >%s<: %v",
			authenData.AccountUser.ID, gameID, err)
		return nil, nil, coreerror.NewUnauthorizedError()
	}

	return authenData, managerSubRec, nil
}

// authorizeManagerModify verifies the authenticated account owns the given game instance
// by confirming their manager subscription is linked to it via game_subscription_instance.
// Used by update, delete, and all lifecycle handlers.
func authorizeManagerModify(l logger.Logger, r *http.Request, mm *domain.Domain, gameID, instanceID string) (*server.AuthenData, error) {
	authenData, managerSubRec, err := requireManagerSubscription(l, r, mm, gameID)
	if err != nil {
		return nil, err
	}

	instanceLinks, err := mm.GetGameSubscriptionInstanceRecsBySubscription(managerSubRec.ID)
	if err != nil {
		l.Warn("failed to get instance links for subscription >%s<: %v", managerSubRec.ID, err)
		return nil, coreerror.NewUnauthorizedError()
	}

	for _, link := range instanceLinks {
		if link.GameInstanceID == instanceID {
			return authenData, nil
		}
	}

	l.Warn("authenticated account_user >%s< does not own game instance >%s<", authenData.AccountUser.ID, instanceID)
	return nil, coreerror.NewUnauthorizedError()
}
