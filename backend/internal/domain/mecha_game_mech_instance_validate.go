package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameMechInstanceArgs struct {
	currRec *mecha_game_record.MechaGameMechInstance
	nextRec *mecha_game_record.MechaGameMechInstance
}

func (m *Domain) validateMechaGameMechInstanceRecForCreate(rec *mecha_game_record.MechaGameMechInstance) error {
	args := &validateMechaGameMechInstanceArgs{nextRec: rec}
	return validateMechaGameMechInstanceRec(args, false)
}

func (m *Domain) validateMechaGameMechInstanceRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameMechInstance) error {
	args := &validateMechaGameMechInstanceArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameMechInstanceRec(args, true)
}

func validateMechaGameMechInstanceRec(args *validateMechaGameMechInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameMechInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameMechInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameMechInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameMechInstanceMechaGameSquadInstanceID, rec.MechaGameSquadInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameMechInstanceMechaGameSectorInstanceID, rec.MechaGameSectorInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameMechInstanceMechaGameChassisID, rec.MechaGameChassisID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_game_record.FieldMechaGameMechInstanceCallsign, rec.Callsign); err != nil {
		return err
	}

	validStatuses := map[string]bool{
		mecha_game_record.MechInstanceStatusOperational: true,
		mecha_game_record.MechInstanceStatusDamaged:     true,
		mecha_game_record.MechInstanceStatusDestroyed:   true,
		mecha_game_record.MechInstanceStatusShutdown:    true,
	}
	if rec.Status == "" {
		rec.Status = mecha_game_record.MechInstanceStatusOperational
	}
	if !validStatuses[rec.Status] {
		return InvalidField(mecha_game_record.FieldMechaGameMechInstanceStatus, rec.Status, "must be one of: operational, damaged, destroyed, shutdown")
	}

	return nil
}
