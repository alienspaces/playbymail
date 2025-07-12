package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameCreatureInstanceRec(cfg GameCreatureInstanceConfig, gameInstanceRec *record.GameInstance) (*record.GameCreatureInstance, error) {
	l := t.Logger("createGameCreatureInstanceRec")

	if cfg.GameCreatureRef == "" {
		return nil, fmt.Errorf("game_creature_instance record must have a GameCreatureRef")
	}
	if cfg.GameLocationRef == "" {
		return nil, fmt.Errorf("game_creature_instance record must have a GameLocationRef")
	}

	var rec *record.GameCreatureInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &record.GameCreatureInstance{}
	}

	rec = t.applyGameCreatureInstanceRecDefaultValues(rec)

	rec.GameID = gameInstanceRec.GameID
	rec.GameInstanceID = gameInstanceRec.ID

	// Resolve foreign keys
	creatureRec, err := t.Data.GetGameCreatureRecByRef(cfg.GameCreatureRef)
	if err != nil {
		l.Error("could not resolve GameCreatureRef >%s< to a valid game creature ID", cfg.GameCreatureRef)
		return nil, fmt.Errorf("could not resolve GameCreatureRef >%s< to a valid game creature ID", cfg.GameCreatureRef)
	}
	rec.GameCreatureID = creatureRec.ID

	locationInstanceRec, err := t.Data.GetGameLocationInstanceRecByLocationRef(cfg.GameLocationRef)
	if err != nil {
		l.Error("could not resolve GameLocationRef >%s< to a valid game location instance ID", cfg.GameLocationRef)
		return nil, fmt.Errorf("could not resolve GameLocationRef >%s< to a valid game location instance ID", cfg.GameLocationRef)
	}
	rec.GameLocationInstanceID = locationInstanceRec.ID

	// Create record
	l.Info("creating game_creature_instance record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateGameCreatureInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_creature_instance record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameCreatureInstanceRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddGameCreatureInstanceRec(createdRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.GameCreatureInstanceRefs[cfg.Reference] = createdRec.ID
	}
	return createdRec, nil
}

func (t *Testing) applyGameCreatureInstanceRecDefaultValues(rec *record.GameCreatureInstance) *record.GameCreatureInstance {
	if rec == nil {
		rec = &record.GameCreatureInstance{}
	}
	rec.IsAlive = true
	return rec
}
