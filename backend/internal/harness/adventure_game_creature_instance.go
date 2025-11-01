package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameCreatureInstanceRec(cfg AdventureGameCreatureInstanceConfig, gameInstanceRec *game_record.GameInstance) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	l := t.Logger("createAdventureGameCreatureInstanceRec")

	if cfg.GameCreatureRef == "" {
		return nil, fmt.Errorf("adventure game creature instance record must have a GameCreatureRef")
	}
	if cfg.GameLocationRef == "" {
		return nil, fmt.Errorf("adventure game creature instance record must have a GameLocationRef")
	}

	var rec *adventure_game_record.AdventureGameCreatureInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameCreatureInstance{}
	}

	rec = t.applyAdventureGameCreatureInstanceRecDefaultValues(rec)

	rec.GameID = gameInstanceRec.GameID
	rec.GameInstanceID = gameInstanceRec.ID

	// Resolve foreign keys
	creatureRec, err := t.Data.GetAdventureGameCreatureRecByRef(cfg.GameCreatureRef)
	if err != nil {
		l.Error("could not resolve GameCreatureRef >%s< to a valid game creature ID", cfg.GameCreatureRef)
		return nil, fmt.Errorf("could not resolve GameCreatureRef >%s< to a valid game creature ID", cfg.GameCreatureRef)
	}
	rec.AdventureGameCreatureID = creatureRec.ID

	locationInstanceRec, err := t.Data.GetAdventureGameLocationInstanceRecByLocationRef(cfg.GameLocationRef)
	if err != nil {
		l.Error("could not resolve GameLocationRef >%s< to a valid game location instance ID", cfg.GameLocationRef)
		return nil, fmt.Errorf("could not resolve GameLocationRef >%s< to a valid game location instance ID", cfg.GameLocationRef)
	}
	rec.AdventureGameLocationInstanceID = locationInstanceRec.ID

	// Create record
	l.Debug("creating game_creature_instance record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameCreatureInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_creature_instance record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddAdventureGameCreatureInstanceRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddAdventureGameCreatureInstanceRec(createdRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameCreatureInstanceRefs[cfg.Reference] = createdRec.ID
	}
	return createdRec, nil
}

func (t *Testing) applyAdventureGameCreatureInstanceRecDefaultValues(rec *adventure_game_record.AdventureGameCreatureInstance) *adventure_game_record.AdventureGameCreatureInstance {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameCreatureInstance{}
	}
	rec.Health = 100
	return rec
}
