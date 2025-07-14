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
	"gitlab.com/alienspaces/playbymail/internal/record"
)

const (
	getManyGameLocationLinks = "get-location-links"
	getOneGameLocationLink   = "get-location-link"
	createGameLocationLink   = "create-location-link"
	deleteGameLocationLink   = "delete-location-link"
)

func (rnr *Runner) gameLocationLinkHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameLocationLinkHandlerConfig")

	l.Debug("Adding location_link handler configuration")

	GameLocationLinkConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "location_link.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "location_link.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "location_link.response.schema.json",
		},
		References: referenceSchemas,
	}

	GameLocationLinkConfig[getManyGameLocationLinks] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/location-links",
		HandlerFunc: rnr.getManyGameLocationLinksHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get location link collection",
		},
	}

	GameLocationLinkConfig[getOneGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/location-links/:location_link_id",
		HandlerFunc: rnr.getGameLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get location link",
		},
	}

	GameLocationLinkConfig[createGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/location-links",
		HandlerFunc: rnr.createGameLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create location link",
		},
	}

	GameLocationLinkConfig[deleteGameLocationLink] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/location-links/:location_link_id",
		HandlerFunc: rnr.deleteGameLocationLinkHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete location link",
		},
	}

	return GameLocationLinkConfig, nil
}

func (rnr *Runner) getManyGameLocationLinksHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyGameLocationLinksHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyGameLocationLinkRecs(opts)
	if err != nil {
		l.Warn("failed getting game_location_link records >%v<", err)
		return err
	}

	res, err := mapper.GameLocationLinkRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) getGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetGameLocationLinkHandler")

	gameLocationLinkID := pp.ByName("location_link_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetGameLocationLinkRec(gameLocationLinkID, nil)
	if err != nil {
		l.Warn("failed getting game_location_link record >%v<", err)
		return err
	}

	res, err := mapper.GameLocationLinkRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game_location_link record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) createGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateGameLocationLinkHandler")

	rec, err := mapper.GameLocationLinkRequestToRecord(l, r, &record.GameLocationLink{})
	if err != nil {
		return err
	}
	if rec.GameID == "" {
		res := map[string]interface{}{"error": map[string]interface{}{"code": "missing_game_id", "detail": "game_id is required"}}
		_ = server.WriteResponse(l, w, http.StatusBadRequest, res)
		return nil
	}

	mm := m.(*domain.Domain)

	rec, err = mm.CreateGameLocationLinkRec(rec)
	if err != nil {
		l.Warn("failed creating game_location_link record >%v<", err)
		return err
	}

	res, err := mapper.GameLocationLinkRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func (rnr *Runner) deleteGameLocationLinkHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteGameLocationLinkHandler")

	GameLocationLinkID := pp.ByName("location_link_id")

	l.Info("deleting location_link record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	if err := mm.DeleteGameLocationLinkRec(GameLocationLinkID); err != nil {
		l.Warn("failed deleting location_link record >%v<", err)
		return err
	}

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
