package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
	adventure_game_record "gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
)

func (t *Testing) createGameItemRec(itemConfig GameItemConfig, gameRec *record.Game) (*adventure_game_record.AdventureGameItem, error) {
	l := t.Logger("createGameItemRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_item record >%#v<", itemConfig)
	}

	var rec *adventure_game_record.AdventureGameItem
	if itemConfig.Record != nil {
		recCopy := *itemConfig.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameItem{}
	}

	rec = t.applyGameItemRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Create record
	l.Info("creating game_item record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateAdventureGameItemRec(rec)
	if err != nil {
		l.Warn("failed creating game_item record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameItemRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameItemRec(rec)

	// Add to references store
	if itemConfig.Reference != "" {
		t.Data.Refs.GameItemRefs[itemConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameItemRecDefaultValues(rec *adventure_game_record.AdventureGameItem) *adventure_game_record.AdventureGameItem {
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
