package game

import (
	"net/http"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// authorizeGameInstanceRead verifies the request is authenticated.
// Used by GET handlers to ensure only authenticated users access game instance data.
func authorizeGameInstanceRead(l logger.Logger, r *http.Request) (*server.AuthenData, error) {
	authenData := server.GetRequestAuthenData(l, r)
	if authenData == nil || authenData.AccountUser.ID == "" {
		l.Warn("authenticated account is required")
		return nil, coreerror.NewUnauthorizedError()
	}
	return authenData, nil
}

// authorizeGameInstanceCreate verifies the authenticated account has an active manager
// subscription for the given game. Returns both the auth data and the manager subscription
// record so callers can link the new instance to the subscription.
func authorizeGameInstanceCreate(l logger.Logger, r *http.Request, mm *domain.Domain, gameID string) (*server.AuthenData, *game_record.GameSubscription, error) {
	authenData, err := authorizeGameInstanceRead(l, r)
	if err != nil {
		return nil, nil, err
	}

	managerSubRec, err := mm.GetGameSubscriptionRecByAccountAndGame(
		authenData.AccountUser.AccountID,
		gameID,
		game_record.GameSubscriptionTypeManager,
	)
	if err != nil {
		l.Warn("failed to find manager subscription for account >%s< and game >%s<: %v",
			authenData.AccountUser.AccountID, gameID, err)
		return nil, nil, coreerror.NewUnauthorizedError()
	}

	return authenData, managerSubRec, nil
}

// authorizeGameInstanceModify verifies the authenticated account owns the given game instance
// by confirming their manager subscription is linked to it via game_subscription_instance.
// Used by update, delete, and all lifecycle handlers.
func authorizeGameInstanceModify(l logger.Logger, r *http.Request, mm *domain.Domain, gameID, instanceID string) (*server.AuthenData, error) {
	authenData, managerSubRec, err := authorizeGameInstanceCreate(l, r, mm, gameID)
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

	l.Warn("authenticated account >%s< does not own game instance >%s<", authenData.AccountUser.AccountID, instanceID)
	return nil, coreerror.NewUnauthorizedError()
}
