package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyLocationLinkRecs -
func (m *Domain) GetManyLocationLinkRecs(opts *sql.Options) ([]*record.LocationLink, error) {
	l := m.Logger("GetManyLocationLinkRecs")
	l.Debug("getting many location link records opts >%#v<", opts)
	r := m.LocationLinkRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetLocationLinkRec -
func (m *Domain) GetLocationLinkRec(recID string, lock *sql.Lock) (*record.LocationLink, error) {
	l := m.Logger("GetLocationLinkRec")
	l.Debug("getting location link record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.LocationLinkRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateLocationLinkRec -
func (m *Domain) CreateLocationLinkRec(rec *record.LocationLink) (*record.LocationLink, error) {
	l := m.Logger("CreateLocationLinkRec")
	l.Debug("creating location link record >%#v<", rec)

	if err := m.validateLocationLinkRecForCreate(rec); err != nil {
		l.Warn("failed to validate location link record >%v<", err)
		return rec, err
	}

	r := m.LocationLinkRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// validateLocationLinkRecForCreate validates the Name field for creation
func (m *Domain) validateLocationLinkRecForCreate(rec *record.LocationLink) error {
	if err := domain.ValidateStringField(record.FieldLocationLinkName, rec.Name); err != nil {
		return err
	}
	if len(rec.Name) > 64 {
		return InvalidFieldValue("name")
	}
	return nil
}

// DeleteLocationLinkRec -
func (m *Domain) DeleteLocationLinkRec(recID string) error {
	l := m.Logger("DeleteLocationLinkRec")
	l.Debug("deleting location link record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return err
	}
	r := m.LocationLinkRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
