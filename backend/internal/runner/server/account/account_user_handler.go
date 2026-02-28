package account

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
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
)

const (
	GetMe                = "get-me"
	GetManyAccountUsers  = "get-many-account-users"
	GetOneAccountUser    = "get-one-account-user"
	CreateOneAccountUser = "create-one-account-user"
	UpdateOneAccountUser = "update-one-account-user"
	DeleteOneAccountUser = "delete-one-account-user"
)

func accountUserHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "accountUserHandlerConfig")

	l.Debug("adding account user handler configuration")

	accountUserConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_schema",
			Name:     "account_user.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/account_schema",
				Name:     "account_user.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_schema",
			Name:     "account_user.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_schema",
			Name:     "account_user.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/account_schema",
				Name:     "account_user.schema.json",
			},
		}...),
	}

	// Lightweight "me" endpoint for the authenticated user
	accountUserConfig[GetMe] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/me",
		HandlerFunc: getMeHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get authenticated account user",
		},
	}

	// Register collection routes first
	// Global collection of account users
	accountUserConfig[GetManyAccountUsers] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/account-users",
		HandlerFunc: getManyAccountUsersHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get account user collection",
		},
	}

	// Collection under account
	accountUserConfig["get-many-account-users-by-account"] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/accounts/:account_id/users",
		HandlerFunc: getManyAccountUsersHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get account user collection by account",
		},
	}

	accountUserConfig[CreateOneAccountUser] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/accounts/:account_id/users",
		HandlerFunc: createAccountUserHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create account user",
		},
	}

	// Now register parameterized routes
	accountUserConfig[GetOneAccountUser] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/accounts/:account_id/users/:account_user_id",
		HandlerFunc: getAccountUserHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get account user",
		},
	}

	accountUserConfig[UpdateOneAccountUser] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/accounts/:account_id/users/:account_user_id",
		HandlerFunc: updateAccountUserHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update account user",
		},
	}

	accountUserConfig[DeleteOneAccountUser] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/accounts/:account_id/users/:account_user_id",
		HandlerFunc: deleteAccountUserHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete account user",
		},
	}

	return accountUserConfig, nil
}

func getMeHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getMeHandler")

	l.Info("querying authenticated account user")

	mm := m.(*domain.Domain)

	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		l.Warn("failed getting authenticated account data")
		return server.WriteResponse(l, w, http.StatusUnauthorized, nil)
	}

	rec, err := mm.GetAccountUserRec(authData.AccountUser.ID, nil)
	if err != nil {
		l.Warn("failed getting account user record >%v<", err)
		return err
	}

	res, err := mapper.AccountUserRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account user record to response >%v<", err)
		return err
	}

	l.Info("responding with authenticated account user record id >%s<", rec.ID)

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func getManyAccountUsersHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAccountUsersHandler")

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	accountID := pp.ByName("account_id")
	if accountID != "" {
		opts.Params = append(opts.Params, coresql.Param{
			Col: account_record.FieldAccountUserAccountID,
			Val: accountID,
		})
	}

	// Default ordering
	if len(opts.OrderBy) == 0 {
		opts.OrderBy = []coresql.OrderBy{
			{Col: account_record.FieldAccountUserCreatedAt, Direction: coresql.OrderDirectionDESC},
		}
	}

	recs, err := mm.GetManyAccountUserRecs(opts)
	if err != nil {
		l.Warn("failed to get account user records >%v<", err)
		return err
	}

	res, err := mapper.AccountUserRecordsToCollectionResponse(l, recs)
	if err != nil {
		l.Warn("failed mapping account user records to collection response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func getAccountUserHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getAccountUserHandler")

	accountUserID := pp.ByName("account_user_id")
	if accountUserID == "" {
		return coreerror.NewInvalidDataError("account_user_id is required")
	}

	accountID := pp.ByName("account_id")
	if accountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAccountUserRec(accountUserID, nil)
	if err != nil {
		l.Warn("failed to get account user record >%v<", err)
		return err
	}

	// Validate hierarchy
	if rec.AccountID != accountID {
		return coreerror.NewNotFoundError(account_record.TableAccountUser, accountUserID)
	}

	res, err := mapper.AccountUserRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account user record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createAccountUserHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createAccountUserHandler")

	accountID := pp.ByName("account_id")
	if accountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	mm := m.(*domain.Domain)

	rec := &account_record.AccountUser{
		AccountID: accountID,
	}

	rec, err := mapper.AccountUserRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping account user request to record >%v<", err)
		return err
	}

	rec, err = mm.CreateAccountUserRec(rec)
	if err != nil {
		l.Warn("failed to create account user record >%v<", err)
		return err
	}

	res, err := mapper.AccountUserRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account user record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateAccountUserHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateAccountUserHandler")

	accountUserID := pp.ByName("account_user_id")
	if accountUserID == "" {
		return coreerror.NewInvalidDataError("account_user_id is required")
	}

	accountID := pp.ByName("account_id")
	if accountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAccountUserRec(accountUserID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed getting account user record >%v<", err)
		return err
	}

	// Validate hierarchy
	if rec.AccountID != accountID {
		return coreerror.NewNotFoundError(account_record.TableAccountUser, accountUserID)
	}

	rec, err = mapper.AccountUserRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping account user request to record >%v<", err)
		return err
	}

	rec, err = mm.UpdateAccountUserRec(rec)
	if err != nil {
		l.Warn("failed to update account user record >%v<", err)
		return err
	}

	res, err := mapper.AccountUserRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account user record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteAccountUserHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteAccountUserHandler")

	accountUserID := pp.ByName("account_user_id")
	if accountUserID == "" {
		return coreerror.NewInvalidDataError("account_user_id is required")
	}

	accountID := pp.ByName("account_id")
	if accountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	mm := m.(*domain.Domain)

	// Verify the account_user belongs to the specified account
	rec, err := mm.GetAccountUserRec(accountUserID, nil)
	if err != nil {
		l.Warn("failed to get account user record >%v<", err)
		return err
	}

	if rec.AccountID != accountID {
		return coreerror.NewNotFoundError(account_record.TableAccountUser, accountUserID)
	}

	err = mm.DeleteAccountUserRec(accountUserID)
	if err != nil {
		l.Warn("failed to delete account user record >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
