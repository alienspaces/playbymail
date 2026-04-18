package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameSquadMechArgs struct {
	currRec *mecha_game_record.MechaGameSquadMech
	nextRec *mecha_game_record.MechaGameSquadMech
}

func (m *Domain) validateMechaGameSquadMechRecForCreate(rec *mecha_game_record.MechaGameSquadMech) error {
	args := &validateMechaGameSquadMechArgs{nextRec: rec}
	return validateMechaGameSquadMechRec(args, false)
}

func (m *Domain) validateMechaGameSquadMechRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameSquadMech) error {
	args := &validateMechaGameSquadMechArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameSquadMechRec(args, true)
}

func validateMechaGameSquadMechRec(args *validateMechaGameSquadMechArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadMechID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadMechGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadMechMechaGameSquadID, rec.MechaGameSquadID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadMechMechaGameChassisID, rec.MechaGameChassisID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_game_record.FieldMechaGameSquadMechCallsign, rec.Callsign); err != nil {
		return err
	}

	return nil
}
