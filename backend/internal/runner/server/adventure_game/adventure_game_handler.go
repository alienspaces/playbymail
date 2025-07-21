package adventure_game

import (
	"maps"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

var referenceSchemas = []jsonschema.Schema{
	{
		Location: "schema",
		Name:     "query.schema.json",
	},
	{
		Location: "schema",
		Name:     "common.schema.json",
	},
}

func AdventureGameHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "AdventureGameHandlerConfig")

	l.Debug("Adding adventure_game handler configuration")

	adventureGameConfig := make(map[string]server.HandlerConfig)

	// Additional handler configurations are added here
	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		adventureGameCharacterHandlerConfig,
		adventureGameCreatureHandlerConfig,
		adventureGameItemHandlerConfig,
		adventureGameLocationHandlerConfig,
		adventureGameLocationInstanceHandlerConfig,
		adventureGameLocationLinkHandlerConfig,
		adventureGameCreatureInstanceHandlerConfig,
		adventureGameItemInstanceHandlerConfig,
		adventureGameLocationLinkRequirementHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		cfg, err := fn(l)
		if err != nil {
			return nil, err
		}
		adventureGameConfig = mergeHandlerConfigs(adventureGameConfig, cfg)
	}

	return adventureGameConfig, nil
}

func mergeHandlerConfigs(hc1 map[string]server.HandlerConfig, hc2 map[string]server.HandlerConfig) map[string]server.HandlerConfig {
	if hc1 == nil {
		hc1 = map[string]server.HandlerConfig{}
	}
	maps.Copy(hc1, hc2)
	return hc1
}

// loggerWithFunctionContext provides a logger with function context
func loggerWithFunctionContext(l logger.Logger, functionName string) logger.Logger {
	return logging.LoggerWithFunctionContext(l, "adventure_game", functionName)
}
