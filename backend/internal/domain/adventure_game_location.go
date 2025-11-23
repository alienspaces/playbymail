package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
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

	// Validate starting location constraint
	if rec.IsStartingLocation {
		if err := m.validateStartingLocationConstraint(rec.GameID, ""); err != nil {
			return rec, err
		}
	}

	r := m.AdventureGameLocationRepository()

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

	curr, err := m.GetAdventureGameLocationRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}

	l.Debug("updating adventure_game_location record >%#v<", next)

	// Validate starting location constraint if setting to true
	if next.IsStartingLocation && !curr.IsStartingLocation {
		if err := m.validateStartingLocationConstraint(next.GameID, next.ID); err != nil {
			return next, err
		}
	}

	r := m.AdventureGameLocationRepository()

	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}

	return next, nil
}

// validateStartingLocationConstraint ensures only one starting location exists per game
// excludeID is used when updating to exclude the current record from the check
func (m *Domain) validateStartingLocationConstraint(gameID, excludeID string) error {
	l := m.Logger("validateStartingLocationConstraint")

	l.Info("validating starting location constraint for game ID >%s< exclude ID >%s<", gameID, excludeID)

	opts := &coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationGameID, Val: gameID},
			{Col: adventure_game_record.FieldAdventureGameLocationIsStartingLocation, Val: true},
		},
		Limit: 2,
	}

	existingRecs, err := m.GetManyAdventureGameLocationRecs(opts)
	if err != nil {
		l.Warn("failed to get many adventure game location records >%v<", err)
		return err
	}

	// Filter out the excluded ID if provided
	var count int
	for _, rec := range existingRecs {
		if rec.ID != excludeID {
			count++
		}
	}

	if count > 0 {
		return InvalidField(adventure_game_record.FieldAdventureGameLocationIsStartingLocation, "true", "only one starting location is allowed per game")
	}

	return nil
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
