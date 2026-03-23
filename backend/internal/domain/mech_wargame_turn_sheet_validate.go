package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameTurnSheetArgs struct {
	nextRec *mech_wargame_record.MechWargameTurnSheet
}

func (m *Domain) validateMechWargameTurnSheetRecForCreate(rec *mech_wargame_record.MechWargameTurnSheet) error {
	args := &validateMechWargameTurnSheetArgs{nextRec: rec}
	return validateMechWargameTurnSheetRec(args)
}

func validateMechWargameTurnSheetRec(args *validateMechWargameTurnSheetArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameTurnSheetGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameTurnSheetMechWargameLanceInstanceID, rec.MechWargameLanceInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameTurnSheetGameTurnSheetID, rec.GameTurnSheetID); err != nil {
		return err
	}

	return nil
}
