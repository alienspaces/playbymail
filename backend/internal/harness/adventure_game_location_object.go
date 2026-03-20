package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameLocationObjectRec(cfg AdventureGameLocationObjectConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameLocationObject, error) {
	l := t.Logger("createAdventureGameLocationObjectRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for adventure game location object record >%#v<", cfg)
	}

	if cfg.LocationRef == "" {
		return nil, fmt.Errorf("location reference is required for adventure game location object record >%#v<", cfg)
	}

	var rec *adventure_game_record.AdventureGameLocationObject
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameLocationObject{}
	}

	rec.GameID = gameRec.ID

	locationRec, err := t.Data.GetAdventureGameLocationRecByRef(cfg.LocationRef)
	if err != nil {
		l.Error("could not resolve LocationRef >%s< to a valid location ID", cfg.LocationRef)
		return nil, fmt.Errorf("could not resolve LocationRef >%s< to a valid location ID", cfg.LocationRef)
	}
	rec.AdventureGameLocationID = locationRec.ID

	// 1. Create object without initial state so state records can reference the object ID.
	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationObjectRec(rec)
	if err != nil {
		l.Warn("failed creating adventure_game_location_object record >%v<", err)
		return nil, err
	}

	t.Data.AddAdventureGameLocationObjectRec(createdRec)
	t.teardownData.AddAdventureGameLocationObjectRec(createdRec)

	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameLocationObjectRefs[cfg.Reference] = createdRec.ID
	}

	l.Debug("created adventure_game_location_object record ID >%s<", createdRec.ID)

	// 2. Create state records now that the object ID is known.
	for _, stateCfg := range cfg.AdventureGameLocationObjectStateConfigs {
		if _, err := t.createAdventureGameLocationObjectStateRec(stateCfg, gameRec, createdRec.ID); err != nil {
			l.Warn("failed creating adventure_game_location_object_state record >%v<", err)
			return nil, err
		}
	}

	// 3. Update the object with its initial state now that states are registered.
	if cfg.InitialStateRef != "" {
		stateRec, err := t.Data.GetAdventureGameLocationObjectStateRecByRef(cfg.InitialStateRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve InitialStateRef >%s< for object: %w", cfg.InitialStateRef, err)
		}
		createdRec.InitialAdventureGameLocationObjectStateID = nullstring.FromString(stateRec.ID)
		updatedRec, err := t.Domain.(*domain.Domain).UpdateAdventureGameLocationObjectRec(createdRec)
		if err != nil {
			l.Warn("failed updating adventure_game_location_object record with initial state >%v<", err)
			return nil, err
		}
		createdRec = updatedRec
	}

	// 4. Create effect records.
	for i := range cfg.AdventureGameLocationObjectEffectConfigs {
		_, err := t.createAdventureGameLocationObjectEffectRec(cfg.AdventureGameLocationObjectEffectConfigs[i], createdRec)
		if err != nil {
			l.Warn("failed creating adventure_game_location_object_effect record >%v<", err)
			return nil, err
		}
	}

	return createdRec, nil
}

func (t *Testing) createAdventureGameLocationObjectStateRec(cfg AdventureGameLocationObjectStateConfig, gameRec *game_record.Game, objectID string) (*adventure_game_record.AdventureGameLocationObjectState, error) {
	l := t.Logger("createAdventureGameLocationObjectStateRec")

	var rec *adventure_game_record.AdventureGameLocationObjectState
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameLocationObjectState{}
	}

	rec.GameID = gameRec.ID
	if objectID != "" {
		rec.AdventureGameLocationObjectID = objectID
	}

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationObjectStateRec(rec)
	if err != nil {
		l.Warn("failed creating adventure_game_location_object_state record >%v<", err)
		return nil, err
	}

	t.Data.AddAdventureGameLocationObjectStateRec(createdRec)
	t.teardownData.AddAdventureGameLocationObjectStateRec(createdRec)

	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameLocationObjectStateRefs[cfg.Reference] = createdRec.ID
	}

	l.Debug("created adventure_game_location_object_state record ID >%s<", createdRec.ID)

	return createdRec, nil
}

func (t *Testing) createAdventureGameLocationObjectEffectRec(cfg AdventureGameLocationObjectEffectConfig, objectRec *adventure_game_record.AdventureGameLocationObject) (*adventure_game_record.AdventureGameLocationObjectEffect, error) {
	l := t.Logger("createAdventureGameLocationObjectEffectRec")

	if objectRec == nil {
		return nil, fmt.Errorf("object record is nil for adventure game location object effect record >%#v<", cfg)
	}

	var rec *adventure_game_record.AdventureGameLocationObjectEffect
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameLocationObjectEffect{}
	}

	rec.GameID = objectRec.GameID
	rec.AdventureGameLocationObjectID = objectRec.ID

	if cfg.RequiredStateRef != "" {
		stateRec, err := t.Data.GetAdventureGameLocationObjectStateRecByRef(cfg.RequiredStateRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve RequiredStateRef >%s<: %w", cfg.RequiredStateRef, err)
		}
		rec.RequiredAdventureGameLocationObjectStateID = nullstring.FromString(stateRec.ID)
	}

	if cfg.ResultStateRef != "" {
		stateRec, err := t.Data.GetAdventureGameLocationObjectStateRecByRef(cfg.ResultStateRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve ResultStateRef >%s<: %w", cfg.ResultStateRef, err)
		}
		rec.ResultAdventureGameLocationObjectStateID = nullstring.FromString(stateRec.ID)
	}

	if cfg.ResultObjectRef != "" {
		objRec, err := t.Data.GetAdventureGameLocationObjectRecByRef(cfg.ResultObjectRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve ResultObjectRef >%s<: %w", cfg.ResultObjectRef, err)
		}
		rec.ResultAdventureGameLocationObjectID = nullstring.FromString(objRec.ID)
	}

	if cfg.ResultItemRef != "" {
		itemRec, err := t.Data.GetAdventureGameItemRecByRef(cfg.ResultItemRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve ResultItemRef >%s<: %w", cfg.ResultItemRef, err)
		}
		rec.ResultAdventureGameItemID = nullstring.FromString(itemRec.ID)
	}

	if cfg.ResultLocationRef != "" {
		locRec, err := t.Data.GetAdventureGameLocationRecByRef(cfg.ResultLocationRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve ResultLocationRef >%s<: %w", cfg.ResultLocationRef, err)
		}
		rec.ResultAdventureGameLocationID = nullstring.FromString(locRec.ID)
	}

	if cfg.ResultCreatureRef != "" {
		creatureRec, err := t.Data.GetAdventureGameCreatureRecByRef(cfg.ResultCreatureRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve ResultCreatureRef >%s<: %w", cfg.ResultCreatureRef, err)
		}
		rec.ResultAdventureGameCreatureID = nullstring.FromString(creatureRec.ID)
	}

	if cfg.RequiredItemRef != "" {
		itemRec, err := t.Data.GetAdventureGameItemRecByRef(cfg.RequiredItemRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve RequiredItemRef >%s<: %w", cfg.RequiredItemRef, err)
		}
		rec.RequiredAdventureGameItemID = nullstring.FromString(itemRec.ID)
	}

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationObjectEffectRec(rec)
	if err != nil {
		l.Warn("failed creating adventure_game_location_object_effect record >%v<", err)
		return nil, err
	}

	t.Data.AddAdventureGameLocationObjectEffectRec(createdRec)
	t.teardownData.AddAdventureGameLocationObjectEffectRec(createdRec)

	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameLocationObjectEffectRefs[cfg.Reference] = createdRec.ID
	}

	l.Debug("created adventure_game_location_object_effect record ID >%s<", createdRec.ID)

	return createdRec, nil
}
