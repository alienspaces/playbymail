package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-location-links

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/location-links
// GET (document)    /api/v1/adventure-games/{game_id}/location-links/{location_link_id}
// POST (document)   /api/v1/adventure-games/{game_id}/location-links
// PUT (document)    /api/v1/adventure-games/{game_id}/location-links/{location_link_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/location-links/{location_link_id}

const (
	// API Resource Search Path
	searchManyAdventureGameLocationLinks = "search-many-adventure-game-location-links"

	// API Resource CRUD Paths
	getManyAdventureGameLocationLinks  = "get-many-adventure-game-location-links"
	getOneAdventureGameLocationLink    = "get-one-adventure-game-location-link"
	createOneAdventureGameLocationLink = "create-one-adventure-game-location-link"
	updateOneAdventureGameLocationLink = "update-one-adventure-game-location-link"
	deleteOneAdventureGameLocationLink = "delete-one-adventure-game-location-link"
)

func (rnr *Runner) adventureGameLocationLinkHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "adventureGameLocationLinkHandlerConfig")

	l.Debug("Adding adventure_game_location_link handler configuration")

	GameLocationLinkConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_location_link.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_location_link.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "adventure_game_location_link.response.schema.json",
		},
		References: referenceSchemas,
	}

	// New Adventure Game Location Link API paths
	GameLocationLinkConfig[searchManyAdventureGameLocationLinks] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-location-links",
		HandlerFunc: rnr.searchManyAdventureGameLocationLinksHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game location links",
		},
	}

	GameLocationLinkConfig[getManyAdventureGameLocationLinks] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-links",
		HandlerFunc: rnr.getManyAdventureGameLocationLinksHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game location links",
		},
	}

	GameLocationLinkConfig[getOneAdventureGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-links/:location_link_id",
		HandlerFunc: rnr.getOneAdventureGameLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game location link",
		},
	}

	GameLocationLinkConfig[createOneAdventureGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/location-links",
		HandlerFunc: rnr.createOneAdventureGameLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create adventure game location link",
		},
	}

	GameLocationLinkConfig[updateOneAdventureGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/location-links/:location_link_id",
		HandlerFunc: rnr.updateOneAdventureGameLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update adventure game location link",
		},
	}

	GameLocationLinkConfig[deleteOneAdventureGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/location-links/:location_link_id",
		HandlerFunc: rnr.deleteOneAdventureGameLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete adventure game location link",
		},
	}

	return GameLocationLinkConfig, nil
}

func (rnr *Runner) searchManyAdventureGameLocationLinksHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "SearchManyAdventureGameLocationLinksHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter to only return adventure game location links
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_type",
		Val: "adventure",
	})

	recs, err := mm.GetManyAdventureGameLocationLinkRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location link records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationLinkRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getManyAdventureGameLocationLinksHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyAdventureGameLocationLinksHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: "game_id",
		Val: gameID,
	})

	mm := m.(*domain.Domain)

	recs, err := mm.GetManyAdventureGameLocationLinkRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location link records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationLinkRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getOneAdventureGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetOneAdventureGameLocationLinkHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	locationLinkID := pp.ByName("location_link_id")
	if locationLinkID == "" {
		l.Warn("location link id is required")
		return coreerror.NewNotFoundError("location link", locationLinkID)
	}

	mm := m.(*domain.Domain)

	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed getting game record >%v<", err)
		return err
	}

	rec, err := mm.GetAdventureGameLocationLinkRec(locationLinkID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location link record >%v<", err)
		return err
	}

	// Verify the location link belongs to the specified game
	if rec.GameID != gameRec.ID {
		l.Warn("location link does not belong to specified game >%s< != >%s<", rec.GameID, gameRec.ID)
		return coreerror.NewNotFoundError("location link", locationLinkID)
	}

	res, err := mapper.AdventureGameLocationLinkRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game location link record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) createOneAdventureGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateOneAdventureGameLocationLinkHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	rec, err := mapper.AdventureGameLocationLinkRequestToRecord(l, r, &record.AdventureGameLocationLink{})
	if err != nil {
		return err
	}

	// Set the game ID from the path parameter
	rec.GameID = gameID

	mm := m.(*domain.Domain)

	rec, err = mm.CreateAdventureGameLocationLinkRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location link record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationLinkRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) updateOneAdventureGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "updateOneAdventureGameLocationLinkHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	locationLinkID := pp.ByName("location_link_id")
	if locationLinkID == "" {
		l.Warn("location link id is required")
		return coreerror.NewNotFoundError("location link", locationLinkID)
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationLinkRec(locationLinkID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location link record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		l.Warn("location link does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location link", locationLinkID)
	}

	rec, err = mapper.AdventureGameLocationLinkRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameLocationLinkRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game location link record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationLinkRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) deleteOneAdventureGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteOneAdventureGameLocationLinkHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	locationLinkID := pp.ByName("location_link_id")
	if locationLinkID == "" {
		l.Warn("location link id is required")
		return coreerror.NewNotFoundError("location link", locationLinkID)
	}

	l.Info("deleting adventure game location link record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed getting game record >%v<", err)
		return err
	}

	// First get the record to verify it belongs to the specified game
	rec, err := mm.GetAdventureGameLocationLinkRec(locationLinkID, nil)
	if err != nil {
		return err
	}

	// Verify the location link belongs to the specified game
	if rec.GameID != gameRec.ID {
		l.Warn("location link does not belong to specified game >%s< != >%s<", rec.GameID, gameRec.ID)
		return coreerror.NewNotFoundError("location link", locationLinkID)
	}

	if err := mm.DeleteAdventureGameLocationLinkRec(locationLinkID); err != nil {
		l.Warn("failed deleting adventure game location link record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
