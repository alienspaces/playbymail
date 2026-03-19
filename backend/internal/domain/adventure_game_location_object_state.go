package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameLocationObjectStateRecs -
func (m *Domain) GetManyAdventureGameLocationObjectStateRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameLocationObjectState, error) {
	l := m.Logger("GetManyAdventureGameLocationObjectStateRecs")
	l.Debug("getting many adventure_game_location_object_state records opts >%#v<", opts)
	r := m.AdventureGameLocationObjectStateRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameLocationObjectStateRec -
func (m *Domain) GetAdventureGameLocationObjectStateRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameLocationObjectState, error) {
	l := m.Logger("GetAdventureGameLocationObjectStateRec")
	l.Debug("getting adventure_game_location_object_state record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameLocationObjectStateRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameLocationObjectState, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameLocationObjectStateRec -
func (m *Domain) CreateAdventureGameLocationObjectStateRec(rec *adventure_game_record.AdventureGameLocationObjectState) (*adventure_game_record.AdventureGameLocationObjectState, error) {
	l := m.Logger("CreateAdventureGameLocationObjectStateRec")
	l.Debug("creating adventure_game_location_object_state record >%#v<", rec)
	if err := m.validateAdventureGameLocationObjectStateRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_location_object_state record >%v<", err)
		return rec, err
	}
	r := m.AdventureGameLocationObjectStateRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameLocationObjectStateRec -
func (m *Domain) UpdateAdventureGameLocationObjectStateRec(rec *adventure_game_record.AdventureGameLocationObjectState) (*adventure_game_record.AdventureGameLocationObjectState, error) {
	l := m.Logger("UpdateAdventureGameLocationObjectStateRec")

	currRec, err := m.GetAdventureGameLocationObjectStateRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_location_object_state record >%#v<", rec)

	if err := m.validateAdventureGameLocationObjectStateRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate adventure_game_location_object_state record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameLocationObjectStateRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameLocationObjectStateRec -
func (m *Domain) DeleteAdventureGameLocationObjectStateRec(recID string) error {
	l := m.Logger("DeleteAdventureGameLocationObjectStateRec")
	l.Debug("deleting adventure_game_location_object_state record ID >%s<", recID)
	_, err := m.GetAdventureGameLocationObjectStateRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameLocationObjectStateRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameLocationObjectStateRec -
func (m *Domain) RemoveAdventureGameLocationObjectStateRec(recID string) error {
	l := m.Logger("RemoveAdventureGameLocationObjectStateRec")
	l.Debug("removing adventure_game_location_object_state record ID >%s<", recID)
	r := m.AdventureGameLocationObjectStateRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
