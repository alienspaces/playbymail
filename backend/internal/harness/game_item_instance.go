package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameItemInstanceRec(cfg GameItemInstanceConfig, gameInstanceRec *record.GameInstance) (*record.GameItemInstance, error) {
	l := t.Logger("createGameItemInstanceRec")

	if gameInstanceRec == nil {
		return nil, fmt.Errorf("game instance record is nil for game_item_instance record >%#v<", cfg)
	}

	if cfg.GameItemRef == "" {
		return nil, fmt.Errorf("game_item_instance record must have a GameItemRef")
	}

	var rec *record.GameItemInstance
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &record.GameItemInstance{}
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
	rec.GameItemID = gameItemRec.ID

	// The game location instance the item is attached to is retrieved by reference
	if cfg.GameLocationRef != "" {
		gameLocationInstanceRec, err := t.Data.GetGameLocationInstanceRecByLocationRef(cfg.GameLocationRef)
		if err != nil {
			l.Error("could not resolve GameLocationRef >%s< to a valid game location ID", cfg.GameLocationRef)
			return nil, fmt.Errorf("could not resolve GameLocationRef >%s< to a valid game location ID", cfg.GameLocationRef)
		}
		rec.GameLocationInstanceID = nullstring.FromString(gameLocationInstanceRec.ID)
	}

	l.Info("creating game_item_instance record >%#v<", rec)
	createdRec, err := t.Domain.(*domain.Domain).CreateGameItemInstanceRec(rec)
	if err != nil {
		l.Warn("failed creating game_item_instance record >%v<", err)
		return nil, err
	}

	t.Data.AddGameItemInstanceRec(createdRec)
	t.teardownData.AddGameItemInstanceRec(createdRec)

	if cfg.Reference != "" {
		t.Data.Refs.GameItemInstanceRefs[cfg.Reference] = createdRec.ID
	}
	return createdRec, nil
}

func (t *Testing) applyGameItemInstanceRecDefaultValues(rec *record.GameItemInstance) *record.GameItemInstance {
	if rec == nil {
		rec = &record.GameItemInstance{}
	}
	return rec
}
