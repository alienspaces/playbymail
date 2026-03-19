package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameLocationObjectEffectRecs -
func (m *Domain) GetManyAdventureGameLocationObjectEffectRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameLocationObjectEffect, error) {
	l := m.Logger("GetManyAdventureGameLocationObjectEffectRecs")
	l.Debug("getting many adventure_game_location_object_effect records opts >%#v<", opts)
	r := m.AdventureGameLocationObjectEffectRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameLocationObjectEffectRec -
func (m *Domain) GetAdventureGameLocationObjectEffectRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameLocationObjectEffect, error) {
	l := m.Logger("GetAdventureGameLocationObjectEffectRec")
	l.Debug("getting adventure_game_location_object_effect record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameLocationObjectEffectRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameLocationObjectEffect, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameLocationObjectEffectRec -
func (m *Domain) CreateAdventureGameLocationObjectEffectRec(rec *adventure_game_record.AdventureGameLocationObjectEffect) (*adventure_game_record.AdventureGameLocationObjectEffect, error) {
	l := m.Logger("CreateAdventureGameLocationObjectEffectRec")
	l.Debug("creating adventure_game_location_object_effect record >%#v<", rec)
	if err := m.validateAdventureGameLocationObjectEffectRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_location_object_effect record >%v<", err)
		return rec, err
	}
	r := m.AdventureGameLocationObjectEffectRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameLocationObjectEffectRec -
func (m *Domain) UpdateAdventureGameLocationObjectEffectRec(rec *adventure_game_record.AdventureGameLocationObjectEffect) (*adventure_game_record.AdventureGameLocationObjectEffect, error) {
	l := m.Logger("UpdateAdventureGameLocationObjectEffectRec")

	currRec, err := m.GetAdventureGameLocationObjectEffectRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_location_object_effect record >%#v<", rec)

	if err := m.validateAdventureGameLocationObjectEffectRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate adventure_game_location_object_effect record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameLocationObjectEffectRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameLocationObjectEffectRec -
func (m *Domain) DeleteAdventureGameLocationObjectEffectRec(recID string) error {
	l := m.Logger("DeleteAdventureGameLocationObjectEffectRec")
	l.Debug("deleting adventure_game_location_object_effect record ID >%s<", recID)
	_, err := m.GetAdventureGameLocationObjectEffectRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameLocationObjectEffectRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameLocationObjectEffectRec -
func (m *Domain) RemoveAdventureGameLocationObjectEffectRec(recID string) error {
	l := m.Logger("RemoveAdventureGameLocationObjectEffectRec")
	l.Debug("removing adventure_game_location_object_effect record ID >%s<", recID)
	r := m.AdventureGameLocationObjectEffectRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
