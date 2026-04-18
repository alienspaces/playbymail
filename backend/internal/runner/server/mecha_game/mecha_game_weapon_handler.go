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
	GetManyMechaGameWeapons  = "get-many-mecha-weapons"
	GetOneMechaGameWeapon    = "get-one-mecha-weapon"
	CreateOneMechaGameWeapon = "create-one-mecha-weapon"
	UpdateOneMechaGameWeapon = "update-one-mecha-weapon"
	DeleteOneMechaGameWeapon = "delete-one-mecha-weapon"
)

func mechaGameWeaponHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaGameWeaponHandlerConfig")

	l.Debug("Adding mecha weapon handler configuration")

	weaponConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_weapon.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_weapon.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_weapon.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_weapon.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_weapon.schema.json",
		}),
	}

	weaponConfig[GetManyMechaGameWeapons] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/weapons",
		HandlerFunc: getManyMechaGameWeaponsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mecha weapons"},
	}

	weaponConfig[GetOneMechaGameWeapon] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/weapons/:weapon_id",
		HandlerFunc: getOneMechaGameWeaponHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mecha weapon"},
	}

	weaponConfig[CreateOneMechaGameWeapon] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/weapons",
		HandlerFunc: createOneMechaGameWeaponHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mecha weapon"},
	}

	weaponConfig[UpdateOneMechaGameWeapon] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/weapons/:weapon_id",
		HandlerFunc: updateOneMechaGameWeaponHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mecha weapon"},
	}

	weaponConfig[DeleteOneMechaGameWeapon] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/weapons/:weapon_id",
		HandlerFunc: deleteOneMechaGameWeaponHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mecha weapon"},
	}

	return weaponConfig, nil
}

func getManyMechaGameWeaponsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaGameWeaponsHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{Col: mecha_game_record.FieldMechaGameWeaponGameID, Val: gameID})

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechaGameWeaponRecs(opts)
	if err != nil {
		l.Warn("failed getting mecha weapon records >%v<", err)
		return err
	}

	res, err := mapper.MechaGameWeaponRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechaGameWeaponHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaGameWeaponHandler")

	gameID := pp.ByName("game_id")
	weaponID := pp.ByName("weapon_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaGameWeaponRec(weaponID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("weapon", weaponID)
	}

	res, err := mapper.MechaGameWeaponRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechaGameWeaponHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaGameWeaponHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_game_record.MechaGameWeapon{GameID: gameID}
	rec, err := mapper.MechaGameWeaponRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaGameWeaponRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameWeaponRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechaGameWeaponHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaGameWeaponHandler")

	gameID := pp.ByName("game_id")
	weaponID := pp.ByName("weapon_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameWeaponRec(weaponID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("weapon", weaponID)
	}

	rec, err = mapper.MechaGameWeaponRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaGameWeaponRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameWeaponRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechaGameWeaponHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaGameWeaponHandler")

	gameID := pp.ByName("game_id")
	weaponID := pp.ByName("weapon_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameWeaponRec(weaponID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("weapon", weaponID)
	}

	if err := mm.DeleteMechaGameWeaponRec(weaponID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
