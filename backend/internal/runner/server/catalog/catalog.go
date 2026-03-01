package catalog

import (
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	packageName = "catalog"
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

func CatalogHandlerConfig(cfg config.Config, l logger.Logger, scnr turnsheet.TurnSheetScanner) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "CatalogHandlerConfig")

	l.Debug("Adding catalog handler configuration")

	catalogConfig := make(map[string]server.HandlerConfig)

	handlerConfigFuncs := []func(logger.Logger) (map[string]server.HandlerConfig, error){
		catalogGameHandlerConfig,
	}

	for _, fn := range handlerConfigFuncs {
		cfg, err := fn(l)
		if err != nil {
			return nil, err
		}
		catalogConfig = server.MergeHandlerConfigs(catalogConfig, cfg)
	}

	return catalogConfig, nil
}
