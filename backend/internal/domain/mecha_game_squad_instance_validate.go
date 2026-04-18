package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameSquadInstanceArgs struct {
	currRec *mecha_game_record.MechaGameSquadInstance
	nextRec *mecha_game_record.MechaGameSquadInstance
}

func (m *Domain) validateMechaGameSquadInstanceRecForCreate(rec *mecha_game_record.MechaGameSquadInstance) error {
	args := &validateMechaGameSquadInstanceArgs{nextRec: rec}
	return validateMechaGameSquadInstanceRec(args, false)
}

func (m *Domain) validateMechaGameSquadInstanceRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameSquadInstance) error {
	args := &validateMechaGameSquadInstanceArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameSquadInstanceRec(args, true)
}

func validateMechaGameSquadInstanceRec(args *validateMechaGameSquadInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadInstanceMechaGameSquadID, rec.MechaGameSquadID); err != nil {
		return err
	}

	// GameSubscriptionInstanceID is NULL for computer-opponent squads.
	if rec.GameSubscriptionInstanceID.Valid {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSquadInstanceGameSubscriptionInstanceID, rec.GameSubscriptionInstanceID.String); err != nil {
			return err
		}
	}

	return nil
}
