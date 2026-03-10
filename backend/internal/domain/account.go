package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

// Creates or upodates an account, account user, account user contact and basic subscriptions records for users
// registering to play a game or logging in for the first time.
func (m *Domain) UpsertAccount(
	accountRec *account_record.Account,
	accountUserRec *account_record.AccountUser,
	accountUserContactRec *account_record.AccountUserContact,
) (*account_record.Account, *account_record.AccountUser, *account_record.AccountUserContact, []*account_record.AccountSubscription, error) {

	l := m.Logger("UpsertAccount")

	existingAccountUserRec, err := m.GetAccountUserRecByEmail(accountUserRec.Email)
	if err != nil {
		l.Warn("failed to get account user by email >%v<", err)
		return nil, nil, nil, nil, err
	}

	// When we have an existing account user record, we need to update the account, account user, and account user contact records.
	if existingAccountUserRec != nil {
		// Get the existing account record
		existingAccountRec, err := m.GetAccountRec(existingAccountUserRec.AccountID, coresql.ForUpdate)
		if err != nil {
			l.Warn("failed to get account record >%v<", err)
			return nil, nil, nil, nil, err
		}

		// Update the existing account record
		existingAccountRec.Name = accountRec.Name
		if _, err = m.UpdateAccountRec(existingAccountRec); err != nil {
			l.Warn("failed to update account >%v<", err)
			return nil, nil, nil, nil, err
		}

		// Update the existing account user record
		accountUserRec.ID = existingAccountUserRec.ID
		accountUserRec.CreatedAt = existingAccountUserRec.CreatedAt
		if accountUserRec, err = m.UpdateAccountUserRec(accountUserRec); err != nil {
			l.Warn("failed to update account user >%v<", err)
			return nil, nil, nil, nil, err
		}

		// Get the existing account user contact record
		existingAccountUserContactRec, err := m.GetAccountUserContactRecByAccountUserID(existingAccountUserRec.ID, coresql.ForUpdate)
		if err != nil {
			l.Warn("failed to get account user contact record >%v<", err)
			return nil, nil, nil, nil, err
		}

		// Update the existing account user contact record
		if existingAccountUserContactRec != nil {
			if accountUserContactRec != nil {
				existingAccountUserContactRec.Name = accountUserContactRec.Name
				existingAccountUserContactRec.PostalAddressLine1 = accountUserContactRec.PostalAddressLine1
				existingAccountUserContactRec.PostalAddressLine2 = accountUserContactRec.PostalAddressLine2
				existingAccountUserContactRec.StateProvince = accountUserContactRec.StateProvince
				existingAccountUserContactRec.Country = accountUserContactRec.Country
				existingAccountUserContactRec.PostalCode = accountUserContactRec.PostalCode
			}
			if _, err = m.UpdateAccountUserContactRec(existingAccountUserContactRec); err != nil {
				l.Warn("failed to update account user contact >%v<", err)
				return nil, nil, nil, nil, err
			}
		} else {
			if accountUserContactRec == nil {
				accountUserContactRec = &account_record.AccountUserContact{}
			}
			accountUserContactRec.AccountUserID = existingAccountUserRec.ID
			if accountUserContactRec, err = m.CreateAccountUserContactRec(accountUserContactRec); err != nil {
				l.Warn("failed to create account user contact >%v<", err)
				return nil, nil, nil, nil, err
			}
		}
	} else {
		// Create a new account, account user, and account user contact records
		accountRec, err = m.CreateAccountRec(accountRec)
		if err != nil {
			l.Warn("failed to create account >%v<", err)
			return nil, nil, nil, nil, err
		}

		accountUserRec.AccountID = accountRec.ID
		accountUserRec, err = m.CreateAccountUserRec(accountUserRec)
		if err != nil {
			l.Warn("failed to create account user >%v<", err)
			return nil, nil, nil, nil, err
		}

		contactRec := accountUserContactRec
		if contactRec == nil {
			contactRec = &account_record.AccountUserContact{}
		}
		contactRec.AccountUserID = accountUserRec.ID
		accountUserContactRec, err = m.CreateAccountUserContactRec(contactRec)
		if err != nil {
			l.Warn("failed to create account user contact >%v<", err)
			return nil, nil, nil, nil, err
		}
	}

	// We now need to verify the user has a basic deisgner, manager and player subscriptions.
	designerAccountSubscriptionRec, err := m.GetAccountSubscriptionRecByAccountUserID(accountUserRec.ID, account_record.AccountSubscriptionTypeBasicGameDesigner)
	if err != nil {
		l.Warn("failed to get basic game designer subscription >%v<", err)
		return nil, nil, nil, nil, err
	}

	if designerAccountSubscriptionRec == nil {
		l.Warn("no basic game designer subscription found for account user >%s<", accountUserRec.ID)

		designerAccountSubscriptionRec, err = m.CreateAccountSubscriptionRec(&account_record.AccountSubscription{
			AccountID:          nullstring.FromString(accountRec.ID),
			AccountUserID:      nullstring.FromString(accountUserRec.ID),
			SubscriptionType:   account_record.AccountSubscriptionTypeBasicGameDesigner,
			SubscriptionPeriod: account_record.AccountSubscriptionPeriodEternal,
			Status:             account_record.AccountSubscriptionStatusActive,
			AutoRenew:          true,
		})
		if err != nil {
			l.Warn("failed to create basic game designer subscription >%v<", err)
			return nil, nil, nil, nil, err
		}
	}

	managerAccountSubscriptionRec, err := m.GetAccountSubscriptionRecByAccountUserID(accountUserRec.ID, account_record.AccountSubscriptionTypeBasicManager)
	if err != nil {
		l.Warn("failed to get basic manager subscription >%v<", err)
		return nil, nil, nil, nil, err
	}

	if managerAccountSubscriptionRec == nil {
		l.Warn("no basic game manager subscription found for account user >%s<", accountUserRec.ID)

		managerAccountSubscriptionRec, err = m.CreateAccountSubscriptionRec(&account_record.AccountSubscription{
			AccountID:          nullstring.FromString(accountRec.ID),
			AccountUserID:      nullstring.FromString(accountUserRec.ID),
			SubscriptionType:   account_record.AccountSubscriptionTypeBasicManager,
			SubscriptionPeriod: account_record.AccountSubscriptionPeriodEternal,
			Status:             account_record.AccountSubscriptionStatusActive,
			AutoRenew:          true,
		})
		if err != nil {
			l.Warn("failed to create basic game manager subscription >%v<", err)
			return nil, nil, nil, nil, err
		}
	}

	playerAccountSubscriptionRec, err := m.GetAccountSubscriptionRecByAccountUserID(accountUserRec.ID, account_record.AccountSubscriptionTypeBasicPlayer)
	if err != nil {
		l.Warn("failed to get basic player subscription >%v<", err)
		return nil, nil, nil, nil, err
	}

	if playerAccountSubscriptionRec == nil {
		l.Warn("no basic player subscription found for account user >%s<", accountUserRec.ID)

		playerAccountSubscriptionRec, err = m.CreateAccountSubscriptionRec(&account_record.AccountSubscription{
			AccountID:          nullstring.FromString(accountRec.ID),
			AccountUserID:      nullstring.FromString(accountUserRec.ID),
			SubscriptionType:   account_record.AccountSubscriptionTypeBasicPlayer,
			SubscriptionPeriod: account_record.AccountSubscriptionPeriodEternal,
			Status:             account_record.AccountSubscriptionStatusActive,
			AutoRenew:          true,
		})
		if err != nil {
			l.Warn("failed to create basic player subscription >%v<", err)
			return nil, nil, nil, nil, err
		}
	}

	return accountRec, accountUserRec, accountUserContactRec, []*account_record.AccountSubscription{designerAccountSubscriptionRec, managerAccountSubscriptionRec, playerAccountSubscriptionRec}, nil
}

// GetAccountRec returns the account (parent/tenant) record by ID.
func (m *Domain) GetAccountRec(recID string, lock *coresql.Lock) (*account_record.Account, error) {
	l := m.Logger("GetAccountRec")

	l.Debug("getting account record ID >%s<", recID)

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

// GetManyAccountRecs returns account (parent/tenant) records.
func (m *Domain) GetManyAccountRecs(opts *coresql.Options) ([]*account_record.Account, error) {
	l := m.Logger("GetManyAccountRecs")

	l.Info("getting many account records opts >%#v<", opts)

	r := m.AccountRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// CreateAccountRec creates a new account (parent/tenant) record.
func (m *Domain) CreateAccountRec(rec *account_record.Account) (*account_record.Account, error) {
	l := m.Logger("CreateAccountRec")

	l.Debug("creating account record >%#v<", rec)

	if rec != nil && rec.Status == "" {
		rec.Status = account_record.AccountStatusActive
	}

	r := m.AccountRepository()

	if err := m.validateAccountRecForCreate(rec); err != nil {
		l.Warn("failed to validate account record >%v<", err)
		return rec, err
	}

	createdRec, err := r.CreateOne(rec)
	if err != nil {
		l.Warn("failed to create account record >%v<", err)
		return nil, databaseError(err)
	}

	return createdRec, nil
}

// UpdateAccountRec updates an account (parent/tenant) record.
func (m *Domain) UpdateAccountRec(rec *account_record.Account) (*account_record.Account, error) {
	l := m.Logger("UpdateAccountRec")

	l.Debug("updating account record ID >%s<", rec.ID)

	curr, err := m.GetAccountRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	rec.CreatedAt = curr.CreatedAt

	if rec.Status == "" {
		rec.Status = curr.Status
	}

	r := m.AccountRepository()

	if err := m.validateAccountRecForUpdate(curr, rec); err != nil {
		l.Warn("failed to validate account record >%v<", err)
		return rec, err
	}

	rec, err = r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
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

// RemoveAccountRec hard-deletes an account and all dependent records in FK order.
func (m *Domain) RemoveAccountRec(recID string) error {
	l := m.Logger("RemoveAccountRec")

	l.Debug("removing account record ID >%s< and all dependents", recID)

	accountFilter := &coresql.Options{
		Params: []coresql.Param{
			{Col: "account_id", Val: recID},
		},
	}

	// 1. game_subscription_instance (references account_id, game_subscription_id)
	gameSubscriptionInstanceRecs, err := m.GetManyGameSubscriptionInstanceRecs(accountFilter)
	if err != nil {
		return databaseError(err)
	}
	for _, gameSubscriptionInstanceRec := range gameSubscriptionInstanceRecs {
		if err := m.RemoveGameSubscriptionInstanceRec(gameSubscriptionInstanceRec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 2. game_subscription (references account_id)
	gameSubscriptionRecs, err := m.GetManyGameSubscriptionRecs(accountFilter)
	if err != nil {
		return databaseError(err)
	}
	for _, rec := range gameSubscriptionRecs {
		if err := m.RemoveGameSubscriptionRec(rec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 3. account_subscription (by account_id for designer/manager subs, by account_user_id for player subs)
	accountSubscriptionRecs, err := m.GetManyAccountSubscriptionRecs(accountFilter)
	if err != nil {
		return databaseError(err)
	}
	for _, rec := range accountSubscriptionRecs {
		if err := m.RemoveAccountSubscriptionRec(rec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 4. account_user (needed for player subscription and contact cleanup)
	accountUserRecs, err := m.GetManyAccountUserRecs(accountFilter)
	if err != nil {
		return databaseError(err)
	}

	// 5. account_user
	for _, accountUserRec := range accountUserRecs {
		if err := m.RemoveAccountUserRec(accountUserRec.ID); err != nil {
			return databaseError(err)
		}
	}

	// 6. account
	r := m.AccountRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	l.Info("removed account >%s< and all dependents (%d users, %d subscriptions, %d game subscriptions, %d game subscription instances)",
		recID, len(accountUserRecs), len(accountSubscriptionRecs), len(gameSubscriptionRecs), len(gameSubscriptionInstanceRecs))

	return nil
}
