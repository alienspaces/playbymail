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

const (
	GetManyGameAdministrations = "get-game-administrations"
	GetOneGameAdministration   = "get-game-administration"
	CreateGameAdministration   = "create-game-administration"
	UpdateGameAdministration   = "update-game-administration"
	DeleteGameAdministration   = "delete-game-administration"
)

func (rnr *Runner) gameAdministrationHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameAdministrationHandlerConfig")

	l.Debug("Adding game administration handler configuration")

	config := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main:       jsonschema.Schema{Name: "game_administration.collection.response.schema.json"},
		References: referenceSchemas,
	}
	requestSchema := jsonschema.SchemaWithReferences{
		Main:       jsonschema.Schema{Name: "game_administration.request.schema.json"},
		References: referenceSchemas,
	}
	responseSchema := jsonschema.SchemaWithReferences{
		Main:       jsonschema.Schema{Name: "game_administration.response.schema.json"},
		References: referenceSchemas,
	}

	config[GetManyGameAdministrations] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-administrations",
		HandlerFunc: rnr.getManyGameAdministrationsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
	}
	config[GetOneGameAdministration] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-administrations/:game_administration_id",
		HandlerFunc: rnr.getGameAdministrationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
	}
	config[CreateGameAdministration] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/game-administrations",
		HandlerFunc: rnr.createGameAdministrationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
	}
	config[UpdateGameAdministration] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/game-administrations/:game_administration_id",
		HandlerFunc: rnr.updateGameAdministrationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
	}
	config[DeleteGameAdministration] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/game-administrations/:game_administration_id",
		HandlerFunc: rnr.deleteGameAdministrationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
		},
	}

	return config, nil
}

func (rnr *Runner) getManyGameAdministrationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "getManyGameAdministrationsHandler")
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

func (rnr *Runner) getGameAdministrationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "getGameAdministrationHandler")
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

func (rnr *Runner) createGameAdministrationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "createGameAdministrationHandler")
	mm := m.(*domain.Domain)
	rec, err := mapper.GameAdministrationRequestToRecord(l, r, &record.GameAdministration{})
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

func (rnr *Runner) updateGameAdministrationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "updateGameAdministrationHandler")
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

func (rnr *Runner) deleteGameAdministrationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "deleteGameAdministrationHandler")
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
