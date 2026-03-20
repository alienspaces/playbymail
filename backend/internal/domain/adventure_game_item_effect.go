package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameItemEffectRecs -
func (m *Domain) GetManyAdventureGameItemEffectRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameItemEffect, error) {
	l := m.Logger("GetManyAdventureGameItemEffectRecs")
	l.Debug("getting many adventure_game_item_effect records opts >%#v<", opts)
	r := m.AdventureGameItemEffectRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameItemEffectRec -
func (m *Domain) GetAdventureGameItemEffectRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameItemEffect, error) {
	l := m.Logger("GetAdventureGameItemEffectRec")
	l.Debug("getting adventure_game_item_effect record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameItemEffectRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameItemEffect, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameItemEffectRec -
func (m *Domain) CreateAdventureGameItemEffectRec(rec *adventure_game_record.AdventureGameItemEffect) (*adventure_game_record.AdventureGameItemEffect, error) {
	l := m.Logger("CreateAdventureGameItemEffectRec")
	l.Debug("creating adventure_game_item_effect record >%#v<", rec)
	if err := m.validateAdventureGameItemEffectRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_item_effect record >%v<", err)
		return rec, err
	}
	r := m.AdventureGameItemEffectRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameItemEffectRec -
func (m *Domain) UpdateAdventureGameItemEffectRec(rec *adventure_game_record.AdventureGameItemEffect) (*adventure_game_record.AdventureGameItemEffect, error) {
	l := m.Logger("UpdateAdventureGameItemEffectRec")

	currRec, err := m.GetAdventureGameItemEffectRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_item_effect record >%#v<", rec)

	if err := m.validateAdventureGameItemEffectRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate adventure_game_item_effect record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameItemEffectRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameItemEffectRec -
func (m *Domain) DeleteAdventureGameItemEffectRec(recID string) error {
	l := m.Logger("DeleteAdventureGameItemEffectRec")
	l.Debug("deleting adventure_game_item_effect record ID >%s<", recID)
	_, err := m.GetAdventureGameItemEffectRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameItemEffectRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameItemEffectRec -
func (m *Domain) RemoveAdventureGameItemEffectRec(recID string) error {
	l := m.Logger("RemoveAdventureGameItemEffectRec")
	l.Debug("removing adventure_game_item_effect record ID >%s<", recID)
	r := m.AdventureGameItemEffectRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
