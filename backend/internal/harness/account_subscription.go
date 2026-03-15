package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

func (t *Testing) createAccountUserSubscription(subscriptionConfig AccountUserSubscriptionConfig, accountUserRec *account_record.AccountUser) (*account_record.AccountSubscription, error) {
	l := t.Logger("createAccountSubscriptionRec")

	if accountUserRec == nil {
		return nil, fmt.Errorf("account user record is nil for account subscription")
	}

	var rec *account_record.AccountSubscription
	if subscriptionConfig.Record != nil {
		recCopy := *subscriptionConfig.Record
		rec = &recCopy
	}

	rec = t.applyAccountSubscriptionRecDefaultValues(rec)

	// Set subscription type if provided
	if subscriptionConfig.SubscriptionType != "" {
		rec.SubscriptionType = subscriptionConfig.SubscriptionType
	}

	// All subscription types use both account_id and account_user_id.
	rec.AccountID = nullstring.FromString(accountUserRec.AccountID)
	rec.AccountUserID = nullstring.FromString(accountUserRec.ID)

	// Create record
	l.Debug("creating account user subscription record >%#v<", rec)

	accountSubscriptionRec, err := t.Domain.(*domain.Domain).CreateAccountSubscriptionRec(rec)
	if err != nil {
		l.Warn("failed creating account user subscription record >%v<", err)
		return nil, err
	}

	// Add to data store and teardown data so it is removed during teardown
	t.Data.AddAccountSubscriptionRec(accountSubscriptionRec)
	t.teardownData.AddAccountSubscriptionRec(accountSubscriptionRec)

	return accountSubscriptionRec, nil
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
