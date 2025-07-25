package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameInstanceRec(cfg GameInstanceConfig, gameRec *record.Game) (*record.AdventureGameInstance, error) {
	l := t.Logger("createGameInstanceRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_instance record >%#v<", cfg)
	}

	var rec *record.AdventureGameInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &record.AdventureGameInstance{}
	}

	rec = t.applyGameInstanceRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	l.Info("creating game_instance record >%#v<", rec)

	// Create record
	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameInstanceRec(rec)
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

func (t *Testing) applyGameInstanceRecDefaultValues(rec *record.AdventureGameInstance) *record.AdventureGameInstance {
	if rec == nil {
		rec = &record.AdventureGameInstance{}
	}
	return rec
}
