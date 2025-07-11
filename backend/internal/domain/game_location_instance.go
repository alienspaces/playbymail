package domain

import (
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (m *Domain) GetGameLocationInstanceRec(recID string, lock *sql.Lock) (*record.GameLocationInstance, error) {
	r := m.GameLocationInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateGameLocationInstanceRec(rec *record.GameLocationInstance) (*record.GameLocationInstance, error) {
	r := m.GameLocationInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateGameLocationInstanceRec(next *record.GameLocationInstance) (*record.GameLocationInstance, error) {
	r := m.GameLocationInstanceRepository()
	next, err := r.UpdateOne(next)
	if err != nil {
		return next, err
	}
	return next, nil
}

func (m *Domain) DeleteGameLocationInstanceRec(recID string) error {
	r := m.GameLocationInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) RemoveGameLocationInstanceRec(recID string) error {
	r := m.GameLocationInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) ValidateGameLocationInstance(rec *record.GameLocationInstance) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) GetManyGameLocationInstanceRecs(opts *sql.Options) ([]*record.GameLocationInstance, error) {
	r := m.GameLocationInstanceRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}
