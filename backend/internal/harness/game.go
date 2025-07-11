package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameRec(gameConfig GameConfig) (*record.Game, error) {
	l := t.Logger("createGameRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *record.Game
	if gameConfig.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *gameConfig.Record
		rec = &recCopy
	} else {
		rec = &record.Game{}
	}

	rec = t.applyGameRecDefaultValues(rec)

	l.Info("creating game record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameRec(rec)
	if err != nil {
		l.Warn("failed creating game record >%v<", err)
		return nil, err
	}

	// Add the game record to the data store
	t.Data.AddGameRec(rec)

	// Add the game record to the teardown data store
	t.teardownData.AddGameRec(rec)

	// Add the game record to the data store refs
	if gameConfig.Reference != "" {
		t.Data.Refs.GameRefs[gameConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameRecDefaultValues(rec *record.Game) *record.Game {
	if rec == nil {
		rec = &record.Game{}
	}

	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}

	if rec.GameType == "" {
		rec.GameType = record.GameTypeAdventure
	}

	return rec
}
