package domain

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameCreaturePlacementRecForCreate(rec *adventure_game_record.AdventureGameCreaturePlacement) error {
	return validateAdventureGameCreaturePlacementRec(rec, false)
}

func (m *Domain) validateAdventureGameCreaturePlacementRecForUpdate(rec *adventure_game_record.AdventureGameCreaturePlacement) error {
	return validateAdventureGameCreaturePlacementRec(rec, true)
}

func validateAdventureGameCreaturePlacementRec(rec *adventure_game_record.AdventureGameCreaturePlacement, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreaturePlacementID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreaturePlacementGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreaturePlacementAdventureGameCreatureID, rec.AdventureGameCreatureID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreaturePlacementAdventureGameLocationID, rec.AdventureGameLocationID); err != nil {
		return err
	}

	if rec.InitialCount < 0 {
		return InvalidField(
			adventure_game_record.FieldAdventureGameCreaturePlacementInitialCount,
			fmt.Sprintf("%d", rec.InitialCount),
			"initial_count must be zero or greater",
		)
	}

	return nil
}

