package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameCharacterInstanceRecs -
func (m *Domain) GetManyAdventureGameCharacterInstanceRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := m.Logger("GetManyAdventureGameCharacterInstanceRecs")
	l.Debug("getting many adventure_game_character_instance records opts >%#v<", opts)
	r := m.AdventureGameCharacterInstanceRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameCharacterInstanceRec -
func (m *Domain) GetAdventureGameCharacterInstanceRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := m.Logger("GetAdventureGameCharacterInstanceRec")
	l.Debug("getting adventure_game_character_instance record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameCharacterInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameCharacterInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameCharacterInstanceRec -
func (m *Domain) CreateAdventureGameCharacterInstanceRec(rec *adventure_game_record.AdventureGameCharacterInstance) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := m.Logger("CreateAdventureGameCharacterInstanceRec")
	l.Debug("creating adventure_game_character_instance record >%#v<", rec)
	r := m.AdventureGameCharacterInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameCharacterInstanceRec -
func (m *Domain) UpdateAdventureGameCharacterInstanceRec(next *adventure_game_record.AdventureGameCharacterInstance) (*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := m.Logger("UpdateAdventureGameCharacterInstanceRec")
	_, err := m.GetAdventureGameCharacterInstanceRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating adventure_game_character_instance record >%#v<", next)
	r := m.AdventureGameCharacterInstanceRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteAdventureGameCharacterInstanceRec -
func (m *Domain) DeleteAdventureGameCharacterInstanceRec(recID string) error {
	l := m.Logger("DeleteAdventureGameCharacterInstanceRec")
	l.Debug("deleting adventure_game_character_instance record ID >%s<", recID)
	_, err := m.GetAdventureGameCharacterInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameCharacterInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameCharacterInstanceRec -
func (m *Domain) RemoveAdventureGameCharacterInstanceRec(recID string) error {
	l := m.Logger("RemoveAdventureGameCharacterInstanceRec")
	l.Debug("removing adventure_game_character_instance record ID >%s<", recID)
	_, err := m.GetAdventureGameCharacterInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameCharacterInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
