package player

import (
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// The set of player API's are APIs meant to be accessed by individual players.

const (
	packageName = "player"
)

var referenceSchemas = []jsonschema.Schema{
	{
		Location: "api/common_schema",
		Name:     "query.schema.json",
	},
	{
		Location: "api/common_schema",
		Name:     "common.schema.json",
	},
}

func PlayerHandlerConfig(cfg config.Config, l logger.Logger, scnr turnsheet.TurnSheetScanner) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "PlayerHandlerConfig")

	l.Debug("Adding player handler configuration")

	playerConfig := make(map[string]server.HandlerConfig)

	// Additional handler configurations are added here
	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		playerTurnSheetHandlerConfig,
		playerJoinGameHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		cfg, err := fn(l)
		if err != nil {
			return nil, err
		}
		playerConfig = server.MergeHandlerConfigs(playerConfig, cfg)
	}

	return playerConfig, nil
}
