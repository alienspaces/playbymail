package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameCreatureInstanceRecs -
func (m *Domain) GetManyAdventureGameCreatureInstanceRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameCreatureInstance, error) {
	l := m.Logger("GetManyAdventureGameCreatureInstanceRecs")

	l.Debug("getting many adventure_game_creature_instance records opts >%#v<", opts)

	r := m.AdventureGameCreatureInstanceRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAdventureGameCreatureInstanceRec -
func (m *Domain) GetAdventureGameCreatureInstanceRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	l := m.Logger("GetAdventureGameCreatureInstanceRec")

	l.Debug("getting adventure_game_creature_instance record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AdventureGameCreatureInstanceRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameCreatureInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAdventureGameCreatureInstanceRec -
func (m *Domain) CreateAdventureGameCreatureInstanceRec(rec *adventure_game_record.AdventureGameCreatureInstance) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	l := m.Logger("CreateAdventureGameCreatureInstanceRec")

	l.Debug("creating adventure_game_creature_instance record >%#v<", rec)

	if err := m.validateAdventureGameCreatureInstanceRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_creature_instance record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameCreatureInstanceRepository()

	var err error
	createdRec, err := r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return createdRec, nil
}

// UpdateAdventureGameCreatureInstanceRec -
func (m *Domain) UpdateAdventureGameCreatureInstanceRec(rec *adventure_game_record.AdventureGameCreatureInstance) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	l := m.Logger("UpdateAdventureGameCreatureInstanceRec")

	currRec, err := m.GetAdventureGameCreatureInstanceRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_creature_instance record >%#v<", rec)

	if err := m.validateAdventureGameCreatureInstanceRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate adventure_game_creature_instance record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameCreatureInstanceRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameCreatureInstanceRec -
func (m *Domain) DeleteAdventureGameCreatureInstanceRec(recID string) error {
	l := m.Logger("DeleteAdventureGameCreatureInstanceRec")

	l.Debug("deleting adventure_game_creature_instance record ID >%s<", recID)

	_, err := m.GetAdventureGameCreatureInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameCreatureInstanceRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAdventureGameCreatureInstanceRec -
func (m *Domain) RemoveAdventureGameCreatureInstanceRec(recID string) error {
	l := m.Logger("RemoveAdventureGameCreatureInstanceRec")

	l.Debug("removing adventure_game_creature_instance record ID >%s<", recID)

	_, err := m.GetAdventureGameCreatureInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameCreatureInstanceRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
