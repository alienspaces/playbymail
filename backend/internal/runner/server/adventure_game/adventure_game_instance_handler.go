package adventure_game

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
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-instances
//
// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/instances
// GET (document)    /api/v1/adventure-games/{game_id}/instances/{instance_id}
// POST (document)   /api/v1/adventure-games/{game_id}/instances
// PUT (document)    /api/v1/adventure-games/{game_id}/instances/{instance_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/instances/{instance_id}

const (
	searchManyAdventureGameInstances = "search-many-adventure-game-instances"
	getManyAdventureGameInstances    = "get-many-adventure-game-instances"
	getOneAdventureGameInstance      = "get-one-adventure-game-instance"
	createOneAdventureGameInstance   = "create-one-adventure-game-instance"
	updateOneAdventureGameInstance   = "update-one-adventure-game-instance"
	deleteOneAdventureGameInstance   = "delete-one-adventure-game-instance"
)

func adventureGameInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameInstanceHandlerConfig")

	l.Debug("Adding adventure_game_instance handler configuration")

	gameInstanceConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_instance.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_instance.response.schema.json",
		},
		References: referenceSchemas,
	}

	gameInstanceConfig[searchManyAdventureGameInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances",
		HandlerFunc: searchManyAdventureGameInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game instances",
		},
	}

	gameInstanceConfig[getManyAdventureGameInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/instances",
		HandlerFunc: getManyAdventureGameInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game instances",
		},
	}

	gameInstanceConfig[getOneAdventureGameInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/instances/:instance_id",
		HandlerFunc: getOneAdventureGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game instance",
		},
	}

	gameInstanceConfig[createOneAdventureGameInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/instances",
		HandlerFunc: createOneAdventureGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game instance",
		},
	}

	gameInstanceConfig[updateOneAdventureGameInstance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/instances/:instance_id",
		HandlerFunc: updateOneAdventureGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game instance",
		},
	}

	gameInstanceConfig[deleteOneAdventureGameInstance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/instances/:instance_id",
		HandlerFunc: deleteOneAdventureGameInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game instance",
		},
	}

	return gameInstanceConfig, nil
}

func searchManyAdventureGameInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameInstancesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.AdventureGameInstanceRepository().GetMany(opts)
	if err != nil {
		l.Warn("failed getting adventure game instance records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameInstancesHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{
		Col: record.FieldAdventureGameInstanceGameID,
		Val: gameID,
	})

	recs, err := mm.AdventureGameInstanceRepository().GetMany(opts)
	if err != nil {
		l.Warn("failed getting adventure game instance records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameInstanceHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetAdventureGameInstanceRec(instanceID, nil)
	if err != nil {
		l.Warn("failed getting adventure game instance record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		l.Warn("instance does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	res, err := mapper.AdventureGameInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameInstanceHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	rec := &record.AdventureGameInstance{GameID: gameID}
	mm := m.(*domain.Domain)
	rec, err := mm.CreateAdventureGameInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game instance record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameInstanceHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetAdventureGameInstanceRec(instanceID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("instance does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	rec, err = mm.UpdateAdventureGameInstanceRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game instance record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameInstanceHandler")

	gameID := pp.ByName("game_id")
	instanceID := pp.ByName("instance_id")
	if gameID == "" || instanceID == "" {
		l.Warn("game id and instance id are required")
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	mm := m.(*domain.Domain)
	rec, err := mm.GetAdventureGameInstanceRec(instanceID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("instance does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("game instance", instanceID)
	}

	if err := mm.DeleteAdventureGameInstanceRec(instanceID); err != nil {
		l.Warn("failed deleting adventure game instance record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
