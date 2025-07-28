package harness

import (
	"github.com/brianvoe/gofakeit"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type GameInstanceConfigurationConfig struct {
	Record    *game_record.GameInstanceConfiguration
	Reference string
}

func (t *Testing) createGameInstanceConfigurationRec(config GameInstanceConfigurationConfig) (*game_record.GameInstanceConfiguration, error) {
	l := t.Logger("createGameInstanceConfigurationRec")

	// Create a new record instance to avoid reusing the same record across tests
	var rec *game_record.GameInstanceConfiguration
	if config.Record != nil {
		// Copy the record to avoid modifying the original
		recCopy := *config.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameInstanceConfiguration{}
	}

	rec = t.applyGameInstanceConfigurationRecDefaultValues(rec)

	// Create record
	l.Info("creating game instance configuration record >%#v<", rec)

	rec, err := t.Domain.(*domain.Domain).CreateGameInstanceConfigurationRec(rec)
	if err != nil {
		l.Warn("failed creating game instance configuration record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameInstanceConfigurationRec(rec)

	// Add to teardown data store
	t.teardownData.AddGameInstanceConfigurationRec(rec)

	// Add to references store
	if config.Reference != "" {
		t.Data.Refs.GameInstanceConfigurationRefs[config.Reference] = rec.ID
	}

	return rec, nil
}

func (t *Testing) applyGameInstanceConfigurationRecDefaultValues(rec *game_record.GameInstanceConfiguration) *game_record.GameInstanceConfiguration {
	if rec == nil {
		rec = &game_record.GameInstanceConfiguration{}
	}

	if rec.GameInstanceID == "" {
		// Try to get a game instance from the data store
		if len(t.Data.GameInstanceRecs) > 0 {
			rec.GameInstanceID = t.Data.GameInstanceRecs[0].ID
		}
	}

	if rec.ConfigKey == "" {
		rec.ConfigKey = gofakeit.Word()
	}

	if rec.ValueType == "" {
		rec.ValueType = "string"
	}

	return rec
}
