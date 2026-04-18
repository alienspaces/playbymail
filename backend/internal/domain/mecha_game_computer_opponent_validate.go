package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameComputerOpponentArgs struct {
	currRec *mecha_game_record.MechaGameComputerOpponent
	nextRec *mecha_game_record.MechaGameComputerOpponent
}

func (m *Domain) validateMechaGameComputerOpponentRecForCreate(rec *mecha_game_record.MechaGameComputerOpponent) error {
	args := &validateMechaGameComputerOpponentArgs{nextRec: rec}
	return validateMechaGameComputerOpponentRec(args, false)
}

func (m *Domain) validateMechaGameComputerOpponentRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameComputerOpponent) error {
	args := &validateMechaGameComputerOpponentArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameComputerOpponentRec(args, true)
}

func validateMechaGameComputerOpponentRec(args *validateMechaGameComputerOpponentArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameComputerOpponentID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameComputerOpponentGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_game_record.FieldMechaGameComputerOpponentName, rec.Name); err != nil {
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
