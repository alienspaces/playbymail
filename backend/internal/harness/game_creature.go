package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameCreatureRec(charConfig GameCreatureConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameCreature, error) {
	l := t.Logger("createGameCreatureRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game_character record >%#v<", charConfig)
	}

	var rec *adventure_game_record.AdventureGameCreature
	if charConfig.Record != nil {
		recCopy := *charConfig.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameCreature{}
	}

	rec = t.applyGameCreatureRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Create record
	l.Debug("creating game_character record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameCreatureRec(rec)
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

func (t *Testing) applyGameCreatureRecDefaultValues(rec *adventure_game_record.AdventureGameCreature) *adventure_game_record.AdventureGameCreature {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameCreature{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	return rec
}
