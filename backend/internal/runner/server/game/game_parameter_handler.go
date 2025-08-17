package game

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

// API Resource Search Path
//
// GET (collection) /api/v1/game-parameters
const (
	SearchManyGameParameters = "search-many-game-parameters"
)

func gameParameterHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameParameterHandlerConfig")

	l.Debug("adding game parameter handler configuration")

	// Create a new map to avoid modifying the passed config
	gameParameterConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_parameter.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_parameter.schema.json",
			},
		}...),
	}

	// Search route
	gameParameterConfig[SearchManyGameParameters] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-parameters",
		HandlerFunc: searchManyGameParametersHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Search game parameter collection",
		},
	}

	return gameParameterConfig, nil
}

// searchManyGameParametersHandler supports searching for game parameters by game type only
func searchManyGameParametersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "searchManyGameParametersHandler")

	// Check if game_type parameter exists and has values before accessing
	gameTypeParams, exists := qp.Params["game_type"]
	if !exists || len(gameTypeParams) == 0 {
		l.Warn("game_type parameter is missing or empty, returning all game parameters")
		recs := domain.GetGameParameters()
		res, err := mapper.GameParameterRecordsToCollectionResponse(l, recs)
		if err != nil {
			l.Warn("failed mapping game parameter records to collection response >%v<", err)
			return err
		}
		if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
			l.Warn("failed writing response >%v<", err)
			return err
		}
		return nil
	}

	gameType := gameTypeParams[0].Val.(string)

	l.Info("searching many game parameters for game type >%s<", gameType)

	var recs []*game_record.GameParameter
	if gameType != "" {
		recs = domain.GetGameParametersByGameType(gameType)
	} else {
		recs = domain.GetGameParameters()
	}

	res, err := mapper.GameParameterRecordsToCollectionResponse(l, recs)
	if err != nil {
		l.Warn("failed mapping game parameter records to collection response >%v<", err)
		return err
	}

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
