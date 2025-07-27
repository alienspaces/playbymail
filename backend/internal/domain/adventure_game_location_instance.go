package domain

import (
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) GetAdventureGameLocationInstanceRec(recID string, lock *sql.Lock) (*adventure_game_record.AdventureGameLocationInstance, error) {
	r := m.AdventureGameLocationInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateAdventureGameLocationInstanceRec(rec *adventure_game_record.AdventureGameLocationInstance) (*adventure_game_record.AdventureGameLocationInstance, error) {
	r := m.AdventureGameLocationInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateAdventureGameLocationInstanceRec(next *adventure_game_record.AdventureGameLocationInstance) (*adventure_game_record.AdventureGameLocationInstance, error) {
	r := m.AdventureGameLocationInstanceRepository()
	next, err := r.UpdateOne(next)
	if err != nil {
		return next, err
	}
	return next, nil
}

func (m *Domain) DeleteAdventureGameLocationInstanceRec(recID string) error {
	r := m.AdventureGameLocationInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) RemoveAdventureGameLocationInstanceRec(recID string) error {
	r := m.AdventureGameLocationInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) ValidateAdventureGameLocationInstance(rec *adventure_game_record.AdventureGameLocationInstance) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) GetManyAdventureGameLocationInstanceRecs(opts *sql.Options) ([]*adventure_game_record.AdventureGameLocationInstance, error) {
	r := m.AdventureGameLocationInstanceRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, err
	}
	return recs, nil
}
