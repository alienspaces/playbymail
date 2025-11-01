package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (t *Testing) createAdventureGameLocationLinkRequirementRec(cfg AdventureGameLocationLinkRequirementConfig, gameLocationLinkRec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	l := t.Logger("createAdventureGameLocationLinkRequirementRec")

	if gameLocationLinkRec == nil {
		return nil, fmt.Errorf("game location link record is nil for adventure game location link requirement record >%#v<", cfg)
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
		rec.AdventureGameItemID = gameItemRec.ID
	}

	// Create record
	l.Debug("creating adventure game location link requirement record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationLinkRequirementRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location link requirement record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddAdventureGameLocationLinkRequirementRec(rec)

	// Add to teardown data store
	t.teardownData.AddAdventureGameLocationLinkRequirementRec(rec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameLocationLinkRequirementRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyAdventureGameLocationLinkRequirementRecDefaultValues(rec *adventure_game_record.AdventureGameLocationLinkRequirement) *adventure_game_record.AdventureGameLocationLinkRequirement {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocationLinkRequirement{}
	}
	if rec.Quantity == 0 {
		rec.Quantity = 1
	}
	return rec
}
