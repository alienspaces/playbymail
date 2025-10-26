package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type GameInstanceParameterConfig struct {
	Reference string // Reference to the game_instance_parameter record
	Record    *game_record.GameInstanceParameter
}

func (t *Testing) createGameInstanceParameterRec(cfg GameInstanceParameterConfig, gameInstanceRec *game_record.GameInstance) (*game_record.GameInstanceParameter, error) {
	l := t.Logger("createGameInstanceParameterRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *game_record.GameInstanceParameter
	if cfg.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameInstanceParameter{}
	}

	rec = t.applyGameInstanceParameterRecDefaultValues(rec)

	rec.GameInstanceID = gameInstanceRec.ID

	// Create record
	l.Debug("creating game instance parameter record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameInstanceParameterRec(rec)
	if err != nil {
		l.Warn("failed creating game instance parameter record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameInstanceParameterRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameInstanceParameterRec(rec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.GameInstanceParameterRefs[cfg.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameInstanceParameterRecDefaultValues(rec *game_record.GameInstanceParameter) *game_record.GameInstanceParameter {
	if rec == nil {
		rec = &game_record.GameInstanceParameter{}
	}

	if rec.GameInstanceID == "" {
		// Try to get a game instance from the data store
		if len(t.Data.GameInstanceRecs) > 0 {
			rec.GameInstanceID = t.Data.GameInstanceRecs[0].ID
		}
	}

	if rec.ParameterKey == "" {
		rec.ParameterKey = gofakeit.Word()
	}

	if !nullstring.IsValid(rec.ParameterValue) {
		rec.ParameterValue = nullstring.FromString(gofakeit.Word())
	}

	return rec
}
