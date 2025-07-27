package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/core/sql"
	adventure_game_record "gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
)

// GetManyAdventureGameLocationLinkRecs -
func (m *Domain) GetManyAdventureGameLocationLinkRecs(opts *sql.Options) ([]*adventure_game_record.AdventureGameLocationLink, error) {
	l := m.Logger("GetManyAdventureGameLocationLinkRecs")
	l.Debug("getting many adventure_game_location_link records opts >%#v<", opts)
	r := m.AdventureGameLocationLinkRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameLocationLinkRec -
func (m *Domain) GetAdventureGameLocationLinkRec(recID string, lock *sql.Lock) (*adventure_game_record.AdventureGameLocationLink, error) {
	l := m.Logger("GetAdventureGameLocationLinkRec")
	l.Debug("getting adventure_game_location_link record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameLocationLinkRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameLocationLinkRec -
func (m *Domain) CreateAdventureGameLocationLinkRec(rec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_record.AdventureGameLocationLink, error) {
	l := m.Logger("CreateAdventureGameLocationLinkRec")
	l.Debug("creating adventure_game_location_link record >%#v<", rec)

	if err := m.validateAdventureGameLocationLinkRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_location_link record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameLocationLinkRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// validateAdventureGameLocationLinkRecForCreate validates the Name field for creation
func (m *Domain) validateAdventureGameLocationLinkRecForCreate(rec *adventure_game_record.AdventureGameLocationLink) error {
	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameLocationLinkName, rec.Name); err != nil {
		return err
	}
	if len(rec.Name) > 64 {
		return InvalidFieldValue("name")
	}
	return nil
}

func (m *Domain) UpdateAdventureGameLocationLinkRec(rec *adventure_game_record.AdventureGameLocationLink) (*adventure_game_record.AdventureGameLocationLink, error) {
	l := m.Logger("UpdateAdventureGameLocationLinkRec")
	l.Debug("updating adventure_game_location_link record >%#v<", rec)

	if err := m.validateAdventureGameLocationLinkRecForUpdate(rec); err != nil {
		l.Warn("failed to validate adventure_game_location_link record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameLocationLinkRepository()

	var err error
	rec, err = r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

func (m *Domain) validateAdventureGameLocationLinkRecForUpdate(rec *adventure_game_record.AdventureGameLocationLink) error {
	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameLocationLinkName, rec.Name); err != nil {
		return err
	}
	return nil
}

// DeleteAdventureGameLocationLinkRec -
func (m *Domain) DeleteAdventureGameLocationLinkRec(recID string) error {
	l := m.Logger("DeleteAdventureGameLocationLinkRec")
	l.Debug("deleting adventure_game_location_link record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return err
	}
	r := m.AdventureGameLocationLinkRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameLocationLinkRec -
func (m *Domain) RemoveAdventureGameLocationLinkRec(recID string) error {
	l := m.Logger("RemoveAdventureGameLocationLinkRec")

	l.Debug("removing adventure_game_location_link record ID >%s<", recID)

	_, err := m.GetAdventureGameLocationLinkRec(recID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameLocationLinkRepository()

	// Add validation here if needed

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
