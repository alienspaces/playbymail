package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
)

// GetManyAdventureGameItemInstanceRecs -
func (m *Domain) GetManyAdventureGameItemInstanceRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("GetManyAdventureGameItemInstanceRecs")
	l.Debug("getting many adventure_game_item_instance records opts >%#v<", opts)
	r := m.AdventureGameItemInstanceRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameItemInstanceRec -
func (m *Domain) GetAdventureGameItemInstanceRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("GetAdventureGameItemInstanceRec")
	l.Debug("getting adventure_game_item_instance record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameItemInstanceRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameItemInstance, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameItemInstanceRec -
func (m *Domain) CreateAdventureGameItemInstanceRec(rec *adventure_game_record.AdventureGameItemInstance) (*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("CreateAdventureGameItemInstanceRec")
	l.Debug("creating adventure_game_item_instance record >%#v<", rec)
	r := m.AdventureGameItemInstanceRepository()
	// Add validation here if needed
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameItemInstanceRec -
func (m *Domain) UpdateAdventureGameItemInstanceRec(next *adventure_game_record.AdventureGameItemInstance) (*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("UpdateAdventureGameItemInstanceRec")
	_, err := m.GetAdventureGameItemInstanceRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}
	l.Debug("updating adventure_game_item_instance record >%#v<", next)
	// Add validation here if needed
	r := m.AdventureGameItemInstanceRepository()
	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}
	return next, nil
}

// DeleteAdventureGameItemInstanceRec -
func (m *Domain) DeleteAdventureGameItemInstanceRec(recID string) error {
	l := m.Logger("DeleteAdventureGameItemInstanceRec")
	l.Debug("deleting adventure_game_item_instance record ID >%s<", recID)
	_, err := m.GetAdventureGameItemInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameItemInstanceRepository()
	// Add validation here if needed
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameItemInstanceRec -
func (m *Domain) RemoveAdventureGameItemInstanceRec(recID string) error {
	l := m.Logger("RemoveAdventureGameItemInstanceRec")
	l.Debug("removing adventure_game_item_instance record ID >%s<", recID)
	_, err := m.GetAdventureGameItemInstanceRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameItemInstanceRepository()
	// Add validation here if needed
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
