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
// GET (collection) /api/v1/adventure-game-item-effects

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/item-effects
// GET (document)    /api/v1/adventure-games/{game_id}/item-effects/{item_effect_id}
// POST (document)   /api/v1/adventure-games/{game_id}/item-effects
// PUT (document)    /api/v1/adventure-games/{game_id}/item-effects/{item_effect_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/item-effects/{item_effect_id}

const (
	SearchManyAdventureGameItemEffects = "searchManyAdventureGameItemEffects"
	GetManyAdventureGameItemEffects    = "getManyAdventureGameItemEffects"
	GetOneAdventureGameItemEffect      = "getOneAdventureGameItemEffect"
	CreateOneAdventureGameItemEffect   = "createOneAdventureGameItemEffect"
	UpdateOneAdventureGameItemEffect   = "updateOneAdventureGameItemEffect"
	DeleteOneAdventureGameItemEffect   = "deleteOneAdventureGameItemEffect"
)

func adventureGameItemEffectHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameItemEffectHandlerConfig")

	l.Debug("Adding adventure_game_item_effect handler configuration")

	itemEffectConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_item_effect.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_item_effect.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_item_effect.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_item_effect.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_item_effect.schema.json",
			},
		}...),
	}

	itemEffectConfig[SearchManyAdventureGameItemEffects] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-item-effects",
		HandlerFunc: searchManyAdventureGameItemEffectsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game item effects",
		},
	}

	itemEffectConfig[GetManyAdventureGameItemEffects] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/item-effects",
		HandlerFunc: getManyAdventureGameItemEffectsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game item effects",
		},
	}

	itemEffectConfig[GetOneAdventureGameItemEffect] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/item-effects/:item_effect_id",
		HandlerFunc: getOneAdventureGameItemEffectHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game item effect",
		},
	}

	itemEffectConfig[CreateOneAdventureGameItemEffect] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/item-effects",
		HandlerFunc: createOneAdventureGameItemEffectHandler,
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
			Title:    "Create adventure game item effect",
		},
	}

	itemEffectConfig[UpdateOneAdventureGameItemEffect] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/item-effects/:item_effect_id",
		HandlerFunc: updateOneAdventureGameItemEffectHandler,
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
			Title:    "Update adventure game item effect",
		},
	}

	itemEffectConfig[DeleteOneAdventureGameItemEffect] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/item-effects/:item_effect_id",
		HandlerFunc: deleteOneAdventureGameItemEffectHandler,
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
			Title:    "Delete adventure game item effect",
		},
	}

	return itemEffectConfig, nil
}

func searchManyAdventureGameItemEffectsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameItemEffectsHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameItemEffectRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game item effect records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemEffectRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameItemEffectsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameItemEffectsHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	opts.Params = append(opts.Params, sql.Param{
		Col: adventure_game_record.FieldAdventureGameItemEffectGameID,
		Val: gameID,
	})

	recs, err := mm.GetManyAdventureGameItemEffectRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game item effect records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemEffectRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameItemEffectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameItemEffectHandler")

	gameID := pp.ByName("game_id")
	itemEffectID := pp.ByName("item_effect_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameItemEffectRec(itemEffectID, nil)
	if err != nil {
		l.Warn("failed getting adventure game item effect record >%v<", err)
		return err
	}

	if rec.GameID != gameID {
		l.Warn("item effect does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item effect", itemEffectID)
	}

	res, err := mapper.AdventureGameItemEffectRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game item effect record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameItemEffectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameItemEffectHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed getting game record >%v<", err)
		return err
	}

	rec := &adventure_game_record.AdventureGameItemEffect{
		GameID: gameRec.ID,
	}

	rec, err = mapper.AdventureGameItemEffectRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateAdventureGameItemEffectRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game item effect record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemEffectRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameItemEffectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameItemEffectHandler")

	gameID := pp.ByName("game_id")
	itemEffectID := pp.ByName("item_effect_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameItemEffectRec(itemEffectID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("item effect does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item effect", itemEffectID)
	}

	rec, err = mapper.AdventureGameItemEffectRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameItemEffectRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game item effect record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameItemEffectRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameItemEffectHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameItemEffectHandler")

	gameID := pp.ByName("game_id")
	itemEffectID := pp.ByName("item_effect_id")

	l.Info("deleting adventure game item effect record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameItemEffectRec(itemEffectID, nil)
	if err != nil {
		return err
	}

	if rec.GameID != gameID {
		l.Warn("item effect does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("item effect", itemEffectID)
	}

	if err := mm.DeleteAdventureGameItemEffectRec(itemEffectID); err != nil {
		l.Warn("failed deleting adventure game item effect record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
