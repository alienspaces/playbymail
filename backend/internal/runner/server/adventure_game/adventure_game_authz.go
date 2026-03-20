package adventure_game

import (
	"net/http"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// authorizeAdventureGameRead verifies the request is authenticated.
// Used by handlers that need to confirm the caller is a known user.
func authorizeAdventureGameRead(l logger.Logger, r *http.Request) (*server.AuthenData, error) {
	authenData := server.GetRequestAuthenData(l, r)
	if authenData == nil || authenData.AccountUser.ID == "" {
		l.Warn("authenticated account is required")
		return nil, coreerror.NewUnauthorizedError()
	}
	return authenData, nil
}

// authorizeAdventureGameDesigner verifies the authenticated account holds an active
// designer subscription for the given game. Used by create, update, and delete
// handlers to ensure only the game's own designer can modify its resources.
func authorizeAdventureGameDesigner(l logger.Logger, r *http.Request, mm *domain.Domain, gameID string) (*server.AuthenData, error) {
	authenData, err := authorizeAdventureGameRead(l, r)
	if err != nil {
		return nil, err
	}

	_, err = mm.GetGameSubscriptionRecByAccountAndGame(
		authenData.AccountUser.AccountID,
		gameID,
		game_record.GameSubscriptionTypeDesigner,
	)
	if err != nil {
		l.Warn("failed to find designer subscription for account >%s< and game >%s<: %v",
			authenData.AccountUser.AccountID, gameID, err)
		return nil, coreerror.NewUnauthorizedError()
	}

	return authenData, nil
}
