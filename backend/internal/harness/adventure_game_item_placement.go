package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameItemPlacementRec(cfg AdventureGameItemPlacementConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameItemPlacement, error) {
	l := t.Logger("createAdventureGameItemPlacementRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for adventure game item placement record >%#v<", cfg)
	}

	var rec *adventure_game_record.AdventureGameItemPlacement
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameItemPlacement{}
	}

	rec.GameID = gameRec.ID

	if cfg.GameItemRef != "" {
		itemID, ok := t.Data.Refs.AdventureGameItemRefs[cfg.GameItemRef]
		if !ok {
			return nil, fmt.Errorf("item ref >%s< not found in data refs", cfg.GameItemRef)
		}
		rec.AdventureGameItemID = itemID
	}

	if cfg.GameLocationRef != "" {
		locationID, ok := t.Data.Refs.AdventureGameLocationRefs[cfg.GameLocationRef]
		if !ok {
			return nil, fmt.Errorf("location ref >%s< not found in data refs", cfg.GameLocationRef)
		}
		rec.AdventureGameLocationID = locationID
	}

	if cfg.InitialCount > 0 {
		rec.InitialCount = cfg.InitialCount
	} else {
		rec.InitialCount = 1
	}

	l.Debug("creating adventure game item placement record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameItemPlacementRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game item placement record >%v<", err)
		return nil, err
	}

	t.Data.AddAdventureGameItemPlacementRec(createdRec)
	t.teardownData.AddAdventureGameItemPlacementRec(createdRec)

	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameItemPlacementRefs[cfg.Reference] = createdRec.ID
	}

	return createdRec, nil
}
