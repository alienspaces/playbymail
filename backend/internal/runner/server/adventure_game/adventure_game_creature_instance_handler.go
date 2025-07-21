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
	"gitlab.com/alienspaces/playbymail/schema"
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-creature-instances

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-game-instances/{game_instance_id}/creatures
// GET (document)    /api/v1/adventure-game-instances/{game_instance_id}/creatures/{creature_instance_id}
// POST (document)   /api/v1/adventure-game-instances/{game_instance_id}/creatures
// PUT (document)    /api/v1/adventure-game-instances/{game_instance_id}/creatures/{creature_instance_id}
// DELETE (document) /api/v1/adventure-game-instances/{game_instance_id}/creatures/{creature_instance_id}

const (
	// API Resource Search Path
	searchManyAdventureGameCreatureInstances = "search-many-adventure-game-creature-instances"

	// API Resource CRUD Paths
	getManyAdventureGameCreatureInstances  = "get-many-adventure-game-creature-instances"
	getOneAdventureGameCreatureInstance    = "get-one-adventure-game-creature-instance"
	createOneAdventureGameCreatureInstance = "create-one-adventure-game-creature-instance"
	updateOneAdventureGameCreatureInstance = "update-one-adventure-game-creature-instance"
	deleteOneAdventureGameCreatureInstance = "delete-one-adventure-game-creature-instance"
)

func adventureGameCreatureInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "adventureGameCreatureInstanceHandlerConfig")

	l.Debug("Adding adventure_game_creature_instance handler configuration")

	gameCreatureInstanceConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_creature_instance.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_creature_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_creature_instance.response.schema.json",
		},
		References: referenceSchemas,
	}

	// New Adventure Game Creature Instance API paths
	gameCreatureInstanceConfig[searchManyAdventureGameCreatureInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-creature-instances",
		HandlerFunc: searchManyAdventureGameCreatureInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game creature instances",
		},
	}

	gameCreatureInstanceConfig[getManyAdventureGameCreatureInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/creature-instances",
		HandlerFunc: getManyAdventureGameCreatureInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game creature instances",
		},
	}

	gameCreatureInstanceConfig[getOneAdventureGameCreatureInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/creature-instances/:creature_instance_id",
		HandlerFunc: getOneAdventureGameCreatureInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game creature instance",
		},
	}

	gameCreatureInstanceConfig[createOneAdventureGameCreatureInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/creature-instances",
		HandlerFunc: createOneAdventureGameCreatureInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game creature instance",
		},
	}

	gameCreatureInstanceConfig[updateOneAdventureGameCreatureInstance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/creature-instances/:creature_instance_id",
		HandlerFunc: updateOneAdventureGameCreatureInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game creature instance",
		},
	}

	gameCreatureInstanceConfig[deleteOneAdventureGameCreatureInstance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/creature-instances/:creature_instance_id",
		HandlerFunc: deleteOneAdventureGameCreatureInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game creature instance",
		},
	}
	return gameCreatureInstanceConfig, nil
}

// New Adventure Game Creature Instance Handlers

func searchManyAdventureGameCreatureInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "SearchManyAdventureGameCreatureInstancesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameCreatureInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game creature instance records >%v<", err)
		return err
	}

	res := mapper.AdventureGameCreatureInstanceRecordsToCollectionResponse(recs)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameCreatureInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyAdventureGameCreatureInstancesHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: record.FieldAdventureGameCreatureInstanceAdventureGameInstanceID,
		Val: gameInstanceID,
	})

	recs, err := mm.GetManyAdventureGameCreatureInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game creature instance records >%v<", err)
		return err
	}

	res := mapper.AdventureGameCreatureInstanceRecordsToCollectionResponse(recs)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameCreatureInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetOneAdventureGameCreatureInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	creatureInstanceID := pp.ByName("creature_instance_id")
	if creatureInstanceID == "" {
		l.Warn("creature instance id is required")
		return coreerror.NewNotFoundError("creature instance", creatureInstanceID)
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetAdventureGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameCreatureInstanceRec(creatureInstanceID, nil)
	if err != nil {
		l.Warn("failed getting adventure game creature instance record >%v<", err)
		return err
	}

	// Verify the creature instance belongs to the specified game
	if rec.AdventureGameInstanceID != gameInstanceRec.ID {
		l.Warn("creature instance does not belong to specified game instance >%s< != >%s<", rec.AdventureGameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("creature instance", creatureInstanceID)
	}

	res := mapper.AdventureGameCreatureInstanceRecordToResponse(rec)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameCreatureInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateOneAdventureGameCreatureInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	var req schema.AdventureGameCreatureInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.AdventureGameCreatureInstanceRequestToRecord(&req, nil)
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

	rec, err = mm.CreateAdventureGameCreatureInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game creature instance record >%v<", err)
		return err
	}

	res := mapper.AdventureGameCreatureInstanceRecordToResponse(rec)

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameCreatureInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateOneAdventureGameCreatureInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	creatureInstanceID := pp.ByName("creature_instance_id")
	if creatureInstanceID == "" {
		l.Warn("creature instance id is required")
		return coreerror.NewNotFoundError("creature instance", creatureInstanceID)
	}

	l.Info("updating adventure game creature instance record with path params >%#v<", pp)

	var req schema.AdventureGameCreatureInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetAdventureGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameCreatureInstanceRec(creatureInstanceID, nil)
	if err != nil {
		return err
	}

	// Verify the creature instance belongs to the specified game
	if rec.AdventureGameInstanceID != gameInstanceRec.ID {
		l.Warn("creature instance does not belong to specified game instance >%s< != >%s<", rec.AdventureGameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("creature instance", creatureInstanceID)
	}

	rec, err = mapper.AdventureGameCreatureInstanceRequestToRecord(&req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameCreatureInstanceRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game creature instance record >%v<", err)
		return err
	}

	data := mapper.AdventureGameCreatureInstanceRecordToResponseData(rec)

	res := schema.AdventureGameCreatureInstanceResponse{
		Data: data,
	}

	l.Info("responding with updated adventure game creature instance record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameCreatureInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteOneAdventureGameCreatureInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	creatureInstanceID := pp.ByName("creature_instance_id")
	if creatureInstanceID == "" {
		l.Warn("creature instance id is required")
		return coreerror.NewNotFoundError("creature instance", creatureInstanceID)
	}

	l.Info("deleting adventure game creature instance record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetAdventureGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	// First get the record to verify it belongs to the specified game
	rec, err := mm.GetAdventureGameCreatureInstanceRec(creatureInstanceID, nil)
	if err != nil {
		return err
	}

	// Verify the creature instance belongs to the specified game
	if rec.AdventureGameInstanceID != gameInstanceRec.ID {
		l.Warn("creature instance does not belong to specified game instance >%s< != >%s<", rec.AdventureGameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("creature instance", creatureInstanceID)
	}

	if err := mm.DeleteAdventureGameCreatureInstanceRec(creatureInstanceID); err != nil {
		l.Warn("failed deleting adventure game creature instance record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
