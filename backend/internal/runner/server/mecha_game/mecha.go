package mecha_game

import (
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	packageName = "mecha"
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

func MechaHandlerConfig(cfg config.Config, l logger.Logger, scnr turnsheet.TurnSheetScanner) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "MechaHandlerConfig")

	l.Debug("Adding mecha handler configuration")

	mechaConfig := make(map[string]server.HandlerConfig)

	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		mechaChassisHandlerConfig,
		mechaWeaponHandlerConfig,
		mechaSectorHandlerConfig,
		mechaSectorLinkHandlerConfig,
		mechaLanceHandlerConfig,
		mechaLanceMechHandlerConfig,
		mechaComputerOpponentHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		cfg, err := fn(l)
		if err != nil {
			return nil, err
		}
		mechaConfig = server.MergeHandlerConfigs(mechaConfig, cfg)
	}

	return mechaConfig, nil
}
