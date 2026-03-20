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
// GET (collection) /api/v1/adventure-game-location-object-effects

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/location-object-effects
// GET (document)    /api/v1/adventure-games/{game_id}/location-object-effects/{location_object_effect_id}
// POST (document)   /api/v1/adventure-games/{game_id}/location-object-effects
// PUT (document)    /api/v1/adventure-games/{game_id}/location-object-effects/{location_object_effect_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/location-object-effects/{location_object_effect_id}

const (
	SearchManyAdventureGameLocationObjectEffects = "searchManyAdventureGameLocationObjectEffects"
	GetManyAdventureGameLocationObjectEffects    = "getManyAdventureGameLocationObjectEffects"
	GetOneAdventureGameLocationObjectEffect      = "getOneAdventureGameLocationObjectEffect"
	CreateOneAdventureGameLocationObjectEffect   = "createOneAdventureGameLocationObjectEffect"
	UpdateOneAdventureGameLocationObjectEffect   = "updateOneAdventureGameLocationObjectEffect"
	DeleteOneAdventureGameLocationObjectEffect   = "deleteOneAdventureGameLocationObjectEffect"
)

func adventureGameLocationObjectEffectHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameLocationObjectEffectHandlerConfig")

	l.Debug("Adding adventure_game_location_object_effect handler configuration")

	locationObjectEffectConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_object_effect.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_object_effect.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_object_effect.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_object_effect.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_object_effect.schema.json",
			},
		}...),
	}

	locationObjectEffectConfig[SearchManyAdventureGameLocationObjectEffects] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-location-object-effects",
		HandlerFunc: searchManyAdventureGameLocationObjectEffectsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game location object effects",
		},
	}

	locationObjectEffectConfig[GetManyAdventureGameLocationObjectEffects] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-object-effects",
		HandlerFunc: getManyAdventureGameLocationObjectEffectsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game location object effects",
		},
	}

	locationObjectEffectConfig[GetOneAdventureGameLocationObjectEffect] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-object-effects/:location_object_effect_id",
		HandlerFunc: getOneAdventureGameLocationObjectEffectHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game location object effect",
		},
	}

	locationObjectEffectConfig[CreateOneAdventureGameLocationObjectEffect] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/location-object-effects",
		HandlerFunc: createOneAdventureGameLocationObjectEffectHandler,
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
			Title:    "Create adventure game location object effect",
		},
	}

	locationObjectEffectConfig[UpdateOneAdventureGameLocationObjectEffect] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/location-object-effects/:location_object_effect_id",
		HandlerFunc: updateOneAdventureGameLocationObjectEffectHandler,
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
			Title:    "Update adventure game location object effect",
		},
	}

	locationObjectEffectConfig[DeleteOneAdventureGameLocationObjectEffect] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/location-object-effects/:location_object_effect_id",
		HandlerFunc: deleteOneAdventureGameLocationObjectEffectHandler,
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
			Title:    "Delete adventure game location object effect",
		},
	}

	return locationObjectEffectConfig, nil
}

func searchManyAdventureGameLocationObjectEffectsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameLocationObjectEffectsHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameLocationObjectEffectRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location object effect records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectEffectRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameLocationObjectEffectsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameLocationObjectEffectsHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	opts.Params = append(opts.Params, sql.Param{
		Col: adventure_game_record.FieldAdventureGameLocationObjectEffectGameID,
		Val: gameID,
	})

	recs, err := mm.GetManyAdventureGameLocationObjectEffectRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location object effect records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectEffectRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameLocationObjectEffectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameLocationObjectEffectHandler")

	gameID := pp.ByName("game_id")
	locationObjectEffectID := pp.ByName("location_object_effect_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationObjectEffectRec(locationObjectEffectID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location object effect record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		l.Warn("location object effect does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location object effect", locationObjectEffectID)
	}

	res, err := mapper.AdventureGameLocationObjectEffectRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game location object effect record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameLocationObjectEffectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameLocationObjectEffectHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeAdventureGameDesigner(l, r, mm, gameID); err != nil {
		return err
	}

	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed getting game record >%v<", err)
		return err
	}

	rec := &adventure_game_record.AdventureGameLocationObjectEffect{
		GameID: gameRec.ID,
	}

	rec, err = mapper.AdventureGameLocationObjectEffectRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateAdventureGameLocationObjectEffectRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location object effect record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectEffectRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameLocationObjectEffectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameLocationObjectEffectHandler")

	gameID := pp.ByName("game_id")
	locationObjectEffectID := pp.ByName("location_object_effect_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeAdventureGameDesigner(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameLocationObjectEffectRec(locationObjectEffectID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("location object effect does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location object effect", locationObjectEffectID)
	}

	rec, err = mapper.AdventureGameLocationObjectEffectRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameLocationObjectEffectRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game location object effect record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectEffectRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameLocationObjectEffectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameLocationObjectEffectHandler")

	gameID := pp.ByName("game_id")
	locationObjectEffectID := pp.ByName("location_object_effect_id")

	l.Info("deleting adventure game location object effect record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	if _, err := authorizeAdventureGameDesigner(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameLocationObjectEffectRec(locationObjectEffectID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("location object effect does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("location object effect", locationObjectEffectID)
	}

	if err := mm.DeleteAdventureGameLocationObjectEffectRec(locationObjectEffectID); err != nil {
		l.Warn("failed deleting adventure game location object effect record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
