package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyAccountRecs -
func (m *Domain) GetManyAccountRecs(opts *coresql.Options) ([]*record.Account, error) {
	l := m.Logger("GetManyAccountRecs")

	l.Debug("getting many client records opts >%#v<", opts)

	r := m.AccountRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAccountRec -
func (m *Domain) GetAccountRec(recID string, lock *coresql.Lock) (*record.Account, error) {
	l := m.Logger("GetAccountRec")

	l.Debug("getting client record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AccountRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableAccount, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAccountRec -
func (m *Domain) CreateAccountRec(rec *record.Account) (*record.Account, error) {
	l := m.Logger("CreateAccountRec")

	l.Debug("creating client record >%#v<", rec)

	r := m.AccountRepository()

	if err := m.validateAccountRecForCreate(rec); err != nil {
		l.Warn("failed to validate client record >%v<", err)
		return rec, err
	}

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateAccountRec -
func (m *Domain) UpdateAccountRec(next *record.Account) (*record.Account, error) {
	l := m.Logger("UpdateAccountRec")

	curr, err := m.GetAccountRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}

	l.Debug("updating client record >%#v<", next)

	if err := m.validateAccountRecForUpdate(next, curr); err != nil {
		l.Warn("failed to validate client record >%v<", err)
		return next, err
	}

	r := m.AccountRepository()

	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}

	return next, nil
}

// DeleteAccountRec -
func (m *Domain) DeleteAccountRec(recID string) error {
	l := m.Logger("DeleteAccountRec")

	l.Debug("deleting client record ID >%s<", recID)

	rec, err := m.GetAccountRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AccountRepository()

	if err := m.validateAccountRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAccountRec -
func (m *Domain) RemoveAccountRec(recID string) error {
	l := m.Logger("RemoveAccountRec")

	l.Debug("removing client record ID >%s<", recID)

	rec, err := m.GetAccountRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AccountRepository()

	if err := m.validateAccountRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
