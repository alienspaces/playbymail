package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// The TagGroup and Tag value is used for the menu headings in the generated API documentation.
const (
	tagGroupGame server.TagGroup = "Games"
	TagGame      server.Tag      = "Games"
)

const (
	getManyGames     = "get-games"
	getOneGame       = "get-game"
	createGame       = "create-game"
	createGameWithID = "create-game-with-id"
	updateGame       = "update-game"
	deleteGame       = "delete-game"
)

func (rnr *Runner) gameHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameHandlerConfig")

	l.Debug("Adding game handler configuration")

	// Create a new map to avoid modifying the passed config
	gameConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game.response.schema.json",
		},
		References: referenceSchemas,
	}

	// Unnested routes
	gameConfig[getManyGames] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/games",
		HandlerFunc: rnr.getManyGamesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get client collection",
		},
	}

	gameConfig[getOneGame] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/games/:game_id",
		HandlerFunc: rnr.getGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game",
		},
	}

	gameConfig[createGame] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/games",
		HandlerFunc: rnr.createGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game",
		},
	}

	gameConfig[createGameWithID] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/games/:game_id",
		HandlerFunc: rnr.createGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game with ID",
		},
	}

	gameConfig[updateGame] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/games/:game_id",
		HandlerFunc: rnr.updateGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game",
		},
	}

	gameConfig[deleteGame] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/games/:game_id",
		HandlerFunc: rnr.deleteGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game",
		},
	}

	return gameConfig, nil
}

// GetManyGamesHandler -
func (rnr *Runner) getManyGamesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyGamesHandler")

	l.Info("querying many game records with params >%#v<", qp)

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyGameRecs(opts)
	if err != nil {
		l.Warn("failed getting game records >%v<", err)
		return err
	}

	res, err := mapper.GameRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	l.Info("responding with >%d< game records", len(res.Data))

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetGameHandler")

	l.Info("querying game record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	recID := pp.ByName("game_id")

	rec, err := mm.GetGameRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed getting game record >%v<", err)
		return err
	}

	res, err := mapper.GameRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game record to response >%v<", err)
		return err
	}

	l.Info("responding with game record id >%s<", rec.ID)

	err = server.WriteResponse(l, w, http.StatusOK, res)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) createGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateGameHandler")

	l.Info("creating game record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	rec, err := mapper.GameRequestToRecord(l, r, &record.Game{})
	if err != nil {
		return err
	}

	rec, err = mm.CreateGameRec(rec)
	if err != nil {
		l.Warn("failed creating game record >%v<", err)
		return err
	}

	res, err := mapper.GameRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	l.Info("responding with created game record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) updateGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateGameHandler")

	l.Info("updating game record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	recID := pp.ByName("game_id")
	rec, err := mm.GetGameRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	rec, err = mapper.GameRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateGameRec(rec)
	if err != nil {
		l.Warn("failed updating game record >%v<", err)
		return err
	}

	res, err := mapper.GameRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	l.Info("responding with updated game record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) deleteGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteGameHandler")

	l.Info("deleting game record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	recID := pp.ByName("game_id")
	rec, err := mm.GetGameRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	err = mm.DeleteGameRec(rec.ID)
	if err != nil {
		l.Warn("failed deleting game record >%v<", err)
		return err
	}

	err = server.WriteResponse(l, w, http.StatusNoContent, nil)
	if err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
