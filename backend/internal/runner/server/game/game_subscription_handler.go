package game

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/riverqueue/river"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
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
)

const (
	GetManyGameSubscriptions  = "get-many-game-subscriptions"
	GetOneGameSubscription    = "get-one-game-subscription"
	CreateOneGameSubscription = "create-one-game-subscription"
	UpdateOneGameSubscription = "update-one-game-subscription"
	DeleteOneGameSubscription = "delete-one-game-subscription"
	ApproveGameSubscription   = "approve-game-subscription"
)

func gameSubscriptionHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "gameSubscriptionHandlerConfig")

	l.Debug("Adding game subscription handler configuration")

	config := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_subscription.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_subscription.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_subscription.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_subscription.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_subscription.schema.json",
			},
		}...),
	}

	config[GetManyGameSubscriptions] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-subscriptions",
		HandlerFunc: getManyGameSubscriptionsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: collectionResponseSchema,
		},
	}
	config[GetOneGameSubscription] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id",
		HandlerFunc: getGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateResponseSchema: responseSchema,
		},
	}
	config[CreateOneGameSubscription] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-subscriptions",
		HandlerFunc: createGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
	}
	config[UpdateOneGameSubscription] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id",
		HandlerFunc: updateGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypeToken},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
	}
	config[DeleteOneGameSubscription] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id",
		HandlerFunc: deleteGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
		},
	}
	config[ApproveGameSubscription] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id/approve",
		HandlerFunc: approveGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes:            []server.AuthenticationType{server.AuthenticationTypePublic},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Approve game subscription",
			Description: "Approve a pending game subscription by verifying the email matches " +
				"the subscription's account and updating the status to active. " +
				"Requires email query parameter.",
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
	rec, err := mapper.GameSubscriptionRequestToRecord(l, r, &game_record.GameSubscription{})
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

func approveGameSubscriptionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "approveGameSubscriptionHandler")

	subscriptionID := pp.ByName("game_subscription_id")

	emailParams, exists := qp.Params["email"]
	if !exists || len(emailParams) == 0 {
		l.Warn("email query parameter is required")
		return coreerror.NewInvalidDataError("email query parameter is required")
	}

	email := emailParams[0].Val.(string)
	if email == "" {
		l.Warn("email query parameter is empty")
		return coreerror.NewInvalidDataError("email query parameter is required")
	}

	mm := m.(*domain.Domain)

	rec, err := mm.ApproveGameSubscription(subscriptionID, email)
	if err != nil {
		l.Warn("failed to approve game subscription >%v<", err)
		return err
	}

	res, err := mapper.GameSubscriptionRecordToResponse(l, rec)
	if err != nil {
		return err
	}

	l.Info("approved game subscription ID >%s< for email >%s<", subscriptionID, email)

	return server.WriteResponse(l, w, http.StatusOK, res)
}
