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

// API Resource Search Path
//
// GET (collection) /api/v1/game-instances
//
// API Resource CRUD Paths
//
// GET (collection)  /api/v1/games/{game_id}/instances
// GET (document)    /api/v1/games/{game_id}/instances/{instance_id}
// POST (document)   /api/v1/games/{game_id}/instances
// PUT (document)    /api/v1/games/{game_id}/instances/{instance_id}
// DELETE (document) /api/v1/games/{game_id}/instances/{instance_id}

const (
	searchManyGameInstances = "search-many-game-instances"
	getManyGameInstances    = "get-many-game-instances"
	getOneGameInstance      = "get-one-game-instance"
	createOneGameInstance   = "create-one-game-instance"
	updateOneGameInstance   = "update-one-game-instance"
	deleteOneGameInstance   = "delete-one-game-instance"
	startGameInstance       = "start-game-instance"
	pauseGameInstance       = "pause-game-instance"
	resumeGameInstance      = "resume-game-instance"
	cancelGameInstance      = "cancel-game-instance"
)

func gameInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameInstanceHandlerConfig")

	l.Debug("Adding game instance handler configuration")

	gameInstanceConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_instance.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_instance.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_instance.schema.json",
			},
		}...),
	}

	gameInstanceConfig[searchManyGameInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-instances",
		HandlerFunc: searchManyGameInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search game instances",
		},
	}

	gameInstanceConfig[getManyGameInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/instances",
		HandlerFunc: getManyGameInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game instances",
		},
	}

	gameInstanceConfig[getOneGameInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id/instances/:instance_id",
		HandlerFunc: getOneGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game instance",
		},
	}

	gameInstanceConfig[createOneGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances",
		HandlerFunc: createOneGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game instance",
		},
	}

	gameInstanceConfig[updateOneGameInstance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/games/:game_id/instances/:instance_id",
		HandlerFunc: updateOneGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game instance",
		},
	}

	gameInstanceConfig[deleteOneGameInstance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/games/:game_id/instances/:instance_id",
		HandlerFunc: deleteOneGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game instance",
		},
	}

	// Runtime management endpoints
	gameInstanceConfig[startGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/start",
		HandlerFunc: startGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Start game instance",
		},
	}

	gameInstanceConfig[pauseGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/pause",
		HandlerFunc: pauseGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Pause game instance",
		},
	}

	gameInstanceConfig[resumeGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/resume",
		HandlerFunc: resumeGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Resume game instance",
		},
	}

	gameInstanceConfig[cancelGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/instances/:instance_id/cancel",
		HandlerFunc: cancelGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Cancel game instance",
		},
	}

	return gameInstanceConfig, nil
}

func searchManyGameInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyGameInstancesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GameInstanceRepository().GetMany(opts)
	if err != nil {
		l.Warn("failed getting game instance records >%v<", err)
		return err
	}

	res, err := mapper.GameInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyGameInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameInstancesHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{
		Col: game_record.FieldGameInstanceGameID,
		Val: gameID,
	})

	recs, err := mm.GameInstanceRepository().GetMany(opts)
	if err != nil {
		l.Warn("failed getting game instance records >%v<", err)
		return err
	}

	res, err := mapper.GameInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneGameInstanceHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed getting game instance record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		l.Warn("instance does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	res, err := mapper.GameInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneGameInstanceHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	// Create record with just the GameID - domain layer will set defaults
	rec := &game_record.GameInstance{
		GameID: gameID,
	}

	mm := m.(*domain.Domain)
	rec, err := mm.CreateGameInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game instance record >%v<", err)
		return err
	}

	res, err := mapper.GameInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneGameInstanceHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("instance does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	rec, err = mm.UpdateGameInstanceRec(rec)
	if err != nil {
		l.Warn("failed updating game instance record >%v<", err)
		return err
	}

	res, err := mapper.GameInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneGameInstanceHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetGameInstanceRec(instanceID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("instance does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	if err := mm.DeleteGameInstanceRec(instanceID); err != nil {
		l.Warn("failed deleting game instance record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// startGameInstanceHandler starts a game instance
func startGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "startGameInstanceHandler")

	l.Info("starting game instance")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Start the game instance
	instance, err := mm.StartGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to start game instance >%v<", err)
		return err
	}

	// Convert to response
	res, err := mapper.GameInstanceRecordToResponse(l, instance)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// pauseGameInstanceHandler pauses a game instance
func pauseGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "pauseGameInstanceHandler")

	l.Info("pausing game instance")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Pause the game instance
	instance, err := mm.PauseGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to pause game instance >%v<", err)
		return err
	}

	// Convert to response
	res, err := mapper.GameInstanceRecordToResponse(l, instance)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// resumeGameInstanceHandler resumes a game instance
func resumeGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "resumeGameInstanceHandler")

	l.Info("resuming game instance")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Resume the game instance
	instance, err := mm.ResumeGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to resume game instance >%v<", err)
		return err
	}

	// Convert to response
	res, err := mapper.GameInstanceRecordToResponse(l, instance)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// cancelGameInstanceHandler cancels a game instance
func cancelGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "cancelGameInstanceHandler")

	l.Info("canceling game instance")

	// Get path parameters
	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")

	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)

	// Cancel the game instance
	instance, err := mm.CancelGameInstance(instanceID)
	if err != nil {
		l.Warn("failed to cancel game instance >%v<", err)
		return err
	}

	// Convert to response
	res, err := mapper.GameInstanceRecordToResponse(l, instance)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
