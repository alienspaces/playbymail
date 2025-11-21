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
	GetManyAccountContacts  = "get-many-account-contacts"
	GetOneAccountContact    = "get-one-account-contact"
	CreateOneAccountContact = "create-one-account-contact"
	UpdateOneAccountContact = "update-one-account-contact"
	DeleteOneAccountContact = "delete-one-account-contact"
)

func accountContactHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "accountContactHandlerConfig")

	l.Debug("adding account contact handler configuration")

	accountContactConfig := make(map[string]server.HandlerConfig)

	collectionResponseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_schema",
			Name:     "account_contact.collection.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/account_schema",
				Name:     "account_contact.schema.json",
			},
		}...),
	}

	requestSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_schema",
			Name:     "account_contact.request.schema.json",
		},
		References: referenceSchemas,
	}

	responseSchema := jsonschema.SchemaWithReferences{
		Main: jsonschema.Schema{
			Location: "api/account_schema",
			Name:     "account_contact.response.schema.json",
		},
		References: append(referenceSchemas, []jsonschema.Schema{
			{
				Location: "api/account_schema",
				Name:     "account_contact.schema.json",
			},
		}...),
	}

	// Register collection routes first
	accountContactConfig[GetManyAccountContacts] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/accounts/:account_id/contacts",
		HandlerFunc: getManyAccountContactsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get account contact collection",
		},
	}

	accountContactConfig[CreateOneAccountContact] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/accounts/:account_id/contacts",
		HandlerFunc: createAccountContactHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Create account contact",
		},
	}

	// Now register parameterized routes
	accountContactConfig[GetOneAccountContact] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/accounts/:account_id/contacts/:account_contact_id",
		HandlerFunc: getAccountContactHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Get account contact",
		},
	}

	accountContactConfig[UpdateOneAccountContact] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/accounts/:account_id/contacts/:account_contact_id",
		HandlerFunc: updateAccountContactHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateRequestSchema:  requestSchema,
			ValidateResponseSchema: responseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Update account contact",
		},
	}

	accountContactConfig[DeleteOneAccountContact] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/accounts/:account_id/contacts/:account_contact_id",
		HandlerFunc: deleteAccountContactHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
		},
		DocumentationConfig: server.DocumentationConfig{
			Document: true,
			Title:    "Delete account contact",
		},
	}

	return accountContactConfig, nil
}

func getManyAccountContactsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAccountContactsHandler")

	accountID := pp.ByName("account_id")
	if accountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)
	opts.Params = append(opts.Params, coresql.Param{
		Col: account_record.FieldAccountContactAccountID,
		Val: accountID,
	})
	// Override default ordering to use created_at descending
	opts.OrderBy = []coresql.OrderBy{
		{Col: account_record.FieldAccountContactCreatedAt, Direction: coresql.OrderDirectionDESC},
	}

	recs, err := mm.GetManyAccountContactRecs(opts)
	if err != nil {
		l.Warn("failed to get account contact records >%v<", err)
		return err
	}

	res, err := mapper.AccountContactRecordsToCollectionResponse(l, recs)
	if err != nil {
		l.Warn("failed mapping account contact records to collection response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func getAccountContactHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getAccountContactHandler")

	accountContactID := pp.ByName("account_contact_id")
	if accountContactID == "" {
		return coreerror.NewInvalidDataError("account_contact_id is required")
	}

	accountID := pp.ByName("account_id")
	if accountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAccountContactRec(accountContactID, nil)
	if err != nil {
		l.Warn("failed to get account contact record >%v<", err)
		return err
	}

	// Verify the account_contact belongs to the specified account
	if rec.AccountID != accountID {
		return coreerror.NewNotFoundError(account_record.TableAccountContact, accountContactID)
	}

	res, err := mapper.AccountContactRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account contact record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createAccountContactHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createAccountContactHandler")

	accountID := pp.ByName("account_id")
	if accountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	mm := m.(*domain.Domain)

	rec := &account_record.AccountContact{
		AccountID: accountID,
	}

	rec, err := mapper.AccountContactRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping account contact request to record >%v<", err)
		return err
	}

	rec, err = mm.CreateAccountContactRec(rec)
	if err != nil {
		l.Warn("failed to create account contact record >%v<", err)
		return err
	}

	res, err := mapper.AccountContactRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account contact record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateAccountContactHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateAccountContactHandler")

	accountContactID := pp.ByName("account_contact_id")
	if accountContactID == "" {
		return coreerror.NewInvalidDataError("account_contact_id is required")
	}

	accountID := pp.ByName("account_id")
	if accountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	mm := m.(*domain.Domain)

	rec := &account_record.AccountContact{
		AccountID: accountID,
	}
	rec.ID = accountContactID

	rec, err := mapper.AccountContactRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping account contact request to record >%v<", err)
		return err
	}

	rec, err = mm.UpdateAccountContactRec(rec)
	if err != nil {
		l.Warn("failed to update account contact record >%v<", err)
		return err
	}

	res, err := mapper.AccountContactRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account contact record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteAccountContactHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteAccountContactHandler")

	accountContactID := pp.ByName("account_contact_id")
	if accountContactID == "" {
		return coreerror.NewInvalidDataError("account_contact_id is required")
	}

	accountID := pp.ByName("account_id")
	if accountID == "" {
		return coreerror.NewInvalidDataError("account_id is required")
	}

	mm := m.(*domain.Domain)

	// Verify the account_contact belongs to the specified account
	rec, err := mm.GetAccountContactRec(accountContactID, nil)
	if err != nil {
		l.Warn("failed to get account contact record >%v<", err)
		return err
	}

	if rec.AccountID != accountID {
		return coreerror.NewNotFoundError(account_record.TableAccountContact, accountContactID)
	}

	err = mm.DeleteAccountContactRec(accountContactID)
	if err != nil {
		l.Warn("failed to delete account contact record >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
