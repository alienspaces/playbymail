package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

// GetManyAccountSubscriptionRecs -
func (m *Domain) GetManyAccountSubscriptionRecs(opts *coresql.Options) ([]*account_record.AccountSubscription, error) {
	l := m.Logger("GetManyAccountSubscriptionRecs")

	l.Debug("getting many account_subscription records opts >%#v<", opts)

	r := m.AccountSubscriptionRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAccountSubscriptionRec -
func (m *Domain) GetAccountSubscriptionRec(recID string, lock *coresql.Lock) (*account_record.AccountSubscription, error) {
	l := m.Logger("GetAccountSubscriptionRec")

	l.Debug("getting account_subscription record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AccountSubscriptionRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(account_record.TableAccountSubscription, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAccountSubscriptionRec -
func (m *Domain) CreateAccountSubscriptionRec(rec *account_record.AccountSubscription) (*account_record.AccountSubscription, error) {
	l := m.Logger("CreateAccountSubscriptionRec")

	l.Debug("creating account_subscription record >%#v<", rec)

	r := m.AccountSubscriptionRepository()

	if err := m.validateAccountSubscriptionRecForCreate(rec); err != nil {
		l.Warn("failed to validate account_subscription record >%v<", err)
		return rec, err
	}

	createdRec, err := r.CreateOne(rec)
	if err != nil {
		l.Warn("failed to create account_subscription record >%v<", err)
		return nil, databaseError(err)
	}

	return createdRec, nil
}

// UpdateAccountSubscriptionRec -
func (m *Domain) UpdateAccountSubscriptionRec(rec *account_record.AccountSubscription) (*account_record.AccountSubscription, error) {
	l := m.Logger("UpdateAccountSubscriptionRec")

	curr, err := m.GetAccountSubscriptionRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating account_subscription record >%#v<", rec)

	if err := m.validateAccountSubscriptionRecForUpdate(curr, rec); err != nil {
		l.Warn("failed to validate account_subscription record >%v<", err)
		return rec, err
	}

	r := m.AccountSubscriptionRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAccountSubscriptionRec -
func (m *Domain) DeleteAccountSubscriptionRec(recID string) error {
	l := m.Logger("DeleteAccountSubscriptionRec")
	l.Debug("deleting account_subscription record ID >%s<", recID)
	rec, err := m.GetAccountSubscriptionRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AccountSubscriptionRepository()
	if err := m.validateAccountSubscriptionRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAccountSubscriptionRec -
func (m *Domain) RemoveAccountSubscriptionRec(recID string) error {
	l := m.Logger("RemoveAccountSubscriptionRec")
	l.Debug("removing account_subscription record ID >%s<", recID)
	rec, err := m.GetAccountSubscriptionRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AccountSubscriptionRepository()
	if err := m.validateAccountSubscriptionRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
