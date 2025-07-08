package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyLocationRecs -
func (m *Domain) GetManyLocationRecs(opts *coresql.Options) ([]*record.Location, error) {
	l := m.Logger("GetManyLocationRecs")

	l.Debug("getting many location records opts >%#v<", opts)

	r := m.LocationRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetLocationRec -
func (m *Domain) GetLocationRec(recID string, lock *coresql.Lock) (*record.Location, error) {
	l := m.Logger("GetLocationRec")

	l.Debug("getting location record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.LocationRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableLocation, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateLocationRec -
func (m *Domain) CreateLocationRec(rec *record.Location) (*record.Location, error) {
	l := m.Logger("CreateLocationRec")

	l.Debug("creating location record >%#v<", rec)

	r := m.LocationRepository()

	// Add validation here if needed

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateLocationRec -
func (m *Domain) UpdateLocationRec(next *record.Location) (*record.Location, error) {
	l := m.Logger("UpdateLocationRec")

	_, err := m.GetLocationRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}

	l.Debug("updating location record >%#v<", next)

	// Add validation here if needed

	r := m.LocationRepository()

	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}

	return next, nil
}

// DeleteLocationRec -
func (m *Domain) DeleteLocationRec(recID string) error {
	l := m.Logger("DeleteLocationRec")

	l.Debug("deleting location record ID >%s<", recID)

	_, err := m.GetLocationRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.LocationRepository()

	// Add validation here if needed

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveLocationRec -
func (m *Domain) RemoveLocationRec(recID string) error {
	l := m.Logger("RemoveLocationRec")

	l.Debug("removing location record ID >%s<", recID)

	_, err := m.GetLocationRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.LocationRepository()

	// Add validation here if needed

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
