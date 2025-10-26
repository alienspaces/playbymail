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
