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
	tagGroupGameCreature server.TagGroup = "GameCreatures"
	TagGameCreature      server.Tag      = "GameCreatures"
)

const (
	getManyGameCreatures = "get-game-creatures"
	getOneGameCreature   = "get-game-creature"
	createGameCreature   = "create-game-creature"
	updateGameCreature   = "update-game-creature"
	deleteGameCreature   = "delete-game-creature"
)

func (rnr *Runner) gameCreatureHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameCreatureHandlerConfig")

	l.Debug("Adding game_creature handler configuration")

	gameCreatureConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_creature.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_creature.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_creature.response.schema.json",
		},
		References: referenceSchemas,
	}

	gameCreatureConfig[getManyGameCreatures] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-creatures",
		HandlerFunc: rnr.getManyGameCreaturesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
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
				server.AuthenticationTypeAPIKey,
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
				server.AuthenticationTypeAPIKey,
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
				server.AuthenticationTypeAPIKey,
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
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game_creature",
		},
	}

	return gameCreatureConfig, nil
}

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

func (rnr *Runner) updateGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateGameCreatureHandler")

	gameCreatureID := pp.ByName("game_creature_id")

	l.Info("updating game_creature record with path params >%#v<", pp)

	var req schema.GameCreatureRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetGameCreatureRec(gameCreatureID, nil)
	if err != nil {
		return err
	}

	rec, err = mapper.GameCreatureRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateGameCreatureRec(rec)
	if err != nil {
		l.Warn("failed updating game_creature record >%v<", err)
		return err
	}

	data, err := mapper.GameCreatureRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.GameCreatureResponse{
		Data: &data,
	}

	l.Info("responding with updated game_creature record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) deleteGameCreatureHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteGameCreatureHandler")

	gameCreatureID := pp.ByName("game_creature_id")

	l.Info("deleting game_creature record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	if err := mm.DeleteGameCreatureRec(gameCreatureID); err != nil {
		l.Warn("failed deleting game_creature record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
