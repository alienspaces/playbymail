package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameCreaturePlacementRec(cfg AdventureGameCreaturePlacementConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameCreaturePlacement, error) {
	l := t.Logger("createAdventureGameCreaturePlacementRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for adventure game creature placement record >%#v<", cfg)
	}

	var rec *adventure_game_record.AdventureGameCreaturePlacement
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameCreaturePlacement{}
	}

	rec.GameID = gameRec.ID

	if cfg.GameCreatureRef != "" {
		creatureID, ok := t.Data.Refs.AdventureGameCreatureRefs[cfg.GameCreatureRef]
		if !ok {
			return nil, fmt.Errorf("creature ref >%s< not found in data refs", cfg.GameCreatureRef)
		}
		rec.AdventureGameCreatureID = creatureID
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

	l.Debug("creating adventure game creature placement record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameCreaturePlacementRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game creature placement record >%v<", err)
		return nil, err
	}

	t.Data.AddAdventureGameCreaturePlacementRec(createdRec)
	t.teardownData.AddAdventureGameCreaturePlacementRec(createdRec)

	if cfg.Reference != "" {
		t.Data.Refs.AdventureGameCreaturePlacementRefs[cfg.Reference] = createdRec.ID
	}

	return createdRec, nil
}
