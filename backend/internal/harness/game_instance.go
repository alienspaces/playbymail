package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameInstanceRec(cfg GameInstanceConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameInstance, error) {
	l := t.Logger("createGameInstanceRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_instance record >%#v<", cfg)
	}

	var rec *adventure_game_record.AdventureGameInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameInstance{}
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

func (t *Testing) applyGameInstanceRecDefaultValues(rec *adventure_game_record.AdventureGameInstance) *adventure_game_record.AdventureGameInstance {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameInstance{}
	}

	// Set default status if not already set
	if rec.Status == "" {
		rec.Status = adventure_game_record.GameInstanceStatusCreated
	}

	// Set default turn deadline if not already set
	if rec.TurnDeadlineHours == 0 {
		rec.TurnDeadlineHours = 168 // 7 days default
	}

	return rec
}
