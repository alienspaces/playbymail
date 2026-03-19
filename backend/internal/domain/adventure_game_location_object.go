package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameLocationObjectRecs -
func (m *Domain) GetManyAdventureGameLocationObjectRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameLocationObject, error) {
	l := m.Logger("GetManyAdventureGameLocationObjectRecs")
	l.Debug("getting many adventure_game_location_object records opts >%#v<", opts)
	r := m.AdventureGameLocationObjectRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameLocationObjectRec -
func (m *Domain) GetAdventureGameLocationObjectRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameLocationObject, error) {
	l := m.Logger("GetAdventureGameLocationObjectRec")
	l.Debug("getting adventure_game_location_object record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameLocationObjectRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameLocationObject, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameLocationObjectRec -
func (m *Domain) CreateAdventureGameLocationObjectRec(rec *adventure_game_record.AdventureGameLocationObject) (*adventure_game_record.AdventureGameLocationObject, error) {
	l := m.Logger("CreateAdventureGameLocationObjectRec")
	l.Debug("creating adventure_game_location_object record >%#v<", rec)
	if err := m.validateAdventureGameLocationObjectRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_location_object record >%v<", err)
		return rec, err
	}
	r := m.AdventureGameLocationObjectRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameLocationObjectRec -
func (m *Domain) UpdateAdventureGameLocationObjectRec(rec *adventure_game_record.AdventureGameLocationObject) (*adventure_game_record.AdventureGameLocationObject, error) {
	l := m.Logger("UpdateAdventureGameLocationObjectRec")

	currRec, err := m.GetAdventureGameLocationObjectRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_location_object record >%#v<", rec)

	if err := m.validateAdventureGameLocationObjectRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate adventure_game_location_object record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameLocationObjectRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameLocationObjectRec -
func (m *Domain) DeleteAdventureGameLocationObjectRec(recID string) error {
	l := m.Logger("DeleteAdventureGameLocationObjectRec")
	l.Debug("deleting adventure_game_location_object record ID >%s<", recID)
	_, err := m.GetAdventureGameLocationObjectRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameLocationObjectRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameLocationObjectRec -
func (m *Domain) RemoveAdventureGameLocationObjectRec(recID string) error {
	l := m.Logger("RemoveAdventureGameLocationObjectRec")
	l.Debug("removing adventure_game_location_object record ID >%s<", recID)
	r := m.AdventureGameLocationObjectRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
