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

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-games/{game_id}/location-objects/{location_object_id}/states
// GET (document)    /api/v1/adventure-games/{game_id}/location-objects/{location_object_id}/states/{state_id}
// POST (document)   /api/v1/adventure-games/{game_id}/location-objects/{location_object_id}/states
// PUT (document)    /api/v1/adventure-games/{game_id}/location-objects/{location_object_id}/states/{state_id}
// DELETE (document) /api/v1/adventure-games/{game_id}/location-objects/{location_object_id}/states/{state_id}

const (
	GetManyAdventureGameLocationObjectStates  = "getManyAdventureGameLocationObjectStates"
	GetOneAdventureGameLocationObjectState    = "getOneAdventureGameLocationObjectState"
	CreateOneAdventureGameLocationObjectState = "createOneAdventureGameLocationObjectState"
	UpdateOneAdventureGameLocationObjectState = "updateOneAdventureGameLocationObjectState"
	DeleteOneAdventureGameLocationObjectState = "deleteOneAdventureGameLocationObjectState"
)

func adventureGameLocationObjectStateHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameLocationObjectStateHandlerConfig")

	l.Debug("Adding adventure_game_location_object_state handler configuration")

	stateConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_object_state.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_object_state.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_object_state.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_object_state.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_object_state.schema.json",
			},
		}...),
	}

	stateConfig[GetManyAdventureGameLocationObjectStates] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-objects/:location_object_id/states",
		HandlerFunc: getManyAdventureGameLocationObjectStatesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game location object states",
		},
	}

	stateConfig[GetOneAdventureGameLocationObjectState] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/location-objects/:location_object_id/states/:state_id",
		HandlerFunc: getOneAdventureGameLocationObjectStateHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game location object state",
		},
	}

	stateConfig[CreateOneAdventureGameLocationObjectState] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/location-objects/:location_object_id/states",
		HandlerFunc: createOneAdventureGameLocationObjectStateHandler,
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
			Title:    "Create adventure game location object state",
		},
	}

	stateConfig[UpdateOneAdventureGameLocationObjectState] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/location-objects/:location_object_id/states/:state_id",
		HandlerFunc: updateOneAdventureGameLocationObjectStateHandler,
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
			Title:    "Update adventure game location object state",
		},
	}

	stateConfig[DeleteOneAdventureGameLocationObjectState] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/location-objects/:location_object_id/states/:state_id",
		HandlerFunc: deleteOneAdventureGameLocationObjectStateHandler,
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
			Title:    "Delete adventure game location object state",
		},
	}

	return stateConfig, nil
}

func getManyAdventureGameLocationObjectStatesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameLocationObjectStatesHandler")

	locationObjectID := pp.ByName("location_object_id")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	opts.Params = append(opts.Params, sql.Param{
		Col: adventure_game_record.FieldAdventureGameLocationObjectStateAdventureGameLocationObjectID,
		Val: locationObjectID,
	})

	recs, err := mm.GetManyAdventureGameLocationObjectStateRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location object state records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectStateRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameLocationObjectStateHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameLocationObjectStateHandler")

	locationObjectID := pp.ByName("location_object_id")
	stateID := pp.ByName("state_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameLocationObjectStateRec(stateID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location object state record >%v<", err)
		return err
	}

	if rec.AdventureGameLocationObjectID != locationObjectID {
		l.Warn("state does not belong to specified location object >%s< != >%s<", rec.AdventureGameLocationObjectID, locationObjectID)
		return coreerror.NewNotFoundError("location object state", stateID)
	}

	res, err := mapper.AdventureGameLocationObjectStateRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game location object state record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createOneAdventureGameLocationObjectStateHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameLocationObjectStateHandler")

	gameID := pp.ByName("game_id")
	locationObjectID := pp.ByName("location_object_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	objectRec, err := mm.GetAdventureGameLocationObjectRec(locationObjectID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location object record >%v<", err)
		return err
	}

	if objectRec.GameID != gameID {
		l.Warn("location object does not belong to specified game >%s< != >%s<", objectRec.GameID, gameID)
		return coreerror.NewNotFoundError("location object", locationObjectID)
	}

	rec := &adventure_game_record.AdventureGameLocationObjectState{
		GameID:                        gameID,
		AdventureGameLocationObjectID: locationObjectID,
	}

	rec, err = mapper.AdventureGameLocationObjectStateRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateAdventureGameLocationObjectStateRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location object state record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectStateRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateOneAdventureGameLocationObjectStateHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameLocationObjectStateHandler")

	gameID := pp.ByName("game_id")
	locationObjectID := pp.ByName("location_object_id")
	stateID := pp.ByName("state_id")
	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameLocationObjectStateRec(stateID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if rec.AdventureGameLocationObjectID != locationObjectID {
		l.Warn("state does not belong to specified location object >%s< != >%s<", rec.AdventureGameLocationObjectID, locationObjectID)
		return coreerror.NewNotFoundError("location object state", stateID)
	}

	rec, err = mapper.AdventureGameLocationObjectStateRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameLocationObjectStateRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game location object state record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationObjectStateRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteOneAdventureGameLocationObjectStateHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameLocationObjectStateHandler")

	gameID := pp.ByName("game_id")
	locationObjectID := pp.ByName("location_object_id")
	stateID := pp.ByName("state_id")

	l.Info("deleting adventure game location object state record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	if _, err := authorizeDesignerModify(l, r, mm, gameID); err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameLocationObjectStateRec(stateID, nil)
	if err != nil {
		return err
	}

	if rec.AdventureGameLocationObjectID != locationObjectID {
		l.Warn("state does not belong to specified location object >%s< != >%s<", rec.AdventureGameLocationObjectID, locationObjectID)
		return coreerror.NewNotFoundError("location object state", stateID)
	}

	if err := mm.DeleteAdventureGameLocationObjectStateRec(stateID); err != nil {
		l.Warn("failed deleting adventure game location object state record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
