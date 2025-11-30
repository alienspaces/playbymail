package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	corerecord "gitlab.com/alienspaces/playbymail/core/record"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

// Utility for HMAC-SHA256 hashing
func hmacSHA256(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// isEmailAddress checks if a string looks like an email address (simple check)
func isEmailAddress(s string) bool {
	// Simple check: contains @ and has text before and after it
	atIndex := strings.Index(s, "@")
	if atIndex <= 0 || atIndex >= len(s)-1 {
		return false
	}
	// Check for at least one dot after @
	afterAt := s[atIndex+1:]
	return strings.Contains(afterAt, ".") && !strings.HasSuffix(afterAt, ".")
}

// Session token expiry duration
const sessionTokenExpiryDuration = 15 * time.Minute

// SessionTokenExpiryDuration returns the duration after which session tokens expire.
func (m *Domain) SessionTokenExpiryDuration() time.Duration {
	return sessionTokenExpiryDuration
}

// SessionTokenExpirySeconds returns the number of seconds until session tokens expire.
func (m *Domain) SessionTokenExpirySeconds() int {
	return int(sessionTokenExpiryDuration.Seconds())
}

// GetManyAccountRecs -
func (m *Domain) GetManyAccountRecs(opts *coresql.Options) ([]*account_record.Account, error) {
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
func (m *Domain) GetAccountRec(recID string, lock *coresql.Lock) (*account_record.Account, error) {
	l := m.Logger("GetAccountRec")

	l.Debug("getting client record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AccountRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(account_record.TableAccount, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAccountRec -
func (m *Domain) CreateAccountRec(rec *account_record.Account) (*account_record.Account, error) {
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
func (m *Domain) UpdateAccountRec(rec *account_record.Account) (*account_record.Account, error) {
	l := m.Logger("UpdateAccountRec")

	curr, err := m.GetAccountRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating client record >%#v<", rec)

	if err := m.validateAccountRecForUpdate(rec, curr); err != nil {
		l.Warn("failed to validate client record >%v<", err)
		return rec, err
	}

	r := m.AccountRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
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

func (m *Domain) GenerateAccountVerificationToken(rec *account_record.Account) (string, error) {
	l := m.Logger("GenerateAccountVerificationToken")

	l.Debug("generating verification token for account ID >%s<", rec.ID)

	// Generate a new UUID for the token
	token := corerecord.NewRecordID()

	// HMAC hash the token
	hash := hmacSHA256(m.config.TokenHMACKey, token)

	rec.VerificationToken = nullstring.FromString(hash)
	rec.VerificationTokenExpiresAt = nulltime.FromTime(corerecord.NewRecordTimestamp().Add(15 * time.Minute))

	_, err := m.UpdateAccountRec(rec)
	if err != nil {
		l.Warn("failed to update account >%v<", err)
		return "", err
	}

	l.Info("generated verification token >%s< for account ID >%s<", token, rec.ID)

	return token, nil
}

// VerifyAccountVerificationToken verifies a verification token and returns a session token.
// If testBypassEnabled is true, the token is treated as an email address for test authentication.
func (m *Domain) VerifyAccountVerificationToken(token string, testBypassEnabled bool) (string, error) {
	l := m.Logger("VerifyAccountVerificationToken")

	l.Info("verifying verification token >%s<", token)

	repo := m.AccountRepository()

	var rec *account_record.Account

	// Test bypass mode: allow using email address as the verification code
	// This is enabled when the caller has verified the test bypass header/token
	if testBypassEnabled {
		// Check if the token looks like an email address
		if isEmailAddress(token) {
			l.Info("test bypass: attempting magic auth with email >%s<", token)

			recs, err := repo.GetMany(&coresql.Options{
				Params: []coresql.Param{
					{Col: account_record.FieldAccountEmail, Val: token},
				},
				Limit: 1,
			})
			if err != nil {
				l.Warn("test bypass: failed to get account by email >%s< >%v<", token, err)
				return "", err
			}

			if len(recs) > 0 {
				rec = recs[0]
				l.Info("test bypass: magic auth successful for email >%s<", token)
			}
		}
	}

	// If no test bypass match, try the normal verification token lookup
	if rec == nil {
		// HMAC hash the provided token
		hash := hmacSHA256(m.config.TokenHMACKey, token)

		recs, err := repo.GetMany(&coresql.Options{
			Params: []coresql.Param{
				{Col: account_record.FieldAccountVerificationToken, Val: hash},
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

		rec = recs[0]
	}

	l.Info("account found for verification token >%s<", token)

	// Generate session token
	sessionToken := corerecord.NewRecordID()

	// Hash the session token (keep using bcrypt or switch to HMAC as well if desired)
	hashedSessionToken := hmacSHA256(m.config.TokenHMACKey, sessionToken)

	rec.SessionToken = nullstring.FromString(hashedSessionToken)
	rec.SessionTokenExpiresAt = nulltime.FromTime(corerecord.NewRecordTimestamp().Add(sessionTokenExpiryDuration))

	_, err := m.UpdateAccountRec(rec)
	if err != nil {
		l.Warn("failed to update account >%v<", err)
		return "", err
	}

	l.Info("generated session token >%s< for account ID >%s<", sessionToken, rec.ID)

	return sessionToken, nil
}

// IsTestBypassEnabled checks if test bypass authentication is enabled based on
// the configured header name, value, and the provided header value.
func (m *Domain) IsTestBypassEnabled(headerValue string) bool {
	// Both config values must be set
	if m.config.TestBypassHeaderName == "" || m.config.TestBypassHeaderValue == "" {
		return false
	}

	// Header value must match the configured value exactly
	return headerValue == m.config.TestBypassHeaderValue
}

// GetTestBypassHeaderName returns the configured test bypass header name.
func (m *Domain) GetTestBypassHeaderName() string {
	return m.config.TestBypassHeaderName
}

func (m *Domain) VerifyAccountSessionToken(token string) (*account_record.Account, error) {
	l := m.Logger("VerifyAccountSessionToken")

	l.Info("verifying account session token >%s<", token)

	// HMAC hash the provided token
	hash := hmacSHA256(m.config.TokenHMACKey, token)

	// Look up account by session token
	repo := m.AccountRepository()

	recs, err := repo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountSessionToken, Val: hash},
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

	// Now get the account record with a lock; we want to extend the
	// session token expiration time and need to wait for any concurrent
	// updates to the account record.
	rec, err = m.GetAccountRec(rec.ID, coresql.ForUpdate)
	if err != nil {
		l.Warn("failed to get account >%v<", err)
		return nil, err
	}

	// Extend the expiration time of the session token
	rec.SessionTokenExpiresAt = nulltime.FromTime(corerecord.NewRecordTimestamp().Add(sessionTokenExpiryDuration))

	_, err = m.UpdateAccountRec(rec)
	if err != nil {
		l.Warn("failed to update account >%v<", err)
		return nil, err
	}

	l.Info("account found for session token >%s<", token)

	return rec, nil
}

func (m *Domain) GetAccountRecByEmail(email string) (*account_record.Account, error) {
	l := m.Logger("GetAccountRecByEmail")

	l.Debug("getting account record by email >%s<", email)

	if email == "" {
		return nil, coreerror.NewInvalidDataError("email is required")
	}

	repo := m.AccountRepository()

	recs, err := repo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountEmail, Val: email},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get account by email >%v<", err)
		return nil, databaseError(err)
	}

	if len(recs) == 0 {
		return nil, nil
	}

	return recs[0], nil
}
