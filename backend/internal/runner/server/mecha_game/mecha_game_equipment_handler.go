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
	GetManyMechaGameEquipments  = "get-many-mecha-equipments"
	GetOneMechaGameEquipment    = "get-one-mecha-equipment"
	CreateOneMechaGameEquipment = "create-one-mecha-equipment"
	UpdateOneMechaGameEquipment = "update-one-mecha-equipment"
	DeleteOneMechaGameEquipment = "delete-one-mecha-equipment"
)

func mechaGameEquipmentHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaGameEquipmentHandlerConfig")

	l.Debug("Adding mecha equipment handler configuration")

	equipmentConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_equipment.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_equipment.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_equipment.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_equipment.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_equipment.schema.json",
		}),
	}

	equipmentConfig[GetManyMechaGameEquipments] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/equipments",
		HandlerFunc: getManyMechaGameEquipmentsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mecha equipments"},
	}

	equipmentConfig[GetOneMechaGameEquipment] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/equipments/:equipment_id",
		HandlerFunc: getOneMechaGameEquipmentHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mecha equipment"},
	}

	equipmentConfig[CreateOneMechaGameEquipment] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/equipments",
		HandlerFunc: createOneMechaGameEquipmentHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mecha equipment"},
	}

	equipmentConfig[UpdateOneMechaGameEquipment] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/equipments/:equipment_id",
		HandlerFunc: updateOneMechaGameEquipmentHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mecha equipment"},
	}

	equipmentConfig[DeleteOneMechaGameEquipment] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/equipments/:equipment_id",
		HandlerFunc: deleteOneMechaGameEquipmentHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mecha equipment"},
	}

	return equipmentConfig, nil
}

func getManyMechaGameEquipmentsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaGameEquipmentsHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{Col: mecha_game_record.FieldMechaGameEquipmentGameID, Val: gameID})

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechaGameEquipmentRecs(opts)
	if err != nil {
		l.Warn("failed getting mecha equipment records >%v<", err)
		return err
	}

	res, err := mapper.MechaGameEquipmentRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechaGameEquipmentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaGameEquipmentHandler")

	gameID := pp.ByName("game_id")
	equipmentID := pp.ByName("equipment_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaGameEquipmentRec(equipmentID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("equipment", equipmentID)
	}

	res, err := mapper.MechaGameEquipmentRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechaGameEquipmentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaGameEquipmentHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_game_record.MechaGameEquipment{GameID: gameID}
	rec, err := mapper.MechaGameEquipmentRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaGameEquipmentRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameEquipmentRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechaGameEquipmentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaGameEquipmentHandler")

	gameID := pp.ByName("game_id")
	equipmentID := pp.ByName("equipment_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameEquipmentRec(equipmentID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("equipment", equipmentID)
	}

	rec, err = mapper.MechaGameEquipmentRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaGameEquipmentRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameEquipmentRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechaGameEquipmentHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaGameEquipmentHandler")

	gameID := pp.ByName("game_id")
	equipmentID := pp.ByName("equipment_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameEquipmentRec(equipmentID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("equipment", equipmentID)
	}

	if err := mm.DeleteMechaGameEquipmentRec(equipmentID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
