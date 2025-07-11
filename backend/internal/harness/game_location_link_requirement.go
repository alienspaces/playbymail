package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameLocationLinkRequirementRec(cfg GameLocationLinkRequirementConfig, gameLocationLinkRec *record.GameLocationLink) (*record.GameLocationLinkRequirement, error) {
	l := t.Logger("createGameLocationLinkRequirementRec")

	if gameLocationLinkRec == nil {
		return nil, fmt.Errorf("game location link record is nil for game_location_link_requirement record >%#v<", cfg)
	}

	var rec *record.GameLocationLinkRequirement
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &record.GameLocationLinkRequirement{}
	}

	rec = t.applyGameLocationLinkRequirementRecDefaultValues(rec)

	rec.GameLocationLinkID = gameLocationLinkRec.ID

	if cfg.GameItemRef != "" {
		gameItemRec, err := t.Data.GetGameItemRecByRef(cfg.GameItemRef)
		if err != nil {
			l.Error("could not resolve GameItemRef >%s< to a valid game item ID", cfg.GameItemRef)
			return nil, fmt.Errorf("could not resolve GameItemRef >%s< to a valid game item ID", cfg.GameItemRef)
		}
		rec.GameItemID = gameItemRec.ID
	}

	l.Info("creating game_location_link_requirement record >%#v<", rec)

	// Create record
	rec, err := t.Domain.(*domain.Domain).CreateGameLocationLinkRequirementRec(rec)
	if err != nil {
		l.Warn("failed creating game_location_link_requirement record >%v<", err)
		return nil, err
	}

	// Add to data
	t.Data.AddGameLocationLinkRequirementRec(rec)

	// Add to teardown data
	t.teardownData.AddGameLocationLinkRequirementRec(rec)

	if cfg.Reference != "" {
		t.Data.Refs.GameLocationLinkRequirementRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameLocationLinkRequirementRecDefaultValues(rec *record.GameLocationLinkRequirement) *record.GameLocationLinkRequirement {
	if rec == nil {
		rec = &record.GameLocationLinkRequirement{}
	}
	if rec.Quantity == 0 {
		rec.Quantity = 1
	}
	return rec
}
