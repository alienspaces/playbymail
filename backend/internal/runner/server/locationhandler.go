package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/schema"
)

const (
	tagGroupGameLocation server.TagGroup = "GameLocations"
	TagGameLocation      server.Tag      = "GameLocations"
)

const (
	getManyGameLocations = "get-game-locations"
	getOneGameLocation   = "get-game-location"
	createGameLocation   = "create-game-location"
	updateGameLocation   = "update-game-location"
	deleteGameLocation   = "delete-game-location"
)

func (rnr *Runner) gameLocationHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameLocationHandlerConfig")

	l.Debug("Adding game_location handler configuration")

	gameLocationConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_location.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_location.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_location.response.schema.json",
		},
		References: referenceSchemas,
	}

	gameLocationConfig[getManyGameLocations] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-locations",
		HandlerFunc: rnr.getManyGameLocationsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game_location collection",
		},
	}

	gameLocationConfig[getOneGameLocation] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-locations/:game_location_id",
		HandlerFunc: rnr.getGameLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game_location",
		},
	}

	gameLocationConfig[createGameLocation] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/game-locations",
		HandlerFunc: rnr.createGameLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game_location",
		},
	}

	gameLocationConfig[updateGameLocation] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/game-locations/:game_location_id",
		HandlerFunc: rnr.updateGameLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game_location",
		},
	}

	gameLocationConfig[deleteGameLocation] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/game-locations/:game_location_id",
		HandlerFunc: rnr.deleteGameLocationHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game_location",
		},
	}

	return gameLocationConfig, nil
}

func (rnr *Runner) getManyGameLocationsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyGameLocationsHandler")

	l.Info("querying many game_location records with params >%#v<", qp)

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyGameLocationRecs(opts)
	if err != nil {
		l.Warn("failed getting game_location records >%v<", err)
		return err
	}

	data, err := server.Paginate(l, recs, mapper.GameLocationRecordToResponseData, qp.PageSize)
	if err != nil {
		return err
	}

	l.Info("responding with >%d< game_location records", len(data))

	res := schema.GameLocationCollectionResponse(data)

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getGameLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetGameLocationHandler")

	gameLocationID := pp.ByName("game_location_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetGameLocationRec(gameLocationID, nil)
	if err != nil {
		l.Warn("failed getting game_location record >%v<", err)
		return err
	}

	data, err := mapper.GameLocationRecordToResponseData(l, rec)
	if err != nil {
		l.Warn("failed mapping game_location record to response data >%v<", err)
		return err
	}

	l.Info("responding with game_location record id >%s<", rec.ID)

	res := schema.GameLocationResponse{
		GameLocationResponseData: &data,
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) createGameLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateGameLocationHandler")

	l.Info("creating game_location record with path params >%#v<", pp)

	var req schema.GameLocationRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	rec, err := mapper.GameLocationRequestToRecord(l, &req, nil)
	if err != nil {
		return err
	}

	mm := m.(*domain.Domain)

	rec, err = mm.CreateGameLocationRec(rec)
	if err != nil {
		l.Warn("failed creating game_location record >%v<", err)
		return err
	}

	respData, err := mapper.GameLocationRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.GameLocationResponse{
		GameLocationResponseData: &respData,
	}

	l.Info("responding with created game_location record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) updateGameLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateGameLocationHandler")

	gameLocationID := pp.ByName("game_location_id")

	l.Info("updating game_location record with path params >%#v<", pp)

	var req schema.GameLocationRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetGameLocationRec(gameLocationID, nil)
	if err != nil {
		return err
	}

	rec, err = mapper.GameLocationRequestToRecord(l, &req, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateGameLocationRec(rec)
	if err != nil {
		l.Warn("failed updating game_location record >%v<", err)
		return err
	}

	data, err := mapper.GameLocationRecordToResponseData(l, rec)
	if err != nil {
		return err
	}

	res := schema.GameLocationResponse{
		GameLocationResponseData: &data,
	}

	l.Info("responding with updated game_location record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) deleteGameLocationHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteGameLocationHandler")

	gameLocationID := pp.ByName("game_location_id")

	l.Info("deleting game_location record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	if err := mm.DeleteGameLocationRec(gameLocationID); err != nil {
		l.Warn("failed deleting game_location record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
