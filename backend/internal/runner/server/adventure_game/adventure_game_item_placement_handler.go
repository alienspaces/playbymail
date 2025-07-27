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
	"gitlab.com/alienspaces/playbymail/schema"
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-item-placements

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/item-placements
// GET (document)    /api/v1/adventure-games/{game_id}/item-placements/{placement_id}
// POST (document)   /api/v1/adventure-games/{game_id}/item-placements
// PUT (document)    /api/v1/adventure-games/{game_id}/item-placements/{placement_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/item-placements/{placement_id}

const (
	// API Resource Search Path
	SearchManyAdventureGameItemPlacements = "search-many-adventure-game-item-placements"

	// API Resource CRUD Paths
	GetManyAdventureGameItemPlacements  = "get-many-adventure-game-item-placements"
	GetOneAdventureGameItemPlacement    = "get-one-adventure-game-item-placement"
	CreateOneAdventureGameItemPlacement = "create-one-adventure-game-item-placement"
	UpdateOneAdventureGameItemPlacement = "update-one-adventure-game-item-placement"
	DeleteOneAdventureGameItemPlacement = "delete-one-adventure-game-item-placement"
)

func adventureGameItemPlacementHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameItemPlacementHandlerConfig")

	l.Debug("Adding adventure_game_item_placement handler configuration")

	itemPlacementConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_item_placement.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_item_placement.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_item_placement.response.schema.json",
		},
		References: referenceSchemas,
	}

	// New Adventure Game Item Placement API paths
	itemPlacementConfig[SearchManyAdventureGameItemPlacements] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-item-placements",
		HandlerFunc: searchManyAdventureGameItemPlacementsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game item placements",
		},
	}

	itemPlacementConfig[GetManyAdventureGameItemPlacements] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/item-placements",
		HandlerFunc: getManyAdventureGameItemPlacementsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get many adventure game item placements",
		},
	}

	itemPlacementConfig[GetOneAdventureGameItemPlacement] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/item-placements/:placement_id",
		HandlerFunc: getOneAdventureGameItemPlacementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get one adventure game item placement",
		},
	}

	itemPlacementConfig[CreateOneAdventureGameItemPlacement] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/item-placements",
		HandlerFunc: createOneAdventureGameItemPlacementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create one adventure game item placement",
		},
	}

	itemPlacementConfig[UpdateOneAdventureGameItemPlacement] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/item-placements/:placement_id",
		HandlerFunc: updateOneAdventureGameItemPlacementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update one adventure game item placement",
		},
	}

	itemPlacementConfig[DeleteOneAdventureGameItemPlacement] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/item-placements/:placement_id",
		HandlerFunc: deleteOneAdventureGameItemPlacementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete one adventure game item placement",
		},
	}

	return itemPlacementConfig, nil
}

// searchManyAdventureGameItemPlacementsHandler -
func searchManyAdventureGameItemPlacementsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameItemPlacementsHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameItemPlacementRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game item placement records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemPlacementRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// getManyAdventureGameItemPlacementsHandler -
func getManyAdventureGameItemPlacementsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameItemPlacementsHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_id",
		Val: gameID,
	})

	recs, err := mm.GetManyAdventureGameItemPlacementRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game item placement records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemPlacementRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// getOneAdventureGameItemPlacementHandler -
func getOneAdventureGameItemPlacementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameItemPlacementHandler")

	gameID := pp.ByName("game_id")
	placementID := pp.ByName("placement_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameItemPlacementRec(placementID, nil)
	if err != nil {
		l.Warn("failed getting adventure game item placement record >%v<", err)
		return err
	}

	// Verify the placement belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("placement does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item_placement", placementID)
	}

	res, err := mapper.AdventureGameItemPlacementRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game item placement record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// createOneAdventureGameItemPlacementHandler -
func createOneAdventureGameItemPlacementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameItemPlacementHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	var request schema.AdventureGameItemPlacementRequest
	if _, err := server.ReadRequest(l, r, &request); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.AdventureGameItemPlacementRequestToRecord(l, &request, &adventure_game_record.AdventureGameItemPlacement{})
	if err != nil {
		return err
	}

	// Set the game ID from the path parameter
	rec.GameID = gameID

	rec, err = mm.CreateAdventureGameItemPlacementRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game item placement record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemPlacementRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game item placement record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// updateOneAdventureGameItemPlacementHandler -
func updateOneAdventureGameItemPlacementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameItemPlacementHandler")

	gameID := pp.ByName("game_id")
	placementID := pp.ByName("placement_id")
	mm := m.(*domain.Domain)

	var request schema.AdventureGameItemPlacementRequest
	if _, err := server.ReadRequest(l, r, &request); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mm.GetAdventureGameItemPlacementRec(placementID, sql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed getting adventure game item placement record >%v<", err)
		return err
	}

	// Verify the placement belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("placement does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item_placement", placementID)
	}

	rec, err = mapper.AdventureGameItemPlacementRequestToRecord(l, &request, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameItemPlacementRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game item placement record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemPlacementRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game item placement record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// deleteOneAdventureGameItemPlacementHandler -
func deleteOneAdventureGameItemPlacementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameItemPlacementHandler")

	gameID := pp.ByName("game_id")
	placementID := pp.ByName("placement_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameItemPlacementRec(placementID, sql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed getting adventure game item placement record >%v<", err)
		return err
	}

	// Verify the placement belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("placement does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item_placement", placementID)
	}

	if err := mm.DeleteAdventureGameItemPlacementRec(placementID); err != nil {
		l.Warn("failed deleting adventure game item placement record >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
