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
	tagGroupGameCreatureInstance server.TagGroup = "GameCreatureInstances"
	TagGameCreatureInstance      server.Tag      = "GameCreatureInstances"
)

const (
	getManyGameCreatureInstances = "get-game-creature-instances"
	getOneGameCreatureInstance   = "get-game-creature-instance"
	createGameCreatureInstance   = "create-game-creature-instance"
	updateGameCreatureInstance   = "update-game-creature-instance"
	deleteGameCreatureInstance   = "delete-game-creature-instance"
)

func (rnr *Runner) gameCreatureInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameCreatureInstanceHandlerConfig")

	l.Debug("Adding game_creature_instance handler configuration")

	cfg := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_creature_instance.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_creature_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_creature_instance.response.schema.json",
		},
		References: referenceSchemas,
	}

	cfg[getManyGameCreatureInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-creature-instances",
		HandlerFunc: rnr.getManyGameCreatureInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game_creature_instance collection",
		},
	}

	cfg[getOneGameCreatureInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-creature-instances/:game_creature_instance_id",
		HandlerFunc: rnr.getGameCreatureInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game_creature_instance",
		},
	}

	cfg[createGameCreatureInstance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/game-creature-instances",
		HandlerFunc: rnr.createGameCreatureInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game_creature_instance",
		},
	}

	cfg[updateGameCreatureInstance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/game-creature-instances/:game_creature_instance_id",
		HandlerFunc: rnr.updateGameCreatureInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game_creature_instance",
		},
	}

	cfg[deleteGameCreatureInstance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/game-creature-instances/:game_creature_instance_id",
		HandlerFunc: rnr.deleteGameCreatureInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game_creature_instance",
		},
	}

	return cfg, nil
}

func (rnr *Runner) getManyGameCreatureInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyGameCreatureInstancesHandler")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	recs, err := mm.GetManyGameCreatureInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting game_creature_instance records >%v<", err)
		return err
	}
	res := mapper.GameCreatureInstanceRecordsToCollectionResponse(recs)
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) getGameCreatureInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetGameCreatureInstanceHandler")
	id := pp.ByName("game_creature_instance_id")
	mm := m.(*domain.Domain)
	rec, err := mm.GetGameCreatureInstanceRec(id, nil)
	if err != nil {
		l.Warn("failed getting game_creature_instance record >%v<", err)
		return err
	}
	res := mapper.GameCreatureInstanceRecordToResponse(rec)
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) createGameCreatureInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateGameCreatureInstanceHandler")
	var req schema.GameCreatureInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}
	mm := m.(*domain.Domain)
	rec := mapper.GameCreatureInstanceRequestToRecord(&req)
	rec, err := mm.CreateGameCreatureInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_creature_instance record >%v<", err)
		return err
	}
	res := mapper.GameCreatureInstanceRecordToResponse(rec)
	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) updateGameCreatureInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateGameCreatureInstanceHandler")
	id := pp.ByName("game_creature_instance_id")
	mm := m.(*domain.Domain)
	existing, err := mm.GetGameCreatureInstanceRec(id, nil)
	if err != nil {
		l.Warn("failed getting game_creature_instance record >%v<", err)
		return err
	}
	var req schema.GameCreatureInstanceRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}
	// Update fields from request
	updated := mapper.GameCreatureInstanceRequestToRecord(&req)
	updated.ID = existing.ID
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedAt = existing.UpdatedAt
	updated.DeletedAt = existing.DeletedAt
	rec, err := mm.UpdateGameCreatureInstanceRec(updated)
	if err != nil {
		l.Warn("failed updating game_creature_instance record >%v<", err)
		return err
	}
	res := mapper.GameCreatureInstanceRecordToResponse(rec)
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) deleteGameCreatureInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteGameCreatureInstanceHandler")
	id := pp.ByName("game_creature_instance_id")
	mm := m.(*domain.Domain)
	if err := mm.DeleteGameCreatureInstanceRec(id); err != nil {
		l.Warn("failed deleting game_creature_instance record >%v<", err)
		return err
	}
	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}
