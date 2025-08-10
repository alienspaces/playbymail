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
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/game-instances/{game_instance_id}/parameters
// GET (document)    /api/v1/game-instances/{game_instance_id}/parameters/{parameter_id}
// POST (document)   /api/v1/game-instances/{game_instance_id}/parameters
// PUT (document)    /api/v1/game-instances/{game_instance_id}/parameters/{parameter_id}
// DELETE (document) /api/v1/game-instances/{game_instance_id}/parameters/{parameter_id}

const (
	getManyGameInstanceParameters  = "get-many-game-instance-parameters"
	getOneGameInstanceParameter    = "get-one-game-instance-parameter"
	createOneGameInstanceParameter = "create-one-game-instance-parameter"
	updateOneGameInstanceParameter = "update-one-game-instance-parameter"
	deleteOneGameInstanceParameter = "delete-one-game-instance-parameter"
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

	gameInstanceParameterConfig[getManyGameInstanceParameters] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-instances/:game_instance_id/parameters",
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

	gameInstanceParameterConfig[getOneGameInstanceParameter] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-instances/:game_instance_id/parameters/:parameter_id",
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

	gameInstanceParameterConfig[createOneGameInstanceParameter] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-instances/:game_instance_id/parameters",
		HandlerFunc: createOneGameInstanceParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game instance parameter",
		},
	}

	gameInstanceParameterConfig[updateOneGameInstanceParameter] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/game-instances/:game_instance_id/parameters/:parameter_id",
		HandlerFunc: updateOneGameInstanceParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game instance parameter",
		},
	}

	gameInstanceParameterConfig[deleteOneGameInstanceParameter] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/game-instances/:game_instance_id/parameters/:parameter_id",
		HandlerFunc: deleteOneGameInstanceParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
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

	gameInstanceID := pp.ByName("game_instance_id")

	l.Info("getting many game instance parameters for game instance >%s<", gameInstanceID)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{
		Col: game_record.FieldGameInstanceParameterGameInstanceID,
		Val: gameInstanceID,
	})

	recs, err := m.(*domain.Domain).GetGameInstanceParameterRecs(opts)
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

	gameInstanceID := pp.ByName("game_instance_id")
	parameterID := pp.ByName("parameter_id")

	l.Info("getting game instance parameter >%s< for game instance >%s<", parameterID, gameInstanceID)

	rec, err := m.(*domain.Domain).GetGameInstanceParameterRec(parameterID, nil)
	if err != nil {
		l.Warn("failed getting game instance parameter >%v<", err)
		return err
	}

	// Validate the parameter belongs to the game instance
	if rec.GameInstanceID != gameInstanceID {
		l.Warn("parameter >%s< does not belong to game instance >%s<", parameterID, gameInstanceID)
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

	gameInstanceID := pp.ByName("game_instance_id")

	l.Info("creating game instance parameter for game instance >%s<", gameInstanceID)

	rec := &game_record.GameInstanceParameter{}
	rec, err := mapper.GameInstanceParameterRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping game instance parameter request to record >%v<", err)
		return err
	}

	// Validate the game_instance_id matches the URL parameter
	if rec.GameInstanceID != gameInstanceID {
		l.Warn("game_instance_id in body >%s< does not match URL parameter >%s<", rec.GameInstanceID, gameInstanceID)
		return coreerror.NewInvalidError("game_instance_id", "game_instance_id mismatch")
	}

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

	gameInstanceID := pp.ByName("game_instance_id")
	parameterID := pp.ByName("parameter_id")

	l.Info("updating game instance parameter >%s< for game instance >%s<", parameterID, gameInstanceID)

	rec, err := m.(*domain.Domain).GetGameInstanceParameterRec(parameterID, nil)
	if err != nil {
		l.Warn("failed getting game instance parameter >%v<", err)
		return err
	}

	// Validate the parameter belongs to the game instance
	if rec.GameInstanceID != gameInstanceID {
		l.Warn("parameter >%s< does not belong to game instance >%s<", parameterID, gameInstanceID)
		return coreerror.NewNotFoundError("parameter", parameterID)
	}

	rec, err = mapper.GameInstanceParameterRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping game instance parameter request to record >%v<", err)
		return err
	}

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

	gameInstanceID := pp.ByName("game_instance_id")
	parameterID := pp.ByName("parameter_id")

	l.Info("deleting game instance parameter >%s< for game instance >%s<", parameterID, gameInstanceID)

	rec, err := m.(*domain.Domain).GetGameInstanceParameterRec(parameterID, nil)
	if err != nil {
		l.Warn("failed getting game instance parameter >%v<", err)
		return err
	}

	// Validate the parameter belongs to the game instance
	if rec.GameInstanceID != gameInstanceID {
		l.Warn("parameter >%s< does not belong to game instance >%s<", parameterID, gameInstanceID)
		return coreerror.NewNotFoundError("parameter", parameterID)
	}

	err = m.(*domain.Domain).DeleteGameInstanceParameterRec(parameterID)
	if err != nil {
		l.Warn("failed deleting game instance parameter >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
