package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (t *Testing) createGameCreatureRec(charConfig GameCreatureConfig, gameRec *record.Game) (*record.GameCreature, error) {
	l := t.Logger("createGameCreatureRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_character record >%#v<", charConfig)
	}

	var rec *record.GameCreature
	if charConfig.Record != nil {
		recCopy := *charConfig.Record
		rec = &recCopy
	} else {
		rec = &record.GameCreature{}
	}

	rec = t.applyGameCreatureRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Create record
	l.Info("creating game_character record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateGameCreatureRec(rec)
	if err != nil {
		l.Warn("failed creating game_character record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameCreatureRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddGameCreatureRec(createdRec)

	// Add to references store
	if charConfig.Reference != "" {
		t.Data.Refs.GameCreatureRefs[charConfig.Reference] = createdRec.ID
	}

	return createdRec, nil
}

func (t *Testing) applyGameCreatureRecDefaultValues(rec *record.GameCreature) *record.GameCreature {
	if rec == nil {
		rec = &record.GameCreature{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	return rec
}
