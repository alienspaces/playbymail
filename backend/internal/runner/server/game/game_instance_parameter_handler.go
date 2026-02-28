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
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/games/{game_id}/instances/{instance_id}/parameters
// GET (document)    /api/v1/games/{game_id}/instances/{instance_id}/parameters/{parameter_id}
// POST (document)   /api/v1/games/{game_id}/instances/{instance_id}/parameters
// PUT (document)    /api/v1/games/{game_id}/instances/{instance_id}/parameters/{parameter_id}
// DELETE (document) /api/v1/games/{game_id}/instances/{instance_id}/parameters/{parameter_id}

const (
	GetManyGameInstanceParameters  = "get-many-game-instance-parameters"
	GetOneGameInstanceParameter    = "get-one-game-instance-parameter"
	CreateOneGameInstanceParameter = "create-one-game-instance-parameter"
	UpdateOneGameInstanceParameter = "update-one-game-instance-parameter"
	DeleteOneGameInstanceParameter = "delete-one-game-instance-parameter"
)

func gameInstanceParameterHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameInstanceParameterHandlerConfig")

	l.Debug("adding game instance parameter handler configuration")

	gameInstanceParameterConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance_parameter.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_instance_parameter.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance_parameter.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance_parameter.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_instance_parameter.schema.json",
			},
		}...),
	}

	gameInstanceParameterConfig[GetManyGameInstanceParameters] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/parameters",
		HandlerFunc: getManyGameInstanceParametersHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game instance parameter collection",
		},
	}

	gameInstanceParameterConfig[GetOneGameInstanceParameter] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/parameters/:parameter_id",
		HandlerFunc: getOneGameInstanceParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game instance parameter",
		},
	}

	gameInstanceParameterConfig[CreateOneGameInstanceParameter] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/parameters",
		HandlerFunc: createOneGameInstanceParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game instance parameter",
		},
	}

	gameInstanceParameterConfig[UpdateOneGameInstanceParameter] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/parameters/:parameter_id",
		HandlerFunc: updateOneGameInstanceParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game instance parameter",
		},
	}

	gameInstanceParameterConfig[DeleteOneGameInstanceParameter] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/parameters/:parameter_id",
		HandlerFunc: deleteOneGameInstanceParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game instance parameter",
		},
	}

	return gameInstanceParameterConfig, nil
}

func getManyGameInstanceParametersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameInstanceParametersHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	l.Info("getting many game instance parameters for game >%s< instance >%s<", gameID, instanceID)

	// Validate that the instance belongs to the specified game
	gameInstance, err := m.(*domain.Domain).GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed getting game instance >%v<", err)
		return err
	}

	if gameInstance.GameID != gameID {
		l.Warn("game instance >%s< does not belong to game >%s<", instanceID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{
		Col: game_record.FieldGameInstanceParameterGameInstanceID,
		Val: instanceID,
	})

	recs, err := m.(*domain.Domain).GetManyGameInstanceParameterRecs(opts)
	if err != nil {
		l.Warn("failed getting game instance parameters >%v<", err)
		return err
	}

	response, err := mapper.GameInstanceParameterRecsToCollectionResponse(l, recs)
	if err != nil {
		l.Warn("failed mapping game instance parameter records to collection response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, response)
}

func getOneGameInstanceParameterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneGameInstanceParameterHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	parameterID := pp.ByName("parameter_id")

	l.Info("getting game instance parameter >%s< for game >%s< instance >%s<", parameterID, gameID, instanceID)

	// Validate that the instance belongs to the specified game
	gameInstance, err := m.(*domain.Domain).GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed getting game instance >%v<", err)
		return err
	}

	if gameInstance.GameID != gameID {
		l.Warn("game instance >%s< does not belong to game >%s<", instanceID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	rec, err := m.(*domain.Domain).GetGameInstanceParameterRec(parameterID, nil)
	if err != nil {
		l.Warn("failed getting game instance parameter >%v<", err)
		return err
	}

	// Validate the parameter belongs to the game instance
	if rec.GameInstanceID != instanceID {
		l.Warn("parameter >%s< does not belong to game instance >%s<", parameterID, instanceID)
		return coreerror.NewNotFoundError("parameter", parameterID)
	}

	response, err := mapper.GameInstanceParameterRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game instance parameter record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, response)
}

func createOneGameInstanceParameterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneGameInstanceParameterHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	l.Info("creating game instance parameter for game >%s< instance >%s<", gameID, instanceID)

	// Validate that the instance belongs to the specified game
	gameInstance, err := m.(*domain.Domain).GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed getting game instance >%v<", err)
		return err
	}

	if gameInstance.GameID != gameID {
		l.Warn("game instance >%s< does not belong to game >%s<", instanceID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	rec := &game_record.GameInstanceParameter{
		GameInstanceID: instanceID,
	}
	rec, err = mapper.GameInstanceParameterRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping game instance parameter request to record >%v<", err)
		return err
	}

	// Ensure the game_instance_id matches the URL parameter
	rec.GameInstanceID = instanceID

	rec, err = m.(*domain.Domain).CreateGameInstanceParameterRec(rec)
	if err != nil {
		l.Warn("failed creating game instance parameter >%v<", err)
		return err
	}

	response, err := mapper.GameInstanceParameterRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game instance parameter record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, response)
}

func updateOneGameInstanceParameterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneGameInstanceParameterHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	parameterID := pp.ByName("parameter_id")

	l.Info("updating game instance parameter >%s< for game >%s< instance >%s<", parameterID, gameID, instanceID)

	// Validate that the instance belongs to the specified game
	gameInstance, err := m.(*domain.Domain).GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed getting game instance >%v<", err)
		return err
	}

	if gameInstance.GameID != gameID {
		l.Warn("game instance >%s< does not belong to game >%s<", instanceID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	rec, err := m.(*domain.Domain).GetGameInstanceParameterRec(parameterID, nil)
	if err != nil {
		l.Warn("failed getting game instance parameter >%v<", err)
		return err
	}

	// Validate the parameter belongs to the game instance
	if rec.GameInstanceID != instanceID {
		l.Warn("parameter >%s< does not belong to game instance >%s<", parameterID, instanceID)
		return coreerror.NewNotFoundError("parameter", parameterID)
	}

	rec, err = mapper.GameInstanceParameterRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping game instance parameter request to record >%v<", err)
		return err
	}

	// Ensure the game_instance_id doesn't change
	rec.GameInstanceID = instanceID

	rec, err = m.(*domain.Domain).UpdateGameInstanceParameterRec(rec)
	if err != nil {
		l.Warn("failed updating game instance parameter >%v<", err)
		return err
	}

	response, err := mapper.GameInstanceParameterRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game instance parameter record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, response)
}

func deleteOneGameInstanceParameterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneGameInstanceParameterHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	parameterID := pp.ByName("parameter_id")

	l.Info("deleting game instance parameter >%s< for game >%s< instance >%s<", parameterID, gameID, instanceID)

	// Validate that the instance belongs to the specified game
	gameInstance, err := m.(*domain.Domain).GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed getting game instance >%v<", err)
		return err
	}

	if gameInstance.GameID != gameID {
		l.Warn("game instance >%s< does not belong to game >%s<", instanceID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	rec, err := m.(*domain.Domain).GetGameInstanceParameterRec(parameterID, nil)
	if err != nil {
		l.Warn("failed getting game instance parameter >%v<", err)
		return err
	}

	// Validate the parameter belongs to the game instance
	if rec.GameInstanceID != instanceID {
		l.Warn("parameter >%s< does not belong to game instance >%s<", parameterID, instanceID)
		return coreerror.NewNotFoundError("parameter", parameterID)
	}

	err = m.(*domain.Domain).DeleteGameInstanceParameterRec(parameterID)
	if err != nil {
		l.Warn("failed deleting game instance parameter >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
