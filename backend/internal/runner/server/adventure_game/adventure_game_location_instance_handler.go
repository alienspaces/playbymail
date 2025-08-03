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
// GET (collection) /api/v1/adventure-game-location-instances

// API Resource CRUD Paths
//
// GET (collection)  /api/v1/adventure-game-instances/{game_instance_id}/location-instances
// GET (document)    /api/v1/adventure-game-instances/{game_instance_id}/location-instances/{location_instance_id}
//
// Game location instances are created and managed by the game engine and not through the public API.

const (
	// API Resource Search Path
	SearchManyAdventureGameLocationInstances = "search-many-adventure-game-location-instances"

	// API Resource CRUD Paths
	GetManyAdventureGameLocationInstances = "get-many-adventure-game-location-instances"
	GetOneAdventureGameLocationInstance   = "get-one-adventure-game-location-instance"
)

func adventureGameLocationInstanceHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "adventureGameLocationInstanceHandlerConfig")

	l.Debug("Adding adventure_game_location_instance handler configuration")

	gameLocationInstanceConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_instance.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_instance.schema.json",
			},
		}...),
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/adventure_game_schema",
			Name:     "adventure_game_location_instance.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/adventure_game_schema",
				Name:     "adventure_game_location_instance.schema.json",
			},
		}...),
	}

	// New Adventure Game Location Instance API paths
	gameLocationInstanceConfig[SearchManyAdventureGameLocationInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-location-instances",
		HandlerFunc: searchManyAdventureGameLocationInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search adventure game location instances",
		},
	}

	gameLocationInstanceConfig[GetManyAdventureGameLocationInstances] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/location-instances",
		HandlerFunc: getManyAdventureGameLocationInstancesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get adventure game location instances",
		},
	}

	gameLocationInstanceConfig[GetOneAdventureGameLocationInstance] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/adventure-game-instances/:game_instance_id/location-instances/:location_instance_id",
		HandlerFunc: getOneAdventureGameLocationInstanceHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get adventure game location instance",
		},
	}

	return gameLocationInstanceConfig, nil
}

func searchManyAdventureGameLocationInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyAdventureGameLocationInstancesHandler")

	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAdventureGameLocationInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location instance records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getManyAdventureGameLocationInstancesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAdventureGameLocationInstancesHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Add filter for specific game
	opts.Params = append(opts.Params, sql.Param{
		Col: adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID,
		Val: gameInstanceID,
	})

	mm := m.(*domain.Domain)

	recs, err := mm.GetManyAdventureGameLocationInstanceRecs(opts)
	if err != nil {
		l.Warn("failed getting adventure game location instance records >%v<", err)
		return err
	}

	res, err := mapper.AdventureGameLocationInstanceRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getOneAdventureGameLocationInstanceHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getOneAdventureGameLocationInstanceHandler")

	gameInstanceID := pp.ByName("game_instance_id")
	if gameInstanceID == "" {
		l.Warn("game instance id is required")
		return coreerror.NewNotFoundError("game instance", gameInstanceID)
	}

	locationInstanceID := pp.ByName("location_instance_id")
	if locationInstanceID == "" {
		l.Warn("location instance id is required")
		return coreerror.NewNotFoundError("location instance", locationInstanceID)
	}

	mm := m.(*domain.Domain)

	gameInstanceRec, err := mm.GetGameInstanceRec(gameInstanceID, nil)
	if err != nil {
		return err
	}

	rec, err := mm.GetAdventureGameLocationInstanceRec(locationInstanceID, nil)
	if err != nil {
		l.Warn("failed getting adventure game location instance record >%v<", err)
		return err
	}

	// Verify the location instance belongs to the specified game
	if rec.GameInstanceID != gameInstanceRec.ID {
		l.Warn("location instance does not belong to specified game instance >%s< != >%s<", rec.GameInstanceID, gameInstanceRec.ID)
		return coreerror.NewNotFoundError("location instance", locationInstanceID)
	}

	res, err := mapper.AdventureGameLocationInstanceRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping adventure game location instance record to response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
