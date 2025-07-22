package adventure_game

import (
	"net/http"

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
// GET (collection) /api/v1/adventure-game-locations

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/locations
// GET (document)    /api/v1/adventure-games/{game_id}/locations/{location_id}
// POST (document)   /api/v1/adventure-games/{game_id}/locations
// PUT (document)    /api/v1/adventure-games/{game_id}/locations/{location_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/locations/{location_id}

const (
	// API Resource Search Path
	searchManyAdventureGameLocations = "search-many-adventure-game-locations"

	// API Resource CRUD Paths
	getManyAdventureGameLocations  = "get-many-adventure-game-locations"
	getOneAdventureGameLocation    = "get-one-adventure-game-location"
	createOneAdventureGameLocation = "create-one-adventure-game-location"
	updateOneAdventureGameLocation = "update-one-adventure-game-location"
	deleteOneAdventureGameLocation = "delete-one-adventure-game-location"
)

func adventureGameLocationHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameLocationHandlerConfig")

	l.Debug("Adding adventure_game_location handler configuration")

	gameLocationConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_location.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_location.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_location.response.schema.json",
		},
		References: referenceSchemas,
	}

	// New Adventure Game Location API paths
	gameLocationConfig[searchManyAdventureGameLocations] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-locations",
		HandlerFunc: searchManyAdventureGameLocationsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game locations",
		},
	}

	gameLocationConfig[getManyAdventureGameLocations] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/locations",
		HandlerFunc: getManyAdventureGameLocationsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game locations",
		},
	}

	gameLocationConfig[getOneAdventureGameLocation] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/locations/:location_id",
		HandlerFunc: getOneAdventureGameLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game location",
		},
	}

	gameLocationConfig[createOneAdventureGameLocation] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/locations",
		HandlerFunc: createOneAdventureGameLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game location",
		},
	}

	gameLocationConfig[updateOneAdventureGameLocation] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/locations/:location_id",
		HandlerFunc: updateOneAdventureGameLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game location",
		},
	}

	gameLocationConfig[deleteOneAdventureGameLocation] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/locations/:location_id",
		HandlerFunc: deleteOneAdventureGameLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game location",
		},
	}

	return gameLocationConfig, nil
}

func searchManyAdventureGameLocationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameLocationsHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter to only return adventure game locations
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_type",
		Val: "adventure",
	})

	recs, err := mm.GetManyAdventureGameLocationRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameLocationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameLocationsHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_id",
		Val: gameID,
	})

	recs, err := mm.GetManyAdventureGameLocationRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameLocationHandler")

	gameID := pp.ByName("game_id")
	locationID := pp.ByName("location_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationRec(locationID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location record >%v<", err)
		return err
	}

	// Verify the location belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("location does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location", locationID)
	}

	res, err := mapper.AdventureGameLocationRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game location record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameLocationHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game_id path parameter is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	var req schema.AdventureGameLocationRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.AdventureGameLocationRequestToRecord(l, &req, &record.AdventureGameLocation{})
	if err != nil {
		return err
	}

	// Set the game ID from the path parameter
	rec.GameID = gameID

	mm := m.(*domain.Domain)

	rec, err = mm.CreateAdventureGameLocationRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameLocationHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game_id path parameter is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	locationID := pp.ByName("location_id")
	if locationID == "" {
		l.Warn("location_id path parameter is required")
		return coreerror.NewNotFoundError("location", locationID)
	}

	l.Info("updating adventure game location record with path params >%#v<", pp)

	var req schema.AdventureGameLocationRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationRec(locationID, nil)
	if err != nil {
		return err
	}

	// Verify the location belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("location does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location", locationID)
	}

	rec, err = mapper.AdventureGameLocationRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameLocationRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game location record >%v<", err)
		return err
	}

	data, err := mapper.AdventureGameLocationRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.AdventureGameLocationResponse{
		Data: &data,
	}

	l.Info("responding with updated adventure game location record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameLocationHandler")

	gameID := pp.ByName("game_id")
	locationID := pp.ByName("location_id")

	l.Info("deleting adventure game location record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	// First get the record to verify it belongs to the specified game
	rec, err := mm.GetAdventureGameLocationRec(locationID, nil)
	if err != nil {
		return err
	}

	// Verify the location belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("location does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location", locationID)
	}

	if err := mm.DeleteAdventureGameLocationRec(locationID); err != nil {
		l.Warn("failed deleting adventure game location record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
