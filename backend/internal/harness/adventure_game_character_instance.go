package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameCharacterInstanceRec(cfg AdventureGameCharacterInstanceConfig, gameInstanceRec *game_record.GameInstance) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := t.Logger("createAdventureGameCharacterInstanceRec")

	if cfg.GameCharacterRef == "" {
		return nil, fmt.Errorf("game_character_instance record must have a GameCharacterRef")
	}
	if cfg.GameLocationRef == "" {
		return nil, fmt.Errorf("game_character_instance record must have a GameLocationRef")
	}

	var rec *adventure_game_record.AdventureGameCharacterInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameCharacterInstance{}
	}

	rec = t.applyAdventureGameCharacterInstanceRecDefaultValues(rec)

	rec.GameID = gameInstanceRec.GameID
	rec.GameInstanceID = gameInstanceRec.ID

	// Resolve foreign keys
	characterRec, err := t.Data.GetAdventureGameCharacterRecByRef(cfg.GameCharacterRef)
	if err != nil {
		l.Error("could not resolve GameCharacterRef >%s< to a valid game character ID", cfg.GameCharacterRef)
		return nil, fmt.Errorf("could not resolve GameCharacterRef >%s< to a valid game character ID", cfg.GameCharacterRef)
	}
	rec.AdventureGameCharacterID = characterRec.ID

	locationInstanceRec, err := t.Data.GetAdventureGameLocationInstanceRecByLocationRef(cfg.GameLocationRef)
	if err != nil {
		l.Error("could not resolve GameLocationRef >%s< to a valid game location instance ID", cfg.GameLocationRef)
		return nil, fmt.Errorf("could not resolve GameLocationRef >%s< to a valid game location instance ID", cfg.GameLocationRef)
	}
	rec.AdventureGameLocationInstanceID = locationInstanceRec.ID

	// Create record
	l.Debug("creating game_character_instance record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameCharacterInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_character_instance record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddAdventureGameCharacterInstanceRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddAdventureGameCharacterInstanceRec(createdRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameCharacterInstanceRefs[cfg.Reference] = createdRec.ID
	}

	return createdRec, nil
}

func (t *Testing) applyAdventureGameCharacterInstanceRecDefaultValues(rec *adventure_game_record.AdventureGameCharacterInstance) *adventure_game_record.AdventureGameCharacterInstance {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameCharacterInstance{}
	}
	rec.Health = 100
	if rec.InventoryCapacity == 0 {
		rec.InventoryCapacity = 10
	}
	return rec
}
