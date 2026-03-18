package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (t *Testing) createAdventureGameLocationLinkRequirementRec(cfg AdventureGameLocationLinkRequirementConfig, gameLocationLinkRec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	l := t.Logger("createAdventureGameLocationLinkRequirementRec")

	if gameLocationLinkRec == nil {
		return nil, fmt.Errorf("game location link record is nil for adventure game location link requirement record >%#v<", cfg)
	}

	if cfg.GameItemRef == "" && cfg.GameCreatureRef == "" {
		return nil, fmt.Errorf("one of GameItemRef or GameCreatureRef is required for adventure game location link requirement record >%#v<", cfg)
	}

	if cfg.GameItemRef != "" && cfg.GameCreatureRef != "" {
		return nil, fmt.Errorf("only one of GameItemRef or GameCreatureRef may be set for adventure game location link requirement record >%#v<", cfg)
	}

	var rec *adventure_game_record.AdventureGameLocationLinkRequirement
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameLocationLinkRequirement{}
	}

	rec = t.applyAdventureGameLocationLinkRequirementRecDefaultValues(rec)

	rec.GameID = gameLocationLinkRec.GameID
	rec.AdventureGameLocationLinkID = gameLocationLinkRec.ID

	if cfg.GameItemRef != "" {
		gameItemRec, err := t.Data.GetAdventureGameItemRecByRef(cfg.GameItemRef)
		if err != nil {
			l.Error("could not resolve GameItemRef >%s< to a valid game item ID", cfg.GameItemRef)
			return nil, fmt.Errorf("could not resolve GameItemRef >%s< to a valid game item ID", cfg.GameItemRef)
		}
		rec.AdventureGameItemID = nullstring.FromString(gameItemRec.ID)
	}

	if cfg.GameCreatureRef != "" {
		gameCreatureRec, err := t.Data.GetAdventureGameCreatureRecByRef(cfg.GameCreatureRef)
		if err != nil {
			l.Error("could not resolve GameCreatureRef >%s< to a valid game creature ID", cfg.GameCreatureRef)
			return nil, fmt.Errorf("could not resolve GameCreatureRef >%s< to a valid game creature ID", cfg.GameCreatureRef)
		}
		rec.AdventureGameCreatureID = nullstring.FromString(gameCreatureRec.ID)
	}

	// Create record
	l.Debug("creating adventure game location link requirement record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationLinkRequirementRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location link requirement record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddAdventureGameLocationLinkRequirementRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddAdventureGameLocationLinkRequirementRec(createdRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameLocationLinkRequirementRefs[cfg.Reference] = createdRec.ID
	}

	return createdRec, nil
}

func (t *Testing) applyAdventureGameLocationLinkRequirementRecDefaultValues(rec *adventure_game_record.AdventureGameLocationLinkRequirement) *adventure_game_record.AdventureGameLocationLinkRequirement {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocationLinkRequirement{}
	}
	if rec.Quantity == 0 {
		rec.Quantity = 1
	}
	if rec.Purpose == "" {
		rec.Purpose = adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse
	}
	if rec.Condition == "" {
		rec.Condition = adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory
	}
	return rec
}
