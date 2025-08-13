package domain

import (
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) GetGameParameterRec(recID string, lock *sql.Lock) (*game_record.GameParameter, error) {
	r := m.GameParameterRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateGameParameterRec(rec *game_record.GameParameter) (*game_record.GameParameter, error) {
	r := m.GameParameterRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateGameParameterRec(next *game_record.GameParameter) (*game_record.GameParameter, error) {
	r := m.GameParameterRepository()
	rec, err := r.UpdateOne(next)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) DeleteGameParameterRec(recID string) error {
	r := m.GameParameterRepository()
	return r.DeleteOne(recID)
}

func (m *Domain) GetGameParameterRecs(opts *sql.Options) ([]*game_record.GameParameter, error) {
	r := m.GameParameterRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (m *Domain) ValidateGameParameter(rec *game_record.GameParameter) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) RemoveGameParameterRec(recID string) error {
	r := m.GameParameterRepository()
	if err := r.RemoveOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) GetManyGameParameterRecs(opts *sql.Options) ([]*game_record.GameParameter, error) {
	r := m.GameParameterRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

// GetGameParametersByGameType gets all parameters for a specific game type
func (m *Domain) GetGameParametersByGameType(gameType string) ([]*game_record.GameParameter, error) {
	r := m.GameParameterRepository()
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
