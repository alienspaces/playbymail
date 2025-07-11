package domain

import (
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (m *Domain) GetGameInstanceRec(recID string, lock *sql.Lock) (*record.GameInstance, error) {
	r := m.GameInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateGameInstanceRec(rec *record.GameInstance) (*record.GameInstance, error) {
	r := m.GameInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateGameInstanceRec(next *record.GameInstance) (*record.GameInstance, error) {
	r := m.GameInstanceRepository()
	next, err := r.UpdateOne(next)
	if err != nil {
		return next, err
	}
	return next, nil
}

func (m *Domain) DeleteGameInstanceRec(recID string) error {
	r := m.GameInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) Validate(rec *record.GameInstance) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) RemoveGameInstanceRec(recID string) error {
	r := m.GameInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return err
	}
	return nil
}
