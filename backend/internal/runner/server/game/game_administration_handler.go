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
	GetManyGameAdministrations = "get-game-administrations"
	GetOneGameAdministration   = "get-game-administration"
	CreateGameAdministration   = "create-game-administration"
	UpdateGameAdministration   = "update-game-administration"
	DeleteGameAdministration   = "delete-game-administration"
)

func gameAdministrationHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameAdministrationHandlerConfig")

	l.Debug("Adding game administration handler configuration")

	config := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "game_administration.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "game_administration.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api",
			Name:     "game_administration.response.schema.json",
		},
		References: referenceSchemas,
	}

	config[GetManyGameAdministrations] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-administrations",
		HandlerFunc: getManyGameAdministrationsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game administration collection",
		},
	}
	config[GetOneGameAdministration] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-administrations/:game_administration_id",
		HandlerFunc: getGameAdministrationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game administration",
		},
	}
	config[CreateGameAdministration] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-administrations",
		HandlerFunc: createGameAdministrationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game administration",
		},
	}
	config[UpdateGameAdministration] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/game-administrations/:game_administration_id",
		HandlerFunc: updateGameAdministrationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game administration",
		},
	}
	config[DeleteGameAdministration] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/game-administrations/:game_administration_id",
		HandlerFunc: deleteGameAdministrationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game administration",
		},
	}

	return config, nil
}

func getManyGameAdministrationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameAdministrationsHandler")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	recs, err := mm.GetManyGameAdministrationRecs(opts)
	if err != nil {
		return err
	}
	res, err := mapper.GameAdministrationRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getGameAdministrationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameAdministrationHandler")
	mm := m.(*domain.Domain)
	recID := pp.ByName("game_administration_id")
	rec, err := mm.GetGameAdministrationRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	res, err := mapper.GameAdministrationRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createGameAdministrationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createGameAdministrationHandler")
	mm := m.(*domain.Domain)
	rec, err := mapper.GameAdministrationRequestToRecord(l, r, &game_record.GameAdministration{})
	if err != nil {
		return err
	}
	rec, err = mm.CreateGameAdministrationRec(rec)
	if err != nil {
		return err
	}
	res, err := mapper.GameAdministrationRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateGameAdministrationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateGameAdministrationHandler")
	mm := m.(*domain.Domain)
	recID := pp.ByName("game_administration_id")
	rec, err := mm.GetGameAdministrationRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	rec, err = mapper.GameAdministrationRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}
	rec, err = mm.UpdateGameAdministrationRec(rec)
	if err != nil {
		return err
	}
	res, err := mapper.GameAdministrationRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteGameAdministrationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteGameAdministrationHandler")
	mm := m.(*domain.Domain)
	recID := pp.ByName("game_administration_id")
	rec, err := mm.GetGameAdministrationRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if err := mm.DeleteGameAdministrationRec(rec.ID); err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
