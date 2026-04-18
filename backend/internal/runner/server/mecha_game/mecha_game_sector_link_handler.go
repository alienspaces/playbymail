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
	GetManyMechaGameSectorLinks  = "get-many-mecha-sector-links"
	GetOneMechaGameSectorLink    = "get-one-mecha-sector-link"
	CreateOneMechaGameSectorLink = "create-one-mecha-sector-link"
	UpdateOneMechaGameSectorLink = "update-one-mecha-sector-link"
	DeleteOneMechaGameSectorLink = "delete-one-mecha-sector-link"
)

func mechaGameSectorLinkHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "mechaGameSectorLinkHandlerConfig")

	l.Debug("Adding mecha sector link handler configuration")

	sectorLinkConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector_link.collection.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector_link.schema.json",
		}),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector_link.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector_link.response.schema.json",
		},
		References: append(referenceSchemas, jsonschema.Schema{
			Location: "api/mecha_game_schema",
			Name:     "mecha_game_sector_link.schema.json",
		}),
	}

	sectorLinkConfig[GetManyMechaGameSectorLinks] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/sector-links",
		HandlerFunc: getManyMechaGameSectorLinksHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Collection: true, Title: "Get mecha sector links"},
	}

	sectorLinkConfig[GetOneMechaGameSectorLink] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/mecha-games/:game_id/sector-links/:sector_link_id",
		HandlerFunc: getOneMechaGameSectorLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Get mecha sector link"},
	}

	sectorLinkConfig[CreateOneMechaGameSectorLink] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/mecha-games/:game_id/sector-links",
		HandlerFunc: createOneMechaGameSectorLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Create mecha sector link"},
	}

	sectorLinkConfig[UpdateOneMechaGameSectorLink] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/mecha-games/:game_id/sector-links/:sector_link_id",
		HandlerFunc: updateOneMechaGameSectorLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions:       []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Update mecha sector link"},
	}

	sectorLinkConfig[DeleteOneMechaGameSectorLink] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/mecha-games/:game_id/sector-links/:sector_link_id",
		HandlerFunc: deleteOneMechaGameSectorLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:      []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{handler_auth.PermissionGameDesign},
		},
		DocumentationConfig: server.DocumentationConfig{Document: true, Title: "Delete mecha sector link"},
	}

	return sectorLinkConfig, nil
}

func getManyMechaGameSectorLinksHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyMechaGameSectorLinksHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		return coreerror.NewParamError("game_id is required")
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, sql.Param{Col: mecha_game_record.FieldMechaGameSectorLinkGameID, Val: gameID})

	mm := m.(*domain.Domain)
	recs, err := mm.GetManyMechaGameSectorLinkRecs(opts)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameSectorLinkRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getOneMechaGameSectorLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneMechaGameSectorLinkHandler")

	gameID := pp.ByName("game_id")
	sectorLinkID := pp.ByName("sector_link_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetMechaGameSectorLinkRec(sectorLinkID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("sector_link", sectorLinkID)
	}

	res, err := mapper.MechaGameSectorLinkRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createOneMechaGameSectorLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneMechaGameSectorLinkHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec := &mecha_game_record.MechaGameSectorLink{GameID: gameID}
	rec, err := mapper.MechaGameSectorLinkRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateMechaGameSectorLinkRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameSectorLinkRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateOneMechaGameSectorLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneMechaGameSectorLinkHandler")

	gameID := pp.ByName("game_id")
	sectorLinkID := pp.ByName("sector_link_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameSectorLinkRec(sectorLinkID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("sector_link", sectorLinkID)
	}

	rec, err = mapper.MechaGameSectorLinkRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateMechaGameSectorLinkRec(rec)
	if err != nil {
		return err
	}

	res, err := mapper.MechaGameSectorLinkRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteOneMechaGameSectorLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneMechaGameSectorLinkHandler")

	gameID := pp.ByName("game_id")
	sectorLinkID := pp.ByName("sector_link_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetMechaGameSectorLinkRec(sectorLinkID, nil)
	if err != nil {
		return err
	}
	if rec.GameID != gameID {
		return coreerror.NewNotFoundError("sector_link", sectorLinkID)
	}

	if err := mm.DeleteMechaGameSectorLinkRec(sectorLinkID); err != nil {
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
