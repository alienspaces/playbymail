package game

import (
	"net/http"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// requireDesignerSubscription verifies the authenticated account user holds an active
// designer subscription for the given game. Returns both the auth data and the subscription record.
// Authentication is already guaranteed by the token middleware before any handler runs.
func requireDesignerSubscription(l logger.Logger, r *http.Request, mm *domain.Domain, gameID string) (*server.AuthenData, *game_record.GameSubscription, error) {
	authenData := server.GetRequestAuthenData(l, r)

	designerSubRec, err := mm.GetGameSubscriptionRecByAccountUserAndGame(
		authenData.AccountUser.ID,
		gameID,
		game_record.GameSubscriptionTypeDesigner,
	)
	if err != nil {
		l.Warn("failed to find designer subscription for account_user >%s< and game >%s<: %v",
			authenData.AccountUser.ID, gameID, err)
		return nil, nil, coreerror.NewUnauthorizedError()
	}

	return authenData, designerSubRec, nil
}

// authorizeDesignerModify verifies the authenticated account user holds an active
// designer subscription for the given game. Used by create, update, and delete
// handlers to ensure only the game's own designer can modify its resources.
func authorizeDesignerModify(l logger.Logger, r *http.Request, mm *domain.Domain, gameID string) (*server.AuthenData, *game_record.GameSubscription, error) {
	return requireDesignerSubscription(l, r, mm, gameID)
}
