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
	"gitlab.com/alienspaces/playbymail/schema"
)

const (
	getManyGameItemInstances = "get-game-item-instances"
	getOneGameItemInstance   = "get-game-item-instance"
	createGameItemInstance   = "create-game-item-instance"
	updateGameItemInstance   = "update-game-item-instance"
	deleteGameItemInstance   = "delete-game-item-instance"
)

func (rnr *Runner) gameItemInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameItemInstanceHandlerConfig")

	l.Debug("Adding game_item_instance handler configuration")

	cfg := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_item_instance.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_item_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_item_instance.response.schema.json",
		},
		References: referenceSchemas,
	}

	cfg[getManyGameItemInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-item-instances",
		HandlerFunc: rnr.getManyGameItemInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game_item_instance collection",
		},
	}

	cfg[getOneGameItemInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-item-instances/:game_item_instance_id",
		HandlerFunc: rnr.getGameItemInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game_item_instance",
		},
	}

	cfg[createGameItemInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/game-item-instances",
		HandlerFunc: rnr.createGameItemInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game_item_instance",
		},
	}

	cfg[updateGameItemInstance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/game-item-instances/:game_item_instance_id",
		HandlerFunc: rnr.updateGameItemInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game_item_instance",
		},
	}

	cfg[deleteGameItemInstance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/game-item-instances/:game_item_instance_id",
		HandlerFunc: rnr.deleteGameItemInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game_item_instance",
		},
	}

	return cfg, nil
}

func (rnr *Runner) getManyGameItemInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyGameItemInstancesHandler")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	recs, err := mm.GetManyGameItemInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting game_item_instance records >%v<", err)
		return err
	}
	res, err := mapper.GameItemInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) getGameItemInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetGameItemInstanceHandler")
	id := pp.ByName("game_item_instance_id")
	mm := m.(*domain.Domain)
	rec, err := mm.GetGameItemInstanceRec(id, nil)
	if err != nil {
		l.Warn("failed getting game_item_instance record >%v<", err)
		return err
	}
	res, err := mapper.GameItemInstanceRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game_item_instance record to response >%v<", err)
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) createGameItemInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateGameItemInstanceHandler")
	var req schema.GameItemInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}
	mm := m.(*domain.Domain)
	rec, err := mapper.GameItemInstanceRequestToRecord(l, &req, nil)
	if err != nil {
		return err
	}
	rec, err = mm.CreateGameItemInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_item_instance record >%v<", err)
		return err
	}
	res, err := mapper.GameItemInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) updateGameItemInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateGameItemInstanceHandler")
	id := pp.ByName("game_item_instance_id")
	mm := m.(*domain.Domain)
	existing, err := mm.GetGameItemInstanceRec(id, nil)
	if err != nil {
		l.Warn("failed getting game_item_instance record >%v<", err)
		return err
	}
	var req schema.GameItemInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}
	rec, err := mapper.GameItemInstanceRequestToRecord(l, &req, existing)
	if err != nil {
		return err
	}
	rec, err = mm.UpdateGameItemInstanceRec(rec)
	if err != nil {
		l.Warn("failed updating game_item_instance record >%v<", err)
		return err
	}
	res, err := mapper.GameItemInstanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) deleteGameItemInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteGameItemInstanceHandler")
	id := pp.ByName("game_item_instance_id")
	mm := m.(*domain.Domain)
	if err := mm.DeleteGameItemInstanceRec(id); err != nil {
		l.Warn("failed deleting game_item_instance record >%v<", err)
		return err
	}
	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}
