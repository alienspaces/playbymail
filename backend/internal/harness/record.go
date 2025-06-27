package harness

import (
	"fmt"

	"github.com/brianvoe/gofakeit"

	corerecord "gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// UniqueName appends a UUID4 to the end of the name to make it unique
// for parallel test execution.
func UniqueName(name string) string {
	if name == "" {
		name = gofakeit.Color()
	}
	return fmt.Sprintf("%s (%s)", name, corerecord.NewRecordID())
}

// removes the unique UUID4 from the end of the name to make it normal for
// test harness functions that return a record based on its non unique name.
func NormalName(name string) string {
	return name[:len(name)-39]
}

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

	l.Info("Creating game record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameRec(rec)
	if err != nil {
		l.Warn("failed creating game record >%v<", err)
		return nil, err
	}

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
		rec.Name = UniqueName(rec.Name)
	}

	return rec
}
