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
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
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
	GetManyMechaComputerOpponents  = "get-many-mecha-computer-opponents"
	GetOneMechaComputerOpponent    = "get-one-mecha-computer-opponent"
	CreateOneMechaComputerOpponent = "create-one-mecha-computer-opponent"
	UpdateOneMechaComputerOpponent = "update-one-mecha-computer-opponent"
	DeleteOneMechaComputerOpponent = "delete-one-mecha-computer-opponent"
)

func mechaComputerOpponentHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaComputerOpponentHandlerConfig")
	l.Debug("Adding mecha computer opponent handler configuration")

	config := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_computer_opponent.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_computer_opponent.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_computer_opponent.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_computer_opponent.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_computer_opponent.schema.json",
		}),
	}

	config[GetManyMechaComputerOpponents] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents",
		HandlerFunc: getManyMechaComputerOpponentsHandler,
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

	config[GetOneMechaComputerOpponent] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents/:computer_opponent_id",
		HandlerFunc: getOneMechaComputerOpponentHandler,
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

	config[CreateOneMechaComputerOpponent] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents",
		HandlerFunc: createOneMechaComputerOpponentHandler,
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

	config[UpdateOneMechaComputerOpponent] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents/:computer_opponent_id",
		HandlerFunc: updateOneMechaComputerOpponentHandler,
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

	config[DeleteOneMechaComputerOpponent] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/computer-opponents/:computer_opponent_id",
		HandlerFunc: deleteOneMechaComputerOpponentHandler,
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

func getManyMechaComputerOpponentsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaComputerOpponentsHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{
		Col: mecha_record.FieldMechaComputerOpponentGameID,
		Val: gameID,
	})

	mm := m.(*domain.Domain)

	recs, err := mm.GetManyMechaComputerOpponentRecs(opts)
	if err != nil {
		l.Warn("failed getting mecha computer opponent records >%v<", err)
		return err
	}

	res, err := mapper.MechaComputerOpponentRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneMechaComputerOpponentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaComputerOpponentHandler")

	gameID := pp.ByName("game_id")
	opponentID := pp.ByName("computer_opponent_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaComputerOpponentRec(opponentID, nil)
	if err != nil {
		l.Warn("failed getting mecha computer opponent record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("computer opponent", opponentID)
	}

	res, err := mapper.MechaComputerOpponentRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneMechaComputerOpponentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaComputerOpponentHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_record.MechaComputerOpponent{
		GameID: gameID,
	}

	rec, err := mapper.MechaComputerOpponentRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaComputerOpponentRec(rec)
	if err != nil {
		l.Warn("failed creating mecha computer opponent record >%v<", err)
		return err
	}

	res, err := mapper.MechaComputerOpponentRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneMechaComputerOpponentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaComputerOpponentHandler")

	gameID := pp.ByName("game_id")
	opponentID := pp.ByName("computer_opponent_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaComputerOpponentRec(opponentID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("computer opponent", opponentID)
	}

	rec, err = mapper.MechaComputerOpponentRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaComputerOpponentRec(rec)
	if err != nil {
		l.Warn("failed updating mecha computer opponent record >%v<", err)
		return err
	}

	res, err := mapper.MechaComputerOpponentRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneMechaComputerOpponentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaComputerOpponentHandler")

	gameID := pp.ByName("game_id")
	opponentID := pp.ByName("computer_opponent_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaComputerOpponentRec(opponentID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("computer opponent", opponentID)
	}

	if err := mm.DeleteMechaComputerOpponentRec(opponentID); err != nil {
		l.Warn("failed deleting mecha computer opponent record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
