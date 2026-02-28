package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameTurnSheetArgs struct {
	nextRec *adventure_game_record.AdventureGameTurnSheet
	currRec *adventure_game_record.AdventureGameTurnSheet
}

func (m *Domain) populateAdventureGameTurnSheetValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameTurnSheet) (*validateAdventureGameTurnSheetArgs, error) {
	args := &validateAdventureGameTurnSheetArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameTurnSheetRecForCreate(rec *adventure_game_record.AdventureGameTurnSheet) error {
	args, err := m.populateAdventureGameTurnSheetValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameTurnSheetRecForCreate(args)
}

func (m *Domain) validateAdventureGameTurnSheetRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameTurnSheet) error {
	args, err := m.populateAdventureGameTurnSheetValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameTurnSheetRecForUpdate(args)
}

func validateAdventureGameTurnSheetRecForCreate(args *validateAdventureGameTurnSheetArgs) error {
	return validateAdventureGameTurnSheetRec(args, false)
}

func validateAdventureGameTurnSheetRecForUpdate(args *validateAdventureGameTurnSheetArgs) error {
	return validateAdventureGameTurnSheetRec(args, true)
}

func validateAdventureGameTurnSheetRec(args *validateAdventureGameTurnSheetArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameTurnSheetID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameTurnSheetGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID, rec.AdventureGameCharacterInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameTurnSheetGameTurnSheetID, rec.GameTurnSheetID); err != nil {
		return err
	}

	return nil
}
