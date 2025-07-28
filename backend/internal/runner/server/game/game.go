package game

import (
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	packageName = "game"
)

var referenceSchemas = []jsonschema.Schema{
	{
		Location: "api",
		Name:     "query.schema.json",
	},
	{
		Location: "api",
		Name:     "common.schema.json",
	},
}

func GameHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "GameHandlerConfig")

	l.Debug("Adding game handler configuration")

	gameConfig := make(map[string]server.HandlerConfig)

	// Additional handler configurations are added here
	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		gameHandlerConfig,
		gameConfigurationHandlerConfig,
		gameSubscriptionHandlerConfig,
		gameAdministrationHandlerConfig,
		gameInstanceHandlerConfig,
		gameInstanceConfigurationHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		cfg, err := fn(l)
		if err != nil {
			return nil, err
		}
		gameConfig = server.MergeHandlerConfigs(gameConfig, cfg)
	}

	return gameConfig, nil
}
