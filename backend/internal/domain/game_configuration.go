package domain

import (
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) GetGameConfigurationRec(recID string, lock *sql.Lock) (*game_record.GameConfiguration, error) {
	r := m.GameConfigurationRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateGameConfigurationRec(rec *game_record.GameConfiguration) (*game_record.GameConfiguration, error) {
	r := m.GameConfigurationRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateGameConfigurationRec(next *game_record.GameConfiguration) (*game_record.GameConfiguration, error) {
	r := m.GameConfigurationRepository()
	rec, err := r.UpdateOne(next)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) DeleteGameConfigurationRec(recID string) error {
	r := m.GameConfigurationRepository()
	return r.DeleteOne(recID)
}

func (m *Domain) GetGameConfigurationRecs(opts *sql.Options) ([]*game_record.GameConfiguration, error) {
	r := m.GameConfigurationRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (m *Domain) ValidateGameConfiguration(rec *game_record.GameConfiguration) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) RemoveGameConfigurationRec(recID string) error {
	r := m.GameConfigurationRepository()
	if err := r.RemoveOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) GetManyGameConfigurationRecs(opts *sql.Options) ([]*game_record.GameConfiguration, error) {
	r := m.GameConfigurationRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

// GetGameConfigurationsByGameType gets all configurations for a specific game type
func (m *Domain) GetGameConfigurationsByGameType(gameType string) ([]*game_record.GameConfiguration, error) {
	r := m.GameConfigurationRepository()
	opts := &sql.Options{
		Params: []sql.Param{
			{
				Col: "game_type",
				Val: gameType,
			},
		},
	}
	return r.GetMany(opts)
}
