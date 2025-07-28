package domain

import (
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) GetGameInstanceConfigurationRec(recID string, lock *sql.Lock) (*game_record.GameInstanceConfiguration, error) {
	r := m.GameInstanceConfigurationRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateGameInstanceConfigurationRec(rec *game_record.GameInstanceConfiguration) (*game_record.GameInstanceConfiguration, error) {
	r := m.GameInstanceConfigurationRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateGameInstanceConfigurationRec(next *game_record.GameInstanceConfiguration) (*game_record.GameInstanceConfiguration, error) {
	r := m.GameInstanceConfigurationRepository()
	rec, err := r.UpdateOne(next)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) DeleteGameInstanceConfigurationRec(recID string) error {
	r := m.GameInstanceConfigurationRepository()
	return r.DeleteOne(recID)
}

func (m *Domain) GetGameInstanceConfigurationRecs(opts *sql.Options) ([]*game_record.GameInstanceConfiguration, error) {
	r := m.GameInstanceConfigurationRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (m *Domain) ValidateGameInstanceConfiguration(rec *game_record.GameInstanceConfiguration) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) GetManyGameInstanceConfigurationRecs(opts *sql.Options) ([]*game_record.GameInstanceConfiguration, error) {
	r := m.GameInstanceConfigurationRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

// GetGameInstanceConfigurationsByGameInstanceID gets all configurations for a specific game instance
func (m *Domain) GetGameInstanceConfigurationsByGameInstanceID(gameInstanceID string) ([]*game_record.GameInstanceConfiguration, error) {
	r := m.GameInstanceConfigurationRepository()
	opts := &sql.Options{
		Params: []sql.Param{
			{
				Col: "game_instance_id",
				Val: gameInstanceID,
			},
		},
	}
	return r.GetMany(opts)
}
