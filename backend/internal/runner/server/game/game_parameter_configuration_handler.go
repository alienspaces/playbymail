package game

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
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// API Resource Search Path
//
// GET (collection) /api/v1/game-parameter-configurations

const (
	GetManyGameParameterConfigurations = "get-many-game-parameter-configurations"
)

func gameParameterConfigurationHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameParameterConfigurationHandlerConfig")

	l.Debug("Adding game parameter configuration handler configuration")

	gameParameterConfigurationConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_parameter_configuration.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_parameter_configuration.schema.json",
			},
		}...),
	}

	// Collection endpoint
	gameParameterConfigurationConfig[GetManyGameParameterConfigurations] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-parameter-configurations",
		HandlerFunc: getManyGameParameterConfigurationsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game parameter configuration collection",
		},
	}

	return gameParameterConfigurationConfig, nil
}

func getManyGameParameterConfigurationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameParameterConfigurationsHandler")

	l.Info("querying many game parameter configuration records with params >%#v<", qp)

	// Get game_type query parameter if provided
	var gameType string
	if gameTypeParams, exists := qp.Params["game_type"]; exists && len(gameTypeParams) > 0 {
		if val, ok := gameTypeParams[0].Val.(string); ok {
			gameType = val
		}
	}

	var recs []*game_record.GameParameter
	if gameType != "" {
		l.Info("filtering configurations by game_type >%s<", gameType)
		recs = domain.GetGameParameterConfigurationsByGameType(gameType)
	} else {
		l.Info("getting all configurations")
		recs = domain.GetGameParameterConfigurations()
	}

	res, err := mapper.GameParameterConfigurationRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	l.Info("responding with >%d< game parameter configuration records", len(res.Data))

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
