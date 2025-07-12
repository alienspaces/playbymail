package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameLocationInstanceRec(cfg GameLocationInstanceConfig, gameInstanceRec *record.GameInstance) (*record.GameLocationInstance, error) {
	l := t.Logger("createGameLocationInstanceRec")

	if gameInstanceRec == nil {
		return nil, fmt.Errorf("game instance record is nil for game_location_instance record >%#v<", cfg)
	}

	if cfg.GameLocationRef == "" {
		return nil, fmt.Errorf("game location reference is required for game_location_instance record >%#v<", cfg)
	}

	var rec *record.GameLocationInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &record.GameLocationInstance{}
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
	rec.GameLocationID = gameLocationRec.ID

	// Create record
	l.Info("creating game_location_instance record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateGameLocationInstanceRec(rec)
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

func (t *Testing) applyGameLocationInstanceRecDefaultValues(rec *record.GameLocationInstance) *record.GameLocationInstance {
	if rec == nil {
		rec = &record.GameLocationInstance{}
	}
	return rec
}
