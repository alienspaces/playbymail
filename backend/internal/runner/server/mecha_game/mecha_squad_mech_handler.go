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
	GetManyMechaSquadMechs  = "get-many-mecha-squad-mechs"
	GetOneMechaSquadMech    = "get-one-mecha-squad-mech"
	CreateOneMechaSquadMech = "create-one-mecha-squad-mech"
	UpdateOneMechaSquadMech = "update-one-mecha-squad-mech"
	DeleteOneMechaSquadMech = "delete-one-mecha-squad-mech"
)

func mechaSquadMechHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaSquadMechHandlerConfig")

	l.Debug("Adding mecha squad mech handler configuration")

	squadMechConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_squad_mech.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_squad_mech.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_squad_mech.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_squad_mech.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_schema",
			Name:     "mecha_squad_mech.schema.json",
		}),
	}

	squadMechConfig[GetManyMechaSquadMechs] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs",
		HandlerFunc: getManyMechaSquadMechsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mecha squad mechs"},
	}

	squadMechConfig[GetOneMechaSquadMech] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs/:mech_id",
		HandlerFunc: getOneMechaSquadMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mecha squad mech"},
	}

	squadMechConfig[CreateOneMechaSquadMech] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs",
		HandlerFunc: createOneMechaSquadMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mecha squad mech"},
	}

	squadMechConfig[UpdateOneMechaSquadMech] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs/:mech_id",
		HandlerFunc: updateOneMechaSquadMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mecha squad mech"},
	}

	squadMechConfig[DeleteOneMechaSquadMech] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs/:mech_id",
		HandlerFunc: deleteOneMechaSquadMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mecha squad mech"},
	}

	return squadMechConfig, nil
}

func getManyMechaSquadMechsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaSquadMechsHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	if gameID == "" || squadID == "" {
		return coreerror.NewParamError("game_id and squad_id are required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params,
		sql.Param{Col: mecha_record.FieldMechaSquadMechGameID, Val: gameID},
		sql.Param{Col: mecha_record.FieldMechaSquadMechMechaSquadID, Val: squadID},
	)

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechaSquadMechRecs(opts)
	if err != nil {
		return err
	}

	res, err := mapper.MechaSquadMechRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechaSquadMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaSquadMechHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaSquadMechRec(mechID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechaSquadID != squadID {
		return coreerror.NewNotFoundError("squad_mech", mechID)
	}

	res, err := mapper.MechaSquadMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechaSquadMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaSquadMechHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_record.MechaSquadMech{
		GameID:       gameID,
		MechaSquadID: squadID,
	}

	rec, err := mapper.MechaSquadMechRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaSquadMechRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaSquadMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechaSquadMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaSquadMechHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaSquadMechRec(mechID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechaSquadID != squadID {
		return coreerror.NewNotFoundError("squad_mech", mechID)
	}

	rec, err = mapper.MechaSquadMechRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaSquadMechRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaSquadMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechaSquadMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaSquadMechHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaSquadMechRec(mechID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechaSquadID != squadID {
		return coreerror.NewNotFoundError("squad_mech", mechID)
	}

	if err := mm.DeleteMechaSquadMechRec(mechID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
