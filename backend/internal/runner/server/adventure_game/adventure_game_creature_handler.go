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
	SearchManyAdventureGameCreatures = "search-many-adventure-game-creatures"

	// API Resource CRUD Paths
	GetManyAdventureGameCreatures  = "get-many-adventure-game-creatures"
	GetOneAdventureGameCreature    = "get-one-adventure-game-creature"
	CreateOneAdventureGameCreature = "create-one-adventure-game-creature"
	UpdateOneAdventureGameCreature = "update-one-adventure-game-creature"
	DeleteOneAdventureGameCreature = "delete-one-adventure-game-creature"
)

func adventureGameCreatureHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameCreatureHandlerConfig")

	l.Debug("Adding adventure_game_creature handler configuration")

	gameCreatureConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_creature.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_creature.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_creature.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_creature.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_creature.schema.json",
			},
		}...),
	}

	// New Adventure Game Creature API paths
	gameCreatureConfig[SearchManyAdventureGameCreatures] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-creatures",
		HandlerFunc: searchManyAdventureGameCreaturesHandler,
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

	gameCreatureConfig[GetManyAdventureGameCreatures] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/creatures",
		HandlerFunc: getManyAdventureGameCreaturesHandler,
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

	gameCreatureConfig[GetOneAdventureGameCreature] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/creatures/:creature_id",
		HandlerFunc: getOneAdventureGameCreatureHandler,
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

	gameCreatureConfig[CreateOneAdventureGameCreature] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/creatures",
		HandlerFunc: createOneAdventureGameCreatureHandler,
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

	gameCreatureConfig[UpdateOneAdventureGameCreature] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/creatures/:creature_id",
		HandlerFunc: updateOneAdventureGameCreatureHandler,
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

	gameCreatureConfig[DeleteOneAdventureGameCreature] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/creatures/:creature_id",
		HandlerFunc: deleteOneAdventureGameCreatureHandler,
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

	return gameCreatureConfig, nil
}

func searchManyAdventureGameCreaturesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameCreaturesHandler")

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

func getManyAdventureGameCreaturesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameCreaturesHandler")

	// Create SQL options from query parameters
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts.Params = append(opts.Params, sql.Param{
		Col: "game_id",
		Val: gameID,
	})

	mm := m.(*domain.Domain)

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

func getOneAdventureGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameCreatureHandler")

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

func createOneAdventureGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameCreatureHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed getting game record >%v<", err)
		return err
	}

	rec := &adventure_game_record.AdventureGameCreature{
		GameID: gameRec.ID,
	}

	rec, err = mapper.AdventureGameCreatureRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

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

func updateOneAdventureGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameCreatureHandler")

	gameID := pp.ByName("game_id")
	creatureID := pp.ByName("creature_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameCreatureRec(creatureID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	// Verify the creature belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("creature does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("creature", creatureID)
	}

	rec, err = mapper.AdventureGameCreatureRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameCreatureRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game creature record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCreatureRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameCreatureHandler")

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
func getManyGameCreaturesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, "GetManyGameCreaturesHandler")

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
func getGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, "GetGameCreatureHandler")

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
func createGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = logging.LoggerWithFunctionContext(l, "CreateGameCreatureHandler")

	var req adventure_game.GameCreatureRequest
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
// func updateGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
// 	l = logging.LoggerWithFunctionContext(l, "UpdateGameCreatureHandler")

// 	gameCreatureID := pp.ByName("game_creature_id")

// 	l.Info("updating game_creature record with path params >%#v<", pp)

// 	var req adventure_game_schema.AdventureGameCreatureRequest
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

// 	res := adventure_game_schema.AdventureGameCreatureResponse{
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
// func deleteGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
// 	l = logging.LoggerWithFunctionContext(l, "DeleteGameCreatureHandler")

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
