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

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/mech-wargame-games/{game_id}/chassis
// GET (document)    /api/v1/mech-wargame-games/{game_id}/chassis/{chassis_id}
// POST (document)   /api/v1/mech-wargame-games/{game_id}/chassis
// PUT (document)    /api/v1/mech-wargame-games/{game_id}/chassis/{chassis_id}
// DELETE (document) /api/v1/mech-wargame-games/{game_id}/chassis/{chassis_id}

const (
	GetManyMechWargameChassis  = "get-many-mech-wargame-chassis"
	GetOneMechWargameChassis   = "get-one-mech-wargame-chassis"
	CreateOneMechWargameChassis = "create-one-mech-wargame-chassis"
	UpdateOneMechWargameChassis = "update-one-mech-wargame-chassis"
	DeleteOneMechWargameChassis = "delete-one-mech-wargame-chassis"
)

func mechWargameChassisHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechWargameChassisHandlerConfig")

	chassisConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_chassis.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_chassis.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_chassis.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_chassis.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_chassis.schema.json",
		}),
	}

	chassisConfig[GetManyMechWargameChassis] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/chassis",
		HandlerFunc: getManyMechWargameChassisHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get mech wargame chassis",
		},
	}

	chassisConfig[GetOneMechWargameChassis] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/chassis/:chassis_id",
		HandlerFunc: getOneMechWargameChassisHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get mech wargame chassis",
		},
	}

	chassisConfig[CreateOneMechWargameChassis] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mech-wargame-games/:game_id/chassis",
		HandlerFunc: createOneMechWargameChassisHandler,
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
			Title:    "Create mech wargame chassis",
		},
	}

	chassisConfig[UpdateOneMechWargameChassis] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mech-wargame-games/:game_id/chassis/:chassis_id",
		HandlerFunc: updateOneMechWargameChassisHandler,
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
			Title:    "Update mech wargame chassis",
		},
	}

	chassisConfig[DeleteOneMechWargameChassis] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mech-wargame-games/:game_id/chassis/:chassis_id",
		HandlerFunc: deleteOneMechWargameChassisHandler,
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
			Title:    "Delete mech wargame chassis",
		},
	}

	return chassisConfig, nil
}

func getManyMechWargameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechWargameChassisHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{
		Col: mech_wargame_record.FieldMechWargameChassisGameID,
		Val: gameID,
	})

	mm := m.(*domain.Domain)

	recs, err := mm.GetManyMechWargameChassisRecs(opts)
	if err != nil {
		l.Warn("failed getting mech wargame chassis records >%v<", err)
		return err
	}

	res, err := mapper.MechWargameChassisRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneMechWargameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechWargameChassisHandler")

	gameID := pp.ByName("game_id")
	chassisID := pp.ByName("chassis_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechWargameChassisRec(chassisID, nil)
	if err != nil {
		l.Warn("failed getting mech wargame chassis record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("chassis", chassisID)
	}

	res, err := mapper.MechWargameChassisRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneMechWargameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechWargameChassisHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mech_wargame_record.MechWargameChassis{
		GameID: gameID,
	}

	rec, err := mapper.MechWargameChassisRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechWargameChassisRec(rec)
	if err != nil {
		l.Warn("failed creating mech wargame chassis record >%v<", err)
		return err
	}

	res, err := mapper.MechWargameChassisRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneMechWargameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechWargameChassisHandler")

	gameID := pp.ByName("game_id")
	chassisID := pp.ByName("chassis_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameChassisRec(chassisID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("chassis", chassisID)
	}

	rec, err = mapper.MechWargameChassisRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechWargameChassisRec(rec)
	if err != nil {
		l.Warn("failed updating mech wargame chassis record >%v<", err)
		return err
	}

	res, err := mapper.MechWargameChassisRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneMechWargameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechWargameChassisHandler")

	gameID := pp.ByName("game_id")
	chassisID := pp.ByName("chassis_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameChassisRec(chassisID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("chassis", chassisID)
	}

	if err := mm.DeleteMechWargameChassisRec(chassisID); err != nil {
		l.Warn("failed deleting mech wargame chassis record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
