package domain

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameItemPlacementArgs struct {
	nextRec *adventure_game_record.AdventureGameItemPlacement
	currRec *adventure_game_record.AdventureGameItemPlacement
}

func (m *Domain) populateAdventureGameItemPlacementValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameItemPlacement) (*validateAdventureGameItemPlacementArgs, error) {
	args := &validateAdventureGameItemPlacementArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameItemPlacementRecForCreate(rec *adventure_game_record.AdventureGameItemPlacement) error {
	args, err := m.populateAdventureGameItemPlacementValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameItemPlacementRecForCreate(args)
}

func (m *Domain) validateAdventureGameItemPlacementRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameItemPlacement) error {
	args, err := m.populateAdventureGameItemPlacementValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameItemPlacementRecForUpdate(args)
}

func validateAdventureGameItemPlacementRecForCreate(args *validateAdventureGameItemPlacementArgs) error {
	return validateAdventureGameItemPlacementRec(args, false)
}

func validateAdventureGameItemPlacementRecForUpdate(args *validateAdventureGameItemPlacementArgs) error {
	return validateAdventureGameItemPlacementRec(args, true)
}

func validateAdventureGameItemPlacementRec(args *validateAdventureGameItemPlacementArgs, requireID bool) error {
	rec := args.nextRec

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
