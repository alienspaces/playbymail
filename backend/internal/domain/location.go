package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyGameLocationRecs -
func (m *Domain) GetManyGameLocationRecs(opts *coresql.Options) ([]*record.GameLocation, error) {
	l := m.Logger("GetManyGameLocationRecs")

	l.Debug("getting many game_location records opts >%#v<", opts)

	r := m.GameLocationRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameLocationRec -
func (m *Domain) GetGameLocationRec(recID string, lock *coresql.Lock) (*record.GameLocation, error) {
	l := m.Logger("GetGameLocationRec")

	l.Debug("getting game_location record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameLocationRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableGameLocation, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateGameLocationRec -
func (m *Domain) CreateGameLocationRec(rec *record.GameLocation) (*record.GameLocation, error) {
	l := m.Logger("CreateGameLocationRec")

	l.Debug("creating game_location record >%#v<", rec)

	r := m.GameLocationRepository()

	// Add validation here if needed

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateGameLocationRec -
func (m *Domain) UpdateGameLocationRec(next *record.GameLocation) (*record.GameLocation, error) {
	l := m.Logger("UpdateGameLocationRec")

	_, err := m.GetGameLocationRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}

	l.Debug("updating game_location record >%#v<", next)

	// Add validation here if needed

	r := m.GameLocationRepository()

	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}

	return next, nil
}

// DeleteGameLocationRec -
func (m *Domain) DeleteGameLocationRec(recID string) error {
	l := m.Logger("DeleteGameLocationRec")

	l.Debug("deleting game_location record ID >%s<", recID)

	_, err := m.GetGameLocationRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameLocationRepository()

	// Add validation here if needed

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveGameLocationRec -
func (m *Domain) RemoveGameLocationRec(recID string) error {
	l := m.Logger("RemoveGameLocationRec")

	l.Debug("removing game_location record ID >%s<", recID)

	_, err := m.GetGameLocationRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameLocationRepository()

	// Add validation here if needed

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
