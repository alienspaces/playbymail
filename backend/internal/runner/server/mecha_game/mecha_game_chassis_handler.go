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

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/mecha-games/{game_id}/chassis
// GET (document)    /api/v1/mecha-games/{game_id}/chassis/{chassis_id}
// POST (document)   /api/v1/mecha-games/{game_id}/chassis
// PUT (document)    /api/v1/mecha-games/{game_id}/chassis/{chassis_id}
// DELETE (document) /api/v1/mecha-games/{game_id}/chassis/{chassis_id}

const (
	GetManyMechaGameChassis   = "get-many-mecha-chassis"
	GetOneMechaGameChassis    = "get-one-mecha-chassis"
	CreateOneMechaGameChassis = "create-one-mecha-chassis"
	UpdateOneMechaGameChassis = "update-one-mecha-chassis"
	DeleteOneMechaGameChassis = "delete-one-mecha-chassis"
)

func mechaGameChassisHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaGameChassisHandlerConfig")

	l.Debug("Adding mecha chassis handler configuration")

	chassisConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_chassis.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_chassis.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_chassis.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_chassis.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_chassis.schema.json",
		}),
	}

	chassisConfig[GetManyMechaGameChassis] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/chassis",
		HandlerFunc: getManyMechaGameChassisHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get mecha chassis",
		},
	}

	chassisConfig[GetOneMechaGameChassis] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/chassis/:chassis_id",
		HandlerFunc: getOneMechaGameChassisHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get mecha chassis",
		},
	}

	chassisConfig[CreateOneMechaGameChassis] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/chassis",
		HandlerFunc: createOneMechaGameChassisHandler,
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
			Title:    "Create mecha chassis",
		},
	}

	chassisConfig[UpdateOneMechaGameChassis] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/chassis/:chassis_id",
		HandlerFunc: updateOneMechaGameChassisHandler,
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
			Title:    "Update mecha chassis",
		},
	}

	chassisConfig[DeleteOneMechaGameChassis] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/chassis/:chassis_id",
		HandlerFunc: deleteOneMechaGameChassisHandler,
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
			Title:    "Delete mecha chassis",
		},
	}

	return chassisConfig, nil
}

func getManyMechaGameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaGameChassisHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{
		Col: mecha_game_record.FieldMechaGameChassisGameID,
		Val: gameID,
	})

	mm := m.(*domain.Domain)

	recs, err := mm.GetManyMechaGameChassisRecs(opts)
	if err != nil {
		l.Warn("failed getting mecha chassis records >%v<", err)
		return err
	}

	res, err := mapper.MechaGameChassisRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneMechaGameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaGameChassisHandler")

	gameID := pp.ByName("game_id")
	chassisID := pp.ByName("chassis_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaGameChassisRec(chassisID, nil)
	if err != nil {
		l.Warn("failed getting mecha chassis record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("chassis", chassisID)
	}

	res, err := mapper.MechaGameChassisRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneMechaGameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaGameChassisHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_game_record.MechaGameChassis{
		GameID: gameID,
	}

	rec, err := mapper.MechaGameChassisRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaGameChassisRec(rec)
	if err != nil {
		l.Warn("failed creating mecha chassis record >%v<", err)
		return err
	}

	res, err := mapper.MechaGameChassisRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneMechaGameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaGameChassisHandler")

	gameID := pp.ByName("game_id")
	chassisID := pp.ByName("chassis_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameChassisRec(chassisID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("chassis", chassisID)
	}

	rec, err = mapper.MechaGameChassisRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaGameChassisRec(rec)
	if err != nil {
		l.Warn("failed updating mecha chassis record >%v<", err)
		return err
	}

	res, err := mapper.MechaGameChassisRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneMechaGameChassisHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaGameChassisHandler")

	gameID := pp.ByName("game_id")
	chassisID := pp.ByName("chassis_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameChassisRec(chassisID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("chassis", chassisID)
	}

	if err := mm.DeleteMechaGameChassisRec(chassisID); err != nil {
		l.Warn("failed deleting mecha chassis record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
