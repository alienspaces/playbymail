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
	GetManyMechWargameSectors  = "get-many-mech-wargame-sectors"
	GetOneMechWargameSector    = "get-one-mech-wargame-sector"
	CreateOneMechWargameSector = "create-one-mech-wargame-sector"
	UpdateOneMechWargameSector = "update-one-mech-wargame-sector"
	DeleteOneMechWargameSector = "delete-one-mech-wargame-sector"
)

func mechWargameSectorHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechWargameSectorHandlerConfig")

	sectorConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_sector.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_sector.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_sector.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_sector.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mech_wargame_schema",
			Name:     "mech_wargame_sector.schema.json",
		}),
	}

	sectorConfig[GetManyMechWargameSectors] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/sectors",
		HandlerFunc: getManyMechWargameSectorsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mech wargame sectors"},
	}

	sectorConfig[GetOneMechWargameSector] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mech-wargame-games/:game_id/sectors/:sector_id",
		HandlerFunc: getOneMechWargameSectorHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mech wargame sector"},
	}

	sectorConfig[CreateOneMechWargameSector] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mech-wargame-games/:game_id/sectors",
		HandlerFunc: createOneMechWargameSectorHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mech wargame sector"},
	}

	sectorConfig[UpdateOneMechWargameSector] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mech-wargame-games/:game_id/sectors/:sector_id",
		HandlerFunc: updateOneMechWargameSectorHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mech wargame sector"},
	}

	sectorConfig[DeleteOneMechWargameSector] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mech-wargame-games/:game_id/sectors/:sector_id",
		HandlerFunc: deleteOneMechWargameSectorHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mech wargame sector"},
	}

	return sectorConfig, nil
}

func getManyMechWargameSectorsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechWargameSectorsHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{Col: mech_wargame_record.FieldMechWargameSectorGameID, Val: gameID})

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechWargameSectorRecs(opts)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameSectorRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechWargameSectorHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechWargameSectorHandler")

	gameID := pp.ByName("game_id")
	sectorID := pp.ByName("sector_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechWargameSectorRec(sectorID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("sector", sectorID)
	}

	res, err := mapper.MechWargameSectorRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechWargameSectorHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechWargameSectorHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mech_wargame_record.MechWargameSector{GameID: gameID}
	rec, err := mapper.MechWargameSectorRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechWargameSectorRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameSectorRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechWargameSectorHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechWargameSectorHandler")

	gameID := pp.ByName("game_id")
	sectorID := pp.ByName("sector_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameSectorRec(sectorID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("sector", sectorID)
	}

	rec, err = mapper.MechWargameSectorRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechWargameSectorRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechWargameSectorRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechWargameSectorHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechWargameSectorHandler")

	gameID := pp.ByName("game_id")
	sectorID := pp.ByName("sector_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechWargameSectorRec(sectorID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("sector", sectorID)
	}

	if err := mm.DeleteMechWargameSectorRec(sectorID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
