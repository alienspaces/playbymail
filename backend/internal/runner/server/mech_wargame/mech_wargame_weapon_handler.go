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
	GetManyMechWargameWeapons  = "get-many-mech-wargame-weapons"
	GetOneMechWargameWeapon    = "get-one-mech-wargame-weapon"
	CreateOneMechWargameWeapon = "create-one-mech-wargame-weapon"
	UpdateOneMechWargameWeapon = "update-one-mech-wargame-weapon"
	DeleteOneMechWargameWeapon = "delete-one-mech-wargame-weapon"
)

func mechWargameWeaponHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechWargameWeaponHandlerConfig")

	weaponConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_weapon.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_weapon.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_weapon.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_weapon.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_weapon.schema.json",
		}),
	}

	weaponConfig[GetManyMechWargameWeapons] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/weapons",
		HandlerFunc: getManyMechWargameWeaponsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mech wargame weapons"},
	}

	weaponConfig[GetOneMechWargameWeapon] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/weapons/:weapon_id",
		HandlerFunc: getOneMechWargameWeaponHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mech wargame weapon"},
	}

	weaponConfig[CreateOneMechWargameWeapon] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mech-wargame-games/:game_id/weapons",
		HandlerFunc: createOneMechWargameWeaponHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mech wargame weapon"},
	}

	weaponConfig[UpdateOneMechWargameWeapon] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mech-wargame-games/:game_id/weapons/:weapon_id",
		HandlerFunc: updateOneMechWargameWeaponHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mech wargame weapon"},
	}

	weaponConfig[DeleteOneMechWargameWeapon] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mech-wargame-games/:game_id/weapons/:weapon_id",
		HandlerFunc: deleteOneMechWargameWeaponHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mech wargame weapon"},
	}

	return weaponConfig, nil
}

func getManyMechWargameWeaponsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechWargameWeaponsHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{Col: mech_wargame_record.FieldMechWargameWeaponGameID, Val: gameID})

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechWargameWeaponRecs(opts)
	if err != nil {
		l.Warn("failed getting mech wargame weapon records >%v<", err)
		return err
	}

	res, err := mapper.MechWargameWeaponRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechWargameWeaponHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechWargameWeaponHandler")

	gameID := pp.ByName("game_id")
	weaponID := pp.ByName("weapon_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechWargameWeaponRec(weaponID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("weapon", weaponID)
	}

	res, err := mapper.MechWargameWeaponRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechWargameWeaponHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechWargameWeaponHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mech_wargame_record.MechWargameWeapon{GameID: gameID}
	rec, err := mapper.MechWargameWeaponRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechWargameWeaponRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameWeaponRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechWargameWeaponHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechWargameWeaponHandler")

	gameID := pp.ByName("game_id")
	weaponID := pp.ByName("weapon_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameWeaponRec(weaponID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("weapon", weaponID)
	}

	rec, err = mapper.MechWargameWeaponRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechWargameWeaponRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameWeaponRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechWargameWeaponHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechWargameWeaponHandler")

	gameID := pp.ByName("game_id")
	weaponID := pp.ByName("weapon_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameWeaponRec(weaponID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("weapon", weaponID)
	}

	if err := mm.DeleteMechWargameWeaponRec(weaponID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
