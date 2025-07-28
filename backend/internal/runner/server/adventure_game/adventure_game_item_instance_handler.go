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
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-item-instances

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-game-instances/{game_instance_id}/item-instances
// GET (document)    /api/v1/adventure-game-instances/{game_instance_id}/item-instances/{item_instance_id}
// POST (document)   /api/v1/adventure-game-instances/{game_instance_id}/item-instances
// PUT (document)    /api/v1/adventure-game-instances/{game_instance_id}/item-instances/{item_instance_id}
// DELETE (document) /api/v1/adventure-game-instances/{game_instance_id}/item-instances/{item_instance_id}

const (
	// API Resource Search Path
	searchManyAdventureGameItemInstances = "search-many-adventure-game-item-instances"

	// API Resource CRUD Paths
	getManyAdventureGameItemInstances  = "get-many-adventure-game-item-instances"
	getOneAdventureGameItemInstance    = "get-one-adventure-game-item-instance"
	createOneAdventureGameItemInstance = "create-one-adventure-game-item-instance"
	updateOneAdventureGameItemInstance = "update-one-adventure-game-item-instance"
	deleteOneAdventureGameItemInstance = "delete-one-adventure-game-item-instance"
)

func adventureGameItemInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameItemInstanceHandlerConfig")

	l.Debug("Adding adventure_game_item_instance handler configuration")

	gameItemInstanceConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "adventure_game_item_instance.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "adventure_game_item_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "adventure_game_item_instance.response.schema.json",
		},
		References: referenceSchemas,
	}

	// New Adventure Game Item Instance API paths
	gameItemInstanceConfig[searchManyAdventureGameItemInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-item-instances",
		HandlerFunc: searchManyAdventureGameItemInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game item instances",
		},
	}

	gameItemInstanceConfig[getManyAdventureGameItemInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/item-instances",
		HandlerFunc: getManyAdventureGameItemInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game item instances",
		},
	}

	gameItemInstanceConfig[getOneAdventureGameItemInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/item-instances/:item_instance_id",
		HandlerFunc: getOneAdventureGameItemInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game item instance",
		},
	}

	gameItemInstanceConfig[createOneAdventureGameItemInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/item-instances",
		HandlerFunc: createOneAdventureGameItemInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game item instance",
		},
	}

	gameItemInstanceConfig[updateOneAdventureGameItemInstance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/item-instances/:item_instance_id",
		HandlerFunc: updateOneAdventureGameItemInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game item instance",
		},
	}

	gameItemInstanceConfig[deleteOneAdventureGameItemInstance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/item-instances/:item_instance_id",
		HandlerFunc: deleteOneAdventureGameItemInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game item instance",
		},
	}

	return gameItemInstanceConfig, nil
}

func searchManyAdventureGameItemInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameItemInstancesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameItemInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game item instance records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameItemInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameItemInstancesHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: adventure_game_record.FieldAdventureGameItemInstanceGameInstanceID,
		Val: gameInstanceID,
	})

	recs, err := mm.GetManyAdventureGameItemInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game item instance records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameItemInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameItemInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	itemInstanceID := pp.ByName("item_instance_id")
	if itemInstanceID == "" {
		l.Warn("item instance id is required")
		return coreerror.NewNotFoundError("item instance", itemInstanceID)
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameItemInstanceRec(itemInstanceID, nil)
	if err != nil {
		l.Warn("failed getting adventure game item instance record >%v<", err)
		return err
	}

	// Verify the item instance belongs to the specified game
	if rec.GameInstanceID != gameInstanceRec.ID {
		l.Warn("item instance does not belong to specified game instance >%s< != >%s<", rec.GameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("item instance", itemInstanceID)
	}

	res, err := mapper.AdventureGameItemInstanceRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game item instance record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameItemInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameItemInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	var req api.AdventureGameItemInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.AdventureGameItemInstanceRequestToRecord(l, &req, &adventure_game_record.AdventureGameItemInstance{})
	if err != nil {
		return err
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	// Set the game ID and game instance ID from game instance record
	rec.GameID = gameInstanceRec.GameID
	rec.GameInstanceID = gameInstanceRec.ID

	rec, err = mm.CreateAdventureGameItemInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game item instance record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameItemInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameItemInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	itemInstanceID := pp.ByName("item_instance_id")
	if itemInstanceID == "" {
		l.Warn("item instance id is required")
		return coreerror.NewNotFoundError("item instance", itemInstanceID)
	}

	l.Info("updating adventure game item instance record with path params >%#v<", pp)

	var req api.AdventureGameItemInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameItemInstanceRec(itemInstanceID, nil)
	if err != nil {
		return err
	}

	// Verify the item instance belongs to the specified game
	if rec.GameInstanceID != gameInstanceRec.ID {
		l.Warn("item instance does not belong to specified game instance >%s< != >%s<", rec.GameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("item instance", itemInstanceID)
	}

	rec, err = mapper.AdventureGameItemInstanceRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameItemInstanceRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game item instance record >%v<", err)
		return err
	}

	data, err := mapper.AdventureGameItemInstanceRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := api.AdventureGameItemInstanceResponse{
		Data: &data,
	}

	l.Info("responding with updated adventure game item instance record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameItemInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameItemInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	itemInstanceID := pp.ByName("item_instance_id")
	if itemInstanceID == "" {
		l.Warn("item instance id is required")
		return coreerror.NewNotFoundError("item instance", itemInstanceID)
	}

	l.Info("deleting adventure game item instance record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	// First get the record to verify it belongs to the specified game
	rec, err := mm.GetAdventureGameItemInstanceRec(itemInstanceID, nil)
	if err != nil {
		return err
	}

	// Verify the item instance belongs to the specified game
	if rec.GameInstanceID != gameInstanceRec.ID {
		l.Warn("item instance does not belong to specified game instance >%s< != >%s<", rec.GameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("item instance", itemInstanceID)
	}

	if err := mm.DeleteAdventureGameItemInstanceRec(itemInstanceID); err != nil {
		l.Warn("failed deleting adventure game item instance record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
