package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaSquadInstanceArgs struct {
	currRec *mecha_record.MechaSquadInstance
	nextRec *mecha_record.MechaSquadInstance
}

func (m *Domain) validateMechaSquadInstanceRecForCreate(rec *mecha_record.MechaSquadInstance) error {
	args := &validateMechaSquadInstanceArgs{nextRec: rec}
	return validateMechaSquadInstanceRec(args, false)
}

func (m *Domain) validateMechaSquadInstanceRecForUpdate(currRec, nextRec *mecha_record.MechaSquadInstance) error {
	args := &validateMechaSquadInstanceArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaSquadInstanceRec(args, true)
}

func validateMechaSquadInstanceRec(args *validateMechaSquadInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadInstanceMechaSquadID, rec.MechaSquadID); err != nil {
		return err
	}

	// GameSubscriptionInstanceID is NULL for computer-opponent squads.
	if rec.GameSubscriptionInstanceID.Valid {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaSquadInstanceGameSubscriptionInstanceID, rec.GameSubscriptionInstanceID.String); err != nil {
			return err
		}
	}

	return nil
}
