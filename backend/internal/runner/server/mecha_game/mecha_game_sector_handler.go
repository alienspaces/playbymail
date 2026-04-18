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
	GetManyMechaGameSectors  = "get-many-mecha-sectors"
	GetOneMechaGameSector    = "get-one-mecha-sector"
	CreateOneMechaGameSector = "create-one-mecha-sector"
	UpdateOneMechaGameSector = "update-one-mecha-sector"
	DeleteOneMechaGameSector = "delete-one-mecha-sector"
)

func mechaGameSectorHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaGameSectorHandlerConfig")

	l.Debug("Adding mecha sector handler configuration")

	sectorConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector.schema.json",
		}),
	}

	sectorConfig[GetManyMechaGameSectors] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/sectors",
		HandlerFunc: getManyMechaGameSectorsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mecha sectors"},
	}

	sectorConfig[GetOneMechaGameSector] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/sectors/:sector_id",
		HandlerFunc: getOneMechaGameSectorHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mecha sector"},
	}

	sectorConfig[CreateOneMechaGameSector] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/sectors",
		HandlerFunc: createOneMechaGameSectorHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mecha sector"},
	}

	sectorConfig[UpdateOneMechaGameSector] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/sectors/:sector_id",
		HandlerFunc: updateOneMechaGameSectorHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mecha sector"},
	}

	sectorConfig[DeleteOneMechaGameSector] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/sectors/:sector_id",
		HandlerFunc: deleteOneMechaGameSectorHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mecha sector"},
	}

	return sectorConfig, nil
}

func getManyMechaGameSectorsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaGameSectorsHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{Col: mecha_game_record.FieldMechaGameSectorGameID, Val: gameID})

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechaGameSectorRecs(opts)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameSectorRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechaGameSectorHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaGameSectorHandler")

	gameID := pp.ByName("game_id")
	sectorID := pp.ByName("sector_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaGameSectorRec(sectorID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("sector", sectorID)
	}

	res, err := mapper.MechaGameSectorRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechaGameSectorHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaGameSectorHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_game_record.MechaGameSector{GameID: gameID}
	rec, err := mapper.MechaGameSectorRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaGameSectorRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameSectorRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechaGameSectorHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaGameSectorHandler")

	gameID := pp.ByName("game_id")
	sectorID := pp.ByName("sector_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameSectorRec(sectorID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("sector", sectorID)
	}

	rec, err = mapper.MechaGameSectorRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaGameSectorRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameSectorRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechaGameSectorHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaGameSectorHandler")

	gameID := pp.ByName("game_id")
	sectorID := pp.ByName("sector_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameSectorRec(sectorID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("sector", sectorID)
	}

	if err := mm.DeleteMechaGameSectorRec(sectorID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
