package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyAdventureGameItemRecs -
func (m *Domain) GetManyAdventureGameItemRecs(opts *coresql.Options) ([]*record.AdventureGameItem, error) {
	l := m.Logger("GetManyAdventureGameItemRecs")
	l.Debug("getting many adventure_game_item records opts >%#v<", opts)
	r := m.AdventureGameItemRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameItemRec -
func (m *Domain) GetAdventureGameItemRec(recID string, lock *coresql.Lock) (*record.AdventureGameItem, error) {
	l := m.Logger("GetAdventureGameItemRec")
	l.Debug("getting adventure_game_item record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameItemRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(record.TableAdventureGameItem, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameItemRec -
func (m *Domain) CreateAdventureGameItemRec(rec *record.AdventureGameItem) (*record.AdventureGameItem, error) {
	l := m.Logger("CreateAdventureGameItemRec")
	l.Debug("creating adventure_game_item record >%#v<", rec)
	r := m.AdventureGameItemRepository()
	// Add validation here if needed
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameItemRec -
func (m *Domain) UpdateAdventureGameItemRec(next *record.AdventureGameItem) (*record.AdventureGameItem, error) {
	l := m.Logger("UpdateAdventureGameItemRec")
	_, err := m.GetAdventureGameItemRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating adventure_game_item record >%#v<", next)
	// Add validation here if needed
	r := m.AdventureGameItemRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteAdventureGameItemRec -
func (m *Domain) DeleteAdventureGameItemRec(recID string) error {
	l := m.Logger("DeleteAdventureGameItemRec")
	l.Debug("deleting adventure_game_item record ID >%s<", recID)
	_, err := m.GetAdventureGameItemRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameItemRepository()
	// Add validation here if needed
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameItemRec -
func (m *Domain) RemoveAdventureGameItemRec(recID string) error {
	l := m.Logger("RemoveAdventureGameItemRec")
	l.Debug("removing adventure_game_item record ID >%s<", recID)
	_, err := m.GetAdventureGameItemRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameItemRepository()
	// Add validation here if needed
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
