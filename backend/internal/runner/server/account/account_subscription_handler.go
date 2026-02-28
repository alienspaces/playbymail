package account

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
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/runner/server/handler_auth"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	GetMyAccountSubscriptions = "get-my-account-subscriptions"
)

func accountSubscriptionHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "accountSubscriptionHandlerConfig")

	l.Debug("adding account subscription handler configuration")

	accountSubscriptionConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_subscription_schema",
			Name:     "account_subscription.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/account_subscription_schema",
				Name:     "account_subscription.schema.json",
			},
		}...),
	}

	// Register "my account subscriptions" route
	accountSubscriptionConfig[GetMyAccountSubscriptions] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/account/subscriptions",
		HandlerFunc: getMyAccountSubscriptionsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			AuthzPermissions: []server.AuthorizedPermission{
				handler_auth.PermissionGameDesign,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get my account subscriptions",
		},
	}

	return accountSubscriptionConfig, nil
}

func getMyAccountSubscriptionsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getMyAccountSubscriptionsHandler")

	l.Info("getting authenticated user account subscriptions")

	mm := m.(*domain.Domain)

	// Get the authenticated account from the request context
	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		l.Warn("failed getting authenticated account data")
		return server.WriteResponse(l, w, http.StatusUnauthorized, nil)
	}

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, coresql.Param{
		Col: account_record.FieldAccountSubscriptionAccountID,
		Val: authData.AccountUser.AccountID,
	})
	// Order by created_at descending
	opts.OrderBy = []coresql.OrderBy{
		{Col: account_record.FieldAccountSubscriptionCreatedAt, Direction: coresql.OrderDirectionDESC},
	}

	recs, err := mm.GetManyAccountSubscriptionRecs(opts)
	if err != nil {
		l.Warn("failed to get account subscription records >%v<", err)
		return err
	}

	res, err := mapper.AccountSubscriptionRecordsToCollectionResponse(l, recs)
	if err != nil {
		l.Warn("failed mapping account subscription records to collection response >%v<", err)
		return err
	}

	l.Info("responding with account subscription records count >%d<", len(recs))

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}
