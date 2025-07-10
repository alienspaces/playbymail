package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

const (
	tagGroupGameItem server.TagGroup = "GameItems"
	TagGameItem      server.Tag      = "GameItems"
)

const (
	getManyGameItems = "get-game-items"
	getOneGameItem   = "get-game-item"
	createGameItem   = "create-game-item"
	updateGameItem   = "update-game-item"
	deleteGameItem   = "delete-game-item"
)

func (rnr *Runner) gameItemHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameItemHandlerConfig")

	l.Debug("Adding game_item handler configuration")

	gameItemConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_item.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_item.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_item.response.schema.json",
		},
		References: referenceSchemas,
	}

	gameItemConfig[getManyGameItems] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-items",
		HandlerFunc: rnr.getManyGameItemsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
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
		HandlerFunc: rnr.getGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
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
		HandlerFunc: rnr.createGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
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
		HandlerFunc: rnr.updateGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
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
		HandlerFunc: rnr.deleteGameItemHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game_item",
		},
	}

	return gameItemConfig, nil
}

func (rnr *Runner) getManyGameItemsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
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

func (rnr *Runner) getGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
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

func (rnr *Runner) createGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
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

func (rnr *Runner) updateGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
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

func (rnr *Runner) deleteGameItemHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
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
