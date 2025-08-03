package game

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
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

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/game-instances/{game_instance_id}/configurations
// GET (document)    /api/v1/game-instances/{game_instance_id}/configurations/{configuration_id}
// POST (document)   /api/v1/game-instances/{game_instance_id}/configurations
// PUT (document)    /api/v1/game-instances/{game_instance_id}/configurations/{configuration_id}
// DELETE (document) /api/v1/game-instances/{game_instance_id}/configurations/{configuration_id}

const (
	getManyGameInstanceConfigurations  = "get-many-game-instance-configurations"
	getOneGameInstanceConfiguration    = "get-one-game-instance-configuration"
	createOneGameInstanceConfiguration = "create-one-game-instance-configuration"
	updateOneGameInstanceConfiguration = "update-one-game-instance-configuration"
	deleteOneGameInstanceConfiguration = "delete-one-game-instance-configuration"
)

func gameInstanceConfigurationHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameInstanceConfigurationHandlerConfig")

	l.Debug("Adding game instance configuration handler configuration")

	gameInstanceConfigurationConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance_configuration.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_instance_configuration.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance_configuration.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance_configuration.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_instance_configuration.schema.json",
			},
		}...),
	}

	gameInstanceConfigurationConfig[getManyGameInstanceConfigurations] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-instances/:game_instance_id/configurations",
		HandlerFunc: getManyGameInstanceConfigurationsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game instance configurations",
		},
	}

	gameInstanceConfigurationConfig[getOneGameInstanceConfiguration] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-instances/:game_instance_id/configurations/:configuration_id",
		HandlerFunc: getOneGameInstanceConfigurationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game instance configuration",
		},
	}

	gameInstanceConfigurationConfig[createOneGameInstanceConfiguration] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-instances/:game_instance_id/configurations",
		HandlerFunc: createOneGameInstanceConfigurationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game instance configuration",
		},
	}

	gameInstanceConfigurationConfig[updateOneGameInstanceConfiguration] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/game-instances/:game_instance_id/configurations/:configuration_id",
		HandlerFunc: updateOneGameInstanceConfigurationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game instance configuration",
		},
	}

	gameInstanceConfigurationConfig[deleteOneGameInstanceConfiguration] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/game-instances/:game_instance_id/configurations/:configuration_id",
		HandlerFunc: deleteOneGameInstanceConfigurationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game instance configuration",
		},
	}

	return gameInstanceConfigurationConfig, nil
}

func getManyGameInstanceConfigurationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameInstanceConfigurationsHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game_instance_id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	mm := m.(*domain.Domain)
	recs, err := mm.GetGameInstanceConfigurationsByGameInstanceID(gameInstanceID)
	if err != nil {
		l.Warn("failed getting game instance configuration records >%v<", err)
		return err
	}

	res, err := mapper.GameInstanceConfigurationRecsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneGameInstanceConfigurationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneGameInstanceConfigurationHandler")

	configurationID := pp.ByName("configuration_id")
	if configurationID == "" {
		l.Warn("configuration_id is required")
		return coreerror.NewNotFoundError("game instance configuration", configurationID)
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetGameInstanceConfigurationRec(configurationID, nil)
	if err != nil {
		l.Warn("failed getting game instance configuration record >%v<", err)
		return err
	}

	res, err := mapper.GameInstanceConfigurationRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneGameInstanceConfigurationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneGameInstanceConfigurationHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game_instance_id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	var request game_record.GameInstanceConfiguration
	_, err := server.ReadRequest(l, r, &request)
	if err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	// Set the game instance ID from the URL parameter
	request.GameInstanceID = gameInstanceID

	mm := m.(*domain.Domain)
	if err := mm.ValidateGameInstanceConfiguration(&request); err != nil {
		l.Warn("validation failed >%v<", err)
		return err
	}

	rec, err := mm.CreateGameInstanceConfigurationRec(&request)
	if err != nil {
		l.Warn("failed creating game instance configuration record >%v<", err)
		return err
	}

	res, err := mapper.GameInstanceConfigurationRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneGameInstanceConfigurationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneGameInstanceConfigurationHandler")

	configurationID := pp.ByName("configuration_id")
	if configurationID == "" {
		l.Warn("configuration_id is required")
		return coreerror.NewNotFoundError("game instance configuration", configurationID)
	}

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game_instance_id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	var request game_record.GameInstanceConfiguration
	_, err := server.ReadRequest(l, r, &request)
	if err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	// Set the IDs from the URL parameters
	request.ID = configurationID
	request.GameInstanceID = gameInstanceID

	mm := m.(*domain.Domain)
	if err := mm.ValidateGameInstanceConfiguration(&request); err != nil {
		l.Warn("validation failed >%v<", err)
		return err
	}

	rec, err := mm.UpdateGameInstanceConfigurationRec(&request)
	if err != nil {
		l.Warn("failed updating game instance configuration record >%v<", err)
		return err
	}

	res, err := mapper.GameInstanceConfigurationRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneGameInstanceConfigurationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneGameInstanceConfigurationHandler")

	configurationID := pp.ByName("configuration_id")
	if configurationID == "" {
		l.Warn("configuration_id is required")
		return coreerror.NewNotFoundError("game instance configuration", configurationID)
	}

	mm := m.(*domain.Domain)
	if err := mm.DeleteGameInstanceConfigurationRec(configurationID); err != nil {
		l.Warn("failed deleting game instance configuration record >%v<", err)
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
