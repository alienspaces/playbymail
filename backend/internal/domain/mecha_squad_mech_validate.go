package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaSquadMechArgs struct {
	currRec *mecha_record.MechaSquadMech
	nextRec *mecha_record.MechaSquadMech
}

func (m *Domain) validateMechaSquadMechRecForCreate(rec *mecha_record.MechaSquadMech) error {
	args := &validateMechaSquadMechArgs{nextRec: rec}
	return validateMechaSquadMechRec(args, false)
}

func (m *Domain) validateMechaSquadMechRecForUpdate(currRec, nextRec *mecha_record.MechaSquadMech) error {
	args := &validateMechaSquadMechArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaSquadMechRec(args, true)
}

func validateMechaSquadMechRec(args *validateMechaSquadMechArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadMechID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadMechGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadMechMechaSquadID, rec.MechaSquadID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadMechMechaChassisID, rec.MechaChassisID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaSquadMechCallsign, rec.Callsign); err != nil {
		return err
	}

	return nil
}
