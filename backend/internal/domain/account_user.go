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

// GetAccountUserRecByAccountID returns the account_user record for the given parent account ID.
// Since each account has exactly one account_user (1:1 mapping), this returns the unique user.
func (m *Domain) GetAccountUserRecByAccountID(accountID string, lock *coresql.Lock) (*account_record.AccountUser, error) {
	l := m.Logger("GetAccountUserRecByAccountID")

	l.Debug("getting account user record by account ID >%s<", accountID)

	if err := domain.ValidateUUIDField("account_id", accountID); err != nil {
		return nil, err
	}

	r := m.AccountUserRepository()

	opts := &coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountUserAccountID, Val: accountID},
		},
		Lock:  lock,
		Limit: 1,
	}

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	if len(recs) == 0 {
		return nil, coreerror.NewNotFoundError(account_record.TableAccountUser, accountID)
	}

	return recs[0], nil
}

func (m *Domain) GetAccountUserRecByEmail(email string) (*account_record.AccountUser, error) {
	l := m.Logger("GetAccountUserRecByEmail")

	l.Debug("getting account user record by email >%s<", email)

	if email == "" {
		return nil, coreerror.NewInvalidDataError("email is required")
	}

	repo := m.AccountUserRepository()

	recs, err := repo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountUserEmail, Val: email},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get account user by email >%v<", err)
		return nil, databaseError(err)
	}

	if len(recs) == 0 {
		return nil, nil
	}

	return recs[0], nil
}

func (m *Domain) GetAccountUserRecByVerificationToken(verificationToken string) (*account_record.AccountUser, error) {
	l := m.Logger("GetAccountUserRecByVerificationToken")

	l.Debug("getting account user record by verification token >%s<", verificationToken)

	if verificationToken == "" {
		return nil, coreerror.NewInvalidDataError("verification token is required")
	}

	repo := m.AccountUserRepository()

	recs, err := repo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountUserVerificationToken, Val: verificationToken},
		},
		Limit: 1,
	})
	if err != nil {
		l.Warn("failed to get account user by verification token >%v<", err)
		return nil, databaseError(err)
	}

	if len(recs) == 0 {
		return nil, nil
	}

	return recs[0], nil
}

func (m *Domain) GetAccountUserRecBySessionToken(token string) (*account_record.AccountUser, error) {
	l := m.Logger("GetAccountUserRecBySessionToken")

	l.Debug("getting account user record by session token >%s<", token)

	if token == "" {
		return nil, coreerror.NewInvalidDataError("session token is required")
	}

	// HMAC hash the provided token
	hash := hmacSHA256(m.config.TokenHMACKey, token)

	// Look up user account by session token
	repo := m.AccountUserRepository()

	recs, err := repo.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{Col: account_record.FieldAccountUserSessionToken, Val: hash},
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

	return rec, nil
}

// CreateAccountUserRec creates an account user record
func (m *Domain) CreateAccountUserRec(rec *account_record.AccountUser) (*account_record.AccountUser, error) {
	l := m.Logger("CreateAccountUserRec")

	l.Debug("creating account user record >%#v<", rec)

	if rec != nil && rec.Status == "" {
		rec.Status = account_record.AccountUserStatusPendingApproval
	}

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

// UpdateAccountRec uses FOR UPDATE (wait) rather than NOWAIT because brief
// contention is expected between River job workers and API requests operating
// on the same account row (e.g. verification token generation vs session
// token generation).
func (m *Domain) UpdateAccountUserRec(rec *account_record.AccountUser) (*account_record.AccountUser, error) {
	l := m.Logger("UpdateAccountUserRec")

	accountUserRec, err := m.GetAccountUserRec(rec.ID, coresql.ForUpdate)
	if err != nil {
		return rec, err
	}

	l.Debug("updating account user record >%#v<", rec)

	if rec.Status == "" {
		rec.Status = accountUserRec.Status
	}

	if err := m.validateAccountUserRecForUpdate(accountUserRec, rec); err != nil {
		l.Warn("failed to validate account user record >%v<", err)
		return rec, err
	}

	r := m.AccountUserRepository()

	updatedAccountUserRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedAccountUserRec, nil
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

// RemoveAccountUserRec hard-deletes an account_user record by ID.
func (m *Domain) RemoveAccountUserRec(recID string) error {
	l := m.Logger("RemoveAccountUserRec")

	l.Debug("removing account user record ID >%s<", recID)

	accountUserFilter := &coresql.Options{
		Params: []coresql.Param{
			{Col: "account_user_id", Val: recID},
		},
	}

	// 1. game_subscription_instance (references account_user_id)
	gameSubscriptionInstanceRecs, err := m.GetManyGameSubscriptionInstanceRecs(accountUserFilter)
	if err != nil {
		return databaseError(err)
	}
	for _, gameSubscriptionInstanceRec := range gameSubscriptionInstanceRecs {
		if err := m.RemoveGameSubscriptionInstanceRec(gameSubscriptionInstanceRec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 2. game_subscription (references account_user_id)
	gameSubscriptionRecs, err := m.GetManyGameSubscriptionRecs(accountUserFilter)
	if err != nil {
		return databaseError(err)
	}
	for _, rec := range gameSubscriptionRecs {
		if err := m.RemoveGameSubscriptionRec(rec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 3. account_subscription (by account_user_id for player subs)
	accountSubscriptionRecs, err := m.GetManyAccountSubscriptionRecs(accountUserFilter)
	if err != nil {
		return databaseError(err)
	}
	for _, rec := range accountSubscriptionRecs {
		if err := m.RemoveAccountSubscriptionRec(rec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 4. account_user
	r := m.AccountUserRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	l.Info("removed account user >%s< and all dependents (%d subscriptions, %d game subscriptions, %d game subscription instances)",
		recID, len(accountSubscriptionRecs), len(gameSubscriptionRecs), len(gameSubscriptionInstanceRecs))

	return nil
}

// User account session functions that should move out of this file..

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

func (m *Domain) GenerateAccountUserVerificationToken(rec *account_record.AccountUser) (string, error) {
	l := m.Logger("GenerateAccountUserVerificationToken")

	l.Debug("generating verification token for user account ID >%s<", rec.ID)

	// Generate a new UUID for the token
	token := corerecord.NewRecordID()

	// HMAC hash the token
	hash := hmacSHA256(m.config.TokenHMACKey, token)

	rec.VerificationToken = nullstring.FromString(hash)
	rec.VerificationTokenExpiresAt = nulltime.FromTime(corerecord.NewRecordTimestamp().Add(15 * time.Minute))

	_, err := m.UpdateAccountUserRec(rec)
	if err != nil {
		l.Warn("failed to update user account >%v<", err)
		return "", err
	}

	l.Info("generated verification token >%s< for user account ID >%s<", token, rec.ID)

	return token, nil
}

// VerifyAccountUserVerificationToken verifies a verification token and returns a session token.
// If testBypassEnabled is true, the token is treated as an email address for test authentication.
func (m *Domain) VerifyAccountUserVerificationToken(token string, testBypassEnabled bool) (string, error) {
	l := m.Logger("VerifyAccountUserVerificationToken")

	l.Info("verifying user account verification token >%s<", token)

	var rec *account_record.AccountUser
	var err error
	// Test bypass mode: allow using email address as the verification code
	// This is enabled when the caller has verified the test bypass header/token
	if testBypassEnabled {
		// Check if the token looks like an email address
		if isEmailAddress(token) {
			l.Info("test bypass: attempting magic auth with email >%s<", token)

			rec, err = m.GetAccountUserRecByEmail(token)
			if err != nil {
				l.Warn("test bypass: failed to get user account by email >%s< >%v<", token, err)
				return "", err
			}

			if rec != nil {
				l.Info("test bypass: magic auth successful for email >%s<", token)
			}
		}
	}

	// If no test bypass match, try the normal verification token lookup
	if rec == nil {
		// HMAC hash the provided token
		hash := hmacSHA256(m.config.TokenHMACKey, token)

		rec, err = m.GetAccountUserRecByVerificationToken(hash)
		if err != nil {
			l.Warn("failed to get user account by verification token >%s< >%v<", token, err)
			return "", err
		}

		if rec != nil {
			l.Info("verification token lookup: auth successful for verification token >%s<", token)
		}
	}

	if rec == nil {
		l.Info("no account found for verification token >%s<", token)
		return "", coreerror.NewInvalidDataError("Invalid verification code.")
	}

	l.Info("user account found for verification token >%s<", token)

	// Generate session token
	sessionToken, err := m.GenerateAccountUserSessionToken(rec)
	if err != nil {
		l.Warn("failed to generate session token >%v<", err)
		return "", err
	}

	return sessionToken, nil
}

// GenerateAccountUserSessionToken generates a session token for an account user record.
func (m *Domain) GenerateAccountUserSessionToken(rec *account_record.AccountUser) (string, error) {
	l := m.Logger("GenerateAccountUserSessionToken")

	l.Debug("generating session token for user account ID >%s<", rec.ID)

	// Generate session token
	sessionToken := corerecord.NewRecordID()

	// Hash the session token
	hashedSessionToken := hmacSHA256(m.config.TokenHMACKey, sessionToken)

	rec.SessionToken = nullstring.FromString(hashedSessionToken)
	rec.SessionTokenExpiresAt = nulltime.FromTime(corerecord.NewRecordTimestamp().Add(sessionTokenExpiryDuration))

	_, err := m.UpdateAccountUserRec(rec)
	if err != nil {
		l.Warn("failed to update user account >%v<", err)
		return "", err
	}

	l.Info("generated session token >%s< for user account ID >%s<", sessionToken, rec.ID)

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

func (m *Domain) VerifyAccountUserSessionToken(token string) (*account_record.AccountUser, error) {
	l := m.Logger("VerifyAccountUserSessionToken")

	l.Info("verifying user account session token >%s<", token)

	rec, err := m.GetAccountUserRecBySessionToken(token)
	if err != nil {
		l.Warn("failed to get account user by session token >%s< >%v<", token, err)
		return nil, err
	}

	if rec == nil {
		l.Info("no account user found for session token >%s<", token)
		return nil, nil
	}

	// Has the session token expired?
	if rec.SessionTokenExpiresAt.Time.Before(corerecord.NewRecordTimestamp()) {
		l.Info("session token >%s< has expired", token)
		return nil, nil
	}

	// Now get the account record with a lock; we want to extend the
	// session token expiration time and need to wait for any concurrent
	// updates to the account record.
	rec, err = m.GetAccountUserRec(rec.ID, coresql.ForUpdate)
	if err != nil {
		l.Warn("failed to get account >%v<", err)
		return nil, err
	}

	// Extend the expiration time of the session token
	rec.SessionTokenExpiresAt = nulltime.FromTime(corerecord.NewRecordTimestamp().Add(sessionTokenExpiryDuration))

	_, err = m.UpdateAccountUserRec(rec)
	if err != nil {
		l.Warn("failed to update account >%v<", err)
		return nil, err
	}

	l.Info("account found for session token >%s<", token)

	return rec, nil
}
