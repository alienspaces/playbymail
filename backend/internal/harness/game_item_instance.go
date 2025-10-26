package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameItemInstanceRec(cfg GameItemInstanceConfig, gameInstanceRec *game_record.GameInstance) (*adventure_game_record.AdventureGameItemInstance, error) {
	l := t.Logger("createGameItemInstanceRec")

	if gameInstanceRec == nil {
		return nil, fmt.Errorf("game instance record is nil for game_item_instance record >%#v<", cfg)
	}

	if cfg.GameItemRef == "" {
		return nil, fmt.Errorf("game_item_instance record must have a GameItemRef")
	}

	var rec *adventure_game_record.AdventureGameItemInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameItemInstance{}
	}

	rec = t.applyGameItemInstanceRecDefaultValues(rec)

	rec.GameID = gameInstanceRec.GameID
	rec.GameInstanceID = gameInstanceRec.ID

	// The game item is retrieved by reference
	gameItemRec, err := t.Data.GetGameItemRecByRef(cfg.GameItemRef)
	if err != nil {
		l.Error("could not resolve GameItemRef >%s< to a valid game item ID", cfg.GameItemRef)
		return nil, fmt.Errorf("could not resolve GameItemRef >%s< to a valid game item ID", cfg.GameItemRef)
	}
	rec.AdventureGameItemID = gameItemRec.ID

	// The game location instance the item is attached to is retrieved by reference
	if cfg.GameLocationRef != "" {
		gameLocationInstanceRec, err := t.Data.GetGameLocationInstanceRecByLocationRef(cfg.GameLocationRef)
		if err != nil {
			l.Error("could not resolve GameLocationRef >%s< to a valid game location ID", cfg.GameLocationRef)
			return nil, fmt.Errorf("could not resolve GameLocationRef >%s< to a valid game location ID", cfg.GameLocationRef)
		}
		rec.AdventureGameLocationInstanceID = nullstring.FromString(gameLocationInstanceRec.ID)
	}

	// Create record
	l.Debug("creating game_item_instance record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameItemInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_item_instance record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameItemInstanceRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddGameItemInstanceRec(createdRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.GameItemInstanceRefs[cfg.Reference] = createdRec.ID
	}
	return createdRec, nil
}

func (t *Testing) applyGameItemInstanceRecDefaultValues(rec *adventure_game_record.AdventureGameItemInstance) *adventure_game_record.AdventureGameItemInstance {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameItemInstance{}
	}
	return rec
}
