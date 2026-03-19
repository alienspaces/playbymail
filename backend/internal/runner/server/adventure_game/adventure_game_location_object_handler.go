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
)

// API Resource Search Path
//
// GET (collection) /api/v1/adventure-game-location-objects

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/location-objects
// GET (document)    /api/v1/adventure-games/{game_id}/location-objects/{location_object_id}
// POST (document)   /api/v1/adventure-games/{game_id}/location-objects
// PUT (document)    /api/v1/adventure-games/{game_id}/location-objects/{location_object_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/location-objects/{location_object_id}

const (
	SearchManyAdventureGameLocationObjects = "searchManyAdventureGameLocationObjects"
	GetManyAdventureGameLocationObjects    = "getManyAdventureGameLocationObjects"
	GetOneAdventureGameLocationObject      = "getOneAdventureGameLocationObject"
	CreateOneAdventureGameLocationObject   = "createOneAdventureGameLocationObject"
	UpdateOneAdventureGameLocationObject   = "updateOneAdventureGameLocationObject"
	DeleteOneAdventureGameLocationObject   = "deleteOneAdventureGameLocationObject"
)

func adventureGameLocationObjectHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameLocationObjectHandlerConfig")

	l.Debug("Adding adventure_game_location_object handler configuration")

	locationObjectConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_object.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_object.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_object.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_object.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_object.schema.json",
			},
		}...),
	}

	locationObjectConfig[SearchManyAdventureGameLocationObjects] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-location-objects",
		HandlerFunc: searchManyAdventureGameLocationObjectsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game location objects",
		},
	}

	locationObjectConfig[GetManyAdventureGameLocationObjects] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-objects",
		HandlerFunc: getManyAdventureGameLocationObjectsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game location objects",
		},
	}

	locationObjectConfig[GetOneAdventureGameLocationObject] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-objects/:location_object_id",
		HandlerFunc: getOneAdventureGameLocationObjectHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game location object",
		},
	}

	locationObjectConfig[CreateOneAdventureGameLocationObject] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/location-objects",
		HandlerFunc: createOneAdventureGameLocationObjectHandler,
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
			Title:    "Create adventure game location object",
		},
	}

	locationObjectConfig[UpdateOneAdventureGameLocationObject] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/location-objects/:location_object_id",
		HandlerFunc: updateOneAdventureGameLocationObjectHandler,
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
			Title:    "Update adventure game location object",
		},
	}

	locationObjectConfig[DeleteOneAdventureGameLocationObject] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/location-objects/:location_object_id",
		HandlerFunc: deleteOneAdventureGameLocationObjectHandler,
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
			Title:    "Delete adventure game location object",
		},
	}

	return locationObjectConfig, nil
}

func searchManyAdventureGameLocationObjectsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameLocationObjectsHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameLocationObjectRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location object records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameLocationObjectsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameLocationObjectsHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	opts.Params = append(opts.Params, sql.Param{
		Col: adventure_game_record.FieldAdventureGameLocationObjectGameID,
		Val: gameID,
	})

	recs, err := mm.GetManyAdventureGameLocationObjectRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location object records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameLocationObjectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameLocationObjectHandler")

	gameID := pp.ByName("game_id")
	locationObjectID := pp.ByName("location_object_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationObjectRec(locationObjectID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location object record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		l.Warn("location object does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location object", locationObjectID)
	}

	res, err := mapper.AdventureGameLocationObjectRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game location object record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameLocationObjectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameLocationObjectHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed getting game record >%v<", err)
		return err
	}

	rec := &adventure_game_record.AdventureGameLocationObject{
		GameID: gameRec.ID,
	}

	rec, err = mapper.AdventureGameLocationObjectRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateAdventureGameLocationObjectRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location object record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameLocationObjectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameLocationObjectHandler")

	gameID := pp.ByName("game_id")
	locationObjectID := pp.ByName("location_object_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationObjectRec(locationObjectID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("location object does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location object", locationObjectID)
	}

	rec, err = mapper.AdventureGameLocationObjectRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameLocationObjectRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game location object record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameLocationObjectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameLocationObjectHandler")

	gameID := pp.ByName("game_id")
	locationObjectID := pp.ByName("location_object_id")

	l.Info("deleting adventure game location object record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationObjectRec(locationObjectID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("location object does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location object", locationObjectID)
	}

	if err := mm.DeleteAdventureGameLocationObjectRec(locationObjectID); err != nil {
		l.Warn("failed deleting adventure game location object record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
