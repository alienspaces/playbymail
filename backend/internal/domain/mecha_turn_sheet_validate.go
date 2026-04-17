package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaTurnSheetArgs struct {
	nextRec *mecha_record.MechaTurnSheet
}

func (m *Domain) validateMechaTurnSheetRecForCreate(rec *mecha_record.MechaTurnSheet) error {
	args := &validateMechaTurnSheetArgs{nextRec: rec}
	return validateMechaTurnSheetRec(args)
}

func validateMechaTurnSheetRec(args *validateMechaTurnSheetArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaTurnSheetGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaTurnSheetMechaSquadInstanceID, rec.MechaSquadInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaTurnSheetGameTurnSheetID, rec.GameTurnSheetID); err != nil {
		return err
	}

	return nil
}
