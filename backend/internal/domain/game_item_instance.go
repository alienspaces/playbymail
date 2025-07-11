package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyGameItemInstanceRecs -
func (m *Domain) GetManyGameItemInstanceRecs(opts *coresql.Options) ([]*record.GameItemInstance, error) {
	l := m.Logger("GetManyGameItemInstanceRecs")
	l.Debug("getting many game_item_instance records opts >%#v<", opts)
	r := m.GameItemInstanceRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetGameItemInstanceRec -
func (m *Domain) GetGameItemInstanceRec(recID string, lock *coresql.Lock) (*record.GameItemInstance, error) {
	l := m.Logger("GetGameItemInstanceRec")
	l.Debug("getting game_item_instance record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.GameItemInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableGameItemInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateGameItemInstanceRec -
func (m *Domain) CreateGameItemInstanceRec(rec *record.GameItemInstance) (*record.GameItemInstance, error) {
	l := m.Logger("CreateGameItemInstanceRec")
	l.Debug("creating game_item_instance record >%#v<", rec)
	r := m.GameItemInstanceRepository()
	// Add validation here if needed
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateGameItemInstanceRec -
func (m *Domain) UpdateGameItemInstanceRec(next *record.GameItemInstance) (*record.GameItemInstance, error) {
	l := m.Logger("UpdateGameItemInstanceRec")
	_, err := m.GetGameItemInstanceRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating game_item_instance record >%#v<", next)
	// Add validation here if needed
	r := m.GameItemInstanceRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteGameItemInstanceRec -
func (m *Domain) DeleteGameItemInstanceRec(recID string) error {
	l := m.Logger("DeleteGameItemInstanceRec")
	l.Debug("deleting game_item_instance record ID >%s<", recID)
	_, err := m.GetGameItemInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameItemInstanceRepository()
	// Add validation here if needed
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveGameItemInstanceRec -
func (m *Domain) RemoveGameItemInstanceRec(recID string) error {
	l := m.Logger("RemoveGameItemInstanceRec")
	l.Debug("removing game_item_instance record ID >%s<", recID)
	_, err := m.GetGameItemInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameItemInstanceRepository()
	// Add validation here if needed
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
