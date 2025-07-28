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
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

const (
	GetManyGameConfigurations     = "get-game-configurations"
	GetOneGameConfiguration       = "get-game-configuration"
	CreateGameConfiguration       = "create-game-configuration"
	CreateGameConfigurationWithID = "create-game-configuration-with-id"
	UpdateGameConfiguration       = "update-game-configuration"
	DeleteGameConfiguration       = "delete-game-configuration"
)

func gameConfigurationHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameConfigurationHandlerConfig")

	l.Debug("Adding game configuration handler configuration")

	// Create a new map to avoid modifying the passed config
	gameConfigurationConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "game_configuration.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "game_configuration.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "game_configuration.response.schema.json",
		},
		References: referenceSchemas,
	}

	// Unnested routes
	gameConfigurationConfig[GetManyGameConfigurations] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-configurations",
		HandlerFunc: getManyGameConfigurationsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game configuration collection",
		},
	}

	gameConfigurationConfig[GetOneGameConfiguration] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-configurations/:game_configuration_id",
		HandlerFunc: getGameConfigurationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game configuration",
		},
	}

	gameConfigurationConfig[CreateGameConfiguration] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-configurations",
		HandlerFunc: createGameConfigurationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game configuration",
		},
	}

	gameConfigurationConfig[UpdateGameConfiguration] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/game-configurations/:game_configuration_id",
		HandlerFunc: updateGameConfigurationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game configuration",
		},
	}

	gameConfigurationConfig[DeleteGameConfiguration] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/game-configurations/:game_configuration_id",
		HandlerFunc: deleteGameConfigurationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game configuration",
		},
	}

	return gameConfigurationConfig, nil
}

func getManyGameConfigurationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameConfigurationsHandler")

	l.Info("getting many game configurations")

	opts := queryparam.ToSQLOptions(qp)

	recs, err := m.(*domain.Domain).GetGameConfigurationRecs(opts)
	if err != nil {
		l.Warn("failed getting game configurations >%v<", err)
		return err
	}

	response := mapper.MapGameConfigurationCollectionResponse(recs)

	return server.WriteResponse(l, w, http.StatusOK, response)
}

func getGameConfigurationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameConfigurationHandler")

	gameConfigurationID := pp.ByName("game_configuration_id")

	l.Info("getting game configuration with id >%s<", gameConfigurationID)

	rec, err := m.(*domain.Domain).GetGameConfigurationRec(gameConfigurationID, nil)
	if err != nil {
		l.Warn("failed getting game configuration >%v<", err)
		return err
	}

	response := mapper.MapGameConfigurationResponse(rec)

	return server.WriteResponse(l, w, http.StatusOK, response)
}

func createGameConfigurationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createGameConfigurationHandler")

	l.Info("creating game configuration")

	var request api.GameConfigurationRequest
	_, err := server.ReadRequest(l, r, &request)
	if err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec := mapper.MapGameConfigurationRequestToRecord(&request)

	rec, err = m.(*domain.Domain).CreateGameConfigurationRec(rec)
	if err != nil {
		l.Warn("failed creating game configuration >%v<", err)
		return err
	}

	response := mapper.MapGameConfigurationResponse(rec)

	return server.WriteResponse(l, w, http.StatusCreated, response)
}

func updateGameConfigurationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateGameConfigurationHandler")

	gameConfigurationID := pp.ByName("game_configuration_id")

	l.Info("updating game configuration with id >%s<", gameConfigurationID)

	var request api.GameConfigurationRequest
	_, err := server.ReadRequest(l, r, &request)
	if err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec := mapper.MapGameConfigurationRequestToRecord(&request)
	rec.ID = gameConfigurationID

	rec, err = m.(*domain.Domain).UpdateGameConfigurationRec(rec)
	if err != nil {
		l.Warn("failed updating game configuration >%v<", err)
		return err
	}

	response := mapper.MapGameConfigurationResponse(rec)

	return server.WriteResponse(l, w, http.StatusOK, response)
}

func deleteGameConfigurationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteGameConfigurationHandler")

	gameConfigurationID := pp.ByName("game_configuration_id")

	l.Info("deleting game configuration with id >%s<", gameConfigurationID)

	err := m.(*domain.Domain).DeleteGameConfigurationRec(gameConfigurationID)
	if err != nil {
		l.Warn("failed deleting game configuration >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
