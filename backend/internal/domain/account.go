package domain

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	corerecord "gitlab.com/alienspaces/playbymail/core/record"
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

func (m *Domain) GenerateAccountVerificationToken(rec *record.Account) (string, error) {
	l := m.Logger("GenerateAccountVerificationToken")

	l.Debug("generating verification token for account ID >%s<", rec.ID)

	// Generate a new UUID for the token
	token := corerecord.NewRecordID()

	// Hash the token
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		l.Warn("failed to hash verification token >%v<", err)
		return "", err
	}

	rec.VerificationToken = nullstring.FromString(string(hashedToken))
	rec.VerificationTokenExpiresAt = nulltime.FromTime(corerecord.NewRecordTimestamp().Add(15 * time.Minute))

	_, err = m.UpdateAccountRec(rec)
	if err != nil {
		l.Warn("failed to update account >%v<", err)
		return "", err
	}

	// TODO: Debug log the token
	l.Info("generated verification token >%s< for account ID >%s<", token, rec.ID)

	return token, nil
}

func (m *Domain) VerifyAccountVerificationToken(token string) (string, error) {
	l := m.Logger("VerifyAccountVerificationToken")

	// TODO: Debug log the token
	l.Info("verifying verification token >%s<", token)

	// Look up account by email
	repo := m.AccountRepository()

	recs, err := repo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: record.FieldAccountVerificationToken, Val: token},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get account by verification token >%s< >%v<", token, err)
		return "", err
	}

	if len(recs) == 0 {
		l.Info("no account found for verification token >%s<", token)
		return "", nil
	}

	rec := recs[0]

	l.Info("account found for verification token >%s<", token)

	// Generate session token
	sessionToken := corerecord.NewRecordID()

	// Hash the session token
	hashedSessionToken, err := bcrypt.GenerateFromPassword([]byte(sessionToken), bcrypt.DefaultCost)
	if err != nil {
		l.Warn("failed to hash session token >%v<", err)
		return "", err
	}

	rec.SessionToken = nullstring.FromString(string(hashedSessionToken))
	rec.SessionTokenExpiresAt = nulltime.FromTime(corerecord.NewRecordTimestamp().Add(15 * time.Minute))

	_, err = m.UpdateAccountRec(rec)
	if err != nil {
		l.Warn("failed to update account >%v<", err)
		return "", err
	}

	// TODO: Debug log the session token
	l.Info("generated session token >%s< for account ID >%s<", sessionToken, rec.ID)

	return sessionToken, nil
}

func (m *Domain) VerifyAccountSessionToken(token string) (*record.Account, error) {
	l := m.Logger("VerifyAccountSessionToken")

	l.Info("verifying account session token >%s<", token)

	// Look up account by session token
	repo := m.AccountRepository()

	recs, err := repo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: record.FieldAccountSessionToken, Val: token},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get account by session token >%s< >%v<", token, err)
		return nil, err
	}

	if len(recs) == 0 {
		l.Info("no account found for session token >%s<", token)
		return nil, nil
	}

	rec := recs[0]

	// Has the session token expired?
	if rec.SessionTokenExpiresAt.Time.Before(corerecord.NewRecordTimestamp()) {
		l.Info("session token >%s< has expired", token)
		return nil, nil
	}

	// Extend the expiration time of the session token
	rec.SessionTokenExpiresAt = nulltime.FromTime(corerecord.NewRecordTimestamp().Add(15 * time.Minute))

	_, err = m.UpdateAccountRec(rec)
	if err != nil {
		l.Warn("failed to update account >%v<", err)
		return nil, err
	}

	l.Info("account found for session token >%s<", token)

	return rec, nil
}
