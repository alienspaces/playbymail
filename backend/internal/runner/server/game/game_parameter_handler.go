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

const (
	GetManyGameParameters     = "get-game-parameters"
	GetOneGameParameter       = "get-game-parameter"
	CreateGameParameter       = "create-game-parameter"
	CreateGameParameterWithID = "create-game-parameter-with-id"
	UpdateGameParameter       = "update-game-parameter"
	DeleteGameParameter       = "delete-game-parameter"
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

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_parameter.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_parameter.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_parameter.schema.json",
			},
		}...),
	}

	// Unnested routes
	gameParameterConfig[GetManyGameParameters] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-parameters",
		HandlerFunc: getManyGameParametersHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get game parameter collection",
		},
	}

	gameParameterConfig[GetOneGameParameter] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-parameters/:game_parameter_id",
		HandlerFunc: getGameParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get game parameter",
		},
	}

	gameParameterConfig[CreateGameParameter] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-parameters",
		HandlerFunc: createGameParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create game parameter",
		},
	}

	gameParameterConfig[UpdateGameParameter] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/game-parameters/:game_parameter_id",
		HandlerFunc: updateGameParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update game parameter",
		},
	}

	gameParameterConfig[DeleteGameParameter] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/game-parameters/:game_parameter_id",
		HandlerFunc: deleteGameParameterHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete game parameter",
		},
	}

	return gameParameterConfig, nil
}

func getManyGameParametersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameParametersHandler")

	l.Info("getting many game parameters")

	opts := queryparam.ToSQLOptions(qp)

	recs, err := m.(*domain.Domain).GetGameParameterRecs(opts)
	if err != nil {
		l.Warn("failed getting game parameters >%v<", err)
		return err
	}

	response, err := mapper.GameParameterRecordsToCollectionResponse(l, recs)
	if err != nil {
		l.Warn("failed mapping game parameter records to collection response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, response)
}

func getGameParameterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameParameterHandler")

	gameParameterID := pp.ByName("game_parameter_id")

	l.Info("getting game parameter with id >%s<", gameParameterID)

	rec, err := m.(*domain.Domain).GetGameParameterRec(gameParameterID, nil)
	if err != nil {
		l.Warn("failed getting game parameter >%v<", err)
		return err
	}

	response, err := mapper.GameParameterRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game parameter record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, response)
}

func createGameParameterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createGameParameterHandler")

	l.Info("creating game parameter")

	rec := &game_record.GameParameter{}
	rec, err := mapper.GameParameterRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping game parameter request to record >%v<", err)
		return err
	}

	rec, err = m.(*domain.Domain).CreateGameParameterRec(rec)
	if err != nil {
		l.Warn("failed creating game parameter >%v<", err)
		return err
	}

	response, err := mapper.GameParameterRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game parameter record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, response)
}

func updateGameParameterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateGameParameterHandler")

	gameParameterID := pp.ByName("game_parameter_id")

	l.Info("updating game parameter with id >%s<", gameParameterID)

	rec, err := m.(*domain.Domain).GetGameParameterRec(gameParameterID, nil)
	if err != nil {
		l.Warn("failed getting game parameter >%v<", err)
		return err
	}

	rec, err = mapper.GameParameterRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping game parameter request to record >%v<", err)
		return err
	}

	rec, err = m.(*domain.Domain).UpdateGameParameterRec(rec)
	if err != nil {
		l.Warn("failed updating game parameter >%v<", err)
		return err
	}

	response, err := mapper.GameParameterRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping game parameter record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, response)
}

func deleteGameParameterHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteGameParameterHandler")

	gameParameterID := pp.ByName("game_parameter_id")

	l.Info("deleting game parameter with id >%s<", gameParameterID)

	_, err := m.(*domain.Domain).GetGameParameterRec(gameParameterID, nil)
	if err != nil {
		l.Warn("failed getting game parameter >%v<", err)
		return err
	}

	err = m.(*domain.Domain).DeleteGameParameterRec(gameParameterID)
	if err != nil {
		l.Warn("failed deleting game parameter >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
