package game

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
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	GetManyGameSubscriptions = "get-game-subscriptions"
	GetOneGameSubscription   = "get-game-subscription"
	CreateGameSubscription   = "create-game-subscription"
	UpdateGameSubscription   = "update-game-subscription"
	DeleteGameSubscription   = "delete-game-subscription"
)

func gameSubscriptionHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameSubscriptionHandlerConfig")

	l.Debug("Adding game subscription handler configuration")

	config := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main:       jsonschema.Schema{Name: "game_subscription.collection.response.schema.json"},
		References: referenceSchemas,
	}
	requestSchema := jsonschema.SchemaWithReferences{
		Main:       jsonschema.Schema{Name: "game_subscription.request.schema.json"},
		References: referenceSchemas,
	}
	responseSchema := jsonschema.SchemaWithReferences{
		Main:       jsonschema.Schema{Name: "game_subscription.response.schema.json"},
		References: referenceSchemas,
	}

	config[GetManyGameSubscriptions] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-subscriptions",
		HandlerFunc: getManyGameSubscriptionsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
	}
	config[GetOneGameSubscription] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/game-subscriptions/:game_subscription_id",
		HandlerFunc: getGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
	}
	config[CreateGameSubscription] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/game-subscriptions",
		HandlerFunc: createGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
	}
	config[UpdateGameSubscription] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/game-subscriptions/:game_subscription_id",
		HandlerFunc: updateGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
	}
	config[DeleteGameSubscription] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/game-subscriptions/:game_subscription_id",
		HandlerFunc: deleteGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
		},
	}

	return config, nil
}

func getManyGameSubscriptionsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameSubscriptionsHandler")
	mm := m.(*domain.Domain)
	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	recs, err := mm.GetManyGameSubscriptionRecs(opts)
	if err != nil {
		return err
	}
	res, err := mapper.GameSubscriptionRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getGameSubscriptionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameSubscriptionHandler")
	mm := m.(*domain.Domain)
	recID := pp.ByName("game_subscription_id")
	rec, err := mm.GetGameSubscriptionRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	res, err := mapper.GameSubscriptionRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createGameSubscriptionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createGameSubscriptionHandler")
	mm := m.(*domain.Domain)
	rec, err := mapper.GameSubscriptionRequestToRecord(l, r, &record.GameSubscription{})
	if err != nil {
		return err
	}
	rec, err = mm.CreateGameSubscriptionRec(rec)
	if err != nil {
		return err
	}
	res, err := mapper.GameSubscriptionRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateGameSubscriptionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateGameSubscriptionHandler")
	mm := m.(*domain.Domain)
	recID := pp.ByName("game_subscription_id")
	rec, err := mm.GetGameSubscriptionRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	rec, err = mapper.GameSubscriptionRequestToRecord(l, r, rec)
	if err != nil {
		return err
	}
	rec, err = mm.UpdateGameSubscriptionRec(rec)
	if err != nil {
		return err
	}
	res, err := mapper.GameSubscriptionRecordToResponse(l, rec)
	if err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteGameSubscriptionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteGameSubscriptionHandler")
	mm := m.(*domain.Domain)
	recID := pp.ByName("game_subscription_id")
	rec, err := mm.GetGameSubscriptionRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	if err := mm.DeleteGameSubscriptionRec(rec.ID); err != nil {
		return err
	}
	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
