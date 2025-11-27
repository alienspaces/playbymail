package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameItemRec(itemConfig AdventureGameItemConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameItem, error) {
	l := t.Logger("createAdventureGameItemRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for adventure game item record >%#v<", itemConfig)
	}

	var rec *adventure_game_record.AdventureGameItem
	if itemConfig.Record != nil {
		recCopy := *itemConfig.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameItem{}
	}

	rec = t.applyAdventureGameItemRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Create record
	l.Debug("creating adventure game item record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateAdventureGameItemRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game item record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddAdventureGameItemRec(rec)

	// Add to teardown data store
	t.teardownData.AddAdventureGameItemRec(rec)

	// Add to references store
	if itemConfig.Reference != "" {
		t.Data.Refs.AdventureGameItemRefs[itemConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyAdventureGameItemRecDefaultValues(rec *adventure_game_record.AdventureGameItem) *adventure_game_record.AdventureGameItem {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameItem{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(10)
	}

	return rec
}
