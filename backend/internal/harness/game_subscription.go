package harness

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullint32"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameSubscriptions() error {
	l := t.Logger("createGameSubscriptions")

	for i := range t.DataConfig.AccountUserGameSubscriptionConfigs {
		subConfig := &t.DataConfig.AccountUserGameSubscriptionConfigs[i]

		var accountUserRec *account_record.AccountUser
		if subConfig.AccountUserRef != "" {
			var err error
			accountUserRec, err = t.Data.GetAccountUserRecByRef(subConfig.AccountUserRef)
			if err != nil {
				l.Warn("failed getting account user by ref >%s< >%v<", subConfig.AccountUserRef, err)
				return err
			}
		} else if subConfig.Record != nil && subConfig.Record.AccountUserID != "" {
			// Account IDs pre-populated by the caller (e.g. demo CLI loader). Build a
			// minimal AccountUser so createGameSubscriptionRec can set rec.AccountID/AccountUserID.
			accountUserRec = &account_record.AccountUser{}
			accountUserRec.ID = subConfig.Record.AccountUserID
			accountUserRec.AccountID = subConfig.Record.AccountID
		} else {
			return fmt.Errorf("subscription config at index >%d< must have AccountUserRef or pre-populated Record.AccountUserID", i)
		}

		_, err := t.createGameSubscriptionRec(*subConfig, accountUserRec)
		if err != nil {
			l.Warn("failed creating game subscription >%v<", err)
			return err
		}
	}

	l.Info("created >%d< game subscription records", len(t.DataConfig.AccountUserGameSubscriptionConfigs))
	return nil
}

func (t *Testing) createGameSubscriptionRec(subscriptionConfig AccountUserGameSubscriptionConfig, accountRec *account_record.AccountUser) (*game_record.GameSubscription, error) {
	l := t.Logger("createGameSubscriptionRec")

	var rec *game_record.GameSubscription
	if subscriptionConfig.Record != nil {
		recCopy := *subscriptionConfig.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameSubscription{}
	}

	rec = t.applyGameSubscriptionRecDefaultValues(rec)

	// Set subscription type if provided
	if subscriptionConfig.SubscriptionType != "" {
		rec.SubscriptionType = subscriptionConfig.SubscriptionType
	}

	// Resolve game ID from GameRef (games are created before subscriptions)
	if subscriptionConfig.GameRef != "" {
		gameID, ok := t.Data.Refs.GameRefs[subscriptionConfig.GameRef]
		if !ok {
			l.Warn("game ref >%s< not found in refs", subscriptionConfig.GameRef)
			return nil, fmt.Errorf("game ref >%s< not found", subscriptionConfig.GameRef)
		}
		rec.GameID = gameID
	}

	// All subscription types are tied to an account user (designer/manager/player).
	rec.AccountID = accountRec.AccountID
	rec.AccountUserID = accountRec.ID

	// Set instance limit if provided
	if subscriptionConfig.InstanceLimit != nil {
		rec.InstanceLimit = nullint32.FromInt32(*subscriptionConfig.InstanceLimit)
	}

	// For player subscriptions, set account_contact_id if not already set
	if rec.SubscriptionType == game_record.GameSubscriptionTypePlayer {
		if !rec.AccountUserContactID.Valid || rec.AccountUserContactID.String == "" {
			accountUserContactRec, err := t.Data.GetAccountUserContactRecByAccountUserID(accountRec.ID)
			if err != nil {
				l.Warn("failed getting account contact for account ID >%s<: %v", accountRec.ID, err)
				return nil, err
			}
			rec.AccountUserContactID = nullstring.FromString(accountUserContactRec.ID)
		}
	}

	// Create record
	l.Debug("creating game_subscription record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameSubscriptionRec(rec)
	if err != nil {
		l.Warn("failed creating game_subscription record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameSubscriptionRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameSubscriptionRec(rec)

	// Add to references store
	if subscriptionConfig.Reference != "" {
		t.Data.Refs.GameSubscriptionRefs[subscriptionConfig.Reference] = rec.ID
	}

	// Process join game subscription if scan data is provided
	if subscriptionConfig.JoinGameScanData != nil {
		ctx := context.Background()
		_, err = t.processJoinGameSubscriptionInSetup(ctx, subscriptionConfig.Reference, subscriptionConfig.JoinGameScanData)
		if err != nil {
			l.Warn("failed processing join game subscription >%v<", err)
			return nil, fmt.Errorf("failed processing join game subscription: %w", err)
		}
		l.Debug("processed join game subscription for subscription >%s<", subscriptionConfig.Reference)
	}

	return rec, nil
}

// createGameSubscriptionInstanceRec creates a GameSubscriptionInstance link (e.g. manager subscription to instance) and adds to teardown.
func (t *Testing) createGameSubscriptionInstanceRec(subscriptionRec *game_record.GameSubscription, instanceRec *game_record.GameInstance) (*game_record.GameSubscriptionInstance, error) {
	l := t.Logger("createGameSubscriptionInstanceRec")

	linkRec := &game_record.GameSubscriptionInstance{
		AccountID:          subscriptionRec.AccountID,
		AccountUserID:      subscriptionRec.AccountUserID,
		GameSubscriptionID: subscriptionRec.ID,
		GameInstanceID:     instanceRec.ID,
	}

	createdRec, err := t.Domain.(*domain.Domain).CreateGameSubscriptionInstanceRec(linkRec)
	if err != nil {
		l.Warn("failed creating game_subscription_instance record >%v<", err)
		return nil, err
	}

	t.Data.AddGameSubscriptionInstanceRec(createdRec)
	t.teardownData.AddGameSubscriptionInstanceRec(createdRec)
	return createdRec, nil
}

func (t *Testing) applyGameSubscriptionRecDefaultValues(rec *game_record.GameSubscription) *game_record.GameSubscription {
	if rec == nil {
		rec = &game_record.GameSubscription{}
	}
	if rec.SubscriptionType == "" {
		rec.SubscriptionType = game_record.GameSubscriptionTypePlayer
	}
	if rec.Status == "" {
		rec.Status = game_record.GameSubscriptionStatusActive
	}
	// Player subscriptions require a delivery method.
	if rec.SubscriptionType == game_record.GameSubscriptionTypePlayer && !rec.DeliveryMethod.Valid {
		rec.DeliveryMethod = nullstring.FromString(game_record.GameSubscriptionDeliveryMethodEmail)
	}
	return rec
}
