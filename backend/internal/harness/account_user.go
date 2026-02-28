package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

func (t *Testing) createAccountUserRec(accountConfig AccountConfig) (*account_record.AccountUser, error) {
	l := t.Logger("createAccountUserRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *account_record.AccountUser
	if accountConfig.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *accountConfig.Record
		rec = &recCopy
	} else {
		rec = &account_record.AccountUser{}
	}

	rec = t.applyAccountUserRecDefaultValues(rec)

	l.Info("creating account record with basic subscriptions for ref >%s<", accountConfig.Reference)
	accountRec, createdRec, accountSubscriptionRecs, err := t.Domain.(*domain.Domain).CreateAccount(rec)
	if err != nil {
		l.Warn("failed creating account record >%v<", err)
		return nil, err
	}

	// Store all basic subscriptions in harness data
	l.Info("storing >%d< basic account subscriptions in harness data for account >%s<", len(accountSubscriptionRecs), createdRec.ID)
	for _, accountSubscriptionRec := range accountSubscriptionRecs {
		t.Data.AddAccountSubscriptionRec(accountSubscriptionRec)
		t.teardownData.AddAccountSubscriptionRec(accountSubscriptionRec)
	}

	// Create additional account subscriptions from config
	for _, subscriptionConfig := range accountConfig.AccountSubscriptionConfigs {
		_, err = t.createAccountSubscriptionRec(subscriptionConfig, createdRec)
		if err != nil {
			l.Warn("failed creating account subscription >%v<", err)
			return nil, err
		}
		l.Debug("created account subscription for account >%s<", createdRec.ID)
	}

	// Always generate session token for all accounts and store in harness data
	sessionToken, err := t.Domain.(*domain.Domain).GenerateAccountSessionToken(createdRec)
	if err != nil {
		l.Warn("failed to generate session token for account >%s< >%v<", createdRec.ID, err)
		return nil, err
	}

	// Store the session token in harness data by account ID
	t.Data.AddAccountSessionToken(createdRec.ID, sessionToken)

	l.Info("generated session token for account >%s< token >%s<", createdRec.ID, sessionToken)

	// Add the account user record to the data store
	t.Data.AddAccountUserRec(createdRec)

	// Add the account user record to the teardown data store
	t.teardownData.AddAccountUserRec(createdRec)

	t.Data.AddAccountRec(accountRec)
	t.teardownData.AddAccountRec(accountRec)

	// Add the account user record to the data store refs
	if accountConfig.Reference != "" {
		t.Data.Refs.AccountUserRefs[accountConfig.Reference] = createdRec.ID
	}

	return createdRec, nil
}

func (t *Testing) createAccountSubscriptionRec(subscriptionConfig AccountSubscriptionConfig, accountRec *account_record.AccountUser) (*account_record.AccountSubscription, error) {
	l := t.Logger("createAccountSubscriptionRec")

	if accountRec == nil {
		return nil, fmt.Errorf("account record is nil for account_subscription")
	}

	var rec *account_record.AccountSubscription
	if subscriptionConfig.Record != nil {
		recCopy := *subscriptionConfig.Record
		rec = &recCopy
	} else {
		rec = &account_record.AccountSubscription{}
	}

	rec = t.applyAccountSubscriptionRecDefaultValues(rec)

	// Set subscription type if provided
	if subscriptionConfig.SubscriptionType != "" {
		rec.SubscriptionType = subscriptionConfig.SubscriptionType
	}

	// Set account or account user ID based on subscription type
	if rec.SubscriptionType == account_record.AccountSubscriptionTypeBasicPlayer ||
		rec.SubscriptionType == account_record.AccountSubscriptionTypeProfessionalPlayer {
		rec.AccountUserID = nullstring.FromString(accountRec.ID)
		rec.AccountID = nullstring.FromString("")
	} else {
		rec.AccountID = nullstring.FromString(accountRec.AccountID)
		rec.AccountUserID = nullstring.FromString("")
	}

	// Create record
	l.Debug("creating account_subscription record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAccountSubscriptionRec(rec)
	if err != nil {
		l.Warn("failed creating account_subscription record >%v<", err)
		return nil, err
	}

	// Add to data store and teardown data so it is removed during teardown
	t.Data.AddAccountSubscriptionRec(createdRec)
	t.teardownData.AddAccountSubscriptionRec(createdRec)

	return createdRec, nil
}

func (t *Testing) applyAccountSubscriptionRecDefaultValues(rec *account_record.AccountSubscription) *account_record.AccountSubscription {
	if rec == nil {
		rec = &account_record.AccountSubscription{}
	}

	if rec.SubscriptionPeriod == "" {
		rec.SubscriptionPeriod = account_record.AccountSubscriptionPeriodEternal
	}

	if rec.Status == "" {
		rec.Status = account_record.AccountSubscriptionStatusActive
	}

	if !rec.AutoRenew {
		rec.AutoRenew = true
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
