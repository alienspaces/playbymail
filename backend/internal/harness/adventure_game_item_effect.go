package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (t *Testing) createAdventureGameItemEffectRec(cfg AdventureGameItemEffectConfig, itemRec *adventure_game_record.AdventureGameItem) (*adventure_game_record.AdventureGameItemEffect, error) {
	l := t.Logger("createAdventureGameItemEffectRec")

	if itemRec == nil {
		return nil, fmt.Errorf("item record is nil for adventure game item effect record >%#v<", cfg)
	}

	var rec *adventure_game_record.AdventureGameItemEffect
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameItemEffect{}
	}

	rec.GameID = itemRec.GameID
	rec.AdventureGameItemID = itemRec.ID

	if cfg.RequiredItemRef != "" {
		requiredItemRec, err := t.Data.GetAdventureGameItemRecByRef(cfg.RequiredItemRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve RequiredItemRef >%s<: %w", cfg.RequiredItemRef, err)
		}
		rec.RequiredAdventureGameItemID = nullstring.FromString(requiredItemRec.ID)
	}

	if cfg.RequiredLocationRef != "" {
		locationRec, err := t.Data.GetAdventureGameLocationRecByRef(cfg.RequiredLocationRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve RequiredLocationRef >%s<: %w", cfg.RequiredLocationRef, err)
		}
		rec.RequiredAdventureGameLocationID = nullstring.FromString(locationRec.ID)
	}

	if cfg.ResultItemRef != "" {
		resultItemRec, err := t.Data.GetAdventureGameItemRecByRef(cfg.ResultItemRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve ResultItemRef >%s<: %w", cfg.ResultItemRef, err)
		}
		rec.ResultAdventureGameItemID = nullstring.FromString(resultItemRec.ID)
	}

	if cfg.ResultLocationRef != "" {
		locationRec, err := t.Data.GetAdventureGameLocationRecByRef(cfg.ResultLocationRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve ResultLocationRef >%s<: %w", cfg.ResultLocationRef, err)
		}
		rec.ResultAdventureGameLocationID = nullstring.FromString(locationRec.ID)
	}

	if cfg.ResultLinkRef != "" {
		linkRec, err := t.Data.GetAdventureGameLocationLinkRecByRef(cfg.ResultLinkRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve ResultLinkRef >%s<: %w", cfg.ResultLinkRef, err)
		}
		rec.ResultAdventureGameLocationLinkID = nullstring.FromString(linkRec.ID)
	}

	if cfg.ResultCreatureRef != "" {
		creatureRec, err := t.Data.GetAdventureGameCreatureRecByRef(cfg.ResultCreatureRef)
		if err != nil {
			return nil, fmt.Errorf("could not resolve ResultCreatureRef >%s<: %w", cfg.ResultCreatureRef, err)
		}
		rec.ResultAdventureGameCreatureID = nullstring.FromString(creatureRec.ID)
	}

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameItemEffectRec(rec)
	if err != nil {
		l.Warn("failed creating adventure_game_item_effect record >%v<", err)
		return nil, err
	}

	t.Data.AddAdventureGameItemEffectRec(createdRec)
	t.teardownData.AddAdventureGameItemEffectRec(createdRec)

	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameItemEffectRefs[cfg.Reference] = createdRec.ID
	}

	l.Debug("created adventure_game_item_effect record ID >%s<", createdRec.ID)

	return createdRec, nil
}
