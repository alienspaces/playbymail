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
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-items

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/items
// GET (document)    /api/v1/adventure-games/{game_id}/items/{item_id}
// POST (document)   /api/v1/adventure-games/{game_id}/items
// PUT (document)    /api/v1/adventure-games/{game_id}/items/{item_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/items/{item_id}

const (
	// API Resource Search Path
	searchManyAdventureGameItems = "search-many-adventure-game-items"

	// API Resource CRUD Paths
	getManyAdventureGameItems  = "get-many-adventure-game-items"
	getOneAdventureGameItem    = "get-one-adventure-game-item"
	createOneAdventureGameItem = "create-one-adventure-game-item"
	updateOneAdventureGameItem = "update-one-adventure-game-item"
	deleteOneAdventureGameItem = "delete-one-adventure-game-item"
)

func adventureGameItemHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameItemHandlerConfig")

	l.Debug("Adding adventure_game_item handler configuration")

	gameItemConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_item.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_item.collection.response.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_item.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_item.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_item.schema.json",
			},
		}...),
	}

	// Get all adventure game items for all adventure games
	gameItemConfig[searchManyAdventureGameItems] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-items",
		HandlerFunc: searchManyAdventureGameItemsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game items",
		},
	}

	// Get all adventure game items for a specific adventure game
	gameItemConfig[getManyAdventureGameItems] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/items",
		HandlerFunc: getManyAdventureGameItemsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game items",
		},
	}

	// Get a specific adventure game item
	gameItemConfig[getOneAdventureGameItem] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/items/:item_id",
		HandlerFunc: getOneAdventureGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game item",
		},
	}

	// Create a new adventure game item
	gameItemConfig[createOneAdventureGameItem] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/items",
		HandlerFunc: createOneAdventureGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameDesign,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game item",
		},
	}

	// Update an existing adventure game item
	gameItemConfig[updateOneAdventureGameItem] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/items/:item_id",
		HandlerFunc: updateOneAdventureGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameDesign,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game item",
		},
	}

	// Delete an existing adventure game item
	gameItemConfig[deleteOneAdventureGameItem] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/items/:item_id",
		HandlerFunc: deleteOneAdventureGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameDesign,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game item",
		},
	}

	return gameItemConfig, nil
}

func searchManyAdventureGameItemsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameItemsHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter to only return adventure game items
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_type",
		Val: "adventure",
	})

	recs, err := mm.GetManyAdventureGameItemRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game item records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameItemsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameItemsHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_id",
		Val: gameID,
	})

	recs, err := mm.GetManyAdventureGameItemRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game item records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameItemHandler")

	gameID := pp.ByName("game_id")
	itemID := pp.ByName("item_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameItemRec(itemID, nil)
	if err != nil {
		l.Warn("failed getting adventure game item record >%v<", err)
		return err
	}

	// Verify the item belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("item does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item", itemID)
	}

	res, err := mapper.AdventureGameItemRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game item record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameItemHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed getting game record >%v<", err)
		return err
	}

	rec := &adventure_game_record.AdventureGameItem{
		GameID: gameRec.ID,
	}

	rec, err = mapper.AdventureGameItemRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateAdventureGameItemRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game item record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameItemHandler")

	gameID := pp.ByName("game_id")
	itemID := pp.ByName("item_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameItemRec(itemID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	// Verify the item belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("item does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item", itemID)
	}

	rec, err = mapper.AdventureGameItemRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameItemRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game item record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameItemHandler")

	gameID := pp.ByName("game_id")
	itemID := pp.ByName("item_id")

	l.Info("deleting adventure game item record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	// First get the record to verify it belongs to the specified game
	rec, err := mm.GetAdventureGameItemRec(itemID, nil)
	if err != nil {
		return err
	}

	// Verify the item belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("item does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item", itemID)
	}

	if err := mm.DeleteAdventureGameItemRec(itemID); err != nil {
		l.Warn("failed deleting adventure game item record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
