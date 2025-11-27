package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (m *Domain) validateGameTurnSheetRecForCreate(rec *game_record.GameTurnSheet) error {
	return validateGameTurnSheetRec(rec, false)
}

func (m *Domain) validateGameTurnSheetRecForUpdate(rec *game_record.GameTurnSheet) error {
	return validateGameTurnSheetRec(rec, true)
}

func validateGameTurnSheetRec(rec *game_record.GameTurnSheet, requireID bool) error {
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
