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

func (t *Testing) createGameSubscriptionRec(subscriptionConfig GameSubscriptionConfig, accountRec *account_record.AccountUser) (*game_record.GameSubscription, error) {
	l := t.Logger("createGameSubscriptionRec")

	if accountRec == nil {
		return nil, fmt.Errorf("account record is nil for game_subscription record >%#v<", subscriptionConfig)
	}

	if subscriptionConfig.GameRef == "" {
		return nil, fmt.Errorf("game_subscription record >%#v< must have a GameRef set", subscriptionConfig)
	}

	// Get game record
	gameRec, err := t.Data.GetGameRecByRef(subscriptionConfig.GameRef)
	if err != nil {
		l.Warn("failed resolving game ref >%s<: %v", subscriptionConfig.GameRef, err)
		return nil, err
	}

	var rec *game_record.GameSubscription
	if subscriptionConfig.Record != nil {
		recCopy := *subscriptionConfig.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameSubscription{}
	}

	rec = t.applyGameSubscriptionRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Set subscription type if provided
	if subscriptionConfig.SubscriptionType != "" {
		rec.SubscriptionType = subscriptionConfig.SubscriptionType
	}

	if rec.SubscriptionType == game_record.GameSubscriptionTypeDesigner || rec.SubscriptionType == game_record.GameSubscriptionTypeManager {
		rec.AccountID = accountRec.AccountID
		rec.AccountUserID = nullstring.FromString("")
	} else {
		// Player subscription
		rec.AccountID = accountRec.AccountID
		// User Subscriptions use account_user_id.
		rec.AccountUserID = nullstring.FromString(accountRec.ID)
	}

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

	rec, err = t.Domain.(*domain.Domain).CreateGameSubscriptionRec(rec)
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

	// Link game instances if provided
	if len(subscriptionConfig.GameInstanceRefs) > 0 {
		for _, instanceRef := range subscriptionConfig.GameInstanceRefs {
			gameInstanceRec, err := t.Data.GetGameInstanceRecByRef(instanceRef)
			if err != nil {
				l.Warn("failed resolving game instance ref >%s<: %v", instanceRef, err)
				return nil, err
			}

			// Validate that game instance belongs to the same game
			if gameInstanceRec.GameID != gameRec.ID {
				return nil, fmt.Errorf("game_instance >%s< does not belong to game >%s<", instanceRef, subscriptionConfig.GameRef)
			}

			// Link instance to subscription
			// account_id will be derived from subscription in validation
			instanceLinkRec := &game_record.GameSubscriptionInstance{
				GameSubscriptionID: rec.ID,
				GameInstanceID:     gameInstanceRec.ID,
			}
			instanceLinkRec, err = t.Domain.(*domain.Domain).CreateGameSubscriptionInstanceRec(instanceLinkRec)
			if err != nil {
				l.Warn("failed linking game instance >%s< to subscription >%s<: %v", instanceRef, rec.ID, err)
				return nil, err
			}

			// Add to data store
			t.Data.AddGameSubscriptionInstanceRec(instanceLinkRec)
			t.teardownData.AddGameSubscriptionInstanceRec(instanceLinkRec)
		}
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

func (t *Testing) applyGameSubscriptionRecDefaultValues(rec *game_record.GameSubscription) *game_record.GameSubscription {
	if rec == nil {
		rec = &game_record.GameSubscription{}
	}
	if rec.SubscriptionType == "" {
		rec.SubscriptionType = game_record.GameSubscriptionTypePlayer // Default subscription type
	}
	if rec.Status == "" {
		rec.Status = game_record.GameSubscriptionStatusActive
	}
	return rec
}
