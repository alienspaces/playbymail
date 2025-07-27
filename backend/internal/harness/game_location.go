package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameLocationRec(gameLocationConfig GameLocationConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameLocation, error) {
	l := t.Logger("createGameLocationRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for game location record >%#v<", gameLocationConfig)
	}

	var rec *adventure_game_record.AdventureGameLocation
	if gameLocationConfig.Record != nil {
		recCopy := *gameLocationConfig.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameLocation{}
	}

	rec = t.applyGameLocationRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Create record
	l.Info("creating game location record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationRec(rec)
	if err != nil {
		l.Warn("failed creating game location record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameLocationRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameLocationRec(rec)

	// Add to references store
	if gameLocationConfig.Reference != "" {
		t.Data.Refs.GameLocationRefs[gameLocationConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameLocationRecDefaultValues(rec *adventure_game_record.AdventureGameLocation) *adventure_game_record.AdventureGameLocation {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocation{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(10)
	}
	return rec
}
