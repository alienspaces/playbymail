package adventure_game

import (
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	packageName = "adventure_game"
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

func AdventureGameHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "AdventureGameHandlerConfig")

	l.Debug("Adding adventure_game handler configuration")

	adventureGameConfig := make(map[string]server.HandlerConfig)

	// Additional handler configurations are added here
	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		adventureGameCharacterHandlerConfig,
		adventureGameCreatureHandlerConfig,
		adventureGameItemHandlerConfig,
		adventureGameItemPlacementHandlerConfig,
		adventureGameCreaturePlacementHandlerConfig,
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
		adventureGameConfig = server.MergeHandlerConfigs(adventureGameConfig, cfg)
	}

	return adventureGameConfig, nil
}
