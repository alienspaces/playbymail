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
	tagGroupGameLocationLinkRequirement server.TagGroup = "GameLocationLinkRequirements"
	TagGameLocationLinkRequirement      server.Tag      = "GameLocationLinkRequirements"
)

const (
	getManyGameLocationLinkRequirements = "get-game-location-link-requirements"
	getOneGameLocationLinkRequirement   = "get-game-location-link-requirement"
	createGameLocationLinkRequirement   = "create-game-location-link-requirement"
	updateGameLocationLinkRequirement   = "update-game-location-link-requirement"
	deleteGameLocationLinkRequirement   = "delete-game-location-link-requirement"
)

func (rnr *Runner) gameLocationLinkRequirementHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "gameLocationLinkRequirementHandlerConfig")

	l.Debug("Adding game_location_link_requirement handler configuration")

	cfg := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_location_link_requirement.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_location_link_requirement.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "game_location_link_requirement.response.schema.json",
		},
		References: referenceSchemas,
	}

	cfg[getManyGameLocationLinkRequirements] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-location-link-requirements",
		HandlerFunc: rnr.getManyGameLocationLinkRequirementsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game_location_link_requirement collection",
		},
	}

	cfg[getOneGameLocationLinkRequirement] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-location-link-requirements/:game_location_link_requirement_id",
		HandlerFunc: rnr.getGameLocationLinkRequirementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game_location_link_requirement",
		},
	}

	cfg[createGameLocationLinkRequirement] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/game-location-link-requirements",
		HandlerFunc: rnr.createGameLocationLinkRequirementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game_location_link_requirement",
		},
	}

	cfg[updateGameLocationLinkRequirement] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/game-location-link-requirements/:game_location_link_requirement_id",
		HandlerFunc: rnr.updateGameLocationLinkRequirementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game_location_link_requirement",
		},
	}

	cfg[deleteGameLocationLinkRequirement] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/game-location-link-requirements/:game_location_link_requirement_id",
		HandlerFunc: rnr.deleteGameLocationLinkRequirementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game_location_link_requirement",
		},
	}

	return cfg, nil
}

func (rnr *Runner) getManyGameLocationLinkRequirementsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyGameLocationLinkRequirementsHandler")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	recs, err := mm.GetManyGameLocationLinkRequirementRecs(opts)
	if err != nil {
		l.Warn("failed getting game_location_link_requirement records >%v<", err)
		return err
	}
	res, err := mapper.GameLocationLinkRequirementRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) getGameLocationLinkRequirementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetGameLocationLinkRequirementHandler")
	id := pp.ByName("game_location_link_requirement_id")
	mm := m.(*domain.Domain)
	rec, err := mm.GetGameLocationLinkRequirementRec(id, nil)
	if err != nil {
		l.Warn("failed getting game_location_link_requirement record >%v<", err)
		return err
	}
	res, err := mapper.GameLocationLinkRequirementRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game_location_link_requirement record to response >%v<", err)
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) createGameLocationLinkRequirementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "CreateGameLocationLinkRequirementHandler")
	var req schema.GameLocationLinkRequirementRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}
	if req.GameID == "" {
		res := map[string]interface{}{"error": map[string]interface{}{"code": "missing_game_id", "detail": "game_id is required"}}
		_ = server.WriteResponse(l, w, http.StatusBadRequest, res)
		return nil
	}
	mm := m.(*domain.Domain)
	rec, err := mapper.GameLocationLinkRequirementRequestToRecord(l, &req, nil)
	if err != nil {
		return err
	}
	rec, err = mm.CreateGameLocationLinkRequirementRec(rec)
	if err != nil {
		l.Warn("failed creating game_location_link_requirement record >%v<", err)
		return err
	}
	res, err := mapper.GameLocationLinkRequirementRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) updateGameLocationLinkRequirementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateGameLocationLinkRequirementHandler")
	id := pp.ByName("game_location_link_requirement_id")
	mm := m.(*domain.Domain)
	existing, err := mm.GetGameLocationLinkRequirementRec(id, nil)
	if err != nil {
		l.Warn("failed getting game_location_link_requirement record >%v<", err)
		return err
	}
	var req schema.GameLocationLinkRequirementRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}
	if req.GameID == "" {
		res := map[string]interface{}{"error": map[string]interface{}{"code": "missing_game_id", "detail": "game_id is required"}}
		_ = server.WriteResponse(l, w, http.StatusBadRequest, res)
		return nil
	}
	rec, err := mapper.GameLocationLinkRequirementRequestToRecord(l, &req, existing)
	if err != nil {
		return err
	}
	rec, err = mm.UpdateGameLocationLinkRequirementRec(rec)
	if err != nil {
		l.Warn("failed updating game_location_link_requirement record >%v<", err)
		return err
	}
	res, err := mapper.GameLocationLinkRequirementRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}

func (rnr *Runner) deleteGameLocationLinkRequirementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteGameLocationLinkRequirementHandler")
	id := pp.ByName("game_location_link_requirement_id")
	mm := m.(*domain.Domain)
	if err := mm.DeleteGameLocationLinkRequirementRec(id); err != nil {
		l.Warn("failed deleting game_location_link_requirement record >%v<", err)
		return err
	}
	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}
	return nil
}
