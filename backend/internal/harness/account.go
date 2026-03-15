package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

func (t *Testing) processAccountConfig(accountConfig AccountConfig) (
	*account_record.Account, []*account_record.AccountUser, []*account_record.AccountUserContact, []*account_record.AccountSubscription, error,
) {
	l := t.Logger("processAccountConfig")

	var allAccountUserRecs []*account_record.AccountUser
	var allAccountUserContactRecs []*account_record.AccountUserContact
	var allAccountSubscriptionRecs []*account_record.AccountSubscription

	// Create a new record instance to avoid reusing the same record across tests
	var accountRec *account_record.Account
	if accountConfig.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *accountConfig.Record
		accountRec = &recCopy
	}

	// Apply account record harness default values
	accountRec = t.applyAccountRecDefaultValues(accountRec)

	for _, accountUserConfig := range accountConfig.AccountUserConfigs {

		var accountUserRec *account_record.AccountUser
		if accountUserConfig.Record != nil {
			recCopy := *accountUserConfig.Record
			accountUserRec = &recCopy
		} else {
			accountUserRec = &account_record.AccountUser{}
		}

		// Apply account user record harness default values
		accountUserRec = t.applyAccountUserRecDefaultValues(accountUserRec)

		l.Info("creating account record with basic subscriptions for ref >%s<", accountConfig.Reference)

		// Upsert account creates or upodates an account, account user, account user contact and
		// basic subscriptions records. If the account user contact or subscriptions are not provided,
		// they are created with default values.
		accountRec, createdRec, accountUserContactRec, accountSubscriptionRecs, err := t.Domain.(*domain.Domain).UpsertAccount(
			accountRec, accountUserRec, &account_record.AccountUserContact{},
		)
		if err != nil {
			l.Warn("failed upserting account record >%v<", err)
			return nil, nil, nil, nil, err
		}

		// Store the account record in harness data
		t.Data.AddAccountRec(accountRec)
		t.teardownData.AddAccountRec(accountRec)

		// Store the account user record in harness data (createdRec has the ID from UpsertAccount)
		t.Data.AddAccountUserRec(createdRec)
		t.teardownData.AddAccountUserRec(createdRec)

		// Add the account user record to the data store refs (keyed by account user ref for lookup by subscription/character configs)
		if accountUserConfig.Reference != "" {
			t.Data.Refs.AccountUserRefs[accountUserConfig.Reference] = createdRec.ID
		}

		// Store the account user contact record in harness data
		t.Data.AddAccountUserContactRec(accountUserContactRec)
		t.teardownData.AddAccountUserContactRec(accountUserContactRec)

		// Store all basic subscriptions in harness data
		for _, accountSubscriptionRec := range accountSubscriptionRecs {
			t.Data.AddAccountSubscriptionRec(accountSubscriptionRec)
			t.teardownData.AddAccountSubscriptionRec(accountSubscriptionRec)
		}

		allAccountUserRecs = append(allAccountUserRecs, createdRec)
		allAccountUserContactRecs = append(allAccountUserContactRecs, accountUserContactRec)
		allAccountSubscriptionRecs = append(allAccountSubscriptionRecs, accountSubscriptionRecs...)

		// Generate a session token for all account users and store in harness data
		// so that it may be used in API handler tests.
		sessionToken, err := t.Domain.(*domain.Domain).GenerateAccountUserSessionToken(createdRec)
		if err != nil {
			l.Warn("failed to generate session token for account user >%s< >%v<", createdRec.ID, err)
			return nil, nil, nil, nil, err
		}
		t.Data.AddAccountSessionToken(createdRec.ID, sessionToken)

		// Process additional user account subscription configurations
		for _, accountUserSubscriptionConfig := range accountUserConfig.AccountUserSubscriptionConfigs {

			accountSubscriptionRec, err := t.createAccountUserSubscription(accountUserSubscriptionConfig, createdRec)
			if err != nil {
				l.Warn("failed creating account user subscription record >%v<", err)
				return nil, nil, nil, nil, err
			}

			allAccountSubscriptionRecs = append(allAccountSubscriptionRecs, accountSubscriptionRec)
		}

	}

	return accountRec, allAccountUserRecs, allAccountUserContactRecs, allAccountSubscriptionRecs, nil
}

func (t *Testing) applyAccountRecDefaultValues(rec *account_record.Account) *account_record.Account {
	if rec == nil {
		rec = &account_record.Account{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Company())
	}

	return rec
}

func (t *Testing) applyAccountUserRecDefaultValues(rec *account_record.AccountUser) *account_record.AccountUser {
	if rec == nil {
		rec = &account_record.AccountUser{}
	}

	if rec.Email == "" {
		rec.Email = UniqueEmail(gofakeit.Email())
	}

	return rec
}
