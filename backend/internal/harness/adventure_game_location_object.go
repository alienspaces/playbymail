package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
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

	if rec.InitialState == "" {
		rec.InitialState = "intact"
	}

	locationRec, err := t.Data.GetAdventureGameLocationRecByRef(cfg.LocationRef)
	if err != nil {
		l.Error("could not resolve LocationRef >%s< to a valid location ID", cfg.LocationRef)
		return nil, fmt.Errorf("could not resolve LocationRef >%s< to a valid location ID", cfg.LocationRef)
	}
	rec.AdventureGameLocationID = locationRec.ID

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

	for _, effectConfig := range cfg.AdventureGameLocationObjectEffectConfigs {
		_, err := t.createAdventureGameLocationObjectEffectRec(effectConfig, createdRec)
		if err != nil {
			l.Warn("failed creating adventure_game_location_object_effect record >%v<", err)
			return nil, err
		}
	}

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
