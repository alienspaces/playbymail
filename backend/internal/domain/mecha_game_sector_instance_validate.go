package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameSectorInstanceArgs struct {
	nextRec *mecha_game_record.MechaGameSectorInstance
}

func (m *Domain) validateMechaGameSectorInstanceRecForCreate(rec *mecha_game_record.MechaGameSectorInstance) error {
	args := &validateMechaGameSectorInstanceArgs{nextRec: rec}
	return validateMechaGameSectorInstanceRec(args)
}

func validateMechaGameSectorInstanceRec(args *validateMechaGameSectorInstanceArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSectorInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSectorInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSectorInstanceMechaGameSectorID, rec.MechaGameSectorID); err != nil {
		return err
	}

	return nil
}
