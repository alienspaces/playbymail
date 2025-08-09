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

const (
	SearchManyAdventureGameCreaturePlacements = "searchManyAdventureGameCreaturePlacements"
	CreateOneAdventureGameCreaturePlacement   = "createOneAdventureGameCreaturePlacement"
	GetOneAdventureGameCreaturePlacement      = "getOneAdventureGameCreaturePlacement"
	UpdateOneAdventureGameCreaturePlacement   = "updateOneAdventureGameCreaturePlacement"
	DeleteOneAdventureGameCreaturePlacement   = "deleteOneAdventureGameCreaturePlacement"
)

// adventureGameCreaturePlacementHandlerConfig -
func adventureGameCreaturePlacementHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameCreaturePlacementHandlerConfig")

	l.Debug("Adding adventure_game_creature_placement handler configuration")

	creaturePlacementConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_creature_placement.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_creature_placement.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_creature_placement.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_creature_placement.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_creature_placement.schema.json",
			},
		}...),
	}

	creaturePlacementConfig[SearchManyAdventureGameCreaturePlacements] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/creature-placements",
		HandlerFunc: searchManyAdventureGameCreaturePlacementsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get many adventure game creature placements",
		},
	}

	creaturePlacementConfig[CreateOneAdventureGameCreaturePlacement] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/adventure-games/:game_id/creature-placements",
		HandlerFunc: createOneAdventureGameCreaturePlacementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create one adventure game creature placement",
		},
	}

	creaturePlacementConfig[GetOneAdventureGameCreaturePlacement] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-games/:game_id/creature-placements/:placement_id",
		HandlerFunc: getOneAdventureGameCreaturePlacementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get one adventure game creature placement",
		},
	}

	creaturePlacementConfig[UpdateOneAdventureGameCreaturePlacement] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/adventure-games/:game_id/creature-placements/:placement_id",
		HandlerFunc: updateOneAdventureGameCreaturePlacementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update one adventure game creature placement",
		},
	}

	creaturePlacementConfig[DeleteOneAdventureGameCreaturePlacement] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/adventure-games/:game_id/creature-placements/:placement_id",
		HandlerFunc: deleteOneAdventureGameCreaturePlacementHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete one adventure game creature placement",
		},
	}

	return creaturePlacementConfig, nil
}

// searchManyAdventureGameCreaturePlacementsHandler -
func searchManyAdventureGameCreaturePlacementsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameCreaturePlacementsHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameCreaturePlacementRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game creature placement records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCreaturePlacementRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// getOneAdventureGameCreaturePlacementHandler -
func getOneAdventureGameCreaturePlacementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameCreaturePlacementHandler")

	gameID := pp.ByName("game_id")
	placementID := pp.ByName("placement_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameCreaturePlacementRec(placementID, nil)
	if err != nil {
		l.Warn("failed getting adventure game creature placement record >%v<", err)
		return err
	}

	// Verify the placement belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("placement does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("creature_placement", placementID)
	}

	res, err := mapper.AdventureGameCreaturePlacementRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game creature placement record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// createOneAdventureGameCreaturePlacementHandler -
func createOneAdventureGameCreaturePlacementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createOneAdventureGameCreaturePlacementHandler")

	gameID := pp.ByName("game_id")
	mm := m.(*domain.Domain)

	// RLS constraints will be applied by the repository
	gameRec, err := mm.GetGameRec(gameID, nil)
	if err != nil {
		l.Warn("failed getting adventure game record >%v<", err)
		return err
	}

	rec := &adventure_game_record.AdventureGameCreaturePlacement{
		GameID: gameRec.ID,
	}

	rec, err = mapper.AdventureGameCreaturePlacementRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.CreateAdventureGameCreaturePlacementRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game creature placement record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCreaturePlacementRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game creature placement record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// updateOneAdventureGameCreaturePlacementHandler -
func updateOneAdventureGameCreaturePlacementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateOneAdventureGameCreaturePlacementHandler")

	gameID := pp.ByName("game_id")
	placementID := pp.ByName("placement_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameCreaturePlacementRec(placementID, sql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed getting adventure game creature placement record >%v<", err)
		return err
	}

	// Verify the placement belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("placement does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("creature_placement", placementID)
	}

	rec, err = mapper.AdventureGameCreaturePlacementRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}

	rec, err = mm.UpdateAdventureGameCreaturePlacementRec(rec)
	if err != nil {
		l.Warn("failed updating adventure game creature placement record >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCreaturePlacementRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game creature placement record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// deleteOneAdventureGameCreaturePlacementHandler -
func deleteOneAdventureGameCreaturePlacementHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteOneAdventureGameCreaturePlacementHandler")

	gameID := pp.ByName("game_id")
	placementID := pp.ByName("placement_id")
	mm := m.(*domain.Domain)

	rec, err := mm.GetAdventureGameCreaturePlacementRec(placementID, sql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed getting adventure game creature placement record >%v<", err)
		return err
	}

	// Verify the placement belongs to the specified game
	if rec.GameID != gameID {
		l.Warn("placement does not belong to specified game >%s< != >%s<", rec.GameID, gameID)
		return coreerror.NewNotFoundError("creature_placement", placementID)
	}

	err = mm.DeleteAdventureGameCreaturePlacementRec(placementID)
	if err != nil {
		l.Warn("failed deleting adventure game creature placement record >%v<", err)
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
