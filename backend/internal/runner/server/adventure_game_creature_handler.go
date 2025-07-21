package runner

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
// GET (collection) /api/v1/adventure-game-creatures

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/creatures
// GET (document)    /api/v1/adventure-games/{game_id}/creatures/{creature_id}
// POST (document)   /api/v1/adventure-games/{game_id}/creatures
// PUT (document)    /api/v1/adventure-games/{game_id}/creatures/{creature_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/creatures/{creature_id}

const (
	// API Resource Search Path
	searchManyAdventureGameCreatures = "search-many-adventure-game-creatures"

	// API Resource CRUD Paths
	getManyAdventureGameCreatures  = "get-many-adventure-game-creatures"
	getOneAdventureGameCreature    = "get-one-adventure-game-creature"
	createOneAdventureGameCreature = "create-one-adventure-game-creature"
	updateOneAdventureGameCreature = "update-one-adventure-game-creature"
	deleteOneAdventureGameCreature = "delete-one-adventure-game-creature"
)

func (rnr *Runner) adventureGameCreatureHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "adventureGameCreatureHandlerConfig")

	l.Debug("Adding adventure_game_creature handler configuration")

	gameCreatureConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_creature.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_creature.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_creature.response.schema.json",
		},
		References: referenceSchemas,
	}

	// New Adventure Game Creature API paths
	gameCreatureConfig[searchManyAdventureGameCreatures] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-creatures",
		HandlerFunc: rnr.searchManyAdventureGameCreaturesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game creatures",
		},
	}

	gameCreatureConfig[getManyAdventureGameCreatures] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/creatures",
		HandlerFunc: rnr.getManyAdventureGameCreaturesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game creatures",
		},
	}

	gameCreatureConfig[getOneAdventureGameCreature] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/creatures/:creature_id",
		HandlerFunc: rnr.getOneAdventureGameCreatureHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game creature",
		},
	}

	gameCreatureConfig[createOneAdventureGameCreature] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/creatures",
		HandlerFunc: rnr.createOneAdventureGameCreatureHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game creature",
		},
	}

	gameCreatureConfig[updateOneAdventureGameCreature] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/creatures/:creature_id",
		HandlerFunc: rnr.updateOneAdventureGameCreatureHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game creature",
		},
	}

	gameCreatureConfig[deleteOneAdventureGameCreature] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/creatures/:creature_id",
		HandlerFunc: rnr.deleteOneAdventureGameCreatureHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game creature",
		},
	}

	// Legacy paths (do not modify)
	// Legacy handlers commented out due to missing handler functions
	/*
		gameCreatureConfig[getManyGameCreatures] = server.HandlerConfig{
			Method:      http.MethodGet,
			Path:        "/v1/game-creatures",
			HandlerFunc: rnr.getManyGameCreaturesHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
				ValidateResponseSchema: collectionResponseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:   true,
				Collection: true,
				Title:      "Get game_creature collection",
			},
		}

		gameCreatureConfig[getOneGameCreature] = server.HandlerConfig{
			Method:      http.MethodGet,
			Path:        "/v1/game-creatures/:game_creature_id",
			HandlerFunc: rnr.getGameCreatureHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document: true,
				Title:    "Get game_creature",
			},
		}

		gameCreatureConfig[createGameCreature] = server.HandlerConfig{
			Method:      http.MethodPost,
			Path:        "/v1/game-creatures",
			HandlerFunc: rnr.createGameCreatureHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
				ValidateRequestSchema:  requestSchema,
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document: true,
				Title:    "Create game_creature",
			},
		}

		gameCreatureConfig[updateGameCreature] = server.HandlerConfig{
			Method:      http.MethodPut,
			Path:        "/v1/game-creatures/:game_creature_id",
			HandlerFunc: rnr.updateGameCreatureHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
				ValidateRequestSchema:  requestSchema,
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document: true,
				Title:    "Update game_creature",
			},
		}

		gameCreatureConfig[deleteGameCreature] = server.HandlerConfig{
			Method:      http.MethodDelete,
			Path:        "/v1/game-creatures/:game_creature_id",
			HandlerFunc: rnr.deleteGameCreatureHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document: true,
				Title:    "Delete game_creature",
			},
		}
	*/

	return gameCreatureConfig, nil
}

// New Adventure Game Creature Handlers

func (rnr *Runner) searchManyAdventureGameCreaturesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "SearchManyAdventureGameCreaturesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter to only return adventure game creatures
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_type",
		Val: "adventure",
	})

	recs, err := mm.GetManyAdventureGameCreatureRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game creature records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCreatureRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getManyAdventureGameCreaturesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyAdventureGameCreaturesHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_id",
		Val: gameID,
	})

	recs, err := mm.GetManyAdventureGameCreatureRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game creature records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCreatureRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getOneAdventureGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetOneAdventureGameCreatureHandler")

	gameID := pp.ByName("game_id")
	creatureID := pp.ByName("creature_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameCreatureRec(creatureID, nil)
	if err != nil {
		l.Warn("failed getting adventure game creature record >%v<", err)
		return err
	}

	// Verify the creature belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("creature does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("creature", creatureID)
	}

	res, err := mapper.AdventureGameCreatureRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game creature record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) createOneAdventureGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateOneAdventureGameCreatureHandler")

	gameID := pp.ByName("game_id")

	var req schema.AdventureGameCreatureRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.AdventureGameCreatureRequestToRecord(l, &req, &record.AdventureGameCreature{})
	if err != nil {
		return err
	}

	// Set the game ID from the path parameter
	rec.GameID = gameID

	mm := m.(*domain.Domain)

	rec, err = mm.CreateAdventureGameCreatureRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game creature record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCreatureRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) updateOneAdventureGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateOneAdventureGameCreatureHandler")

	gameID := pp.ByName("game_id")
	creatureID := pp.ByName("creature_id")

	l.Info("updating adventure game creature record with path params >%#v<", pp)

	var req schema.AdventureGameCreatureRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameCreatureRec(creatureID, nil)
	if err != nil {
		return err
	}

	// Verify the creature belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("creature does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("creature", creatureID)
	}

	rec, err = mapper.AdventureGameCreatureRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameCreatureRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game creature record >%v<", err)
		return err
	}

	data, err := mapper.AdventureGameCreatureRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.AdventureGameCreatureResponse{
		Data: &data,
	}

	l.Info("responding with updated adventure game creature record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) deleteOneAdventureGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteOneAdventureGameCreatureHandler")

	gameID := pp.ByName("game_id")
	creatureID := pp.ByName("creature_id")

	l.Info("deleting adventure game creature record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	// First get the record to verify it belongs to the specified game
	rec, err := mm.GetAdventureGameCreatureRec(creatureID, nil)
	if err != nil {
		return err
	}

	// Verify the creature belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("creature does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("creature", creatureID)
	}

	if err := mm.DeleteAdventureGameCreatureRec(creatureID); err != nil {
		l.Warn("failed deleting adventure game creature record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// Legacy handler - commented out due to Adventure-prefixed naming alignment
/*
func (rnr *Runner) getManyGameCreaturesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyGameCreaturesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyGameCreatureRecs(opts)
	if err != nil {
		l.Warn("failed getting game_creature records >%v<", err)
		return err
	}

	res, err := mapper.GameCreatureRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
*/

// Legacy handler - commented out due to Adventure-prefixed naming alignment
/*
func (rnr *Runner) getGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetGameCreatureHandler")

	gameCreatureID := pp.ByName("game_creature_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetGameCreatureRec(gameCreatureID, nil)
	if err != nil {
		l.Warn("failed getting game_creature record >%v<", err)
		return err
	}

	res, err := mapper.GameCreatureRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game_creature record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
*/

// Legacy handler - commented out due to Adventure-prefixed naming alignment
/*
func (rnr *Runner) createGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateGameCreatureHandler")

	var req schema.GameCreatureRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.GameCreatureRequestToRecord(l, &req, &record.GameCreature{})
	if err != nil {
		return err
	}

	mm := m.(*domain.Domain)

	rec, err = mm.CreateGameCreatureRec(rec)
	if err != nil {
		l.Warn("failed creating game_creature record >%v<", err)
		return err
	}

	res, err := mapper.GameCreatureRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
*/

// // Legacy handler
// func (rnr *Runner) updateGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
// 	l = loggerWithFunctionContext(l, "UpdateGameCreatureHandler")

// 	gameCreatureID := pp.ByName("game_creature_id")

// 	l.Info("updating game_creature record with path params >%#v<", pp)

// 	var req schema.AdventureGameCreatureRequest
// 	if _, err := server.ReadRequest(l, r, &req); err != nil {
// 		l.Warn("failed reading request >%v<", err)
// 		return err
// 	}

// 	mm := m.(*domain.Domain)

// 	rec, err := mm.GetAdventureGameCreatureRec(gameCreatureID, nil)
// 	if err != nil {
// 		return err
// 	}

// 	rec, err = mapper.AdventureGameCreatureRequestToRecord(l, &req, rec)
// 	if err != nil {
// 		return err
// 	}

// 	rec, err = mm.UpdateAdventureGameCreatureRec(rec)
// 	if err != nil {
// 		l.Warn("failed updating game_creature record >%v<", err)
// 		return err
// 	}

// 	data, err := mapper.AdventureGameCreatureRecordToResponseData(l, rec)
// 	if err != nil {
// 		return err
// 	}

// 	res := schema.AdventureGameCreatureResponse{
// 		Data: &data,
// 	}

// 	l.Info("responding with updated game_creature record id >%s<", rec.ID)

// 	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
// 		l.Warn("failed writing response >%v<", err)
// 		return err
// 	}

// 	return nil
// }

// // Legacy handler
// func (rnr *Runner) deleteGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
// 	l = loggerWithFunctionContext(l, "DeleteGameCreatureHandler")

// 	gameCreatureID := pp.ByName("game_creature_id")

// 	l.Info("deleting game_creature record with path params >%#v<", pp)

// 	mm := m.(*domain.Domain)

// 	if err := mm.DeleteAdventureGameCreatureRec(gameCreatureID); err != nil {
// 		l.Warn("failed deleting game_creature record >%v<", err)
// 		return err
// 	}

// 	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
// 		l.Warn("failed writing response >%v<", err)
// 		return err
// 	}

// 	return nil
// }
