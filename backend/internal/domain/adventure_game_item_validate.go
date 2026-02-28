package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameItemArgs struct {
	nextRec *adventure_game_record.AdventureGameItem
	currRec *adventure_game_record.AdventureGameItem
}

func (m *Domain) populateAdventureGameItemValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameItem) (*validateAdventureGameItemArgs, error) {
	args := &validateAdventureGameItemArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameItemRecForCreate(rec *adventure_game_record.AdventureGameItem) error {
	args, err := m.populateAdventureGameItemValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameItemRecForCreate(args)
}

func (m *Domain) validateAdventureGameItemRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameItem) error {
	args, err := m.populateAdventureGameItemValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameItemRecForUpdate(args)
}

func validateAdventureGameItemRecForCreate(args *validateAdventureGameItemArgs) error {
	return validateAdventureGameItemRec(args, false)
}

func validateAdventureGameItemRecForUpdate(args *validateAdventureGameItemArgs) error {
	return validateAdventureGameItemRec(args, true)
}

func validateAdventureGameItemRec(args *validateAdventureGameItemArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameItemName, rec.Name); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameItemDescription, rec.Description); err != nil {
		return err
	}

	return nil
}
