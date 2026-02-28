package domain

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameCreaturePlacementArgs struct {
	nextRec *adventure_game_record.AdventureGameCreaturePlacement
	currRec *adventure_game_record.AdventureGameCreaturePlacement
}

func (m *Domain) populateAdventureGameCreaturePlacementValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameCreaturePlacement) (*validateAdventureGameCreaturePlacementArgs, error) {
	args := &validateAdventureGameCreaturePlacementArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameCreaturePlacementRecForCreate(rec *adventure_game_record.AdventureGameCreaturePlacement) error {
	args, err := m.populateAdventureGameCreaturePlacementValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameCreaturePlacementRecForCreate(args)
}

func (m *Domain) validateAdventureGameCreaturePlacementRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameCreaturePlacement) error {
	args, err := m.populateAdventureGameCreaturePlacementValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameCreaturePlacementRecForUpdate(args)
}

func validateAdventureGameCreaturePlacementRecForCreate(args *validateAdventureGameCreaturePlacementArgs) error {
	return validateAdventureGameCreaturePlacementRec(args, false)
}

func validateAdventureGameCreaturePlacementRecForUpdate(args *validateAdventureGameCreaturePlacementArgs) error {
	return validateAdventureGameCreaturePlacementRec(args, true)
}

func validateAdventureGameCreaturePlacementRec(args *validateAdventureGameCreaturePlacementArgs, requireID bool) error {
	rec := args.nextRec

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

