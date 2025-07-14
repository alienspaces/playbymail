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
	getManyGameLocationInstances = "get-game-location-instances"
	getOneGameLocationInstance   = "get-game-location-instance"
	createGameLocationInstance   = "create-game-location-instance"
	updateGameLocationInstance   = "update-game-location-instance"
	deleteGameLocationInstance   = "delete-game-location-instance"
)

func (rnr *Runner) gameLocationInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameLocationInstanceHandlerConfig")

	l.Debug("adding game_location_instance handler configuration")

	cfg := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_location_instance.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_location_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_location_instance.response.schema.json",
		},
		References: referenceSchemas,
	}

	cfg[getManyGameLocationInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-location-instances",
		HandlerFunc: rnr.getManyGameLocationInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game_location_instance collection",
		},
	}

	cfg[getOneGameLocationInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-location-instances/:game_location_instance_id",
		HandlerFunc: rnr.getGameLocationInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game_location_instance",
		},
	}

	cfg[createGameLocationInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/game-location-instances",
		HandlerFunc: rnr.createGameLocationInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game_location_instance",
		},
	}

	cfg[updateGameLocationInstance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/game-location-instances/:game_location_instance_id",
		HandlerFunc: rnr.updateGameLocationInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game_location_instance",
		},
	}

	cfg[deleteGameLocationInstance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/game-location-instances/:game_location_instance_id",
		HandlerFunc: rnr.deleteGameLocationInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game_location_instance",
		},
	}

	return cfg, nil
}

func (rnr *Runner) getManyGameLocationInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "getManyGameLocationInstancesHandler")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	recs, err := mm.GetManyGameLocationInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting game_location_instance records >%v<", err)
		return err
	}
	res, err := mapper.GameLocationInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) getGameLocationInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "getGameLocationInstanceHandler")
	id := pp.ByName("game_location_instance_id")
	mm := m.(*domain.Domain)
	rec, err := mm.GetGameLocationInstanceRec(id, nil)
	if err != nil {
		l.Warn("failed getting game_location_instance record >%v<", err)
		return err
	}
	res, err := mapper.GameLocationInstanceRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game_location_instance record to response >%v<", err)
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) createGameLocationInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "createGameLocationInstanceHandler")
	var req schema.GameLocationInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}
	mm := m.(*domain.Domain)
	rec, err := mapper.GameLocationInstanceRequestToRecord(l, &req, nil)
	if err != nil {
		return err
	}
	rec, err = mm.CreateGameLocationInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_location_instance record >%v<", err)
		return err
	}
	res, err := mapper.GameLocationInstanceRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game_location_instance record to response >%v<", err)
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) updateGameLocationInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "updateGameLocationInstanceHandler")
	id := pp.ByName("game_location_instance_id")
	var req schema.GameLocationInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}
	mm := m.(*domain.Domain)
	rec, err := mm.GetGameLocationInstanceRec(id, nil)
	if err != nil {
		l.Warn("failed getting game_location_instance record >%v<", err)
		return err
	}
	rec, err = mapper.GameLocationInstanceRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}
	rec, err = mm.UpdateGameLocationInstanceRec(rec)
	if err != nil {
		l.Warn("failed updating game_location_instance record >%v<", err)
		return err
	}
	res, err := mapper.GameLocationInstanceRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game_location_instance record to response >%v<", err)
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) deleteGameLocationInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "deleteGameLocationInstanceHandler")
	id := pp.ByName("game_location_instance_id")
	mm := m.(*domain.Domain)
	err := mm.DeleteGameLocationInstanceRec(id)
	if err != nil {
		l.Warn("failed deleting game_location_instance record >%v<", err)
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}
