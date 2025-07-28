package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type GameConfigurationConfig struct {
	Record    *game_record.GameConfiguration
	Reference string
}

func (t *Testing) createGameConfigurationRec(config GameConfigurationConfig) (*game_record.GameConfiguration, error) {
	l := t.Logger("createGameConfigurationRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *game_record.GameConfiguration
	if config.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *config.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameConfiguration{}
	}

	rec = t.applyGameConfigurationRecDefaultValues(rec)

	// Create record
	l.Info("creating game configuration record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameConfigurationRec(rec)
	if err != nil {
		l.Warn("failed creating game configuration record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameConfigurationRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameConfigurationRec(rec)

	// Add to references store
	if config.Reference != "" {
		t.Data.Refs.GameConfigurationRefs[config.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameConfigurationRecDefaultValues(rec *game_record.GameConfiguration) *game_record.GameConfiguration {
	if rec == nil {
		rec = &game_record.GameConfiguration{}
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
