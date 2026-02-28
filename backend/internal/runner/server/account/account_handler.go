package account

import (
	"context"
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
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/schema/api/account_schema"
)

const (
	GetManyAccounts     = "get-many-accounts"
	GetOneAccount       = "get-one-account"
	CreateOneAccount    = "create-one-account"
	CreateAccountWithID = "create-account-with-id"
	UpdateOneAccount    = "update-one-account"
	DeleteOneAccount    = "delete-one-account"
)

const (
	RequestAuth    = "request-auth"
	VerifyAuth     = "verify-auth"
	RefreshSession = "refresh-session"
)

func accountHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "accountHandlerConfig")

	l.Debug("adding account handler configuration")

	accountConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_schema",
			Name:     "account.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/account_schema",
				Name:     "account.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_schema",
			Name:     "account.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_schema",
			Name:     "account.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/account_schema",
				Name:     "account.schema.json",
			},
		}...),
	}

	// Register collection routes first
	accountConfig[GetManyAccounts] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/accounts",
		HandlerFunc: getManyAccountsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get account collection",
		},
	}

	accountConfig[CreateOneAccount] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/accounts",
		HandlerFunc: createAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create account",
		},
	}

	// Now register parameterized routes
	accountConfig[GetOneAccount] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/accounts/:account_id",
		HandlerFunc: getAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get account",
		},
	}

	accountConfig[CreateAccountWithID] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/accounts/:account_id",
		HandlerFunc: createAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create account with ID",
		},
	}

	accountConfig[UpdateOneAccount] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/accounts/:account_id",
		HandlerFunc: updateAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update account",
		},
	}

	accountConfig[DeleteOneAccount] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/accounts/:account_id",
		HandlerFunc: deleteAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete account",
		},
	}

	// Register auth routes
	accountConfig[RequestAuth] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/request-auth",
		HandlerFunc: requestAuthHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/account_schema",
					Name:     "account.request-auth.request.schema.json",
				},
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/account_schema",
					Name:     "account.request-auth.response.schema.json",
				},
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Request authentication",
		},
	}

	accountConfig[VerifyAuth] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/verify-auth",
		HandlerFunc: verifyAuthHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/account_schema",
					Name:     "account.verify-auth.request.schema.json",
				},
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/account_schema",
					Name:     "account.verify-auth.response.schema.json",
				},
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Verify authentication",
		},
	}

	accountConfig[RefreshSession] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/refresh-session",
		HandlerFunc: refreshSessionHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Location: "api/account_schema",
					Name:     "account.refresh-session.response.schema.json",
				},
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Refresh session",
		},
	}

	return accountConfig, nil
}

func getManyAccountsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAccountsHandler")

	l.Info("querying many account records with params >%#v<", qp)

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	// Filter by authenticated account ID if present
	authData := server.GetRequestAuthenData(l, r)
	if authData != nil {
		opts.Params = append(opts.Params, coresql.Param{
			Col: "id",
			Val: authData.AccountUser.AccountID,
		})
	}

	recs, err := mm.GetManyAccountParentRecs(opts)
	if err != nil {
		l.Warn("failed getting account records >%v<", err)
		return err
	}

	res, err := mapper.AccountRecordsToCollectionResponse(l, recs)
	if err != nil {
		return err
	}

	l.Info("responding with >%d< account records", len(res.Data))

	if err = server.WriteResponse(l, w, http.StatusOK, res, server.XPaginationHeader(len(recs), qp.PageSize)); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func getAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getAccountHandler")

	accountID := pp.ByName("account_id")
	if accountID == "" {
		l.Warn("account ID is empty")
		return coreerror.RequiredPathParameter("account_id")
	}
	l.Info("querying account record with account_id >%s<", accountID)

	mm := m.(*domain.Domain)

	rec, err := mm.GetAccountParentRec(accountID, nil)
	if err != nil {
		l.Warn("failed getting account record >%v<", err)
		return err
	}

	res, err := mapper.AccountRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account record to response >%v<", err)
		return err
	}

	l.Info("responding with account record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func createAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createAccountHandler")

	l.Info("creating account record with request >%#v<", r)

	mm := m.(*domain.Domain)

	accountRec, err := mapper.AccountRequestToRecord(l, r, &account_record.Account{})
	if err != nil {
		l.Warn("failed mapping account request to record >%v<", err)
		return err
	}

	userRec := &account_record.AccountUser{
		Email: accountRec.Name,
	}

	createdAccountRec, _, _, err := mm.CreateAccount(userRec)
	if err != nil {
		l.Warn("failed creating account record >%v<", err)
		return err
	}

	// The domain sets the account name from the email; override with the requested name
	if accountRec.Name != "" {
		createdAccountRec.Name = accountRec.Name
		createdAccountRec, err = mm.UpdateAccountParentRec(createdAccountRec)
		if err != nil {
			l.Warn("failed updating account name >%v<", err)
			return err
		}
	}

	res, err := mapper.AccountRecordToResponse(l, createdAccountRec)
	if err != nil {
		l.Warn("failed mapping account record to response >%v<", err)
		return err
	}

	l.Info("responding with created account record id >%s<", createdAccountRec.ID)

	if err = server.WriteResponse(l, w, http.StatusCreated, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateAccountHandler")

	l.Info("updating account record with path params >%#v<", pp)

	mm := m.(*domain.Domain)

	accountID := pp.ByName("account_id")
	if accountID == "" {
		l.Warn("account ID is empty")
		return coreerror.RequiredPathParameter("account_id")
	}

	rec, err := mm.GetAccountParentRec(accountID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed getting account record >%v<", err)
		return err
	}

	rec, err = mapper.AccountRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping account request to record >%v<", err)
		return err
	}

	updatedRec, err := mm.UpdateAccountParentRec(rec)
	if err != nil {
		l.Warn("failed updating account record >%v<", err)
		return err
	}

	res, err := mapper.AccountRecordToResponse(l, updatedRec)
	if err != nil {
		l.Warn("failed mapping account record to response >%v<", err)
		return err
	}

	l.Info("responding with updated account record id >%s<", updatedRec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteAccountHandler")

	accountID := pp.ByName("account_id")
	l.Info("deleting account record with account_id >%s<", accountID)

	mm := m.(*domain.Domain)

	if err := mm.RemoveAccountRec(accountID); err != nil {
		l.Warn("failed deleting account record >%v<", err)
		return err
	}

	l.Info("deleted account record id >%s<", accountID)

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// requestAuthHandler handles POST /request-auth
func requestAuthHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "requestAuthHandler")

	l.Info("requesting authentication token for email >%s<", r.Header.Get("Authorization"))

	var req account_schema.RequestAuthRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return server.WriteResponse(l, w, http.StatusOK, mapper.MapRequestAuthResponse("ok"))
	}
	if req.Email == "" {
		return server.WriteResponse(l, w, http.StatusOK, mapper.MapRequestAuthResponse("ok"))
	}

	mm := m.(*domain.Domain)

	// Within API handler context we must use the job client passed down through
	// the handler chain.
	err := sendAccountVerificationEmail(mm, jc, req.Email)
	if err != nil {
		l.Warn("failed sending account verification email >%v<", err)
		return server.WriteResponse(l, w, http.StatusOK, mapper.MapRequestAuthResponse("ok"))
	}

	return server.WriteResponse(l, w, http.StatusOK, mapper.MapRequestAuthResponse("ok"))
}

func verifyAuthHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "verifyAuthHandler")

	var req account_schema.VerifyAuthRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	// Check if test bypass authentication is enabled via header
	testBypassEnabled := false
	bypassHeaderName := mm.GetTestBypassHeaderName()
	if bypassHeaderName != "" {
		headerValue := r.Header.Get(bypassHeaderName)
		if headerValue != "" {
			testBypassEnabled = mm.IsTestBypassEnabled(headerValue)
			if testBypassEnabled {
				l.Info("test bypass authentication enabled via header >%s<", bypassHeaderName)
			}
		}
	}

	sessionToken, err := mm.VerifyAccountVerificationToken(req.VerificationToken, testBypassEnabled)
	if err != nil {
		l.Warn("failed verifying account verification token >%v<", err)
		return err
	}

	l.Info("verified account verification token for email >%s<, session token >%s<", req.Email, sessionToken)

	return server.WriteResponse(l, w, http.StatusOK, mapper.MapVerifyAuthResponse(sessionToken))
}

func refreshSessionHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "refreshSessionHandler")

	l.Info("refreshing session")

	mm := m.(*domain.Domain)

	// The session is already refreshed by the authentication middleware
	// (VerifyAccountSessionToken extends the expiry on each call).
	// This endpoint just confirms the session is valid and returns the new expiry.

	l.Info("session refreshed successfully")

	return server.WriteResponse(l, w, http.StatusOK, mapper.MapRefreshSessionResponse("ok", mm.SessionTokenExpirySeconds()))
}

// SendAccountVerificationEmail generates, stores, and emails a verification token for the given email address.
func sendAccountVerificationEmail(m *domain.Domain, jc *river.Client[pgx.Tx], emailAddr string) error {
	l := m.Logger("SendAccountVerificationEmail")

	// Look up account by email
	repo := m.AccountUserRepository()

	recs, err := repo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountUserEmail, Val: emailAddr},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get account by email >%s< >%v<", emailAddr, err)
		return err
	}

	var rec *account_record.AccountUser
	if len(recs) == 0 {
		l.Info("no account found for email >%s<, creating account", emailAddr)
		rec = &account_record.AccountUser{
			Email: emailAddr,
		}
		_, rec, _, err = m.CreateAccount(rec)
		if err != nil {
			l.Warn("failed to create account >%v<", err)
			return err
		}
	} else {
		rec = recs[0]
	}

	// Register job to send verification token
	l.Info("registering job to send verification token for account ID >%s<", rec.ID)

	// Within API handler context we must use the transaction from the domain model so
	// if there is an error the entire API request transaction is rolled back.
	if _, err := jc.InsertTx(context.Background(), m.Tx, &jobworker.SendAccountVerificationEmailWorkerArgs{
		AccountID: rec.ID,
	}, &river.InsertOpts{
		Queue: jobqueue.QueueDefault,
	}); err != nil {
		l.Warn("failed to enqueue account verification email job >%v<", err)
		return err
	}

	return nil
}
