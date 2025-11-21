package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

func (t *Testing) createAccountContactRec(accountID string) (*account_record.AccountContact, error) {
	l := t.Logger("createAccountContactRec")

	rec := &account_record.AccountContact{
		AccountID:          accountID,
		Name:               gofakeit.Name(),
		PostalAddressLine1: gofakeit.Address().Address,
		StateProvince:      gofakeit.Address().State,
		Country:            gofakeit.Address().Country,
		PostalCode:         gofakeit.Address().Zip,
	}

	l.Debug("creating account contact record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateAccountContactRec(rec)
	if err != nil {
		l.Warn("failed creating account contact record >%v<", err)
		return nil, err
	}

	// Add the account contact record to the data store
	t.Data.AddAccountContactRec(rec)

	// Add the account contact record to the teardown data store
	t.teardownData.AddAccountContactRec(rec)

	return rec, nil
}
