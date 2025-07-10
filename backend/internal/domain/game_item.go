package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyGameItemRecs -
func (m *Domain) GetManyGameItemRecs(opts *coresql.Options) ([]*record.GameItem, error) {
	l := m.Logger("GetManyGameItemRecs")
	l.Debug("getting many game_item records opts >%#v<", opts)
	r := m.GameItemRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetGameItemRec -
func (m *Domain) GetGameItemRec(recID string, lock *coresql.Lock) (*record.GameItem, error) {
	l := m.Logger("GetGameItemRec")
	l.Debug("getting game_item record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.GameItemRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableGameItem, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateGameItemRec -
func (m *Domain) CreateGameItemRec(rec *record.GameItem) (*record.GameItem, error) {
	l := m.Logger("CreateGameItemRec")
	l.Debug("creating game_item record >%#v<", rec)
	r := m.GameItemRepository()
	// Add validation here if needed
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateGameItemRec -
func (m *Domain) UpdateGameItemRec(next *record.GameItem) (*record.GameItem, error) {
	l := m.Logger("UpdateGameItemRec")
	_, err := m.GetGameItemRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating game_item record >%#v<", next)
	// Add validation here if needed
	r := m.GameItemRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteGameItemRec -
func (m *Domain) DeleteGameItemRec(recID string) error {
	l := m.Logger("DeleteGameItemRec")
	l.Debug("deleting game_item record ID >%s<", recID)
	_, err := m.GetGameItemRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameItemRepository()
	// Add validation here if needed
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveGameItemRec -
func (m *Domain) RemoveGameItemRec(recID string) error {
	l := m.Logger("RemoveGameItemRec")
	l.Debug("removing game_item record ID >%s<", recID)
	_, err := m.GetGameItemRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.GameItemRepository()
	// Add validation here if needed
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
