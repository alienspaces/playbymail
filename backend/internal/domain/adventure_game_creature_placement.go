package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameCreaturePlacementRecs -
func (m *Domain) GetManyAdventureGameCreaturePlacementRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameCreaturePlacement, error) {
	l := m.Logger("GetManyAdventureGameCreaturePlacementRecs")

	l.Debug("getting many adventure_game_creature_placement records opts >%#v<", opts)

	r := m.AdventureGameCreaturePlacementRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAdventureGameCreaturePlacementRec -
func (m *Domain) GetAdventureGameCreaturePlacementRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameCreaturePlacement, error) {
	l := m.Logger("GetAdventureGameCreaturePlacementRec")

	l.Debug("getting adventure_game_creature_placement record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AdventureGameCreaturePlacementRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameCreaturePlacement, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAdventureGameCreaturePlacementRec -
func (m *Domain) CreateAdventureGameCreaturePlacementRec(rec *adventure_game_record.AdventureGameCreaturePlacement) (*adventure_game_record.AdventureGameCreaturePlacement, error) {
	l := m.Logger("CreateAdventureGameCreaturePlacementRec")

	l.Debug("creating adventure_game_creature_placement record >%#v<", rec)

	if err := m.validateAdventureGameCreaturePlacementRecForCreate(rec); err != nil {
		l.Warn("failed to validate adventure_game_creature_placement record >%v<", err)
		return nil, err
	}

	r := m.AdventureGameCreaturePlacementRepository()

	createdRec, err := r.CreateOne(rec)
	if err != nil {
		return nil, databaseError(err)
	}

	return createdRec, nil
}

// UpdateAdventureGameCreaturePlacementRec -
func (m *Domain) UpdateAdventureGameCreaturePlacementRec(rec *adventure_game_record.AdventureGameCreaturePlacement) (*adventure_game_record.AdventureGameCreaturePlacement, error) {
	l := m.Logger("UpdateAdventureGameCreaturePlacementRec")

	_, err := m.GetAdventureGameCreaturePlacementRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating adventure_game_creature_placement record >%#v<", rec)

	if err := m.validateAdventureGameCreaturePlacementRecForUpdate(rec); err != nil {
		l.Warn("failed to validate adventure_game_creature_placement record >%v<", err)
		return rec, err
	}

	r := m.AdventureGameCreaturePlacementRepository()

	updatedRec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return updatedRec, nil
}

// DeleteAdventureGameCreaturePlacementRec -
func (m *Domain) DeleteAdventureGameCreaturePlacementRec(recID string) error {
	l := m.Logger("DeleteAdventureGameCreaturePlacementRec")

	l.Debug("deleting adventure_game_creature_placement record ID >%s<", recID)

	_, err := m.GetAdventureGameCreaturePlacementRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameCreaturePlacementRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAdventureGameCreaturePlacementRec -
func (m *Domain) RemoveAdventureGameCreaturePlacementRec(recID string) error {
	l := m.Logger("RemoveAdventureGameCreaturePlacementRec")

	l.Debug("removing adventure_game_creature_placement record ID >%s<", recID)

	_, err := m.GetAdventureGameCreaturePlacementRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameCreaturePlacementRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
