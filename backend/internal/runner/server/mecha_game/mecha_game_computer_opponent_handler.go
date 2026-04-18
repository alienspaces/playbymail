package mecha_game

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
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/mecha-games/{game_id}/computer-opponents
// GET (document)    /api/v1/mecha-games/{game_id}/computer-opponents/{computer_opponent_id}
// POST (document)   /api/v1/mecha-games/{game_id}/computer-opponents
// PUT (document)    /api/v1/mecha-games/{game_id}/computer-opponents/{computer_opponent_id}
// DELETE (document) /api/v1/mecha-games/{game_id}/computer-opponents/{computer_opponent_id}

const (
	GetManyMechaGameComputerOpponents  = "get-many-mecha-computer-opponents"
	GetOneMechaGameComputerOpponent    = "get-one-mecha-computer-opponent"
	CreateOneMechaGameComputerOpponent = "create-one-mecha-computer-opponent"
	UpdateOneMechaGameComputerOpponent = "update-one-mecha-computer-opponent"
	DeleteOneMechaGameComputerOpponent = "delete-one-mecha-computer-opponent"
)

func mechaGameComputerOpponentHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaGameComputerOpponentHandlerConfig")
	l.Debug("Adding mecha computer opponent handler configuration")

	config := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_computer_opponent.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_computer_opponent.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_computer_opponent.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_computer_opponent.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_computer_opponent.schema.json",
		}),
	}

	config[GetManyMechaGameComputerOpponents] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents",
		HandlerFunc: getManyMechaGameComputerOpponentsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get mecha computer opponents",
		},
	}

	config[GetOneMechaGameComputerOpponent] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents/:computer_opponent_id",
		HandlerFunc: getOneMechaGameComputerOpponentHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get mecha computer opponent",
		},
	}

	config[CreateOneMechaGameComputerOpponent] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents",
		HandlerFunc: createOneMechaGameComputerOpponentHandler,
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
			Title:    "Create mecha computer opponent",
		},
	}

	config[UpdateOneMechaGameComputerOpponent] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents/:computer_opponent_id",
		HandlerFunc: updateOneMechaGameComputerOpponentHandler,
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
			Title:    "Update mecha computer opponent",
		},
	}

	config[DeleteOneMechaGameComputerOpponent] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents/:computer_opponent_id",
		HandlerFunc: deleteOneMechaGameComputerOpponentHandler,
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
			Title:    "Delete mecha computer opponent",
		},
	}

	return config, nil
}

func getManyMechaGameComputerOpponentsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaGameComputerOpponentsHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{
		Col: mecha_game_record.FieldMechaGameComputerOpponentGameID,
		Val: gameID,
	})

	mm := m.(*domain.Domain)

	recs, err := mm.GetManyMechaGameComputerOpponentRecs(opts)
	if err != nil {
		l.Warn("failed getting mecha computer opponent records >%v<", err)
		return err
	}

	res, err := mapper.MechaGameComputerOpponentRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneMechaGameComputerOpponentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaGameComputerOpponentHandler")

	gameID := pp.ByName("game_id")
	opponentID := pp.ByName("computer_opponent_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaGameComputerOpponentRec(opponentID, nil)
	if err != nil {
		l.Warn("failed getting mecha computer opponent record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("computer opponent", opponentID)
	}

	res, err := mapper.MechaGameComputerOpponentRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneMechaGameComputerOpponentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaGameComputerOpponentHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_game_record.MechaGameComputerOpponent{
		GameID: gameID,
	}

	rec, err := mapper.MechaGameComputerOpponentRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaGameComputerOpponentRec(rec)
	if err != nil {
		l.Warn("failed creating mecha computer opponent record >%v<", err)
		return err
	}

	res, err := mapper.MechaGameComputerOpponentRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneMechaGameComputerOpponentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaGameComputerOpponentHandler")

	gameID := pp.ByName("game_id")
	opponentID := pp.ByName("computer_opponent_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameComputerOpponentRec(opponentID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("computer opponent", opponentID)
	}

	rec, err = mapper.MechaGameComputerOpponentRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaGameComputerOpponentRec(rec)
	if err != nil {
		l.Warn("failed updating mecha computer opponent record >%v<", err)
		return err
	}

	res, err := mapper.MechaGameComputerOpponentRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneMechaGameComputerOpponentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaGameComputerOpponentHandler")

	gameID := pp.ByName("game_id")
	opponentID := pp.ByName("computer_opponent_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameComputerOpponentRec(opponentID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("computer opponent", opponentID)
	}

	if err := mm.DeleteMechaGameComputerOpponentRec(opponentID); err != nil {
		l.Warn("failed deleting mecha computer opponent record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
