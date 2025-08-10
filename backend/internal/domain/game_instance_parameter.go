package domain

import (
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) GetGameInstanceParameterRec(recID string, lock *sql.Lock) (*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateGameInstanceParameterRec(rec *game_record.GameInstanceParameter) (*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateGameInstanceParameterRec(next *game_record.GameInstanceParameter) (*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
	rec, err := r.UpdateOne(next)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) DeleteGameInstanceParameterRec(recID string) error {
	r := m.GameInstanceParameterRepository()
	return r.DeleteOne(recID)
}

func (m *Domain) GetGameInstanceParameterRecs(opts *sql.Options) ([]*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (m *Domain) ValidateGameInstanceParameter(rec *game_record.GameInstanceParameter) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) GetManyGameInstanceParameterRecs(opts *sql.Options) ([]*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

// GetGameInstanceParametersByGameInstanceID gets all parameters for a specific game instance
func (m *Domain) GetGameInstanceParametersByGameInstanceID(gameInstanceID string) ([]*game_record.GameInstanceParameter, error) {
	r := m.GameInstanceParameterRepository()
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
