package mech_wargame

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
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	GetManyMechWargameLanceMechs  = "get-many-mech-wargame-lance-mechs"
	GetOneMechWargameLanceMech    = "get-one-mech-wargame-lance-mech"
	CreateOneMechWargameLanceMech = "create-one-mech-wargame-lance-mech"
	UpdateOneMechWargameLanceMech = "update-one-mech-wargame-lance-mech"
	DeleteOneMechWargameLanceMech = "delete-one-mech-wargame-lance-mech"
)

func mechWargameLanceMechHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechWargameLanceMechHandlerConfig")

	lanceMechConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance_mech.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance_mech.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance_mech.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance_mech.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance_mech.schema.json",
		}),
	}

	lanceMechConfig[GetManyMechWargameLanceMechs] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances/:lance_id/mechs",
		HandlerFunc: getManyMechWargameLanceMechsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mech wargame lance mechs"},
	}

	lanceMechConfig[GetOneMechWargameLanceMech] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances/:lance_id/mechs/:mech_id",
		HandlerFunc: getOneMechWargameLanceMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mech wargame lance mech"},
	}

	lanceMechConfig[CreateOneMechWargameLanceMech] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances/:lance_id/mechs",
		HandlerFunc: createOneMechWargameLanceMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mech wargame lance mech"},
	}

	lanceMechConfig[UpdateOneMechWargameLanceMech] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances/:lance_id/mechs/:mech_id",
		HandlerFunc: updateOneMechWargameLanceMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mech wargame lance mech"},
	}

	lanceMechConfig[DeleteOneMechWargameLanceMech] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances/:lance_id/mechs/:mech_id",
		HandlerFunc: deleteOneMechWargameLanceMechHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mech wargame lance mech"},
	}

	return lanceMechConfig, nil
}

func getManyMechWargameLanceMechsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechWargameLanceMechsHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	if gameID == "" || lanceID == "" {
		return coreerror.NewParamError("game_id and lance_id are required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params,
		sql.Param{Col: mech_wargame_record.FieldMechWargameLanceMechGameID, Val: gameID},
		sql.Param{Col: mech_wargame_record.FieldMechWargameLanceMechMechWargameLanceID, Val: lanceID},
	)

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechWargameLanceMechRecs(opts)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameLanceMechRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechWargameLanceMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechWargameLanceMechHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechWargameLanceMechRec(mechID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechWargameLanceID != lanceID {
		return coreerror.NewNotFoundError("lance_mech", mechID)
	}

	res, err := mapper.MechWargameLanceMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechWargameLanceMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechWargameLanceMechHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mech_wargame_record.MechWargameLanceMech{
		GameID:             gameID,
		MechWargameLanceID: lanceID,
	}

	rec, err := mapper.MechWargameLanceMechRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechWargameLanceMechRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameLanceMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechWargameLanceMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechWargameLanceMechHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameLanceMechRec(mechID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechWargameLanceID != lanceID {
		return coreerror.NewNotFoundError("lance_mech", mechID)
	}

	rec, err = mapper.MechWargameLanceMechRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechWargameLanceMechRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameLanceMechRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechWargameLanceMechHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechWargameLanceMechHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mechID := pp.ByName("mech_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameLanceMechRec(mechID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID || rec.MechWargameLanceID != lanceID {
		return coreerror.NewNotFoundError("lance_mech", mechID)
	}

	if err := mm.DeleteMechWargameLanceMechRec(mechID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
