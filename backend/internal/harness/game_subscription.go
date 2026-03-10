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

// createGameSubscriptionRecOnly creates the subscription record without linking
// instances. Use this when instances don't exist yet; call linkSubscriptionInstances later.
func (t *Testing) createGameSubscriptionRecOnly(subscriptionConfig GameSubscriptionConfig, accountRec *account_record.AccountUser) (*game_record.GameSubscription, error) {
	stripped := subscriptionConfig
	stripped.GameInstanceRefs = nil
	return t.createGameSubscriptionRec(stripped, accountRec)
}

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

	// Link game instances if refs are provided and instances already exist
	if err := t.linkSubscriptionInstances(subscriptionConfig, rec, gameRec); err != nil {
		return nil, err
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

// linkSubscriptionInstances creates game_subscription_instance records for each
// GameInstanceRef in the config. It is a no-op when GameInstanceRefs is empty.
// Called from createGameSubscriptionRec when instances exist, or deferred to a
// later pass when subscriptions must be created before instances.
func (t *Testing) linkSubscriptionInstances(subscriptionConfig GameSubscriptionConfig, gameSubscriptionRec *game_record.GameSubscription, gameRec *game_record.Game) error {
	l := t.Logger("linkSubscriptionInstances")

	for _, instanceRef := range subscriptionConfig.GameInstanceRefs {
		gameInstanceRec, err := t.Data.GetGameInstanceRecByRef(instanceRef)
		if err != nil {
			l.Warn("failed resolving game instance ref >%s<: %v", instanceRef, err)
			return err
		}

		if gameInstanceRec.GameID != gameRec.ID {
			return fmt.Errorf("game_instance >%s< does not belong to game >%s<", instanceRef, subscriptionConfig.GameRef)
		}

		instanceLinkRec := &game_record.GameSubscriptionInstance{
			AccountID:          gameSubscriptionRec.AccountID,
			AccountUserID:      gameSubscriptionRec.AccountUserID,
			GameSubscriptionID: gameSubscriptionRec.ID,
			GameInstanceID:     gameInstanceRec.ID,
		}
		instanceLinkRec, err = t.Domain.(*domain.Domain).CreateGameSubscriptionInstanceRec(instanceLinkRec)
		if err != nil {
			l.Warn("failed linking game instance >%s< to subscription >%s<: %v", instanceRef, gameSubscriptionRec.ID, err)
			return err
		}

		t.Data.AddGameSubscriptionInstanceRec(instanceLinkRec)
		t.teardownData.AddGameSubscriptionInstanceRec(instanceLinkRec)
	}

	return nil
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
