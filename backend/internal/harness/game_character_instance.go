package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameCharacterInstanceRec(cfg GameCharacterInstanceConfig, gameInstanceRec *record.AdventureGameInstance) (*record.AdventureGameCharacterInstance, error) {
	l := t.Logger("createGameCharacterInstanceRec")

	if cfg.GameCharacterRef == "" {
		return nil, fmt.Errorf("game_character_instance record must have a GameCharacterRef")
	}
	if cfg.GameLocationRef == "" {
		return nil, fmt.Errorf("game_character_instance record must have a GameLocationRef")
	}

	var rec *record.AdventureGameCharacterInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &record.AdventureGameCharacterInstance{}
	}

	rec = t.applyGameCharacterInstanceRecDefaultValues(rec)

	rec.GameID = gameInstanceRec.GameID
	rec.AdventureGameInstanceID = gameInstanceRec.ID

	// Resolve foreign keys
	characterRec, err := t.Data.GetGameCharacterRecByRef(cfg.GameCharacterRef)
	if err != nil {
		l.Error("could not resolve GameCharacterRef >%s< to a valid game character ID", cfg.GameCharacterRef)
		return nil, fmt.Errorf("could not resolve GameCharacterRef >%s< to a valid game character ID", cfg.GameCharacterRef)
	}
	rec.AdventureGameCharacterID = characterRec.ID

	locationInstanceRec, err := t.Data.GetGameLocationInstanceRecByLocationRef(cfg.GameLocationRef)
	if err != nil {
		l.Error("could not resolve GameLocationRef >%s< to a valid game location instance ID", cfg.GameLocationRef)
		return nil, fmt.Errorf("could not resolve GameLocationRef >%s< to a valid game location instance ID", cfg.GameLocationRef)
	}
	rec.AdventureGameLocationInstanceID = locationInstanceRec.ID

	// Create record
	l.Info("creating game_character_instance record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameCharacterInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_character_instance record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameCharacterInstanceRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddGameCharacterInstanceRec(createdRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.GameCharacterInstanceRefs[cfg.Reference] = createdRec.ID
	}
	return createdRec, nil
}

func (t *Testing) applyGameCharacterInstanceRecDefaultValues(rec *record.AdventureGameCharacterInstance) *record.AdventureGameCharacterInstance {
	if rec == nil {
		rec = &record.AdventureGameCharacterInstance{}
	}
	rec.Health = 100
	return rec
}
