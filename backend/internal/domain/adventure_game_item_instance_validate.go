package domain

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameItemInstanceArgs struct {
	nextRec *adventure_game_record.AdventureGameItemInstance
	currRec *adventure_game_record.AdventureGameItemInstance
}

func (m *Domain) populateAdventureGameItemInstanceValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameItemInstance) (*validateAdventureGameItemInstanceArgs, error) {
	args := &validateAdventureGameItemInstanceArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameItemInstanceRecForCreate(rec *adventure_game_record.AdventureGameItemInstance) error {
	args, err := m.populateAdventureGameItemInstanceValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameItemInstanceRecForCreate(args)
}

func (m *Domain) validateAdventureGameItemInstanceRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameItemInstance) error {
	args, err := m.populateAdventureGameItemInstanceValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameItemInstanceRecForUpdate(args)
}

func validateAdventureGameItemInstanceRecForCreate(args *validateAdventureGameItemInstanceArgs) error {
	return validateAdventureGameItemInstanceRec(args, false)
}

func validateAdventureGameItemInstanceRecForUpdate(args *validateAdventureGameItemInstanceArgs) error {
	return validateAdventureGameItemInstanceRec(args, true)
}

func validateAdventureGameItemInstanceRec(args *validateAdventureGameItemInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, rec.AdventureGameItemID); err != nil {
		return err
	}

	// Validate that exactly one owner ID is set
	ownerCount := 0
	if nullstring.IsValid(rec.AdventureGameLocationInstanceID) {
		if err := domain.ValidateNullUUIDField(adventure_game_record.FieldAdventureGameItemInstanceAdventureGameLocationInstanceID, rec.AdventureGameLocationInstanceID); err != nil {
			return err
		}
		ownerCount++
	}
	if nullstring.IsValid(rec.AdventureGameCharacterInstanceID) {
		if err := domain.ValidateNullUUIDField(adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, rec.AdventureGameCharacterInstanceID); err != nil {
			return err
		}
		ownerCount++
	}
	if nullstring.IsValid(rec.AdventureGameCreatureInstanceID) {
		if err := domain.ValidateNullUUIDField(adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCreatureInstanceID, rec.AdventureGameCreatureInstanceID); err != nil {
			return err
		}
		ownerCount++
	}

	if ownerCount != 1 {
		return InvalidField(
			"owner",
			"",
			"exactly one owner (location, character, or creature instance) must be specified",
		)
	}

	if rec.UsesRemaining < 0 {
		return InvalidField(
			adventure_game_record.FieldAdventureGameItemInstanceUsesRemaining,
			fmt.Sprintf("%d", rec.UsesRemaining),
			"uses_remaining must be zero or greater",
		)
	}

	return nil
}
