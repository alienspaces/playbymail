package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameCreatureRec(creatureConfig AdventureGameCreatureConfig, gameRec *game_record.Game) (*adventure_game_record.AdventureGameCreature, error) {
	l := t.Logger("createAdventureGameCreatureRec")

	if gameRec == nil {
		return nil, fmt.Errorf("game record is nil for adventure game creature record >%#v<", creatureConfig)
	}

	var rec *adventure_game_record.AdventureGameCreature
	if creatureConfig.Record != nil {
		recCopy := *creatureConfig.Record
		rec = &recCopy
	} else {
		rec = &adventure_game_record.AdventureGameCreature{}
	}

	rec = t.applyAdventureGameCreatureRecDefaultValues(rec)

	rec.GameID = gameRec.ID

	// Create record
	l.Debug("creating adventure game creature record >%#v<", rec)

	createdRec, err := t.Domain.(*domain.Domain).CreateAdventureGameCreatureRec(rec)
	if err != nil {
		l.Warn("failed creating adventure game creature record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddAdventureGameCreatureRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddAdventureGameCreatureRec(createdRec)

	// Add to references store
	if creatureConfig.Reference != "" {
		t.Data.Refs.AdventureGameCreatureRefs[creatureConfig.Reference] = createdRec.ID
	}

	return createdRec, nil
}

func (t *Testing) applyAdventureGameCreatureRecDefaultValues(rec *adventure_game_record.AdventureGameCreature) *adventure_game_record.AdventureGameCreature {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameCreature{}
	}
	if rec.Name == "" {
		rec.Name = UniqueName(gofakeit.Name())
	}
	if rec.Description == "" {
		rec.Description = gofakeit.Sentence(10)
	}

	return rec
}
