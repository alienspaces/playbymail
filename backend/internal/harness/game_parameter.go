package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type GameParameterConfig struct {
	Record    *game_record.GameParameter
	Reference string
}

func (t *Testing) createGameParameterRec(config GameParameterConfig) (*game_record.GameParameter, error) {
	l := t.Logger("createGameParameterRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *game_record.GameParameter
	if config.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *config.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameParameter{}
	}

	rec = t.applyGameParameterRecDefaultValues(rec)

	// Create record
	l.Info("creating game parameter record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameParameterRec(rec)
	if err != nil {
		l.Warn("failed creating game parameter record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameParameterRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameParameterRec(rec)

	// Add to references store
	if config.Reference != "" {
		t.Data.Refs.GameParameterRefs[config.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameParameterRecDefaultValues(rec *game_record.GameParameter) *game_record.GameParameter {
	if rec == nil {
		rec = &game_record.GameParameter{}
	}

	if rec.GameType == "" {
		rec.GameType = "adventure"
	}

	if rec.ConfigKey == "" {
		rec.ConfigKey = gofakeit.Word()
	}

	if rec.ValueType == "" {
		rec.ValueType = "string"
	}

	if !rec.IsRequired {
		rec.IsRequired = false
	}

	return rec
}
