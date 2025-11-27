package domain

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameItemPlacementRecForCreate(rec *adventure_game_record.AdventureGameItemPlacement) error {
	return validateAdventureGameItemPlacementRec(rec, false)
}

func (m *Domain) validateAdventureGameItemPlacementRecForUpdate(rec *adventure_game_record.AdventureGameItemPlacement) error {
	return validateAdventureGameItemPlacementRec(rec, true)
}

func validateAdventureGameItemPlacementRec(rec *adventure_game_record.AdventureGameItemPlacement, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemPlacementID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemPlacementGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemPlacementAdventureGameItemID, rec.AdventureGameItemID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemPlacementAdventureGameLocationID, rec.AdventureGameLocationID); err != nil {
		return err
	}

	if rec.InitialCount < 0 {
		return InvalidField(
			adventure_game_record.FieldAdventureGameItemPlacementInitialCount,
			fmt.Sprintf("%d", rec.InitialCount),
			"initial_count must be zero or greater",
		)
	}

	return nil
}
