package harness

import (
	"fmt"
	"time"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameInstanceRec(cfg GameInstanceConfig, gameRec *game_record.Game) (*game_record.GameInstance, error) {
	l := t.Logger("createGameInstanceRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_instance record >%#v<", cfg)
	}

	var rec *game_record.GameInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameInstance{}
	}

	rec = t.applyGameInstanceRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Set game_subscription_id from the first Manager subscription for this game
	// or fall back to the first subscription of any type
	if rec.GameSubscriptionID == "" {
		for _, subRec := range t.Data.GameSubscriptionRecs {
			if subRec.GameID == gameRec.ID {
				if subRec.SubscriptionType == game_record.GameSubscriptionTypeManager {
					rec.GameSubscriptionID = subRec.ID
					break
				}
				// Fall back to first subscription if no Manager found yet
				if rec.GameSubscriptionID == "" {
					rec.GameSubscriptionID = subRec.ID
				}
			}
		}
	}

	if rec.GameSubscriptionID == "" {
		return nil, fmt.Errorf("no subscription found for game >%s<, cannot create game instance without game_subscription_id", gameRec.ID)
	}

	l.Debug("creating game_instance record >%#v<", rec)

	// Create record
	createdRec, err := t.Domain.(*domain.Domain).CreateGameInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_instance record >%v<", err)
		return rec, err
	}

	// Add to data store
	t.Data.AddGameInstanceRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddGameInstanceRec(createdRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.GameInstanceRefs[cfg.Reference] = createdRec.ID
	}

	return createdRec, nil
}

func (t *Testing) applyGameInstanceRecDefaultValues(rec *game_record.GameInstance) *game_record.GameInstance {
	if rec == nil {
		rec = &game_record.GameInstance{}
	}

	// Set default status if not already set
	if rec.Status == "" {
		rec.Status = game_record.GameInstanceStatusCreated
	}

	// Set timestamps if not already set
	now := time.Now()
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = now
	}
	if !rec.UpdatedAt.Valid {
		rec.UpdatedAt = nulltime.FromTime(now)
	}

	return rec
}
