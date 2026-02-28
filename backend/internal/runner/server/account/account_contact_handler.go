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
	GetManyAccountUserContacts  = "get-many-account-contacts"
	GetOneAccountUserContact    = "get-one-account-contact"
	CreateOneAccountUserContact = "create-one-account-contact"
	UpdateOneAccountUserContact = "update-one-account-contact"
	DeleteOneAccountUserContact = "delete-one-account-contact"
)

func accountUserContactHandlerConfig(l logger.Logger) (map[string]server.HandlerConfig, error) {
	l = logging.LoggerWithFunctionContext(l, packageName, "accountUserContactHandlerConfig")

	l.Debug("adding account contact handler configuration")

	accountUserContactConfig := make(map[string]server.HandlerConfig)

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
	accountUserContactConfig[GetManyAccountUserContacts] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/account-user-contacts",
		HandlerFunc: getManyAccountUserContactsHandler,
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

	accountUserContactConfig["get-many-account-user-contacts-by-user"] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/accounts/:account_id/users/:account_user_id/contacts",
		HandlerFunc: getManyAccountUserContactsHandler,
		MiddlewareConfig: server.MiddlewareConfig{
			AuthenTypes: []server.AuthenticationType{
				server.AuthenticationTypeToken,
			},
			ValidateResponseSchema: collectionResponseSchema,
		},
		DocumentationConfig: server.DocumentationConfig{
			Document:   true,
			Collection: true,
			Title:      "Get account contact collection by user",
		},
	}

	accountUserContactConfig[CreateOneAccountUserContact] = server.HandlerConfig{
		Method:      http.MethodPost,
		Path:        "/api/v1/accounts/:account_id/users/:account_user_id/contacts",
		HandlerFunc: createAccountUserContactHandler,
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
	accountUserContactConfig[GetOneAccountUserContact] = server.HandlerConfig{
		Method:      http.MethodGet,
		Path:        "/api/v1/accounts/:account_id/users/:account_user_id/contacts/:account_user_contact_id",
		HandlerFunc: getAccountUserContactHandler,
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

	accountUserContactConfig[UpdateOneAccountUserContact] = server.HandlerConfig{
		Method:      http.MethodPut,
		Path:        "/api/v1/accounts/:account_id/users/:account_user_id/contacts/:account_user_contact_id",
		HandlerFunc: updateAccountUserContactHandler,
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

	accountUserContactConfig[DeleteOneAccountUserContact] = server.HandlerConfig{
		Method:      http.MethodDelete,
		Path:        "/api/v1/accounts/:account_id/users/:account_user_id/contacts/:account_user_contact_id",
		HandlerFunc: deleteAccountUserContactHandler,
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

	return accountUserContactConfig, nil
}

func getManyAccountUserContactsHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getManyAccountUserContactsHandler")

	mm := m.(*domain.Domain)

	opts := queryparam.ToSQLOptionsWithDefaults(qp)

	accountUserID := pp.ByName("account_user_id")
	if accountUserID != "" {
		opts.Params = append(opts.Params, coresql.Param{
			Col: account_record.FieldAccountUserContactAccountUserID,
			Val: accountUserID,
		})
	}

	// Override default ordering to use created_at descending
	opts.OrderBy = []coresql.OrderBy{
		{Col: account_record.FieldAccountUserContactCreatedAt, Direction: coresql.OrderDirectionDESC},
	}

	recs, err := mm.GetManyAccountUserContactRecs(opts)
	if err != nil {
		l.Warn("failed to get account contact records >%v<", err)
		return err
	}

	res, err := mapper.AccountUserContactRecordsToCollectionResponse(l, recs)
	if err != nil {
		l.Warn("failed mapping account contact records to collection response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func getAccountUserContactHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "getAccountUserContactHandler")

	accountUserContactID := pp.ByName("account_user_contact_id")
	if accountUserContactID == "" {
		return coreerror.NewInvalidDataError("account_user_contact_id is required")
	}

	accountUserID := pp.ByName("account_user_id")
	if accountUserID == "" {
		return coreerror.NewInvalidDataError("account_user_id is required")
	}

	mm := m.(*domain.Domain)

	rec, err := mm.GetAccountUserContactRec(accountUserContactID, nil)
	if err != nil {
		l.Warn("failed to get account contact record >%v<", err)
		return err
	}

	// Verify the account_contact belongs to the specified account user
	if rec.AccountUserID != accountUserID {
		return coreerror.NewNotFoundError(account_record.TableAccountUserContact, accountUserContactID)
	}

	res, err := mapper.AccountUserContactRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account contact record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func createAccountUserContactHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "createAccountUserContactHandler")

	accountUserID := pp.ByName("account_user_id")
	if accountUserID == "" {
		return coreerror.NewInvalidDataError("account_user_id is required")
	}

	mm := m.(*domain.Domain)

	rec := &account_record.AccountUserContact{
		AccountUserID: accountUserID,
	}

	rec, err := mapper.AccountUserContactRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping account contact request to record >%v<", err)
		return err
	}

	rec, err = mm.CreateAccountUserContactRec(rec)
	if err != nil {
		l.Warn("failed to create account contact record >%v<", err)
		return err
	}

	res, err := mapper.AccountUserContactRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account contact record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusCreated, res)
}

func updateAccountUserContactHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "updateAccountUserContactHandler")

	accountUserContactID := pp.ByName("account_user_contact_id")
	if accountUserContactID == "" {
		return coreerror.NewInvalidDataError("account_user_contact_id is required")
	}

	accountUserID := pp.ByName("account_user_id")
	if accountUserID == "" {
		return coreerror.NewInvalidDataError("account_user_id is required")
	}

	mm := m.(*domain.Domain)

	rec := &account_record.AccountUserContact{
		AccountUserID: accountUserID,
	}
	rec.ID = accountUserContactID

	rec, err := mapper.AccountUserContactRequestToRecord(l, r, rec)
	if err != nil {
		l.Warn("failed mapping account contact request to record >%v<", err)
		return err
	}

	rec, err = mm.UpdateAccountUserContactRec(rec)
	if err != nil {
		l.Warn("failed to update account contact record >%v<", err)
		return err
	}

	res, err := mapper.AccountUserContactRecordToResponse(l, rec)
	if err != nil {
		l.Warn("failed mapping account contact record to response >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusOK, res)
}

func deleteAccountUserContactHandler(w http.ResponseWriter, r *http.Request, pp httprouter.Params, qp *queryparam.QueryParams, l logger.Logger, m domainer.Domainer, jc *river.Client[pgx.Tx]) error {
	l = logging.LoggerWithFunctionContext(l, packageName, "deleteAccountUserContactHandler")

	accountUserContactID := pp.ByName("account_user_contact_id")
	if accountUserContactID == "" {
		return coreerror.NewInvalidDataError("account_user_contact_id is required")
	}

	accountUserID := pp.ByName("account_user_id")
	if accountUserID == "" {
		return coreerror.NewInvalidDataError("account_user_id is required")
	}

	mm := m.(*domain.Domain)

	// Verify the account_contact belongs to the specified account
	rec, err := mm.GetAccountUserContactRec(accountUserContactID, nil)
	if err != nil {
		l.Warn("failed to get account contact record >%v<", err)
		return err
	}

	if rec.AccountUserID != accountUserID {
		return coreerror.NewNotFoundError(account_record.TableAccountUserContact, accountUserContactID)
	}

	err = mm.DeleteAccountUserContactRec(accountUserContactID)
	if err != nil {
		l.Warn("failed to delete account contact record >%v<", err)
		return err
	}

	return server.WriteResponse(l, w, http.StatusNoContent, nil)
}
