package player

import (
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

func playerJoinGameHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "playerJoinGameHandlerConfig")

	l.Debug("Adding player join game handler configuration")

	playerJoinGameConfig := make(map[string]server.HandlerConfig)

	// TODO: Add join game endpoints
	// GET /api/v1/player/join-game/:join_game_key
	// POST /api/v1/player/join-game/:join_game_key/verify-email
	// GET /api/v1/player/join-game/:join_game_key/sheet
	// PUT /api/v1/player/join-game/:join_game_key/sheet
	// POST /api/v1/player/join-game/:join_game_key/submit

	return playerJoinGameConfig, nil
}

