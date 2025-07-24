package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

// GetManyAdventureGameCreaturePlacementRecs -
func (m *Domain) GetManyAdventureGameCreaturePlacementRecs(opts *coresql.Options) ([]*record.AdventureGameCreaturePlacement, error) {
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
func (m *Domain) GetAdventureGameCreaturePlacementRec(id string, lock *coresql.Lock) (*record.AdventureGameCreaturePlacement, error) {
	l := m.Logger("GetAdventureGameCreaturePlacementRec")
	l.Debug("getting adventure_game_creature_placement record id >%s< lock >%#v<", id, lock)
	r := m.AdventureGameCreaturePlacementRepository()
	rec, err := r.GetOne(id, lock)
	if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// CreateAdventureGameCreaturePlacementRec -
func (m *Domain) CreateAdventureGameCreaturePlacementRec(rec *record.AdventureGameCreaturePlacement) (*record.AdventureGameCreaturePlacement, error) {
	l := m.Logger("CreateAdventureGameCreaturePlacementRec")
	l.Debug("creating adventure_game_creature_placement record >%#v<", rec)

	if rec == nil {
		return nil, coreerror.NewInvalidDataError("record is nil")
	}

	if rec.GameID == "" {
		return nil, coreerror.NewInvalidDataError("game_id is required")
	}

	if rec.AdventureGameCreatureID == "" {
		return nil, coreerror.NewInvalidDataError("adventure_game_creature_id is required")
	}

	if rec.AdventureGameLocationID == "" {
		return nil, coreerror.NewInvalidDataError("adventure_game_location_id is required")
	}

	if rec.InitialCount < 0 {
		return nil, coreerror.NewInvalidDataError("initial_count must be non-negative")
	}

	r := m.AdventureGameCreaturePlacementRepository()
	rec, err := r.CreateOne(rec)
	if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// UpdateAdventureGameCreaturePlacementRec -
func (m *Domain) UpdateAdventureGameCreaturePlacementRec(rec *record.AdventureGameCreaturePlacement) (*record.AdventureGameCreaturePlacement, error) {
	l := m.Logger("UpdateAdventureGameCreaturePlacementRec")
	l.Debug("updating adventure_game_creature_placement record >%#v<", rec)

	if rec == nil {
		return nil, coreerror.NewInvalidDataError("record is nil")
	}

	if rec.ID == "" {
		return nil, coreerror.NewInvalidDataError("id is required")
	}

	if rec.AdventureGameCreatureID == "" {
		return nil, coreerror.NewInvalidDataError("adventure_game_creature_id is required")
	}

	if rec.AdventureGameLocationID == "" {
		return nil, coreerror.NewInvalidDataError("adventure_game_location_id is required")
	}

	if rec.InitialCount < 0 {
		return nil, coreerror.NewInvalidDataError("initial_count must be non-negative")
	}

	r := m.AdventureGameCreaturePlacementRepository()
	rec, err := r.UpdateOne(rec)
	if err != nil {
		return nil, databaseError(err)
	}
	return rec, nil
}

// DeleteAdventureGameCreaturePlacementRec -
func (m *Domain) DeleteAdventureGameCreaturePlacementRec(id string) error {
	l := m.Logger("DeleteAdventureGameCreaturePlacementRec")
	l.Debug("deleting adventure_game_creature_placement record id >%s<", id)

	if id == "" {
		return coreerror.NewInvalidDataError("id is required")
	}

	r := m.AdventureGameCreaturePlacementRepository()
	err := r.DeleteOne(id)
	if err != nil {
		return databaseError(err)
	}
	return nil
}
