package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

// GetManyAccountContactRecs -
func (m *Domain) GetManyAccountContactRecs(opts *coresql.Options) ([]*account_record.AccountContact, error) {
	l := m.Logger("GetManyAccountContactRecs")

	l.Debug("getting many account contact records opts >%#v<", opts)

	r := m.AccountContactRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		l.Warn("failed to get many account contact records >%v<", err)
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAccountContactRec -
func (m *Domain) GetAccountContactRec(recID string, lock *coresql.Lock) (*account_record.AccountContact, error) {
	l := m.Logger("GetAccountContactRec")

	l.Debug("getting account contact record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {

		return nil, err
	}

	r := m.AccountContactRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(account_record.TableAccountContact, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAccountContactRec -
func (m *Domain) CreateAccountContactRec(rec *account_record.AccountContact) (*account_record.AccountContact, error) {
	l := m.Logger("CreateAccountContactRec")

	l.Debug("creating account contact record >%#v<", rec)

	r := m.AccountContactRepository()

	if err := m.validateAccountContactRecForCreate(rec); err != nil {
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

// UpdateAccountContactRec -
func (m *Domain) UpdateAccountContactRec(rec *account_record.AccountContact) (*account_record.AccountContact, error) {
	l := m.Logger("UpdateAccountContactRec")

	curr, err := m.GetAccountContactRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating account contact record >%#v<", rec)

	if err := m.validateAccountContactRecForUpdate(rec, curr); err != nil {
		l.Warn("failed to validate account contact record >%v<", err)
		return rec, err
	}

	r := m.AccountContactRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		l.Warn("failed to update account contact record >%v<", err)
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAccountContactRec -
func (m *Domain) DeleteAccountContactRec(recID string) error {
	l := m.Logger("DeleteAccountContactRec")

	l.Debug("deleting account contact record ID >%s<", recID)

	rec, err := m.GetAccountContactRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get account contact record ID >%s< >%v<", recID, err)
		return err
	}

	r := m.AccountContactRepository()

	if err := m.validateAccountContactRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.DeleteOne(recID); err != nil {
		l.Warn("failed to delete account contact record ID >%s< >%v<", recID, err)
		return databaseError(err)
	}

	return nil
}

// RemoveAccountContactRec -
func (m *Domain) RemoveAccountContactRec(recID string) error {
	l := m.Logger("RemoveAccountContactRec")

	l.Debug("removing account contact record ID >%s<", recID)

	rec, err := m.GetAccountContactRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		l.Warn("failed to get account contact record ID >%s< >%v<", recID, err)
		return err
	}

	r := m.AccountContactRepository()

	if err := m.validateAccountContactRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.RemoveOne(recID); err != nil {
		l.Warn("failed to remove account contact record ID >%s< >%v<", recID, err)
		return databaseError(err)
	}

	return nil
}
