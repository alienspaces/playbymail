package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyGameCreatureRecs -
func (m *Domain) GetManyGameCreatureRecs(opts *coresql.Options) ([]*record.GameCreature, error) {
	l := m.Logger("GetManyGameCreatureRecs")
	l.Debug("getting many game_creature records opts >%#v<", opts)
	r := m.GameCreatureRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetGameCreatureRec -
func (m *Domain) GetGameCreatureRec(recID string, lock *coresql.Lock) (*record.GameCreature, error) {
	l := m.Logger("GetGameCreatureRec")
	l.Debug("getting game_creature record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.GameCreatureRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableGameCreature, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateGameCreatureRec -
func (m *Domain) CreateGameCreatureRec(rec *record.GameCreature) (*record.GameCreature, error) {
	l := m.Logger("CreateGameCreatureRec")
	l.Debug("creating game_creature record >%#v<", rec)
	r := m.GameCreatureRepository()
	// Add validation here if needed
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateGameCreatureRec -
func (m *Domain) UpdateGameCreatureRec(next *record.GameCreature) (*record.GameCreature, error) {
	l := m.Logger("UpdateGameCreatureRec")
	_, err := m.GetGameCreatureRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating game_creature record >%#v<", next)
	// Add validation here if needed
	r := m.GameCreatureRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteGameCreatureRec -
func (m *Domain) DeleteGameCreatureRec(recID string) error {
	l := m.Logger("DeleteGameCreatureRec")
	l.Debug("deleting game_creature record ID >%s<", recID)
	_, err := m.GetGameCreatureRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameCreatureRepository()
	// Add validation here if needed
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveGameCreatureRec -
func (m *Domain) RemoveGameCreatureRec(recID string) error {
	l := m.Logger("RemoveGameCreatureRec")
	l.Debug("removing game_creature record ID >%s<", recID)
	_, err := m.GetGameCreatureRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameCreatureRepository()
	// Add validation here if needed
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
