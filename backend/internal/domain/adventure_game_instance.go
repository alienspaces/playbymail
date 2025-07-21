package domain

import (
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func (m *Domain) GetAdventureGameInstanceRec(recID string, lock *sql.Lock) (*record.AdventureGameInstance, error) {
	r := m.AdventureGameInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (m *Domain) CreateAdventureGameInstanceRec(rec *record.AdventureGameInstance) (*record.AdventureGameInstance, error) {
	r := m.AdventureGameInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (m *Domain) UpdateAdventureGameInstanceRec(next *record.AdventureGameInstance) (*record.AdventureGameInstance, error) {
	r := m.AdventureGameInstanceRepository()
	next, err := r.UpdateOne(next)
	if err != nil {
		return next, err
	}
	return next, nil
}

func (m *Domain) DeleteAdventureGameInstanceRec(recID string) error {
	r := m.AdventureGameInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return err
	}
	return nil
}

func (m *Domain) ValidateAdventureGameInstance(rec *record.AdventureGameInstance) error {
	// Add validation logic as needed
	return nil
}

func (m *Domain) RemoveAdventureGameInstanceRec(recID string) error {
	r := m.AdventureGameInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return err
	}
	return nil
}
