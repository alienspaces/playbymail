package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameLocationInstanceRecs -
func (m *Domain) GetManyAdventureGameLocationInstanceRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameLocationInstance, error) {
	l := m.Logger("GetManyAdventureGameLocationInstanceRecs")

	l.Debug("getting many adventure_game_location_instance records opts >%#v<", opts)

	r := m.AdventureGameLocationInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAdventureGameLocationInstanceRec -
func (m *Domain) GetAdventureGameLocationInstanceRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameLocationInstance, error) {
	l := m.Logger("GetAdventureGameLocationInstanceRec")

	l.Debug("getting adventure_game_location_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AdventureGameLocationInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameLocationInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAdventureGameLocationInstanceRec -
func (m *Domain) CreateAdventureGameLocationInstanceRec(rec *adventure_game_record.AdventureGameLocationInstance) (*adventure_game_record.AdventureGameLocationInstance, error) {
	l := m.Logger("CreateAdventureGameLocationInstanceRec")

	l.Debug("creating adventure_game_location_instance record >%#v<", rec)

	if err := m.validateAdventureGameLocationInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_location_instance record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameLocationInstanceRepository()

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateAdventureGameLocationInstanceRec -
func (m *Domain) UpdateAdventureGameLocationInstanceRec(rec *adventure_game_record.AdventureGameLocationInstance) (*adventure_game_record.AdventureGameLocationInstance, error) {
	l := m.Logger("UpdateAdventureGameLocationInstanceRec")

	currRec, err := m.GetAdventureGameLocationInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_location_instance record >%#v<", rec)

	if err := m.validateAdventureGameLocationInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate adventure_game_location_instance record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameLocationInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameLocationInstanceRec -
func (m *Domain) DeleteAdventureGameLocationInstanceRec(recID string) error {
	l := m.Logger("DeleteAdventureGameLocationInstanceRec")

	l.Debug("deleting adventure_game_location_instance record ID >%s<", recID)

	_, err := m.GetAdventureGameLocationInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameLocationInstanceRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAdventureGameLocationInstanceRec -
func (m *Domain) RemoveAdventureGameLocationInstanceRec(recID string) error {
	l := m.Logger("RemoveAdventureGameLocationInstanceRec")

	l.Debug("removing adventure_game_location_instance record ID >%s<", recID)

	_, err := m.GetAdventureGameLocationInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameLocationInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
}

	return nil
}
