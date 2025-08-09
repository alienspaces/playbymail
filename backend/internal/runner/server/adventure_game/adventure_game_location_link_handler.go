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
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
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
	SearchManyAdventureGameLocationLinks = "search-many-adventure-game-location-links"

	// API Resource CRUD Paths
	GetManyAdventureGameLocationLinks  = "get-many-adventure-game-location-links"
	GetOneAdventureGameLocationLink    = "get-one-adventure-game-location-link"
	CreateOneAdventureGameLocationLink = "create-one-adventure-game-location-link"
	UpdateOneAdventureGameLocationLink = "update-one-adventure-game-location-link"
	DeleteOneAdventureGameLocationLink = "delete-one-adventure-game-location-link"
)

func adventureGameLocationLinkHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameLocationLinkHandlerConfig")

	l.Debug("Adding adventure_game_location_link handler configuration")

	gameLocationLinkConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_link.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_link.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_link.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_link.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_link.schema.json",
			},
		}...),
	}

	// New Adventure Game Location Link API paths
	gameLocationLinkConfig[SearchManyAdventureGameLocationLinks] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-location-links",
		HandlerFunc: searchManyAdventureGameLocationLinksHandler,
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

	gameLocationLinkConfig[GetManyAdventureGameLocationLinks] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-links",
		HandlerFunc: getManyAdventureGameLocationLinksHandler,
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

	gameLocationLinkConfig[GetOneAdventureGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-links/:location_link_id",
		HandlerFunc: getOneAdventureGameLocationLinkHandler,
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

	gameLocationLinkConfig[CreateOneAdventureGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/location-links",
		HandlerFunc: createOneAdventureGameLocationLinkHandler,
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

	gameLocationLinkConfig[UpdateOneAdventureGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/location-links/:location_link_id",
		HandlerFunc: updateOneAdventureGameLocationLinkHandler,
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

	gameLocationLinkConfig[DeleteOneAdventureGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/location-links/:location_link_id",
		HandlerFunc: deleteOneAdventureGameLocationLinkHandler,
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

	return gameLocationLinkConfig, nil
}

func searchManyAdventureGameLocationLinksHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameLocationLinksHandler")

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

func getManyAdventureGameLocationLinksHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameLocationLinksHandler")

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

func getOneAdventureGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameLocationLinkHandler")

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

func createOneAdventureGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameLocationLinkHandler")

	gameID := pp.ByName("game_id")
	if gameID == "" {
		l.Warn("game id is required")
		return coreerror.NewNotFoundError("game", gameID)
	}

	rec := &adventure_game_record.AdventureGameLocationLink{}
	rec, err := mapper.AdventureGameLocationLinkRequestToRecord(l, r, rec)
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

func updateOneAdventureGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameLocationLinkHandler")

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

func deleteOneAdventureGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameLocationLinkHandler")

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
