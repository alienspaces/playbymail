package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameItemPlacementRecs -
func (m *Domain) GetManyAdventureGameItemPlacementRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameItemPlacement, error) {
	l := m.Logger("GetManyAdventureGameItemPlacementRecs")
	l.Debug("getting many adventure_game_item_placement records opts >%#v<", opts)
	r := m.AdventureGameItemPlacementRepository()
	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}
	return recs, nil
}

// GetAdventureGameItemPlacementRec -
func (m *Domain) GetAdventureGameItemPlacementRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameItemPlacement, error) {
	l := m.Logger("GetAdventureGameItemPlacementRec")
	l.Debug("getting adventure_game_item_placement record ID >%s<", recID)
	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}
	r := m.AdventureGameItemPlacementRepository()
	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameItemPlacement, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameItemPlacementRec -
func (m *Domain) CreateAdventureGameItemPlacementRec(rec *adventure_game_record.AdventureGameItemPlacement) (*adventure_game_record.AdventureGameItemPlacement, error) {
	l := m.Logger("CreateAdventureGameItemPlacementRec")
	l.Debug("creating adventure_game_item_placement record >%#v<", rec)
	if err := m.validateAdventureGameItemPlacementRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_item_placement record >%v<", err)
		return rec, err
	}
	r := m.AdventureGameItemPlacementRepository()
	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameItemPlacementRec -
func (m *Domain) UpdateAdventureGameItemPlacementRec(rec *adventure_game_record.AdventureGameItemPlacement) (*adventure_game_record.AdventureGameItemPlacement, error) {
	l := m.Logger("UpdateAdventureGameItemPlacementRec")

	_, err := m.GetAdventureGameItemPlacementRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_item_placement record >%#v<", rec)

	if err := m.validateAdventureGameItemPlacementRecForUpdate(rec); err != nil {
		l.Warn("failed to validate adventure_game_item_placement record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameItemPlacementRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameItemPlacementRec -
func (m *Domain) DeleteAdventureGameItemPlacementRec(recID string) error {
	l := m.Logger("DeleteAdventureGameItemPlacementRec")
	l.Debug("deleting adventure_game_item_placement record ID >%s<", recID)
	_, err := m.GetAdventureGameItemPlacementRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameItemPlacementRepository()
	// Add validation here if needed
	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}

// RemoveAdventureGameItemPlacementRec -
func (m *Domain) RemoveAdventureGameItemPlacementRec(recID string) error {
	l := m.Logger("RemoveAdventureGameItemPlacementRec")
	l.Debug("removing adventure_game_item_placement record ID >%s<", recID)
	_, err := m.GetAdventureGameItemPlacementRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}
	r := m.AdventureGameItemPlacementRepository()
	// Add validation here if needed
	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}
	return nil
}
