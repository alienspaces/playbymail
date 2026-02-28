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
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	GetManyGameSubscriptions           = "get-many-game-subscriptions"
	GetOneGameSubscription             = "get-one-game-subscription"
	CreateOneGameSubscription          = "create-one-game-subscription"
	UpdateOneGameSubscription          = "update-one-game-subscription"
	DeleteOneGameSubscription          = "delete-one-game-subscription"
	ApproveGameSubscription            = "approve-game-subscription"
	LinkGameInstanceToSubscription     = "link-game-instance-to-subscription"
	UnlinkGameInstanceFromSubscription = "unlink-game-instance-from-subscription"
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
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
				handler_auth.PermissionGamePlaying,
				handler_auth.PermissionGameDesign,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
	}
	config[GetOneGameSubscription] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id",
		HandlerFunc: getGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
				handler_auth.PermissionGamePlaying,
				handler_auth.PermissionGameDesign,
			},
			ValidateResponseSchema: responseSchema,
		},
	}
	config[CreateOneGameSubscription] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-subscriptions",
		HandlerFunc: createGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
				handler_auth.PermissionGamePlaying,
				handler_auth.PermissionGameDesign,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
	}
	config[UpdateOneGameSubscription] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id",
		HandlerFunc: updateGameSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
				handler_auth.PermissionGamePlaying,
				handler_auth.PermissionGameDesign,
			},
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
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
				handler_auth.PermissionGamePlaying,
				handler_auth.PermissionGameDesign,
			},
		},
	}

	// Public route for approving a game subscription that originated from an invitation email
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

	instanceRequestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_subscription_instance.request.schema.json",
		},
		References: referenceSchemas,
	}

	instanceResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/game_schema",
			Name:     "game_subscription_instance.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/game_schema",
				Name:     "game_subscription_instance.schema.json",
			},
		}...),
	}

	// Link instance to subscription
	config[LinkGameInstanceToSubscription] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id/instances",
		HandlerFunc: linkGameInstanceToSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
				handler_auth.PermissionGameDesign,
			},
			ValidateRequestSchema:  instanceRequestSchema,
			ValidateResponseSchema: instanceResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Link game instance to subscription",
			Description: "Link a game instance to a game subscription. " +
				"Validates instance limit and that instance belongs to same game.",
		},
	}

	// Unlink instance from subscription
	config[UnlinkGameInstanceFromSubscription] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/game-subscriptions/:game_subscription_id/instances/:game_instance_id",
		HandlerFunc: unlinkGameInstanceFromSubscriptionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{server.AuthenticationTypeToken},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameManagement,
				handler_auth.PermissionGameDesign,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:    true,
			Title:       "Unlink game instance from subscription",
			Description: "Remove the link between a game instance and a game subscription.",
		},
	}

	return config, nil
}

func getManyGameSubscriptionsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyGameSubscriptionsHandler")

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Use view repository to get subscriptions with aggregated instance IDs
	recs, err := mm.GetManyGameSubscriptionViewRecs(opts)
	if err != nil {
		l.Warn("failed to get many game subscription view records >%v<", err)
		return err
	}

	res, err := mapper.GameSubscriptionViewRecordsToCollectionResponse(l, recs)
	if err != nil {
		l.Warn("failed to map game subscription records to collection response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize))
}

func getGameSubscriptionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getGameSubscriptionHandler")

	mm := m.(*domain.Domain)

	recID := pp.ByName("game_subscription_id")
	if recID == "" {
		l.Warn("game subscription ID is required")
		return coreerror.NewNotFoundError(game_record.TableGameSubscription, recID)
	}

	// Use view repository to get subscription with aggregated instance IDs
	rec, err := mm.GetGameSubscriptionViewRec(recID, nil)
	if err != nil {
		l.Warn("failed to get game subscription view record >%v<", err)
		return err
	}

	res, err := mapper.GameSubscriptionViewRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed to map game subscription record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createGameSubscriptionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createGameSubscriptionHandler")

	mm := m.(*domain.Domain)

	// Get authenticated account
	authenData := server.GetRequestAuthenData(l, r)
	if authenData == nil || authenData.AccountUser.ID == "" {
		l.Warn("authenticated account is required")
		return coreerror.NewUnauthorizedError()
	}

	rec, err := mapper.GameSubscriptionRequestToRecord(l, r, &game_record.GameSubscription{})
	if err != nil {
		l.Warn("failed to map game subscription request to record >%v<", err)
		return err
	}

	// Set accountID from authenticated account (self-subscription) using the tenant account.id
	rec.AccountID = authenData.AccountUser.AccountID

	// Set status to active (self-subscription, no approval required)
	rec.Status = game_record.GameSubscriptionStatusActive

	rec, err = mm.CreateGameSubscriptionRec(rec)
	if err != nil {
		l.Warn("failed to create game subscription record >%v<", err)
		return err
	}

	// New subscription has no instances yet
	res, err := mapper.GameSubscriptionRecordToResponse(l, rec, []string{})
	if err != nil {
		l.Warn("failed to map game subscription record to response >%v<", err)
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
		l.Warn("failed to get game subscription record >%v<", err)
		return err
	}

	rec, err = mapper.GameSubscriptionRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed to map game subscription request to record >%v<", err)
		return err
	}

	rec, err = mm.UpdateGameSubscriptionRec(rec)
	if err != nil {
		l.Warn("failed to update game subscription record >%v<", err)
		return err
	}

	instanceRecs, err := mm.GetGameSubscriptionInstanceRecsBySubscription(recID)
	if err != nil {
		l.Warn("failed to get game subscription instance records >%v<", err)
		return err
	}
	instanceIDs := make([]string, len(instanceRecs))
	for i, instanceRec := range instanceRecs {
		instanceIDs[i] = instanceRec.GameInstanceID
	}

	res, err := mapper.GameSubscriptionRecordToResponse(l, rec, instanceIDs)
	if err != nil {
		l.Warn("failed to map game subscription record to response >%v<", err)
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
		l.Warn("failed to get game subscription record >%v<", err)
		return err
	}

	if err := mm.DeleteGameSubscriptionRec(rec.ID); err != nil {
		l.Warn("failed to delete game subscription record >%v<", err)
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

	// Enqueue game subscription processing job to create game entities
	if _, err := jc.InsertTx(r.Context(), mm.Tx, &jobworker.GameSubscriptionProcessingWorkerArgs{
		GameSubscriptionID: rec.ID,
	}, &river.InsertOpts{Queue: jobqueue.QueueDefault}); err != nil {
		l.Warn("failed to enqueue process subscription job >%v<", err)
		return err
	}

	instanceRecs, err := mm.GetGameSubscriptionInstanceRecsBySubscription(subscriptionID)
	if err != nil {
		l.Warn("failed to get game subscription instance records >%v<", err)
		return err
	}
	instanceIDs := make([]string, len(instanceRecs))
	for i, instanceRec := range instanceRecs {
		instanceIDs[i] = instanceRec.GameInstanceID
	}

	res, err := mapper.GameSubscriptionRecordToResponse(l, rec, instanceIDs)
	if err != nil {
		l.Warn("failed to map game subscription record to response >%v<", err)
		return err
	}

	l.Info("approved game subscription ID >%s< for email >%s<", subscriptionID, email)

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func linkGameInstanceToSubscriptionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "linkGameInstanceToSubscriptionHandler")

	mm := m.(*domain.Domain)

	subscriptionID := pp.ByName("game_subscription_id")
	if subscriptionID == "" {
		l.Warn("game subscription ID is required")
		return coreerror.NewNotFoundError(game_record.TableGameSubscription, subscriptionID)
	}

	rec, err := mapper.GameSubscriptionInstanceRequestToRecord(l, r, &game_record.GameSubscriptionInstance{})
	if err != nil {
		l.Warn("failed to map game subscription instance request to record >%v<", err)
		return err
	}

	// Override subscription ID from path
	rec.GameSubscriptionID = subscriptionID

	// Create the subscription-instance link (account_id will be derived from subscription in validation)
	linkedRec, err := mm.CreateGameSubscriptionInstanceRec(rec)
	if err != nil {
		l.Warn("failed to create game subscription instance link >%v<", err)
		return err
	}

	res, err := mapper.GameSubscriptionInstanceRecordToResponse(l, linkedRec)
	if err != nil {
		l.Warn("failed to map game subscription instance record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func unlinkGameInstanceFromSubscriptionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "unlinkGameInstanceFromSubscriptionHandler")

	mm := m.(*domain.Domain)

	subscriptionID := pp.ByName("game_subscription_id")
	if subscriptionID == "" {
		l.Warn("game subscription ID is required")
		return coreerror.NewNotFoundError(game_record.TableGameSubscription, subscriptionID)
	}

	instanceID := pp.ByName("game_instance_id")
	if instanceID == "" {
		l.Warn("game instance ID is required")
		return coreerror.NewNotFoundError(game_record.TableGameInstance, instanceID)
	}

	// Find the subscription-instance link
	recs, err := mm.GetManyGameSubscriptionInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: game_record.FieldGameSubscriptionInstanceGameSubscriptionID, Val: subscriptionID},
			{Col: game_record.FieldGameSubscriptionInstanceGameInstanceID, Val: instanceID},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to find game subscription instance link >%v<", err)
		return err
	}

	if len(recs) == 0 {
		return coreerror.NewNotFoundError(game_record.TableGameSubscriptionInstance,
			"subscription_id="+subscriptionID+", instance_id="+instanceID)
	}

	// Delete the link
	if err := mm.DeleteGameSubscriptionInstanceRec(recs[0].ID); err != nil {
		l.Warn("failed to delete game subscription instance link >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
