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

// GetAccountParentRec returns the account (parent/tenant) record by ID.
func (m *Domain) GetAccountParentRec(recID string, lock *coresql.Lock) (*account_record.Account, error) {
	l := m.Logger("GetAccountParentRec")

	l.Debug("getting account parent record ID >%s<", recID)

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

// GetManyAccountParentRecs returns account (parent/tenant) records.
func (m *Domain) GetManyAccountParentRecs(opts *coresql.Options) ([]*account_record.Account, error) {
	l := m.Logger("GetManyAccountParentRecs")

	l.Info("getting many account parent records opts >%#v<", opts)

	r := m.AccountRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// UpdateAccountParentRec updates an account (parent/tenant) record.
func (m *Domain) UpdateAccountParentRec(rec *account_record.Account) (*account_record.Account, error) {
	l := m.Logger("UpdateAccountParentRec")

	l.Debug("updating account parent record ID >%s<", rec.ID)

	curr, err := m.GetAccountParentRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	rec.CreatedAt = curr.CreatedAt

	r := m.AccountRepository()

	rec, err = r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// GetManyAccountRecs -
func (m *Domain) GetManyAccountRecs(opts *coresql.Options) ([]*account_record.AccountUser, error) {
	l := m.Logger("GetManyAccountRecs")

	l.Info("getting many account records opts >%#v<", opts)

	r := m.AccountUserRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAccountRec -
func (m *Domain) GetAccountRec(recID string, lock *coresql.Lock) (*account_record.AccountUser, error) {
	l := m.Logger("GetAccountRec")

	l.Debug("getting client record ID >%s<", recID)

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

// CreateAccount creates a new account with an account user and basic subscriptions
func (m *Domain) CreateAccount(rec *account_record.AccountUser) (*account_record.Account, *account_record.AccountUser, []*account_record.AccountSubscription, error) {
	l := m.Logger("CreateAccount")

	l.Info("CreateAccount called for account email >%s<", rec.Email)

	// Create the account record first (1:1 mapping with account_user)
	accountRec := &account_record.Account{
		Name:   rec.Email,
		Status: account_record.AccountUserStatusActive,
	}

	accountRepo := m.AccountRepository()
	accountRec, err := accountRepo.CreateOne(accountRec)
	if err != nil {
		return nil, rec, nil, err
	}

	l.Info("created account >%s< for user >%s<", accountRec.ID, rec.Email)

	// Link user to account
	rec.AccountID = accountRec.ID

	// Create the account user record
	createdRec, err := m.CreateAccountRec(rec)
	if err != nil {
		return accountRec, rec, nil, err
	}

	var subscriptions []*account_record.AccountSubscription

	// Create basic game designer subscription
	l.Info("creating basic game designer subscription for account >%s<", accountRec.ID)
	designerSub := &account_record.AccountSubscription{
		AccountID:          nullstring.FromString(accountRec.ID),
		SubscriptionType:   account_record.AccountSubscriptionTypeBasicGameDesigner,
		SubscriptionPeriod: account_record.AccountSubscriptionPeriodEternal,
		Status:             account_record.AccountSubscriptionStatusActive,
		AutoRenew:          true,
	}
	designerSub, err = m.CreateAccountSubscriptionRec(designerSub)
	if err != nil {
		l.Warn("failed to auto-create basic game designer subscription >%v<", err)
	} else {
		l.Info("created basic game designer subscription ID >%s< Type >%s< Status >%s< for account >%s<",
			designerSub.ID, designerSub.SubscriptionType, designerSub.Status, createdRec.ID)
		subscriptions = append(subscriptions, designerSub)
	}

	// Create basic manager subscription
	l.Info("creating basic manager subscription for account >%s<", accountRec.ID)
	managerSub := &account_record.AccountSubscription{
		AccountID:          nullstring.FromString(accountRec.ID),
		SubscriptionType:   account_record.AccountSubscriptionTypeBasicManager,
		SubscriptionPeriod: account_record.AccountSubscriptionPeriodEternal,
		Status:             account_record.AccountSubscriptionStatusActive,
		AutoRenew:          true,
	}
	managerSub, err = m.CreateAccountSubscriptionRec(managerSub)
	if err != nil {
		l.Warn("failed to auto-create basic manager subscription >%v<", err)
	} else {
		l.Info("created basic manager subscription ID >%s< Type >%s< Status >%s< for account >%s<",
			managerSub.ID, managerSub.SubscriptionType, managerSub.Status, accountRec.ID)
		subscriptions = append(subscriptions, managerSub)
	}

	// Create basic player subscription
	l.Info("creating basic player subscription for account >%s<", accountRec.ID)
	playerSub := &account_record.AccountSubscription{
		AccountUserID:      nullstring.FromString(rec.ID),
		SubscriptionType:   account_record.AccountSubscriptionTypeBasicPlayer,
		SubscriptionPeriod: account_record.AccountSubscriptionPeriodEternal,
		Status:             account_record.AccountSubscriptionStatusActive,
		AutoRenew:          true,
	}
	playerSub, err = m.CreateAccountSubscriptionRec(playerSub)
	if err != nil {
		l.Warn("failed to auto-create basic player subscription >%v<", err)
	} else {
		l.Info("created basic player subscription ID >%s< Type >%s< Status >%s< for account >%s<",
			playerSub.ID, playerSub.SubscriptionType, playerSub.Status, accountRec.ID)
		subscriptions = append(subscriptions, playerSub)
	}

	l.Info("created account >%s< with >%d< subscriptions", createdRec.ID, len(subscriptions))
	return accountRec, createdRec, subscriptions, nil
}

// CreateAccountRec creates an account record without subscriptions
func (m *Domain) CreateAccountRec(rec *account_record.AccountUser) (*account_record.AccountUser, error) {
	l := m.Logger("CreateAccountRec")

	l.Debug("creating client record >%#v<", rec)

	r := m.AccountUserRepository()

	if err := m.validateAccountRecForCreate(rec); err != nil {
		l.Warn("failed to validate client record >%v<", err)
		return rec, err
	}

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
func (m *Domain) UpdateAccountRec(rec *account_record.AccountUser) (*account_record.AccountUser, error) {
	l := m.Logger("UpdateAccountRec")

	curr, err := m.GetAccountRec(rec.ID, coresql.ForUpdate)
	if err != nil {
		return rec, err
	}

	l.Debug("updating client record >%#v<", rec)

	if err := m.validateAccountRecForUpdate(curr, rec); err != nil {
		l.Warn("failed to validate client record >%v<", err)
		return rec, err
	}

	r := m.AccountUserRepository()

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

	r := m.AccountUserRepository()

	if err := m.validateAccountRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAccountUserRec hard-deletes an account_user record by ID.
func (m *Domain) RemoveAccountUserRec(recID string) error {
	l := m.Logger("RemoveAccountUserRec")

	l.Debug("removing account user record ID >%s<", recID)

	rec, err := m.GetAccountRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AccountUserRepository()

	if err := m.validateAccountRecForDelete(rec); err != nil {
		l.Warn("failed domain validation >%v<", err)
		return err
	}

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAccountRec hard-deletes an account and all dependent records in FK order.
func (m *Domain) RemoveAccountRec(recID string) error {
	l := m.Logger("RemoveAccountRec")

	l.Debug("removing account record ID >%s< and all dependents", recID)

	accountFilter := &coresql.Options{Params: []coresql.Param{{Col: "account_id", Val: recID}}}

	// 1. game_subscription_instance (references account_id, game_subscription_id)
	gsiRecs, err := m.GameSubscriptionInstanceRepository().GetMany(accountFilter)
	if err != nil {
		return databaseError(err)
	}
	for _, rec := range gsiRecs {
		if err := m.GameSubscriptionInstanceRepository().RemoveOne(rec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 2. game_subscription (references account_id)
	gsRecs, err := m.GameSubscriptionRepository().GetMany(accountFilter)
	if err != nil {
		return databaseError(err)
	}
	for _, rec := range gsRecs {
		if err := m.GameSubscriptionRepository().RemoveOne(rec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 3. Fetch account_user records (needed for player subscription and contact cleanup)
	auRecs, err := m.AccountUserRepository().GetMany(accountFilter)
	if err != nil {
		return databaseError(err)
	}

	// 4. account_subscription (by account_id for designer/manager subs, by account_user_id for player subs)
	asRecs, err := m.AccountSubscriptionRepository().GetMany(accountFilter)
	if err != nil {
		return databaseError(err)
	}
	for _, rec := range asRecs {
		if err := m.AccountSubscriptionRepository().RemoveOne(rec.ID); err != nil {
			return databaseError(err)
		}
	}
	asCount := len(asRecs)
	for _, auRec := range auRecs {
		userSubFilter := &coresql.Options{Params: []coresql.Param{{Col: account_record.FieldAccountSubscriptionAccountUserID, Val: auRec.ID}}}
		userSubs, err := m.AccountSubscriptionRepository().GetMany(userSubFilter)
		if err != nil {
			return databaseError(err)
		}
		for _, rec := range userSubs {
			if err := m.AccountSubscriptionRepository().RemoveOne(rec.ID); err != nil {
				return databaseError(err)
			}
		}
		asCount += len(userSubs)
	}

	// 5. account_user_contact (by account_user_id)
	for _, auRec := range auRecs {
		contactFilter := &coresql.Options{Params: []coresql.Param{{Col: account_record.FieldAccountUserContactAccountUserID, Val: auRec.ID}}}
		contactRecs, err := m.AccountUserContactRepository().GetMany(contactFilter)
		if err != nil {
			return databaseError(err)
		}
		for _, rec := range contactRecs {
			if err := m.AccountUserContactRepository().RemoveOne(rec.ID); err != nil {
				return databaseError(err)
			}
		}
	}

	// 6. account_user
	for _, rec := range auRecs {
		if err := m.AccountUserRepository().RemoveOne(rec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 7. account
	if err := m.AccountRepository().RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	l.Info("removed account >%s< and all dependents (%d users, %d subscriptions, %d game subscriptions, %d game subscription instances)",
		recID, len(auRecs), asCount, len(gsRecs), len(gsiRecs))

	return nil
}

func (m *Domain) GenerateAccountVerificationToken(rec *account_record.AccountUser) (string, error) {
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

	repo := m.AccountUserRepository()

	var rec *account_record.AccountUser

	// Test bypass mode: allow using email address as the verification code
	// This is enabled when the caller has verified the test bypass header/token
	if testBypassEnabled {
		// Check if the token looks like an email address
		if isEmailAddress(token) {
			l.Info("test bypass: attempting magic auth with email >%s<", token)

			recs, err := repo.GetMany(&coresql.Options{
				Params: []coresql.Param{
					{Col: account_record.FieldAccountUserEmail, Val: token},
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
				{Col: account_record.FieldAccountUserVerificationToken, Val: hash},
			},
			Limit: 1,
		})
		if err != nil {
			l.Warn("failed to get account by verification token >%s< >%v<", token, err)
			return "", err
		}

		if len(recs) > 0 {
			rec = recs[0]
		}
	}

	if rec == nil {
		l.Info("no account found for verification token >%s<", token)
		return "", coreerror.NewInvalidDataError("Invalid verification code.")
	}

	l.Info("account found for verification token >%s<", token)

	// Generate session token
	sessionToken, err := m.GenerateAccountSessionToken(rec)
	if err != nil {
		l.Warn("failed to generate session token >%v<", err)
		return "", err
	}

	return sessionToken, nil
}

// GenerateAccountSessionToken generates a session token for an account record.
func (m *Domain) GenerateAccountSessionToken(rec *account_record.AccountUser) (string, error) {
	l := m.Logger("GenerateAccountSessionToken")

	l.Debug("generating session token for account ID >%s<", rec.ID)

	// Generate session token
	sessionToken := corerecord.NewRecordID()

	// Hash the session token
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

func (m *Domain) VerifyAccountSessionToken(token string) (*account_record.AccountUser, error) {
	l := m.Logger("VerifyAccountSessionToken")

	l.Info("verifying account session token >%s<", token)

	// HMAC hash the provided token
	hash := hmacSHA256(m.config.TokenHMACKey, token)

	// Look up account by session token
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

func (m *Domain) GetAccountRecByEmail(email string) (*account_record.AccountUser, error) {
	l := m.Logger("GetAccountRecByEmail")

	l.Debug("getting account record by email >%s<", email)

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
		l.Warn("failed to get account by email >%v<", err)
		return nil, databaseError(err)
	}

	if len(recs) == 0 {
		return nil, nil
	}

	return recs[0], nil
}
