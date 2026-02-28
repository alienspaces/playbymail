package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type validateGameTurnSheetArgs struct {
	nextRec *game_record.GameTurnSheet
	currRec *game_record.GameTurnSheet
}

func (m *Domain) populateGameTurnSheetValidateArgs(currRec, nextRec *game_record.GameTurnSheet) (*validateGameTurnSheetArgs, error) {
	args := &validateGameTurnSheetArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateGameTurnSheetRecForCreate(rec *game_record.GameTurnSheet) error {
	args, err := m.populateGameTurnSheetValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateGameTurnSheetRecForCreate(args)
}

func (m *Domain) validateGameTurnSheetRecForUpdate(currRec, nextRec *game_record.GameTurnSheet) error {
	args, err := m.populateGameTurnSheetValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateGameTurnSheetRecForUpdate(args)
}

func validateGameTurnSheetRecForCreate(args *validateGameTurnSheetArgs) error {
	return validateGameTurnSheetRec(args, false)
}

func validateGameTurnSheetRecForUpdate(args *validateGameTurnSheetArgs) error {
	return validateGameTurnSheetRec(args, true)
}

func validateGameTurnSheetRec(args *validateGameTurnSheetArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(game_record.FieldGameTurnSheetID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(game_record.FieldGameTurnSheetGameID, rec.GameID); err != nil {
		return err
	}

	if nullstring.IsValid(rec.GameInstanceID) {
		if err := domain.ValidateNullUUIDField(game_record.FieldGameTurnSheetGameInstanceID, rec.GameInstanceID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(game_record.FieldGameTurnSheetAccountID, rec.AccountID); err != nil {
		return err
	}

	if rec.TurnNumber < 0 {
		return InvalidField(
			game_record.FieldGameTurnSheetTurnNumber,
			fmt.Sprintf("%d", rec.TurnNumber),
			"turn_number must be zero or greater",
		)
	}

	return nil
}
