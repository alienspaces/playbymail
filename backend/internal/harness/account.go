package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

func (t *Testing) createAccountRec(accountConfig AccountConfig) (*account_record.Account, error) {
	l := t.Logger("createAccountRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *account_record.Account
	if accountConfig.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *accountConfig.Record
		rec = &recCopy
	} else {
		rec = &account_record.Account{}
	}

	rec = t.applyAccountRecDefaultValues(rec)

	l.Debug("creating account record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateAccountRec(rec)
	if err != nil {
		l.Warn("failed creating account record >%v<", err)
		return nil, err
	}

	// Add the account record to the data store
	t.Data.AddAccountRec(rec)

	// Add the account record to the teardown data store
	t.teardownData.AddAccountRec(rec)

	// Add the account record to the data store refs
	if accountConfig.Reference != "" {
		t.Data.Refs.AccountRefs[accountConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyAccountRecDefaultValues(rec *account_record.Account) *account_record.Account {
	if rec == nil {
		rec = &account_record.Account{}
	}

	if rec.Email == "" {
		rec.Email = UniqueEmail(gofakeit.Email())
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}

	return rec
}
