package domain

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
)

func ValidateAdventureGameCreatureInstance(ctx context.Context, rec *adventure_game_record.AdventureGameCreatureInstance) error {
	// Add validation logic here
	return nil
}

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
		return nil, coreerror.NewNotFoundError("adventure_game_creature_instance", recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameCreatureInstanceRec -
func (m *Domain) CreateAdventureGameCreatureInstanceRec(rec *adventure_game_record.AdventureGameCreatureInstance) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	l := m.Logger("CreateAdventureGameCreatureInstanceRec")
	l.Debug("creating adventure_game_creature_instance record >%#v<", rec)
	r := m.AdventureGameCreatureInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameCreatureInstanceRec -
func (m *Domain) UpdateAdventureGameCreatureInstanceRec(next *adventure_game_record.AdventureGameCreatureInstance) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	l := m.Logger("UpdateAdventureGameCreatureInstanceRec")
	_, err := m.GetAdventureGameCreatureInstanceRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating adventure_game_creature_instance record >%#v<", next)
	r := m.AdventureGameCreatureInstanceRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
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
