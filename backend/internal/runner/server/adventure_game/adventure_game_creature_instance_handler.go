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
// GET (collection) /api/v1/adventure-game-creature-instances

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-game-instances/{game_instance_id}/creatures
// GET (document)    /api/v1/adventure-game-instances/{game_instance_id}/creatures/{creature_instance_id}
//
// Creature instances are created and managed by the game engine and not through the public API.

const (
	// API Resource Search Path
	searchManyAdventureGameCreatureInstances = "search-many-adventure-game-creature-instances"

	// API Resource CRUD Paths
	getManyAdventureGameCreatureInstances = "get-many-adventure-game-creature-instances"
	getOneAdventureGameCreatureInstance   = "get-one-adventure-game-creature-instance"
)

func adventureGameCreatureInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameCreatureInstanceHandlerConfig")

	l.Debug("Adding adventure_game_creature_instance handler configuration")

	gameCreatureInstanceConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_creature_instance.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_creature_instance.schema.json",
			},
		}...),
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_creature_instance.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_creature_instance.schema.json",
			},
		}...),
	}

	// New Adventure Game Creature Instance API paths
	gameCreatureInstanceConfig[searchManyAdventureGameCreatureInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-creature-instances",
		HandlerFunc: searchManyAdventureGameCreatureInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game creature instances",
		},
	}

	gameCreatureInstanceConfig[getManyAdventureGameCreatureInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/creature-instances",
		HandlerFunc: getManyAdventureGameCreatureInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game creature instances",
		},
	}

	gameCreatureInstanceConfig[getOneAdventureGameCreatureInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/creature-instances/:creature_instance_id",
		HandlerFunc: getOneAdventureGameCreatureInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game creature instance",
		},
	}

	return gameCreatureInstanceConfig, nil
}

func searchManyAdventureGameCreatureInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameCreatureInstancesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameCreatureInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game creature instance records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameCreatureInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameCreatureInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameCreatureInstancesHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	mm := m.(*domain.Domain)

	// Create SQL options from query parameters
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID,
		Val: gameInstanceID,
	})

	recs, err := mm.GetManyAdventureGameCreatureInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game creature instance records >%v<", err)
		return err
	}
	res, err := mapper.AdventureGameCreatureInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameCreatureInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameCreatureInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	creatureInstanceID := pp.ByName("creature_instance_id")
	if creatureInstanceID == "" {
		l.Warn("creature instance id is required")
		return coreerror.NewNotFoundError("creature instance", creatureInstanceID)
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameCreatureInstanceRec(creatureInstanceID, nil)
	if err != nil {
		l.Warn("failed getting adventure game creature instance record >%v<", err)
		return err
	}

	// Verify the creature instance belongs to the specified game
	if rec.GameInstanceID != gameInstanceRec.ID {
		l.Warn("creature instance does not belong to specified game instance >%s< != >%s<", rec.GameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("creature instance", creatureInstanceID)
	}

	res, err := mapper.AdventureGameCreatureInstanceRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game creature instance record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
