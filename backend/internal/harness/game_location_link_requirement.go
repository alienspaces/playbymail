package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (t *Testing) createGameLocationLinkRequirementRec(cfg GameLocationLinkRequirementConfig, gameLocationLinkRec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_record.AdventureGameLocationLinkRequirement, error) {
	l := t.Logger("createGameLocationLinkRequirementRec")

	if gameLocationLinkRec == nil {
		return nil, fmt.Errorf("game location link record is nil for game_location_link_requirement record >%#v<", cfg)
	}

	var rec *adventure_game_record.AdventureGameLocationLinkRequirement
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameLocationLinkRequirement{}
	}

	rec = t.applyGameLocationLinkRequirementRecDefaultValues(rec)

	rec.GameID = gameLocationLinkRec.GameID
	rec.AdventureGameLocationLinkID = gameLocationLinkRec.ID

	if cfg.GameItemRef != "" {
		gameItemRec, err := t.Data.GetGameItemRecByRef(cfg.GameItemRef)
		if err != nil {
			l.Error("could not resolve GameItemRef >%s< to a valid game item ID", cfg.GameItemRef)
			return nil, fmt.Errorf("could not resolve GameItemRef >%s< to a valid game item ID", cfg.GameItemRef)
		}
		rec.AdventureGameItemID = gameItemRec.ID
	}

	// Create record
	l.Info("creating game_location_link_requirement record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationLinkRequirementRec(rec)
	if err != nil {
		l.Warn("failed creating game_location_link_requirement record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameLocationLinkRequirementRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameLocationLinkRequirementRec(rec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.GameLocationLinkRequirementRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameLocationLinkRequirementRecDefaultValues(rec *adventure_game_record.AdventureGameLocationLinkRequirement) *adventure_game_record.AdventureGameLocationLinkRequirement {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocationLinkRequirement{}
	}
	if rec.Quantity == 0 {
		rec.Quantity = 1
	}
	return rec
}
