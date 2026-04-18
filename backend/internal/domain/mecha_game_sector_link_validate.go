package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameSectorLinkArgs struct {
	currRec *mecha_game_record.MechaGameSectorLink
	nextRec *mecha_game_record.MechaGameSectorLink
}

func (m *Domain) validateMechaGameSectorLinkRecForCreate(rec *mecha_game_record.MechaGameSectorLink) error {
	args := &validateMechaGameSectorLinkArgs{nextRec: rec}
	return validateMechaGameSectorLinkRec(args, false)
}

func (m *Domain) validateMechaGameSectorLinkRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameSectorLink) error {
	args := &validateMechaGameSectorLinkArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameSectorLinkRec(args, true)
}

func validateMechaGameSectorLinkRec(args *validateMechaGameSectorLinkArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSectorLinkID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSectorLinkGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSectorLinkFromMechaGameSectorID, rec.FromMechaGameSectorID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameSectorLinkToMechaGameSectorID, rec.ToMechaGameSectorID); err != nil {
		return err
	}

	if rec.FromMechaGameSectorID == rec.ToMechaGameSectorID {
		return InvalidField(mecha_game_record.FieldMechaGameSectorLinkFromMechaGameSectorID, rec.FromMechaGameSectorID, "sector link cannot point to the same sector")
	}

	return nil
}
