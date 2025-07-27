package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameCreatureRecs -
func (m *Domain) GetManyAdventureGameCreatureRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameCreature, error) {
	l := m.Logger("GetManyAdventureGameCreatureRecs")
	l.Debug("getting many adventure_game_creature records opts >%#v<", opts)
	r := m.AdventureGameCreatureRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameCreatureRec -
func (m *Domain) GetAdventureGameCreatureRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameCreature, error) {
	l := m.Logger("GetAdventureGameCreatureRec")
	l.Debug("getting adventure_game_creature record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameCreatureRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameCreature, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameCreatureRec -
func (m *Domain) CreateAdventureGameCreatureRec(rec *adventure_game_record.AdventureGameCreature) (*adventure_game_record.AdventureGameCreature, error) {
	l := m.Logger("CreateAdventureGameCreatureRec")
	l.Debug("creating adventure_game_creature record >%#v<", rec)
	r := m.AdventureGameCreatureRepository()
	// Add validation here if needed
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameCreatureRec -
func (m *Domain) UpdateAdventureGameCreatureRec(next *adventure_game_record.AdventureGameCreature) (*adventure_game_record.AdventureGameCreature, error) {
	l := m.Logger("UpdateAdventureGameCreatureRec")
	_, err := m.GetAdventureGameCreatureRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating adventure_game_creature record >%#v<", next)
	// Add validation here if needed
	r := m.AdventureGameCreatureRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteAdventureGameCreatureRec -
func (m *Domain) DeleteAdventureGameCreatureRec(recID string) error {
	l := m.Logger("DeleteAdventureGameCreatureRec")
	l.Debug("deleting adventure_game_creature record ID >%s<", recID)
	_, err := m.GetAdventureGameCreatureRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameCreatureRepository()
	// Add validation here if needed
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameCreatureRec -
func (m *Domain) RemoveAdventureGameCreatureRec(recID string) error {
	l := m.Logger("RemoveAdventureGameCreatureRec")
	l.Debug("removing adventure_game_creature record ID >%s<", recID)
	_, err := m.GetAdventureGameCreatureRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameCreatureRepository()
	// Add validation here if needed
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
