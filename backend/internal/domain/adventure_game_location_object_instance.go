package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameLocationObjectInstanceRecs -
func (m *Domain) GetManyAdventureGameLocationObjectInstanceRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameLocationObjectInstance, error) {
	l := m.Logger("GetManyAdventureGameLocationObjectInstanceRecs")
	l.Debug("getting many adventure_game_location_object_instance records opts >%#v<", opts)
	r := m.AdventureGameLocationObjectInstanceRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameLocationObjectInstanceRec -
func (m *Domain) GetAdventureGameLocationObjectInstanceRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameLocationObjectInstance, error) {
	l := m.Logger("GetAdventureGameLocationObjectInstanceRec")
	l.Debug("getting adventure_game_location_object_instance record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameLocationObjectInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameLocationObjectInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameLocationObjectInstanceRec -
func (m *Domain) CreateAdventureGameLocationObjectInstanceRec(rec *adventure_game_record.AdventureGameLocationObjectInstance) (*adventure_game_record.AdventureGameLocationObjectInstance, error) {
	l := m.Logger("CreateAdventureGameLocationObjectInstanceRec")
	l.Debug("creating adventure_game_location_object_instance record >%#v<", rec)
	if err := m.validateAdventureGameLocationObjectInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_location_object_instance record >%v<", err)
		return rec, err
	}
	r := m.AdventureGameLocationObjectInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameLocationObjectInstanceRec -
func (m *Domain) UpdateAdventureGameLocationObjectInstanceRec(rec *adventure_game_record.AdventureGameLocationObjectInstance) (*adventure_game_record.AdventureGameLocationObjectInstance, error) {
	l := m.Logger("UpdateAdventureGameLocationObjectInstanceRec")

	currRec, err := m.GetAdventureGameLocationObjectInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_location_object_instance record >%#v<", rec)

	if err := m.validateAdventureGameLocationObjectInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate adventure_game_location_object_instance record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameLocationObjectInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameLocationObjectInstanceRec -
func (m *Domain) DeleteAdventureGameLocationObjectInstanceRec(recID string) error {
	l := m.Logger("DeleteAdventureGameLocationObjectInstanceRec")
	l.Debug("deleting adventure_game_location_object_instance record ID >%s<", recID)
	_, err := m.GetAdventureGameLocationObjectInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameLocationObjectInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameLocationObjectInstanceRec -
func (m *Domain) RemoveAdventureGameLocationObjectInstanceRec(recID string) error {
	l := m.Logger("RemoveAdventureGameLocationObjectInstanceRec")
	l.Debug("removing adventure_game_location_object_instance record ID >%s<", recID)
	r := m.AdventureGameLocationObjectInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
