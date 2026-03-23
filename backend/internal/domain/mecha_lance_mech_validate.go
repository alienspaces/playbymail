package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaLanceMechArgs struct {
	currRec *mecha_record.MechaLanceMech
	nextRec *mecha_record.MechaLanceMech
}

func (m *Domain) validateMechaLanceMechRecForCreate(rec *mecha_record.MechaLanceMech) error {
	args := &validateMechaLanceMechArgs{nextRec: rec}
	return validateMechaLanceMechRec(args, false)
}

func (m *Domain) validateMechaLanceMechRecForUpdate(currRec, nextRec *mecha_record.MechaLanceMech) error {
	args := &validateMechaLanceMechArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaLanceMechRec(args, true)
}

func validateMechaLanceMechRec(args *validateMechaLanceMechArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceMechID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceMechGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceMechMechaLanceID, rec.MechaLanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceMechMechaChassisID, rec.MechaChassisID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaLanceMechCallsign, rec.Callsign); err != nil {
		return err
	}

	return nil
}
