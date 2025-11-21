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
func (m *Domain) UpdateAccountContactRec(next *account_record.AccountContact) (*account_record.AccountContact, error) {
	l := m.Logger("UpdateAccountContactRec")

	curr, err := m.GetAccountContactRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}

	l.Debug("updating account contact record >%#v<", next)

	if err := m.validateAccountContactRecForUpdate(next, curr); err != nil {
		l.Warn("failed to validate account contact record >%v<", err)
		return next, err
	}

	r := m.AccountContactRepository()

	next, err = r.UpdateOne(next)
	if err != nil {
		l.Warn("failed to update account contact record >%v<", err)
		return next, databaseError(err)
	}

	return next, nil
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

func (m *Domain) validateAccountContactRecForCreate(rec *account_record.AccountContact) error {
	l := m.Logger("validateAccountContactRecForCreate")
	l.Debug("validating account contact record >%#v<", rec)

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField("account_id", rec.AccountID); err != nil {
		return err
	}

	if rec.Name == "" {
		return coreerror.NewInvalidDataError("name is required")
	}

	if rec.PostalAddressLine1 == "" {
		return coreerror.NewInvalidDataError("postal_address_line1 is required")
	}

	if rec.StateProvince == "" {
		return coreerror.NewInvalidDataError("state_province is required")
	}

	if rec.Country == "" {
		return coreerror.NewInvalidDataError("country is required")
	}

	if rec.PostalCode == "" {
		return coreerror.NewInvalidDataError("postal_code is required")
	}

	return nil
}

func (m *Domain) validateAccountContactRecForUpdate(next, curr *account_record.AccountContact) error {
	l := m.Logger("validateAccountContactRecForUpdate")
	l.Debug("validating current account contact record >%#v< against next >%#v<", curr, next)

	if next == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if next.AccountID != curr.AccountID {
		return coreerror.NewInvalidDataError("account_id cannot be updated")
	}

	if next.Name == "" {
		return coreerror.NewInvalidDataError("name is required")
	}

	if next.PostalAddressLine1 == "" {
		return coreerror.NewInvalidDataError("postal_address_line1 is required")
	}

	if next.StateProvince == "" {
		return coreerror.NewInvalidDataError("state_province is required")
	}

	if next.Country == "" {
		return coreerror.NewInvalidDataError("country is required")
	}

	if next.PostalCode == "" {
		return coreerror.NewInvalidDataError("postal_code is required")
	}

	return nil
}

func (m *Domain) validateAccountContactRecForDelete(rec *account_record.AccountContact) error {
	l := m.Logger("validateAccountContactRecForDelete")
	l.Debug("validating account contact record >%#v<", rec)

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	return nil
}
