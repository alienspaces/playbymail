package game

import (
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	packageName = "game"
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

func GameHandlerConfig(cfg config.Config, l logger.Logger, scnr turn_sheet.TurnSheetScanner) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "GameHandlerConfig")

	l.Debug("Adding game handler configuration")

	gameConfig := make(map[string]server.HandlerConfig)

	// Additional handler configurations are added here
	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		gameHandlerConfig,
		gameImageHandlerConfig,
		gameParameterHandlerConfig,
		gameSubscriptionHandlerConfig,
		gameInstanceHandlerConfig,
		gameInstanceParameterHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		cfg, err := fn(l)
		if err != nil {
			return nil, err
		}
		gameConfig = server.MergeHandlerConfigs(gameConfig, cfg)
	}

	turnSheetConfig, err := gameTurnSheetHandlerConfig(l, scnr)
	if err != nil {
		return nil, err
	}
	gameConfig = server.MergeHandlerConfigs(gameConfig, turnSheetConfig)

	return gameConfig, nil
}
