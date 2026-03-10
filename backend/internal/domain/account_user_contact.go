package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

// GetManyAccountUserContactRecs -
func (m *Domain) GetManyAccountUserContactRecs(opts *coresql.Options) ([]*account_record.AccountUserContact, error) {
	l := m.Logger("GetManyAccountUserContactRecs")

	l.Debug("getting many account contact records opts >%#v<", opts)

	r := m.AccountUserContactRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		l.Warn("failed to get many account contact records >%v<", err)
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAccountUserContactRec -
func (m *Domain) GetAccountUserContactRec(recID string, lock *coresql.Lock) (*account_record.AccountUserContact, error) {
	l := m.Logger("GetAccountUserContactRec")

	l.Debug("getting account contact record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {

		return nil, err
	}

	r := m.AccountUserContactRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(account_record.TableAccountUserContact, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// GetAccountUserContactRecByAccountUserID -
func (m *Domain) GetAccountUserContactRecByAccountUserID(accountUserID string, lock *coresql.Lock) (*account_record.AccountUserContact, error) {
	l := m.Logger("GetAccountUserContactRecByAccountUserID")

	l.Debug("getting account user contact by account user ID >%s<", accountUserID)

	recs, err := m.GetManyAccountUserContactRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountUserContactAccountUserID, Val: accountUserID},
		},
		Lock:  lock,
		Limit: 1,
	})
	if err != nil {
		return nil, databaseError(err)
	}

	if len(recs) == 0 {
		return nil, coreerror.NewNotFoundError(account_record.TableAccountUserContact, accountUserID)
	}

	return recs[0], nil
}

// CreateAccountUserContactRec -
func (m *Domain) CreateAccountUserContactRec(rec *account_record.AccountUserContact) (*account_record.AccountUserContact, error) {
	l := m.Logger("CreateAccountUserContactRec")

	l.Debug("creating account contact record >%#v<", rec)

	r := m.AccountUserContactRepository()

	if err := m.validateAccountUserContactRecForCreate(rec); err != nil {
		l.Warn("failed to validate account contact record >%v<", err)
		return rec, err
	}

	rec, err := r.CreateOne(rec)
	if err != nil {
		l.Warn("failed to create account contact record >%v<", err)
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateAccountUserContactRec -
func (m *Domain) UpdateAccountUserContactRec(rec *account_record.AccountUserContact) (*account_record.AccountUserContact, error) {
	l := m.Logger("UpdateAccountUserContactRec")

	curr, err := m.GetAccountUserContactRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating account contact record >%#v<", rec)

	if err := m.validateAccountUserContactRecForUpdate(curr, rec); err != nil {
		l.Warn("failed to validate account contact record >%v<", err)
		return rec, err
	}

	r := m.AccountUserContactRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		l.Warn("failed to update account contact record >%v<", err)
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// UpsertAccountUserContactRec returns the account user contact for the given account user, updating it if it exists or creating it otherwise.
func (m *Domain) UpsertAccountUserContactRec(
	accountUserContactRec *account_record.AccountUserContact,
) (*account_record.AccountUserContact, error) {
	l := m.Logger("UpsertAccountUserContactRec")

	l.Debug("upserting account user contact for account user ID >%s<", accountUserContactRec.AccountUserID)

	existingAccountUserContactRec, err := m.GetAccountUserContactRecByAccountUserID(accountUserContactRec.AccountUserID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get account user contact by account user ID >%v<", err)
		return nil, err
	}

	// Update the existing account user contact record if it exists
	if existingAccountUserContactRec != nil {
		existingAccountUserContactRec.Name = accountUserContactRec.Name
		existingAccountUserContactRec.PostalAddressLine1 = accountUserContactRec.PostalAddressLine1
		existingAccountUserContactRec.PostalAddressLine2 = accountUserContactRec.PostalAddressLine2
		existingAccountUserContactRec.StateProvince = accountUserContactRec.StateProvince
		existingAccountUserContactRec.Country = accountUserContactRec.Country
		existingAccountUserContactRec.PostalCode = accountUserContactRec.PostalCode

		updated, err := m.UpdateAccountUserContactRec(existingAccountUserContactRec)
		if err != nil {
			l.Warn("failed to update account user contact >%v<", err)
			return nil, err
		}

		l.Info("updated account user contact >%s< for account user >%s<", updated.ID, accountUserContactRec.AccountUserID)

		return updated, nil
	}

	createdAccountUserContactRec, err := m.CreateAccountUserContactRec(accountUserContactRec)
	if err != nil {
		l.Warn("failed to create account user contact >%v<", err)
		return nil, err
	}
	l.Info("created account user contact >%s< for account user >%s<", createdAccountUserContactRec.ID, accountUserContactRec.AccountUserID)

	return createdAccountUserContactRec, nil
}

// DeleteAccountUserContactRec -
func (m *Domain) DeleteAccountUserContactRec(recID string) error {
	l := m.Logger("DeleteAccountUserContactRec")

	l.Debug("deleting account contact record ID >%s<", recID)

	rec, err := m.GetAccountUserContactRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get account contact record ID >%s< >%v<", recID, err)
		return err
	}

	r := m.AccountUserContactRepository()

	if err := m.validateAccountUserContactRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.DeleteOne(recID); err != nil {
		l.Warn("failed to delete account contact record ID >%s< >%v<", recID, err)
		return databaseError(err)
	}

	return nil
}

// RemoveAccountUserContactRec -
func (m *Domain) RemoveAccountUserContactRec(recID string) error {
	l := m.Logger("RemoveAccountUserContactRec")

	l.Debug("removing account contact record ID >%s<", recID)

	r := m.AccountUserContactRepository()

	if err := r.RemoveOne(recID); err != nil {
		l.Warn("failed to remove account contact record ID >%s< >%v<", recID, err)
		return databaseError(err)
	}

	return nil
}
