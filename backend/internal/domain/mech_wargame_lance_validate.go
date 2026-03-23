package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameLanceArgs struct {
	currRec *mech_wargame_record.MechWargameLance
	nextRec *mech_wargame_record.MechWargameLance
}

func (m *Domain) validateMechWargameLanceRecForCreate(rec *mech_wargame_record.MechWargameLance) error {
	args := &validateMechWargameLanceArgs{nextRec: rec}
	return validateMechWargameLanceRec(args, false)
}

func (m *Domain) validateMechWargameLanceRecForUpdate(currRec, nextRec *mech_wargame_record.MechWargameLance) error {
	args := &validateMechWargameLanceArgs{currRec: currRec, nextRec: nextRec}
	return validateMechWargameLanceRec(args, true)
}

func validateMechWargameLanceRec(args *validateMechWargameLanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceAccountID, rec.AccountID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceAccountUserID, rec.AccountUserID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mech_wargame_record.FieldMechWargameLanceName, rec.Name); err != nil {
		return err
	}

	return nil
}
