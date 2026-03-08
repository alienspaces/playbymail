package catalog

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	GetCatalogGameInstances = "get-catalog-game-instances"
)

func catalogGameInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "catalogGameInstanceHandlerConfig")

	l.Debug("Adding catalog game instance handler configuration")

	config := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "catalog_game_instance.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "catalog_game_instance.schema.json",
			},
		}...),
	}

	config[GetCatalogGameInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/catalog/game-instances",
		HandlerFunc: getCatalogGameInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypePublic},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Collection:  true,
			Title:       "Get catalog game instances",
			Description: "Returns all game instances open for player enrollment. No authentication required.",
		},
	}

	return config, nil
}

func getCatalogGameInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getCatalogGameInstancesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyCatalogGameInstanceViewRecs(opts)
	if err != nil {
		l.Warn("failed getting catalog game instance view records >%v<", err)
		return err
	}

	l.Info("mapping >%d< catalog game instance view records for response", len(recs))

	res, err := mapper.CatalogGameInstanceViewRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
