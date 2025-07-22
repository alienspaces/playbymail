package adventure_game

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
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
	"gitlab.com/alienspaces/playbymail/schema"
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-location-instances

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-game-instances/{game_instance_id}/location-instances
// GET (document)    /api/v1/adventure-game-instances/{game_instance_id}/location-instances/{location_instance_id}
// POST (document)   /api/v1/adventure-game-instances/{game_instance_id}/location-instances
// PUT (document)    /api/v1/adventure-game-instances/{game_instance_id}/location-instances/{location_instance_id}
// DELETE (document) /api/v1/adventure-game-instances/{game_instance_id}/location-instances/{location_instance_id}

const (
	// API Resource Search Path
	searchManyAdventureGameLocationInstances = "search-many-adventure-game-location-instances"

	// API Resource CRUD Paths
	getManyAdventureGameLocationInstances  = "get-many-adventure-game-location-instances"
	getOneAdventureGameLocationInstance    = "get-one-adventure-game-location-instance"
	createOneAdventureGameLocationInstance = "create-one-adventure-game-location-instance"
	updateOneAdventureGameLocationInstance = "update-one-adventure-game-location-instance"
	deleteOneAdventureGameLocationInstance = "delete-one-adventure-game-location-instance"
)

func adventureGameLocationInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameLocationInstanceHandlerConfig")

	l.Debug("Adding adventure_game_location_instance handler configuration")

	gameLocationInstanceConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_location_instance.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_location_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_location_instance.response.schema.json",
		},
		References: referenceSchemas,
	}

	// New Adventure Game Location Instance API paths
	gameLocationInstanceConfig[searchManyAdventureGameLocationInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-location-instances",
		HandlerFunc: searchManyAdventureGameLocationInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game location instances",
		},
	}

	gameLocationInstanceConfig[getManyAdventureGameLocationInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/location-instances",
		HandlerFunc: getManyAdventureGameLocationInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game location instances",
		},
	}

	gameLocationInstanceConfig[getOneAdventureGameLocationInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/location-instances/:location_instance_id",
		HandlerFunc: getOneAdventureGameLocationInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game location instance",
		},
	}

	gameLocationInstanceConfig[createOneAdventureGameLocationInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/location-instances",
		HandlerFunc: createOneAdventureGameLocationInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game location instance",
		},
	}

	gameLocationInstanceConfig[updateOneAdventureGameLocationInstance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/location-instances/:location_instance_id",
		HandlerFunc: updateOneAdventureGameLocationInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game location instance",
		},
	}

	gameLocationInstanceConfig[deleteOneAdventureGameLocationInstance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/location-instances/:location_instance_id",
		HandlerFunc: deleteOneAdventureGameLocationInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game location instance",
		},
	}

	return gameLocationInstanceConfig, nil
}

func searchManyAdventureGameLocationInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameLocationInstancesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameLocationInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location instance records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameLocationInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameLocationInstancesHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: record.FieldAdventureGameLocationInstanceAdventureGameInstanceID,
		Val: gameInstanceID,
	})

	mm := m.(*domain.Domain)

	recs, err := mm.GetManyAdventureGameLocationInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location instance records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameLocationInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameLocationInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	locationInstanceID := pp.ByName("location_instance_id")
	if locationInstanceID == "" {
		l.Warn("location instance id is required")
		return coreerror.NewNotFoundError("location instance", locationInstanceID)
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetAdventureGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameLocationInstanceRec(locationInstanceID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location instance record >%v<", err)
		return err
	}

	// Verify the location instance belongs to the specified game
	if rec.AdventureGameInstanceID != gameInstanceRec.ID {
		l.Warn("location instance does not belong to specified game instance >%s< != >%s<", rec.AdventureGameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("location instance", locationInstanceID)
	}

	res, err := mapper.AdventureGameLocationInstanceRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game location instance record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameLocationInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameLocationInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	var req schema.AdventureGameLocationInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.AdventureGameLocationInstanceRequestToRecord(l, &req, &record.AdventureGameLocationInstance{})
	if err != nil {
		return err
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetAdventureGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	// Set the game ID and game instance ID from game instance record
	rec.GameID = gameInstanceRec.GameID
	rec.AdventureGameInstanceID = gameInstanceRec.ID

	l.Info("creating adventure game location instance record >%s<", spew.Sdump(rec))

	rec, err = mm.CreateAdventureGameLocationInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location instance record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameLocationInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameLocationInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	locationInstanceID := pp.ByName("location_instance_id")
	if locationInstanceID == "" {
		l.Warn("location instance id is required")
		return coreerror.NewNotFoundError("location instance", locationInstanceID)
	}

	l.Info("updating adventure game location instance record with path params >%#v<", pp)

	var req schema.AdventureGameLocationInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetAdventureGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameLocationInstanceRec(locationInstanceID, nil)
	if err != nil {
		return err
	}

	// Verify the location instance belongs to the specified game
	if rec.AdventureGameInstanceID != gameInstanceRec.ID {
		l.Warn("location instance does not belong to specified game instance >%s< != >%s<", rec.AdventureGameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("location instance", locationInstanceID)
	}

	rec, err = mapper.AdventureGameLocationInstanceRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameLocationInstanceRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game location instance record >%v<", err)
		return err
	}

	data, err := mapper.AdventureGameLocationInstanceRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.AdventureGameLocationInstanceResponse{
		Data: &data,
	}

	l.Info("responding with updated adventure game location instance record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameLocationInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameLocationInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	locationInstanceID := pp.ByName("location_instance_id")
	if locationInstanceID == "" {
		l.Warn("location instance id is required")
		return coreerror.NewNotFoundError("location instance", locationInstanceID)
	}

	l.Info("deleting adventure game location instance record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetAdventureGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	// First get the record to verify it belongs to the specified game
	rec, err := mm.GetAdventureGameLocationInstanceRec(locationInstanceID, nil)
	if err != nil {
		return err
	}

	// Verify the location instance belongs to the specified game
	if rec.AdventureGameInstanceID != gameInstanceRec.ID {
		l.Warn("location instance does not belong to specified game instance >%s< != >%s<", rec.AdventureGameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("location instance", locationInstanceID)
	}

	if err := mm.DeleteAdventureGameLocationInstanceRec(locationInstanceID); err != nil {
		l.Warn("failed deleting adventure game location instance record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
