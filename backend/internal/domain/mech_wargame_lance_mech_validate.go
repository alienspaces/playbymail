package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameLanceMechArgs struct {
	currRec *mech_wargame_record.MechWargameLanceMech
	nextRec *mech_wargame_record.MechWargameLanceMech
}

func (m *Domain) validateMechWargameLanceMechRecForCreate(rec *mech_wargame_record.MechWargameLanceMech) error {
	args := &validateMechWargameLanceMechArgs{nextRec: rec}
	return validateMechWargameLanceMechRec(args, false)
}

func (m *Domain) validateMechWargameLanceMechRecForUpdate(currRec, nextRec *mech_wargame_record.MechWargameLanceMech) error {
	args := &validateMechWargameLanceMechArgs{currRec: currRec, nextRec: nextRec}
	return validateMechWargameLanceMechRec(args, true)
}

func validateMechWargameLanceMechRec(args *validateMechWargameLanceMechArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceMechID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceMechGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceMechMechWargameLanceID, rec.MechWargameLanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceMechMechWargameChassisID, rec.MechWargameChassisID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mech_wargame_record.FieldMechWargameLanceMechCallsign, rec.Callsign); err != nil {
		return err
	}

	return nil
}
