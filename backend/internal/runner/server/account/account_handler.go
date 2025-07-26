package account

import (
	"context"
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
	"gitlab.com/alienspaces/playbymail/internal/jobqueue"
	"gitlab.com/alienspaces/playbymail/internal/jobworker"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/internal/utils/logging"
	"gitlab.com/alienspaces/playbymail/schema"
)

const (
	GetManyAccounts     = "get-accounts"
	GetOneAccount       = "get-account"
	GetMyAccount        = "get-my-account"
	CreateAccount       = "create-account"
	CreateAccountWithID = "create-account-with-id"
	UpdateAccount       = "update-account"
	UpdateMyAccount     = "update-my-account"
	DeleteAccount       = "delete-account"
	DeleteMyAccount     = "delete-my-account"
)

const (
	RequestAuth = "request-auth"
	VerifyAuth  = "verify-auth"
)

func accountHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "accountHandlerConfig")

	l.Debug("adding account handler configuration")

	accountConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "account.collection.response.schema.json",
		},
		References: referenceSchemas,
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "account.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Name: "account.response.schema.json",
		},
		References: referenceSchemas,
	}

	accountConfig[GetManyAccounts] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/accounts",
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

	accountConfig[GetOneAccount] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/accounts/:account_id",
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

	accountConfig[GetMyAccount] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/accounts/me",
		HandlerFunc: getMyAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get my account",
		},
	}

	accountConfig[CreateAccount] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/accounts",
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

	accountConfig[CreateAccountWithID] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/accounts/:account_id",
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

	accountConfig[UpdateAccount] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/accounts/:account_id",
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

	accountConfig[UpdateMyAccount] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/accounts/me",
		HandlerFunc: updateMyAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update my account",
		},
	}

	accountConfig[DeleteAccount] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/accounts/:account_id",
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

	accountConfig[DeleteMyAccount] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/accounts/me",
		HandlerFunc: deleteMyAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete my account",
		},
	}

	accountConfig[RequestAuth] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/request-auth",
		HandlerFunc: requestAuthHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "account.request-auth.request.schema.json",
				},
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "account.request-auth.response.schema.json",
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
		Path:        "/v1/verify-auth",
		HandlerFunc: verifyAuthHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypePublic,
			},
			ValidateRequestSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "account.verify-auth.request.schema.json",
				},
			},
			ValidateResponseSchema: jsonschema.SchemaWithReferences{
				Main: jsonschema.Schema{
					Name: "account.verify-auth.response.schema.json",
				},
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Verify authentication",
		},
	}

	return accountConfig, nil
}

func getManyAccountsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAccountsHandler")

	l.Info("querying many account records with params >%#v<", qp)

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	recs, err := mm.GetManyAccountRecs(opts)
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
	l.Info("querying account record with id >%s<", accountID)

	mm := m.(*domain.Domain)

	rec, err := mm.GetAccountRec(accountID, nil)
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

	rec, err := mapper.AccountRequestToRecord(l, r, &record.Account{})
	if err != nil {
		l.Warn("failed mapping account request to record >%v<", err)
		return err
	}

	createdRec, err := mm.CreateAccountRec(rec)
	if err != nil {
		l.Warn("failed creating account record >%v<", err)
		return err
	}

	res, err := mapper.AccountRecordToResponse(l, createdRec)
	if err != nil {
		l.Warn("failed mapping account record to response >%v<", err)
		return err
	}

	l.Info("responding with created account record id >%s<", createdRec.ID)

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
	rec, err := mm.GetAccountRec(accountID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed getting account record >%v<", err)
		return err
	}

	rec, err = mapper.AccountRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping account request to record >%v<", err)
		return err
	}

	updatedRec, err := mm.UpdateAccountRec(rec)
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
	l.Info("deleting account record with id >%s<", accountID)

	mm := m.(*domain.Domain)

	if err := mm.DeleteAccountRec(accountID); err != nil {
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

	var req schema.RequestAuthRequest
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

	var req schema.VerifyAuthRequest
	if _, err := server.ReadRequest(l, r, &req); err != nil {
		l.Warn("failed reading request >%v<", err)
		return err
	}

	mm := m.(*domain.Domain)

	sessionToken, err := mm.VerifyAccountVerificationToken(req.VerificationToken)
	if err != nil {
		l.Warn("failed verifying account verification token >%v<", err)
		return err
	}

	l.Info("verified account verification token for email >%s<, session token >%s<", req.Email, sessionToken)

	return server.WriteResponse(l, w, http.StatusOK, mapper.MapVerifyAuthResponse(sessionToken))
}

func getMyAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getMyAccountHandler")

	l.Info("querying authenticated user account")

	mm := m.(*domain.Domain)

	// Get the authenticated account from the request context
	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		l.Warn("failed getting authenticated account data")
		return server.WriteResponse(l, w, http.StatusUnauthorized, nil)
	}

	rec, err := mm.GetAccountRec(authData.Account.ID, nil)
	if err != nil {
		l.Warn("failed getting account record >%v<", err)
		return err
	}

	res, err := mapper.AccountRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account record to response >%v<", err)
		return err
	}

	l.Info("responding with authenticated account record id >%s<", rec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func updateMyAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateMyAccountHandler")

	l.Info("updating authenticated user account")

	mm := m.(*domain.Domain)

	// Get the authenticated account from the request context
	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		l.Warn("failed getting authenticated account data")
		return server.WriteResponse(l, w, http.StatusUnauthorized, nil)
	}

	rec, err := mm.GetAccountRec(authData.Account.ID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed getting account record >%v<", err)
		return err
	}

	// Map the request to record, but preserve the email (it cannot be changed)
	originalEmail := rec.Email
	rec, err = mapper.AccountRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping account request to record >%v<", err)
		return err
	}
	// Ensure email cannot be changed
	rec.Email = originalEmail

	updatedRec, err := mm.UpdateAccountRec(rec)
	if err != nil {
		l.Warn("failed updating account record >%v<", err)
		return err
	}

	res, err := mapper.AccountRecordToResponse(l, updatedRec)
	if err != nil {
		l.Warn("failed mapping account record to response >%v<", err)
		return err
	}

	l.Info("responding with updated authenticated account record id >%s<", updatedRec.ID)

	if err = server.WriteResponse(l, w, http.StatusOK, res); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

func deleteMyAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteMyAccountHandler")

	l.Info("deleting authenticated user account")

	mm := m.(*domain.Domain)

	// Get the authenticated account from the request context
	authData := server.GetRequestAuthenData(l, r)
	if authData == nil {
		l.Warn("failed getting authenticated account data")
		return server.WriteResponse(l, w, http.StatusUnauthorized, nil)
	}

	if err := mm.DeleteAccountRec(authData.Account.ID); err != nil {
		l.Warn("failed deleting account record >%v<", err)
		return err
	}

	l.Info("deleted authenticated account record id >%s<", authData.Account.ID)

	if err := server.WriteResponse(l, w, http.StatusNoContent, nil); err != nil {
		l.Warn("failed writing response >%v<", err)
		return err
	}

	return nil
}

// SendAccountVerificationEmail generates, stores, and emails a verification token for the given email address.
func sendAccountVerificationEmail(m *domain.Domain, jc *river.Client[pgx.Tx], emailAddr string) error {
	l := m.Logger("SendAccountVerificationEmail")

	// Look up account by email
	repo := m.AccountRepository()

	recs, err := repo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: record.FieldAccountEmail, Val: emailAddr},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get account by email >%s< >%v<", emailAddr, err)
		return err
	}

	var rec *record.Account
	if len(recs) == 0 {
		l.Info("no account found for email >%s<, creating account", emailAddr)
		rec = &record.Account{
			Email: emailAddr,
		}
		rec, err = m.CreateAccountRec(rec)
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
	jc.InsertTx(context.Background(), m.Tx, &jobworker.SendAccountVerificationEmailWorkerArgs{
		AccountID: rec.ID,
	}, &river.InsertOpts{
		Queue: jobqueue.QueueDefault,
	})

	return nil
}
