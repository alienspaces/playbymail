package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

// GetManyAccountUserRecs -
func (m *Domain) GetManyAccountUserRecs(opts *coresql.Options) ([]*account_record.AccountUser, error) {
	l := m.Logger("GetManyAccountUserRecs")

	l.Info("getting many account user records opts >%#v<", opts)

	r := m.AccountUserRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAccountUserRec -
func (m *Domain) GetAccountUserRec(recID string, lock *coresql.Lock) (*account_record.AccountUser, error) {
	l := m.Logger("GetAccountUserRec")

	l.Debug("getting account user record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AccountUserRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(account_record.TableAccountUser, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAccountUserRec creates an account user record
func (m *Domain) CreateAccountUserRec(rec *account_record.AccountUser) (*account_record.AccountUser, error) {
	l := m.Logger("CreateAccountUserRec")

	l.Debug("creating account user record >%#v<", rec)

	if err := m.validateAccountUserRecForCreate(rec); err != nil {
		l.Warn("failed to validate account user record >%v<", err)
		return rec, err
	}

	r := m.AccountUserRepository()

	rec, err := r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateAccountUserRec -
func (m *Domain) UpdateAccountUserRec(rec *account_record.AccountUser) (*account_record.AccountUser, error) {
	l := m.Logger("UpdateAccountUserRec")

	curr, err := m.GetAccountUserRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating account user record >%#v<", rec)

	if err := m.validateAccountUserRecForUpdate(curr, rec); err != nil {
		l.Warn("failed to validate account user record >%v<", err)
		return rec, err
	}

	r := m.AccountUserRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAccountUserRec -
func (m *Domain) DeleteAccountUserRec(recID string) error {
	l := m.Logger("DeleteAccountUserRec")

	l.Debug("deleting account user record ID >%s<", recID)

	rec, err := m.GetAccountUserRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	if err := m.validateAccountUserRecForDelete(rec); err != nil {
		l.Warn("failed to validate account user record >%v<", err)
		return err
	}

	r := m.AccountUserRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
