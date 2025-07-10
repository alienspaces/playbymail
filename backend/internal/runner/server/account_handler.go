package runner

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"gitlab.com/alienspaces/playbymail/core/jsonschema"
	"gitlab.com/alienspaces/playbymail/core/queryparam"
	"gitlab.com/alienspaces/playbymail/core/server"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/domainer"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/mapper"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

const (
	tagGroupAccount server.TagGroup = "Accounts"
	TagAccount      server.Tag      = "Accounts"
)

const (
	getManyAccounts     = "get-accounts"
	getOneAccount       = "get-account"
	createAccount       = "create-account"
	createAccountWithID = "create-account-with-id"
	updateAccount       = "update-account"
	deleteAccount       = "delete-account"
)

func (rnr *Runner) accountHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = loggerWithFunctionContext(l, "accountHandlerConfig")

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

	accountConfig[getManyAccounts] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/accounts",
		HandlerFunc: rnr.getManyAccountsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get account collection",
		},
	}

	accountConfig[getOneAccount] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/v1/accounts/:account_id",
		HandlerFunc: rnr.getAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get account",
		},
	}

	accountConfig[createAccount] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/accounts",
		HandlerFunc: rnr.createAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create account",
		},
	}

	accountConfig[createAccountWithID] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/v1/accounts/:account_id",
		HandlerFunc: rnr.createAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create account with ID",
		},
	}

	accountConfig[updateAccount] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/v1/accounts/:account_id",
		HandlerFunc: rnr.updateAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update account",
		},
	}

	accountConfig[deleteAccount] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/v1/accounts/:account_id",
		HandlerFunc: rnr.deleteAccountHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeAPIKey,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete account",
		},
	}

	return accountConfig, nil
}

func (rnr *Runner) getManyAccountsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetManyAccountsHandler")

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

func (rnr *Runner) getAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "GetAccountHandler")

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

func (rnr *Runner) createAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	if w == nil {
		panic("DEBUG: http.ResponseWriter w is nil in createAccountHandler")
	}
	l = loggerWithFunctionContext(l, "CreateAccountHandler")

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

func (rnr *Runner) updateAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "UpdateAccountHandler")

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

func (rnr *Runner) deleteAccountHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer) error {
	l = loggerWithFunctionContext(l, "DeleteAccountHandler")

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
