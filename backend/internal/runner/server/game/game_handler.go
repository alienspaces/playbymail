package game

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	GetManyGames     = "get-many-games"
	GetOneGame       = "get-one-game"
	CreateOneGame    = "create-one-game"
	CreateGameWithID = "create-game-with-id"
	UpdateOneGame    = "update-one-game"
	DeleteOneGame    = "delete-one-game"
	PublishGame      = "publish-game"
)

func gameHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameHandlerConfig")

	l.Debug("Adding game handler configuration")

	// Create a new map to avoid modifying the passed config
	gameConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game.schema.json",
			},
		}...),
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

	gameConfig[CreateOneGame] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games",
		HandlerFunc: createGameHandler,
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

	gameConfig[UpdateOneGame] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/games/:game_id",
		HandlerFunc: updateGameHandler,
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
			Title:    "Update game",
		},
	}

	gameConfig[DeleteOneGame] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/games/:game_id",
		HandlerFunc: deleteGameHandler,
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
			Title:    "Delete game",
		},
	}

	gameConfig[PublishGame] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/games/:game_id/publish",
		HandlerFunc: publishGameHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameDesign,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Publish game",
			Description: "Publish a draft game, making it immutable and visible to everyone. Once published, games cannot be modified.",
		},
	}

	return gameConfig, nil
}

// getManyGamesHandler queries the account_game_view for the authenticated user's
// visible games, with optional filters for subscription role and status.
func getManyGamesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "GetManyGamesHandler")

	l.Info("querying many game records with params >%#v<", qp)

	authenData := server.GetRequestAuthenData(l, r)
	if authenData == nil || authenData.AccountUser.ID == "" {
		l.Warn("authenticated account is required to list games")
		return coreerror.NewUnauthorizedError()
	}

	accountID := authenData.AccountUser.AccountID
	l.Info("listing games for account >%s<", accountID)

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, coresql.Param{
		Col: game_record.FieldAccountGameViewAccountID,
		Val: accountID,
	})

	// Boolean view filters: is_designer, is_manager, can_manage
	if vals, ok := qp.Params["is_designer"]; ok && len(vals) > 0 && vals[0].Val.(string) == "true" {
		l.Info("filtering games by is_designer=true")
		opts.Params = append(opts.Params, coresql.Param{
			Col: game_record.FieldAccountGameViewIsDesigner,
			Val: true,
		})
	}
	if vals, ok := qp.Params["is_manager"]; ok && len(vals) > 0 && vals[0].Val.(string) == "true" {
		l.Info("filtering games by is_manager=true")
		opts.Params = append(opts.Params, coresql.Param{
			Col: game_record.FieldAccountGameViewIsManager,
			Val: true,
		})
	}
	if vals, ok := qp.Params["can_manage"]; ok && len(vals) > 0 && vals[0].Val.(string) == "true" {
		l.Info("filtering games by can_manage=true")
		opts.Params = append(opts.Params, coresql.Param{
			Col: game_record.FieldAccountGameViewCanManage,
			Val: true,
		})
	}

	// Support status filter (maps to game_status in the view)
	if statusValues, ok := qp.Params["status"]; ok && len(statusValues) > 0 {
		statusFilter := statusValues[0].Val.(string)
		if statusFilter == game_record.GameStatusDraft || statusFilter == game_record.GameStatusPublished {
			opts.Params = append(opts.Params, coresql.Param{
				Col: game_record.FieldAccountGameViewGameStatus,
				Val: statusFilter,
			})
			l.Info("filtering games by status >%s<", statusFilter)
		}
	}

	recs, err := mm.GetManyAccountGameViewRecs(opts)
	if err != nil {
		l.Warn("failed getting account game view records >%v<", err)
		return err
	}

	res, err := mapper.AccountGameViewRecordsToCollectionResponse(l, recs)
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
	if recID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

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

	// Get authenticated account ID
	authenData := server.GetRequestAuthenData(l, r)
	if authenData == nil || authenData.AccountUser.ID == "" {
		l.Warn("authenticated account is required to create game")
		return coreerror.NewUnauthorizedError()
	}

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

	// Auto-create designer subscription so the creator can access their game
	_, err = mm.CreateDesignerSubscriptionForNewGame(rec, authenData.AccountUser.AccountID)
	if err != nil {
		l.Warn("failed creating designer subscription for game >%s< >%v<", rec.ID, err)
		return err
	}

	// Auto-create manager subscription so the designer can create instances
	_, err = mm.CreateManagerSubscriptionForNewGame(rec, authenData.AccountUser.AccountID)
	if err != nil {
		l.Warn("failed creating manager subscription for game >%s< >%v<", rec.ID, err)
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
	if recID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}
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
	if recID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}
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

func publishGameHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "publishGameHandler")

	l.Info("publishing game with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game ID is empty")
		return coreerror.RequiredPathParameter("game_id")
	}

	// Get the current game record
	currRec, err := mm.GetGameRec(gameID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get game record >%v<", err)
		return err
	}

	// Update status to published
	currRec.Status = game_record.GameStatusPublished

	// Update the game (validation will check status transition)
	rec, err := mm.UpdateGameRec(currRec)
	if err != nil {
		l.Warn("failed publishing game >%v<", err)
		return err
	}

	res, err := mapper.GameRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game record to response >%v<", err)
		return err
	}

	l.Info("responding with published game record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
