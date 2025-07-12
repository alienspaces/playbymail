package domain

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

func ValidateGameCreatureInstance(ctx context.Context, rec *record.GameCreatureInstance) error {
	// Add validation logic here
	return nil
}

// GetManyGameCreatureInstanceRecs -
func (m *Domain) GetManyGameCreatureInstanceRecs(opts *coresql.Options) ([]*record.GameCreatureInstance, error) {
	l := m.Logger("GetManyGameCreatureInstanceRecs")
	l.Debug("getting many game_creature_instance records opts >%#v<", opts)
	r := m.GameCreatureInstanceRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetGameCreatureInstanceRec -
func (m *Domain) GetGameCreatureInstanceRec(recID string, lock *coresql.Lock) (*record.GameCreatureInstance, error) {
	l := m.Logger("GetGameCreatureInstanceRec")
	l.Debug("getting game_creature_instance record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.GameCreatureInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError("game_creature_instance", recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateGameCreatureInstanceRec -
func (m *Domain) CreateGameCreatureInstanceRec(rec *record.GameCreatureInstance) (*record.GameCreatureInstance, error) {
	l := m.Logger("CreateGameCreatureInstanceRec")
	l.Debug("creating game_creature_instance record >%#v<", rec)
	r := m.GameCreatureInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateGameCreatureInstanceRec -
func (m *Domain) UpdateGameCreatureInstanceRec(next *record.GameCreatureInstance) (*record.GameCreatureInstance, error) {
	l := m.Logger("UpdateGameCreatureInstanceRec")
	_, err := m.GetGameCreatureInstanceRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating game_creature_instance record >%#v<", next)
	r := m.GameCreatureInstanceRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteGameCreatureInstanceRec -
func (m *Domain) DeleteGameCreatureInstanceRec(recID string) error {
	l := m.Logger("DeleteGameCreatureInstanceRec")
	l.Debug("deleting game_creature_instance record ID >%s<", recID)
	_, err := m.GetGameCreatureInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameCreatureInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveGameCreatureInstanceRec -
func (m *Domain) RemoveGameCreatureInstanceRec(recID string) error {
	l := m.Logger("RemoveGameCreatureInstanceRec")
	l.Debug("removing game_creature_instance record ID >%s<", recID)
	_, err := m.GetGameCreatureInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameCreatureInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
