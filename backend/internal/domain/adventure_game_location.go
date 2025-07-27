package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
)

// GetManyAdventureGameLocationRecs -
func (m *Domain) GetManyAdventureGameLocationRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameLocation, error) {
	l := m.Logger("GetManyAdventureGameLocationRecs")

	l.Debug("getting many adventure_game_location records opts >%#v<", opts)

	r := m.AdventureGameLocationRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAdventureGameLocationRec -
func (m *Domain) GetAdventureGameLocationRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameLocation, error) {
	l := m.Logger("GetAdventureGameLocationRec")

	l.Debug("getting adventure_game_location record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AdventureGameLocationRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameLocation, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAdventureGameLocationRec -
func (m *Domain) CreateAdventureGameLocationRec(rec *adventure_game_record.AdventureGameLocation) (*adventure_game_record.AdventureGameLocation, error) {
	l := m.Logger("CreateAdventureGameLocationRec")

	l.Debug("creating adventure_game_location record >%#v<", rec)

	r := m.AdventureGameLocationRepository()

	// Add validation here if needed

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateAdventureGameLocationRec -
func (m *Domain) UpdateAdventureGameLocationRec(next *adventure_game_record.AdventureGameLocation) (*adventure_game_record.AdventureGameLocation, error) {
	l := m.Logger("UpdateAdventureGameLocationRec")

	_, err := m.GetAdventureGameLocationRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}

	l.Debug("updating adventure_game_location record >%#v<", next)

	// Add validation here if needed

	r := m.AdventureGameLocationRepository()

	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}

	return next, nil
}

// DeleteAdventureGameLocationRec -
func (m *Domain) DeleteAdventureGameLocationRec(recID string) error {
	l := m.Logger("DeleteAdventureGameLocationRec")

	l.Debug("deleting adventure_game_location record ID >%s<", recID)

	_, err := m.GetAdventureGameLocationRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameLocationRepository()

	// Add validation here if needed

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAdventureGameLocationRec -
func (m *Domain) RemoveAdventureGameLocationRec(recID string) error {
	l := m.Logger("RemoveAdventureGameLocationRec")

	l.Debug("removing adventure_game_location record ID >%s<", recID)

	_, err := m.GetAdventureGameLocationRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameLocationRepository()

	// Add validation here if needed

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
