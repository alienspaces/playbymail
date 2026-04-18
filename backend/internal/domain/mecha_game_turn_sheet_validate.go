package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameTurnSheetArgs struct {
	nextRec *mecha_game_record.MechaGameTurnSheet
}

func (m *Domain) validateMechaGameTurnSheetRecForCreate(rec *mecha_game_record.MechaGameTurnSheet) error {
	args := &validateMechaGameTurnSheetArgs{nextRec: rec}
	return validateMechaGameTurnSheetRec(args)
}

func validateMechaGameTurnSheetRec(args *validateMechaGameTurnSheetArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameTurnSheetGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameTurnSheetMechaGameSquadInstanceID, rec.MechaGameSquadInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameTurnSheetGameTurnSheetID, rec.GameTurnSheetID); err != nil {
		return err
	}

	return nil
}
