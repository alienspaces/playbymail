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
	l = loggerWithFunctionContext(l, "adventureGameItemHandlerConfig")

	l.Debug("Adding adventure_game_item handler configuration")

	gameItemConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_item.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_item.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_item.response.schema.json",
		},
		References: referenceSchemas,
	}

	// New Adventure Game Item API paths
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

	gameItemConfig[createOneAdventureGameItem] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/items",
		HandlerFunc: createOneAdventureGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game item",
		},
	}

	gameItemConfig[updateOneAdventureGameItem] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/items/:item_id",
		HandlerFunc: updateOneAdventureGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game item",
		},
	}

	gameItemConfig[deleteOneAdventureGameItem] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/items/:item_id",
		HandlerFunc: deleteOneAdventureGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game item",
		},
	}

	// Legacy paths (do not modify)
	// Legacy handlers commented out due to Adventure-prefixed naming alignment
	/*
		gameItemConfig[getManyGameItems] = server.HandlerConfig{
			Method:      http.MethodGet,
			Path:        "/v1/game-items",
			HandlerFunc: getManyGameItemsHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
				ValidateResponseSchema: collectionResponseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document:   true,
				Collection: true,
				Title:      "Get game_item collection",
			},
		}

		gameItemConfig[getOneGameItem] = server.HandlerConfig{
			Method:      http.MethodGet,
			Path:        "/v1/game-items/:game_item_id",
			HandlerFunc: getGameItemHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document: true,
				Title:    "Get game_item",
			},
		}

		gameItemConfig[createGameItem] = server.HandlerConfig{
			Method:      http.MethodPost,
			Path:        "/v1/game-items",
			HandlerFunc: createGameItemHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
				ValidateRequestSchema:  requestSchema,
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document: true,
				Title:    "Create game_item",
			},
		}

		gameItemConfig[updateGameItem] = server.HandlerConfig{
			Method:      http.MethodPut,
			Path:        "/v1/game-items/:game_item_id",
			HandlerFunc: updateGameItemHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
				ValidateRequestSchema:  requestSchema,
				ValidateResponseSchema: responseSchema,
			},
			DocumentationConfig: server.DocumentationConfig{
				Document: true,
				Title:    "Update game_item",
			},
		}

		gameItemConfig[deleteGameItem] = server.HandlerConfig{
			Method:      http.MethodDelete,
			Path:        "/v1/game-items/:game_item_id",
			HandlerFunc: deleteGameItemHandler,
			MiddlewareConfig: server.MiddlewareConfig{
				AuthenTypes: []server.AuthenticationType{
					server.AuthenticationTypeToken,
				},
			},
			DocumentationConfig: server.DocumentationConfig{
				Document: true,
				Title:    "Delete game_item",
			},
		}
	*/

	return gameItemConfig, nil
}

// New Adventure Game Item Handlers

func searchManyAdventureGameItemsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "SearchManyAdventureGameItemsHandler")

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

func getManyAdventureGameItemsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyAdventureGameItemsHandler")

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

func getOneAdventureGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetOneAdventureGameItemHandler")

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

func createOneAdventureGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateOneAdventureGameItemHandler")

	gameID := pp.ByName("game_id")

	var req schema.AdventureGameItemRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.AdventureGameItemRequestToRecord(l, &req, &record.AdventureGameItem{})
	if err != nil {
		return err
	}

	// Set the game ID from the path parameter
	rec.GameID = gameID

	mm := m.(*domain.Domain)

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

func updateOneAdventureGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateOneAdventureGameItemHandler")

	gameID := pp.ByName("game_id")
	itemID := pp.ByName("item_id")

	l.Info("updating adventure game item record with path params >%#v<", pp)

	var req schema.AdventureGameItemRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameItemRec(itemID, nil)
	if err != nil {
		return err
	}

	// Verify the item belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("item does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item", itemID)
	}

	rec, err = mapper.AdventureGameItemRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameItemRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game item record >%v<", err)
		return err
	}

	data, err := mapper.AdventureGameItemRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.AdventureGameItemResponse{
		Data: &data,
	}

	l.Info("responding with updated adventure game item record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteOneAdventureGameItemHandler")

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

// Legacy handler - commented out due to Adventure-prefixed naming alignment
/*
func getManyGameItemsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyGameItemsHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyGameItemRecs(opts)
	if err != nil {
		l.Warn("failed getting game_item records >%v<", err)
		return err
	}

	res, err := mapper.GameItemRecordsToCollectionResponse(l, recs)
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
func getGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetGameItemHandler")

	gameItemID := pp.ByName("game_item_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetGameItemRec(gameItemID, nil)
	if err != nil {
		l.Warn("failed getting game_item record >%v<", err)
		return err
	}

	res, err := mapper.GameItemRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game_item record to response >%v<", err)
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
func createGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateGameItemHandler")

	var req schema.GameItemRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.GameItemRequestToRecord(l, &req, &record.GameItem{})
	if err != nil {
		return err
	}

	mm := m.(*domain.Domain)

	rec, err = mm.CreateGameItemRec(rec)
	if err != nil {
		l.Warn("failed creating game_item record >%v<", err)
		return err
	}

	res, err := mapper.GameItemRecordToResponse(l, rec)
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

// Legacy handler - commented out due to Adventure-prefixed naming alignment
/*
func updateGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateGameItemHandler")

	gameItemID := pp.ByName("game_item_id")

	l.Info("updating game_item record with path params >%#v<", pp)

	var req schema.GameItemRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetGameItemRec(gameItemID, nil)
	if err != nil {
		return err
	}

	rec, err = mapper.GameItemRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateGameItemRec(rec)
	if err != nil {
		l.Warn("failed updating game_item record >%v<", err)
		return err
	}

	data, err := mapper.GameItemRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.GameItemResponse{
		Data: &data,
	}

	l.Info("responding with updated game_item record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
*/

// Legacy handler - commented out due to Adventure-prefixed naming alignment
/*
func deleteGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteGameItemHandler")

	gameItemID := pp.ByName("game_item_id")

	l.Info("deleting game_item record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	if err := mm.DeleteGameItemRec(gameItemID); err != nil {
		l.Warn("failed deleting game_item record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
*/
