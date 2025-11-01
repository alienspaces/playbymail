package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameLocationRec(gameLocationConfig AdventureGameLocationConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameLocation, error) {
	l := t.Logger("createAdventureGameLocationRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for adventure game location record >%#v<", gameLocationConfig)
	}

	var rec *adventure_game_record.AdventureGameLocation
	if gameLocationConfig.Record != nil {
		recCopy := *gameLocationConfig.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameLocation{}
	}

	rec = t.applyAdventureGameLocationRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Create record
	l.Debug("creating adventure game location record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateAdventureGameLocationRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game location record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddAdventureGameLocationRec(rec)

	// Add to teardown data store
	t.teardownData.AddAdventureGameLocationRec(rec)

	// Add to references store
	if gameLocationConfig.Reference != "" {
		t.Data.Refs.AdventureGameLocationRefs[gameLocationConfig.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyAdventureGameLocationRecDefaultValues(rec *adventure_game_record.AdventureGameLocation) *adventure_game_record.AdventureGameLocation {
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
