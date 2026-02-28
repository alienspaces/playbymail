package adventure_game

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
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-location-link-requirements

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/location-link-requirements
// GET (document)    /api/v1/adventure-games/{game_id}/location-link-requirements/{location_link_requirement_id}
// POST (document)   /api/v1/adventure-games/{game_id}/location-link-requirements
// PUT (document)    /api/v1/adventure-games/{game_id}/location-link-requirements/{location_link_requirement_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/location-link-requirements/{location_link_requirement_id}

const (
	// API Resource Search Path
	searchManyAdventureGameLocationLinkRequirements = "search-many-adventure-game-location-link-requirements"

	// API Resource CRUD Paths
	getManyAdventureGameLocationLinkRequirements  = "get-many-adventure-game-location-link-requirements"
	getOneAdventureGameLocationLinkRequirement    = "get-one-adventure-game-location-link-requirement"
	createOneAdventureGameLocationLinkRequirement = "create-one-adventure-game-location-link-requirement"
	updateOneAdventureGameLocationLinkRequirement = "update-one-adventure-game-location-link-requirement"
	deleteOneAdventureGameLocationLinkRequirement = "delete-one-adventure-game-location-link-requirement"
)

func adventureGameLocationLinkRequirementHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameLocationLinkRequirementHandlerConfig")

	l.Debug("Adding adventure_game_location_link_requirement handler configuration")

	gameLocationLinkRequirementConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_link_requirement.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_link_requirement.collection.response.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_link_requirement.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_link_requirement.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_link_requirement.schema.json",
			},
		}...),
	}

	// New Adventure Game Location Link Requirement API paths
	gameLocationLinkRequirementConfig[searchManyAdventureGameLocationLinkRequirements] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-location-link-requirements",
		HandlerFunc: searchManyAdventureGameLocationLinkRequirementsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game location link requirements",
		},
	}

	gameLocationLinkRequirementConfig[getManyAdventureGameLocationLinkRequirements] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-link-requirements",
		HandlerFunc: getManyAdventureGameLocationLinkRequirementsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game location link requirements",
		},
	}

	gameLocationLinkRequirementConfig[getOneAdventureGameLocationLinkRequirement] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-link-requirements/:location_link_requirement_id",
		HandlerFunc: getOneAdventureGameLocationLinkRequirementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game location link requirement",
		},
	}

	gameLocationLinkRequirementConfig[createOneAdventureGameLocationLinkRequirement] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/location-link-requirements",
		HandlerFunc: createOneAdventureGameLocationLinkRequirementHandler,
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
			Title:    "Create adventure game location link requirement",
		},
	}

	gameLocationLinkRequirementConfig[updateOneAdventureGameLocationLinkRequirement] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/location-link-requirements/:location_link_requirement_id",
		HandlerFunc: updateOneAdventureGameLocationLinkRequirementHandler,
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
			Title:    "Update adventure game location link requirement",
		},
	}

	gameLocationLinkRequirementConfig[deleteOneAdventureGameLocationLinkRequirement] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/location-link-requirements/:location_link_requirement_id",
		HandlerFunc: deleteOneAdventureGameLocationLinkRequirementHandler,
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
			Title:    "Delete adventure game location link requirement",
		},
	}

	return gameLocationLinkRequirementConfig, nil
}

// New Adventure Game Location Link Requirement Handlers

func searchManyAdventureGameLocationLinkRequirementsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameLocationLinkRequirementsHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameLocationLinkRequirementRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location link requirement records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationLinkRequirementRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameLocationLinkRequirementsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameLocationLinkRequirementsHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_id",
		Val: gameID,
	})

	recs, err := mm.GetManyAdventureGameLocationLinkRequirementRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location link requirement records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationLinkRequirementRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameLocationLinkRequirementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameLocationLinkRequirementHandler")

	gameID := pp.ByName("game_id")
	locationLinkRequirementID := pp.ByName("location_link_requirement_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationLinkRequirementRec(locationLinkRequirementID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location link requirement record >%v<", err)
		return err
	}

	// Verify the location link requirement belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("location link requirement does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location link requirement", locationLinkRequirementID)
	}

	res, err := mapper.AdventureGameLocationLinkRequirementRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game location link requirement record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameLocationLinkRequirementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameLocationLinkRequirementHandler")

	gameID := pp.ByName("game_id")

	rec := &adventure_game_record.AdventureGameLocationLinkRequirement{}
	rec, err := mapper.AdventureGameLocationLinkRequirementRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	// Set the game ID from the path parameter
	rec.GameID = gameID

	mm := m.(*domain.Domain)

	rec, err = mm.CreateAdventureGameLocationLinkRequirementRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location link requirement record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationLinkRequirementRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameLocationLinkRequirementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameLocationLinkRequirementHandler")

	gameID := pp.ByName("game_id")
	locationLinkRequirementID := pp.ByName("location_link_requirement_id")

	l.Info("updating adventure game location link requirement record with path params >%#v<", pp)

	var req adventure_game_schema.AdventureGameLocationLinkRequirementRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationLinkRequirementRec(locationLinkRequirementID, nil)
	if err != nil {
		return err
	}

	// Verify the location link requirement belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("location link requirement does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location link requirement", locationLinkRequirementID)
	}

	rec, err = mapper.AdventureGameLocationLinkRequirementRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameLocationLinkRequirementRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game location link requirement record >%v<", err)
		return err
	}

	data, err := mapper.AdventureGameLocationLinkRequirementRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := adventure_game_schema.AdventureGameLocationLinkRequirementResponse{
		Data: data,
	}

	l.Info("responding with updated adventure game location link requirement record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameLocationLinkRequirementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameLocationLinkRequirementHandler")

	gameID := pp.ByName("game_id")
	locationLinkRequirementID := pp.ByName("location_link_requirement_id")

	l.Info("deleting adventure game location link requirement record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	// First get the record to verify it belongs to the specified game
	rec, err := mm.GetAdventureGameLocationLinkRequirementRec(locationLinkRequirementID, nil)
	if err != nil {
		return err
	}

	// Verify the location link requirement belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("location link requirement does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location link requirement", locationLinkRequirementID)
	}

	if err := mm.DeleteAdventureGameLocationLinkRequirementRec(locationLinkRequirementID); err != nil {
		l.Warn("failed deleting adventure game location link requirement record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
