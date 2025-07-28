package harness

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameLocationInstanceRec(cfg GameLocationInstanceConfig, gameInstanceRec *game_record.GameInstance) (*adventure_game_record.AdventureGameLocationInstance, error) {
	l := t.Logger("createGameLocationInstanceRec")

	if gameInstanceRec == nil {
		return nil, fmt.Errorf("game instance record is nil for game_location_instance record >%#v<", cfg)
	}

	if cfg.GameLocationRef == "" {
		return nil, fmt.Errorf("game location reference is required for game_location_instance record >%#v<", cfg)
	}

	var rec *adventure_game_record.AdventureGameLocationInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameLocationInstance{}
	}

	rec = t.applyGameLocationInstanceRecDefaultValues(rec)

	// Set game_id from parent game instance
	rec.GameID = gameInstanceRec.GameID
	rec.GameInstanceID = gameInstanceRec.ID

	// The game location is retrieved by reference
	gameLocationRec, err := t.Data.GetGameLocationRecByRef(cfg.GameLocationRef)
	if err != nil {
		l.Error("could not resolve GameLocationRef >%s< to a valid game location ID", cfg.GameLocationRef)
		return nil, fmt.Errorf("could not resolve GameLocationRef >%s< to a valid game location ID", cfg.GameLocationRef)
	}
	rec.AdventureGameLocationID = gameLocationRec.ID

	// Create record
	l.Info("creating game_location_instance record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_location_instance record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameLocationInstanceRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddGameLocationInstanceRec(createdRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.GameLocationInstanceRefs[cfg.Reference] = createdRec.ID
	}

	return createdRec, nil
}

func (t *Testing) applyGameLocationInstanceRecDefaultValues(rec *adventure_game_record.AdventureGameLocationInstance) *adventure_game_record.AdventureGameLocationInstance {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocationInstance{}
	}

	// Set timestamps if not already set
	now := time.Now()
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = now
	}
	if !rec.UpdatedAt.Valid {
		rec.UpdatedAt = sql.NullTime{Time: now, Valid: true}
	}

	return rec
}
