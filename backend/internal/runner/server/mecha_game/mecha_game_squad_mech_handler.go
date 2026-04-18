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

const (
	GetManyMechaGameSquadMechs  = "get-many-mecha-squad-mechs"
	GetOneMechaGameSquadMech    = "get-one-mecha-squad-mech"
	CreateOneMechaGameSquadMech = "create-one-mecha-squad-mech"
	UpdateOneMechaGameSquadMech = "update-one-mecha-squad-mech"
	DeleteOneMechaGameSquadMech = "delete-one-mecha-squad-mech"
)

func mechaGameSquadMechHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaGameSquadMechHandlerConfig")

	l.Debug("Adding mecha squad mech handler configuration")

	squadMechConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_squad_mech.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_squad_mech.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_squad_mech.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_squad_mech.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_squad_mech.schema.json",
		}),
	}

	squadMechConfig[GetManyMechaGameSquadMechs] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs",
		HandlerFunc: getManyMechaGameSquadMechsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mecha squad mechs"},
	}

	squadMechConfig[GetOneMechaGameSquadMech] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs/:mech_id",
		HandlerFunc: getOneMechaGameSquadMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mecha squad mech"},
	}

	squadMechConfig[CreateOneMechaGameSquadMech] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs",
		HandlerFunc: createOneMechaGameSquadMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mecha squad mech"},
	}

	squadMechConfig[UpdateOneMechaGameSquadMech] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs/:mech_id",
		HandlerFunc: updateOneMechaGameSquadMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mecha squad mech"},
	}

	squadMechConfig[DeleteOneMechaGameSquadMech] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/squads/:squad_id/mechs/:mech_id",
		HandlerFunc: deleteOneMechaGameSquadMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mecha squad mech"},
	}

	return squadMechConfig, nil
}

func getManyMechaGameSquadMechsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaGameSquadMechsHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	if gameID == "" || squadID == "" {
		return coreerror.NewParamError("game_id and squad_id are required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params,
		sql.Param{Col: mecha_game_record.FieldMechaGameSquadMechGameID, Val: gameID},
		sql.Param{Col: mecha_game_record.FieldMechaGameSquadMechMechaGameSquadID, Val: squadID},
	)

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechaGameSquadMechRecs(opts)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameSquadMechRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechaGameSquadMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaGameSquadMechHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaGameSquadMechRec(mechID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechaGameSquadID != squadID {
		return coreerror.NewNotFoundError("squad_mech", mechID)
	}

	res, err := mapper.MechaGameSquadMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechaGameSquadMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaGameSquadMechHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_game_record.MechaGameSquadMech{
		GameID:       gameID,
		MechaGameSquadID: squadID,
	}

	rec, err := mapper.MechaGameSquadMechRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaGameSquadMechRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameSquadMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechaGameSquadMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaGameSquadMechHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameSquadMechRec(mechID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechaGameSquadID != squadID {
		return coreerror.NewNotFoundError("squad_mech", mechID)
	}

	rec, err = mapper.MechaGameSquadMechRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaGameSquadMechRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameSquadMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechaGameSquadMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaGameSquadMechHandler")

	gameID := pp.ByName("game_id")
	squadID := pp.ByName("squad_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameSquadMechRec(mechID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechaGameSquadID != squadID {
		return coreerror.NewNotFoundError("squad_mech", mechID)
	}

	if err := mm.DeleteMechaGameSquadMechRec(mechID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
