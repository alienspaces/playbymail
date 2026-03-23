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
	GetManyMechWargameLances  = "get-many-mech-wargame-lances"
	GetOneMechWargameLance    = "get-one-mech-wargame-lance"
	CreateOneMechWargameLance = "create-one-mech-wargame-lance"
	UpdateOneMechWargameLance = "update-one-mech-wargame-lance"
	DeleteOneMechWargameLance = "delete-one-mech-wargame-lance"
)

func mechWargameLanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechWargameLanceHandlerConfig")

	lanceConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_lance.schema.json",
		}),
	}

	lanceConfig[GetManyMechWargameLances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances",
		HandlerFunc: getManyMechWargameLancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mech wargame lances"},
	}

	lanceConfig[GetOneMechWargameLance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances/:lance_id",
		HandlerFunc: getOneMechWargameLanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mech wargame lance"},
	}

	lanceConfig[CreateOneMechWargameLance] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances",
		HandlerFunc: createOneMechWargameLanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mech wargame lance"},
	}

	lanceConfig[UpdateOneMechWargameLance] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances/:lance_id",
		HandlerFunc: updateOneMechWargameLanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mech wargame lance"},
	}

	lanceConfig[DeleteOneMechWargameLance] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mech-wargame-games/:game_id/lances/:lance_id",
		HandlerFunc: deleteOneMechWargameLanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mech wargame lance"},
	}

	return lanceConfig, nil
}

func getManyMechWargameLancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechWargameLancesHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{Col: mech_wargame_record.FieldMechWargameLanceGameID, Val: gameID})

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechWargameLanceRecs(opts)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameLanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechWargameLanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechWargameLanceHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechWargameLanceRec(lanceID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("lance", lanceID)
	}

	res, err := mapper.MechWargameLanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechWargameLanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechWargameLanceHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mech_wargame_record.MechWargameLance{GameID: gameID}
	rec, err := mapper.MechWargameLanceRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	// Resolve AccountID from AccountUserID
	if rec.AccountUserID != "" {
		auth := server.GetRequestAuthenData(l, r)
		if auth != nil && auth.AccountUser.ID == rec.AccountUserID {
			rec.AccountID = auth.AccountUser.AccountID
		} else {
			accountUserRec, err := mm.GetAccountUserRec(rec.AccountUserID, nil)
			if err != nil {
				l.Warn("failed getting account user record >%s< >%v<", rec.AccountUserID, err)
				return err
			}
			rec.AccountID = accountUserRec.AccountID
		}
	} else {
		return coreerror.NewParamError("account_user_id is required")
	}

	rec, err = mm.CreateMechWargameLanceRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameLanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechWargameLanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechWargameLanceHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameLanceRec(lanceID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("lance", lanceID)
	}

	rec, err = mapper.MechWargameLanceRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	// Resolve AccountID from AccountUserID
	if rec.AccountUserID != "" {
		auth := server.GetRequestAuthenData(l, r)
		if auth != nil && auth.AccountUser.ID == rec.AccountUserID {
			rec.AccountID = auth.AccountUser.AccountID
		} else {
			accountUserRec, err := mm.GetAccountUserRec(rec.AccountUserID, nil)
			if err != nil {
				l.Warn("failed getting account user record >%s< >%v<", rec.AccountUserID, err)
				return err
			}
			rec.AccountID = accountUserRec.AccountID
		}
	}

	rec, err = mm.UpdateMechWargameLanceRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameLanceRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechWargameLanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechWargameLanceHandler")

	gameID := pp.ByName("game_id")
	lanceID := pp.ByName("lance_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameLanceRec(lanceID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("lance", lanceID)
	}

	if err := mm.DeleteMechWargameLanceRec(lanceID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
