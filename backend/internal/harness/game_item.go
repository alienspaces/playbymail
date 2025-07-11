package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameItemRec(itemConfig GameItemConfig, gameRec *record.Game) (*record.GameItem, error) {
	l := t.Logger("createGameItemRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_item record >%#v<", itemConfig)
	}

	var rec *record.GameItem
	if itemConfig.Record != nil {
		recCopy := *itemConfig.Record
		rec = &recCopy
	} else {
		rec = &record.GameItem{}
	}

	rec = t.applyGameItemRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	l.Info("creating game_item record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameItemRec(rec)
	if err != nil {
		l.Warn("failed creating game_item record >%v<", err)
		return nil, err
	}

	t.Data.AddGameItemRec(rec)
	t.teardownData.AddGameItemRec(rec)

	if itemConfig.Reference != "" {
		t.Data.Refs.GameItemRefs[itemConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameItemRecDefaultValues(rec *record.GameItem) *record.GameItem {
	if rec == nil {
		rec = &record.GameItem{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(10)
	}
	return rec
}
