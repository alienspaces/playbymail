package mech_wargame

import (
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	packageName = "mech_wargame"
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

func MechWargameHandlerConfig(cfg config.Config, l logger.Logger, scnr turnsheet.TurnSheetScanner) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "MechWargameHandlerConfig")

	l.Debug("Adding mech_wargame handler configuration")

	mechWargameConfig := make(map[string]server.HandlerConfig)

	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		mechWargameChassisHandlerConfig,
		mechWargameWeaponHandlerConfig,
		mechWargameSectorHandlerConfig,
		mechWargameSectorLinkHandlerConfig,
		mechWargameLanceHandlerConfig,
		mechWargameLanceMechHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		cfg, err := fn(l)
		if err != nil {
			return nil, err
		}
		mechWargameConfig = server.MergeHandlerConfigs(mechWargameConfig, cfg)
	}

	return mechWargameConfig, nil
}
