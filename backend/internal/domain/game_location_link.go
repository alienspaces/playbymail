package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyGameLocationLinkRecs -
func (m *Domain) GetManyGameLocationLinkRecs(opts *sql.Options) ([]*record.GameLocationLink, error) {
	l := m.Logger("GetManyGameLocationLinkRecs")
	l.Debug("getting many location link records opts >%#v<", opts)
	r := m.GameLocationLinkRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetGameLocationLinkRec -
func (m *Domain) GetGameLocationLinkRec(recID string, lock *sql.Lock) (*record.GameLocationLink, error) {
	l := m.Logger("GetGameLocationLinkRec")
	l.Debug("getting location link record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.GameLocationLinkRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateGameLocationLinkRec -
func (m *Domain) CreateGameLocationLinkRec(rec *record.GameLocationLink) (*record.GameLocationLink, error) {
	l := m.Logger("CreateGameLocationLinkRec")
	l.Debug("creating location link record >%#v<", rec)

	if err := m.validateGameLocationLinkRecForCreate(rec); err != nil {
		l.Warn("failed to validate location link record >%v<", err)
		return rec, err
	}

	r := m.GameLocationLinkRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// validateGameLocationLinkRecForCreate validates the Name field for creation
func (m *Domain) validateGameLocationLinkRecForCreate(rec *record.GameLocationLink) error {
	if err := domain.ValidateStringField(record.FieldGameLocationLinkName, rec.Name); err != nil {
		return err
	}
	if len(rec.Name) > 64 {
		return InvalidFieldValue("name")
	}
	return nil
}

// DeleteGameLocationLinkRec -
func (m *Domain) DeleteGameLocationLinkRec(recID string) error {
	l := m.Logger("DeleteGameLocationLinkRec")
	l.Debug("deleting location link record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return err
	}
	r := m.GameLocationLinkRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveGameLocationLinkRec -
func (m *Domain) RemoveGameLocationLinkRec(recID string) error {
	l := m.Logger("RemoveGameLocationLinkRec")

	l.Debug("removing location link record ID >%s<", recID)

	_, err := m.GetGameLocationLinkRec(recID, sql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameLocationLinkRepository()

	// Add validation here if needed

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
