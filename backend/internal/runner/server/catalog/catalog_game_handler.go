package catalog

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/schema/api/catalog_schema"
)

const (
	GetCatalogGames = "get-catalog-games"
)

func catalogGameHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "catalogGameHandlerConfig")

	l.Debug("Adding catalog game handler configuration")

	config := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/catalog_schema",
			Name:     "catalog_game.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/catalog_schema",
				Name:     "catalog_game.schema.json",
			},
		}...),
	}

	config[GetCatalogGames] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/catalog/games",
		HandlerFunc: getCatalogGamesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypePublic},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Collection:  true,
			Title:       "Get game catalog",
			Description: "Returns all games with active manager subscriptions that have game instances open for player enrollment.",
		},
	}

	return config, nil
}

func getCatalogGamesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getCatalogGamesHandler")

	mm := m.(*domain.Domain)

	// Find all games that have at least one active manager subscription.
	// Use the repository directly so queries are not constrained by per-account RLS
	// on this public endpoint.
	managerSubscriptions, err := mm.GameSubscriptionRepository().GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionSubscriptionType, Val: game_record.GameSubscriptionTypeManager},
			{Col: game_record.FieldGameSubscriptionStatus, Val: game_record.GameSubscriptionStatusActive},
		},
	})
	if err != nil {
		l.Warn("failed getting manager subscriptions >%v<", err)
		return err
	}

	// Collect unique game IDs that have an active manager.
	gameIDSet := map[string]struct{}{}
	for _, sub := range managerSubscriptions {
		gameIDSet[sub.GameID] = struct{}{}
	}

	if len(gameIDSet) == 0 {
		l.Info("no games with active manager subscriptions found")
		res := catalog_schema.CatalogGameCollectionResponse{
			Data: []*catalog_schema.CatalogGameResponseData{},
		}
		return server.WriteResponse(l, w, http.StatusOK, res)
	}

	gameIDs := make([]string, 0, len(gameIDSet))
	for id := range gameIDSet {
		gameIDs = append(gameIDs, id)
	}

	// Fetch the game records.
	gameRecs, err := mm.GetManyGameRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameID, Val: gameIDs},
		},
	})
	if err != nil {
		l.Warn("failed getting game records >%v<", err)
		return err
	}

	// Build the catalog: for each game, collect instances in "created" status
	// that still have available player capacity.
	catalogData := make([]*catalog_schema.CatalogGameResponseData, 0, len(gameRecs))

	for _, gameRec := range gameRecs {
		instanceRecs, err := mm.GetManyGameInstanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: game_record.FieldGameInstanceGameID, Val: gameRec.ID},
				{Col: game_record.FieldGameInstanceStatus, Val: game_record.GameInstanceStatusCreated},
			},
		})
		if err != nil {
			l.Warn("failed getting game instances for game >%s< >%v<", gameRec.ID, err)
			return err
		}

		availableInstances := make([]*catalog_schema.CatalogGameInstanceData, 0)
		for _, inst := range instanceRecs {
			playerCount, err := mm.GetPlayerCountForGameInstance(inst.ID)
			if err != nil {
				l.Warn("failed getting player count for instance >%s< >%v<", inst.ID, err)
				return err
			}

			hasCapacity, err := mm.HasAvailableCapacity(inst.ID)
			if err != nil {
				l.Warn("failed checking capacity for instance >%s< >%v<", inst.ID, err)
				return err
			}

			if hasCapacity {
				availableInstances = append(availableInstances, mapper.CatalogGameInstanceRecordToData(inst, playerCount))
			}
		}

		// Only include games that have at least one open instance.
		if len(availableInstances) > 0 {
			catalogData = append(catalogData, mapper.CatalogGameRecordToResponseData(l, gameRec, availableInstances))
		}
	}

	res := catalog_schema.CatalogGameCollectionResponse{
		Data: catalogData,
	}

	if err := server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
