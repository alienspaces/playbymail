package game

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	GetManyGames     = "get-games"
	GetOneGame       = "get-game"
	CreateGame       = "create-game"
	CreateGameWithID = "create-game-with-id"
	UpdateGame       = "update-game"
	DeleteGame       = "delete-game"
)

func gameHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameHandlerConfig")

	l.Debug("Adding game handler configuration")

	// Create a new map to avoid modifying the passed config
	gameConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "game.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "game.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "game.response.schema.json",
		},
		References: referenceSchemas,
	}

	// Unnested routes
	gameConfig[GetManyGames] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games",
		HandlerFunc: getManyGamesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game collection",
		},
	}

	gameConfig[GetOneGame] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/games/:game_id",
		HandlerFunc: getGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game",
		},
	}

	gameConfig[CreateGame] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games",
		HandlerFunc: createGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game",
		},
	}

	gameConfig[CreateGameWithID] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id",
		HandlerFunc: createGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game with ID",
		},
	}

	gameConfig[UpdateGame] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/games/:game_id",
		HandlerFunc: updateGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game",
		},
	}

	gameConfig[DeleteGame] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/games/:game_id",
		HandlerFunc: deleteGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
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
func getManyGamesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "GetManyGamesHandler")

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

func getGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "GetGameHandler")

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

func createGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "CreateGameHandler")

	l.Info("creating game record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	rec, err := mapper.GameRequestToRecord(l, r, &game_record.Game{})
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

func updateGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "UpdateGameHandler")

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

func deleteGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "DeleteGameHandler")

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
