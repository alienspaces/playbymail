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
		Path:        "/api/v1/catalog/game-subscriptions",
		HandlerFunc: getCatalogGamesHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypePublic},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Collection:  true,
			Title:       "Get game catalog",
			Description: "Returns all active manager subscriptions that have game instances open for player enrollment. Each entry represents a manager's offering of a game.",
		},
	}

	return config, nil
}

func getCatalogGamesHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getCatalogGamesHandler")

	mm := m.(*domain.Domain)

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

	if len(managerSubscriptions) == 0 {
		l.Info("no active manager subscriptions found")
		return server.WriteResponse(l, w, http.StatusOK, catalog_schema.CatalogCollectionResponse{
			Data: []*catalog_schema.CatalogSubscriptionData{},
		})
	}

	// Pre-fetch all game records referenced by subscriptions.
	gameIDSet := map[string]struct{}{}
	for _, sub := range managerSubscriptions {
		gameIDSet[sub.GameID] = struct{}{}
	}
	gameIDs := make([]string, 0, len(gameIDSet))
	for id := range gameIDSet {
		gameIDs = append(gameIDs, id)
	}
	gameRecs, err := mm.GetManyGameRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameID, Val: gameIDs},
		},
	})
	if err != nil {
		l.Warn("failed getting game records >%v<", err)
		return err
	}
	gameByID := make(map[string]*game_record.Game, len(gameRecs))
	for _, g := range gameRecs {
		gameByID[g.ID] = g
	}

	catalogData := make([]*catalog_schema.CatalogSubscriptionData, 0, len(managerSubscriptions))

	for _, sub := range managerSubscriptions {
		gameRec, ok := gameByID[sub.GameID]
		if !ok {
			l.Warn("game >%s< not found for subscription >%s<", sub.GameID, sub.ID)
			continue
		}

		// Find instances linked to this subscription via game_subscription_instance.
		gsiRecs, err := mm.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: game_record.FieldGameSubscriptionInstanceGameSubscriptionID, Val: sub.ID},
			},
		})
		if err != nil {
			l.Warn("failed getting subscription instances for subscription >%s< >%v<", sub.ID, err)
			return err
		}

		// For each linked instance, check status and capacity.
		availableInstances := make([]*game_record.GameInstance, 0)
		playerCounts := make(map[string]int)
		for _, gsi := range gsiRecs {
			instRec, err := mm.GetGameInstanceRec(gsi.GameInstanceID, nil)
			if err != nil {
				l.Warn("failed getting instance >%s< >%v<", gsi.GameInstanceID, err)
				continue
			}
			if instRec.Status != game_record.GameInstanceStatusCreated {
				continue
			}
			hasCapacity, err := mm.HasAvailableCapacity(instRec.ID)
			if err != nil {
				l.Warn("failed checking capacity for instance >%s< >%v<", instRec.ID, err)
				continue
			}
			if !hasCapacity {
				continue
			}
			count, err := mm.GetPlayerCountForGameInstance(instRec.ID)
			if err != nil {
				l.Warn("failed getting player count for instance >%s< >%v<", instRec.ID, err)
				continue
			}
			availableInstances = append(availableInstances, instRec)
			playerCounts[instRec.ID] = count
		}

		if len(availableInstances) == 0 {
			continue
		}

		catalogData = append(catalogData, mapper.CatalogSubscriptionToData(l, sub, gameRec, availableInstances, playerCounts))
	}

	res := catalog_schema.CatalogCollectionResponse{
		Data: catalogData,
	}

	if err := server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
