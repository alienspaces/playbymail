package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaComputerOpponentArgs struct {
	currRec *mecha_record.MechaComputerOpponent
	nextRec *mecha_record.MechaComputerOpponent
}

func (m *Domain) validateMechaComputerOpponentRecForCreate(rec *mecha_record.MechaComputerOpponent) error {
	args := &validateMechaComputerOpponentArgs{nextRec: rec}
	return validateMechaComputerOpponentRec(args, false)
}

func (m *Domain) validateMechaComputerOpponentRecForUpdate(currRec, nextRec *mecha_record.MechaComputerOpponent) error {
	args := &validateMechaComputerOpponentArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaComputerOpponentRec(args, true)
}

func validateMechaComputerOpponentRec(args *validateMechaComputerOpponentArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaComputerOpponentID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaComputerOpponentGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaComputerOpponentName, rec.Name); err != nil {
		return err
	}

	if rec.Aggression < 1 || rec.Aggression > 10 {
		return coreerror.NewInvalidDataError("aggression must be between 1 and 10, got %d", rec.Aggression)
	}

	if rec.IQ < 1 || rec.IQ > 10 {
		return coreerror.NewInvalidDataError("iq must be between 1 and 10, got %d", rec.IQ)
	}

	return nil
}
