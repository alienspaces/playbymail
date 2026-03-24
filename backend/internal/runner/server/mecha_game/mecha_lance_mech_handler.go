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

const (
	GetManyMechaLanceMechs  = "get-many-mecha-lance-mechs"
	GetOneMechaLanceMech    = "get-one-mecha-lance-mech"
	CreateOneMechaLanceMech = "create-one-mecha-lance-mech"
	UpdateOneMechaLanceMech = "update-one-mecha-lance-mech"
	DeleteOneMechaLanceMech = "delete-one-mecha-lance-mech"
)

func mechaLanceMechHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaLanceMechHandlerConfig")

	l.Debug("Adding mecha lance mech handler configuration")

	lanceMechConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_lance_mech.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_lance_mech.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_lance_mech.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_lance_mech.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_lance_mech.schema.json",
		}),
	}

	lanceMechConfig[GetManyMechaLanceMechs] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/lances/:lance_id/mechs",
		HandlerFunc: getManyMechaLanceMechsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mecha lance mechs"},
	}

	lanceMechConfig[GetOneMechaLanceMech] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/lances/:lance_id/mechs/:mech_id",
		HandlerFunc: getOneMechaLanceMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mecha lance mech"},
	}

	lanceMechConfig[CreateOneMechaLanceMech] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/lances/:lance_id/mechs",
		HandlerFunc: createOneMechaLanceMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mecha lance mech"},
	}

	lanceMechConfig[UpdateOneMechaLanceMech] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/lances/:lance_id/mechs/:mech_id",
		HandlerFunc: updateOneMechaLanceMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mecha lance mech"},
	}

	lanceMechConfig[DeleteOneMechaLanceMech] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/lances/:lance_id/mechs/:mech_id",
		HandlerFunc: deleteOneMechaLanceMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mecha lance mech"},
	}

	return lanceMechConfig, nil
}

func getManyMechaLanceMechsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaLanceMechsHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	if gameID == "" || lanceID == "" {
		return coreerror.NewParamError("game_id and lance_id are required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params,
		sql.Param{Col: mecha_record.FieldMechaLanceMechGameID, Val: gameID},
		sql.Param{Col: mecha_record.FieldMechaLanceMechMechaLanceID, Val: lanceID},
	)

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechaLanceMechRecs(opts)
	if err != nil {
		return err
	}

	res, err := mapper.MechaLanceMechRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechaLanceMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaLanceMechHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaLanceMechRec(mechID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechaLanceID != lanceID {
		return coreerror.NewNotFoundError("lance_mech", mechID)
	}

	res, err := mapper.MechaLanceMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechaLanceMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaLanceMechHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_record.MechaLanceMech{
		GameID:       gameID,
		MechaLanceID: lanceID,
	}

	rec, err := mapper.MechaLanceMechRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaLanceMechRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaLanceMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechaLanceMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaLanceMechHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaLanceMechRec(mechID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechaLanceID != lanceID {
		return coreerror.NewNotFoundError("lance_mech", mechID)
	}

	rec, err = mapper.MechaLanceMechRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaLanceMechRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaLanceMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechaLanceMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaLanceMechHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaLanceMechRec(mechID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechaLanceID != lanceID {
		return coreerror.NewNotFoundError("lance_mech", mechID)
	}

	if err := mm.DeleteMechaLanceMechRec(mechID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
