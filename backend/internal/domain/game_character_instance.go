package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyGameCharacterInstanceRecs -
func (m *Domain) GetManyGameCharacterInstanceRecs(opts *coresql.Options) ([]*record.GameCharacterInstance, error) {
	l := m.Logger("GetManyGameCharacterInstanceRecs")
	l.Debug("getting many game_character_instance records opts >%#v<", opts)
	r := m.GameCharacterInstanceRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetGameCharacterInstanceRec -
func (m *Domain) GetGameCharacterInstanceRec(recID string, lock *coresql.Lock) (*record.GameCharacterInstance, error) {
	l := m.Logger("GetGameCharacterInstanceRec")
	l.Debug("getting game_character_instance record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.GameCharacterInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableGameCharacterInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateGameCharacterInstanceRec -
func (m *Domain) CreateGameCharacterInstanceRec(rec *record.GameCharacterInstance) (*record.GameCharacterInstance, error) {
	l := m.Logger("CreateGameCharacterInstanceRec")
	l.Debug("creating game_character_instance record >%#v<", rec)
	r := m.GameCharacterInstanceRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateGameCharacterInstanceRec -
func (m *Domain) UpdateGameCharacterInstanceRec(next *record.GameCharacterInstance) (*record.GameCharacterInstance, error) {
	l := m.Logger("UpdateGameCharacterInstanceRec")
	_, err := m.GetGameCharacterInstanceRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating game_character_instance record >%#v<", next)
	r := m.GameCharacterInstanceRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteGameCharacterInstanceRec -
func (m *Domain) DeleteGameCharacterInstanceRec(recID string) error {
	l := m.Logger("DeleteGameCharacterInstanceRec")
	l.Debug("deleting game_character_instance record ID >%s<", recID)
	_, err := m.GetGameCharacterInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameCharacterInstanceRepository()
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveGameCharacterInstanceRec -
func (m *Domain) RemoveGameCharacterInstanceRec(recID string) error {
	l := m.Logger("RemoveGameCharacterInstanceRec")
	l.Debug("removing game_character_instance record ID >%s<", recID)
	_, err := m.GetGameCharacterInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameCharacterInstanceRepository()
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
