package domain

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetManyAdventureGameTurnSheetRecs -
func (m *Domain) GetManyAdventureGameTurnSheetRecs(opts *coresql.Options) ([]*adventure_game_record.AdventureGameTurnSheet, error) {
	l := m.Logger("GetManyAdventureGameTurnSheetRecs")

	l.Debug("getting many adventure_game_turn_sheet records opts >%#v<", opts)

	r := m.AdventureGameTurnSheetRepository()

	recs, err := r.GetMany(opts)
	if err != nil {
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetAdventureGameTurnSheetRec -
func (m *Domain) GetAdventureGameTurnSheetRec(recID string, lock *coresql.Lock) (*adventure_game_record.AdventureGameTurnSheet, error) {
	l := m.Logger("GetAdventureGameTurnSheetRec")

	l.Debug("getting adventure_game_turn_sheet record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.AdventureGameTurnSheetRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(adventure_game_record.TableAdventureGameTurnSheet, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateAdventureGameTurnSheetRec -
func (m *Domain) CreateAdventureGameTurnSheetRec(rec *adventure_game_record.AdventureGameTurnSheet) (*adventure_game_record.AdventureGameTurnSheet, error) {
	l := m.Logger("CreateAdventureGameTurnSheetRec")

	l.Debug("creating adventure_game_turn_sheet record >%#v<", rec)

	r := m.AdventureGameTurnSheetRepository()

	// Add validation here if needed

	var err error
	rec, err = r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateAdventureGameTurnSheetRec -
func (m *Domain) UpdateAdventureGameTurnSheetRec(next *adventure_game_record.AdventureGameTurnSheet) (*adventure_game_record.AdventureGameTurnSheet, error) {
	l := m.Logger("UpdateAdventureGameTurnSheetRec")

	_, err := m.GetAdventureGameTurnSheetRec(next.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return next, err
	}

	l.Debug("updating adventure_game_turn_sheet record >%#v<", next)

	// Add validation here if needed

	r := m.AdventureGameTurnSheetRepository()

	next, err = r.UpdateOne(next)
	if err != nil {
		return next, databaseError(err)
	}

	return next, nil
}

// DeleteAdventureGameTurnSheetRec -
func (m *Domain) DeleteAdventureGameTurnSheetRec(recID string) error {
	l := m.Logger("DeleteAdventureGameTurnSheetRec")

	l.Debug("deleting adventure_game_turn_sheet record ID >%s<", recID)

	_, err := m.GetAdventureGameTurnSheetRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameTurnSheetRepository()

	// Add validation here if needed

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveAdventureGameTurnSheetRec -
func (m *Domain) RemoveAdventureGameTurnSheetRec(recID string) error {
	l := m.Logger("RemoveAdventureGameTurnSheetRec")

	l.Debug("removing adventure_game_turn_sheet record ID >%s<", recID)

	_, err := m.GetAdventureGameTurnSheetRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.AdventureGameTurnSheetRepository()

	// Add validation here if needed

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}
