package harness

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameSubscriptionRec(subscriptionConfig GameSubscriptionConfig, gameRec *game_record.Game) (*game_record.GameSubscription, error) {
	l := t.Logger("createGameSubscriptionRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_subscription record >%#v<", subscriptionConfig)
	}

	if subscriptionConfig.AccountRef == "" {
		return nil, fmt.Errorf("game_subscription record >%#v< must have an AccountRef set", subscriptionConfig)
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

	// Get account record
	accountRec, err := t.Data.GetAccountRecByRef(subscriptionConfig.AccountRef)
	if err != nil {
		l.Warn("failed resolving account ref >%s<: %v", subscriptionConfig.AccountRef, err)
		return nil, err
	}
	rec.AccountID = accountRec.ID

	// Set subscription type if provided
	if subscriptionConfig.SubscriptionType != "" {
		rec.SubscriptionType = subscriptionConfig.SubscriptionType
	}

	// For player subscriptions, set account_contact_id if not already set
	if rec.SubscriptionType == game_record.GameSubscriptionTypePlayer {
		if !rec.AccountContactID.Valid || rec.AccountContactID.String == "" {
			accountContactRec, err := t.Data.GetAccountContactRecByAccountID(accountRec.ID)
			if err != nil {
				l.Warn("failed getting account contact for account ID >%s<: %v", accountRec.ID, err)
				return nil, err
			}
			rec.AccountContactID = nullstring.FromString(accountContactRec.ID)
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
		rec.SubscriptionType = "Player" // Default subscription type
	}
	return rec
}
