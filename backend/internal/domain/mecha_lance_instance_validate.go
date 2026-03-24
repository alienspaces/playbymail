package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaLanceInstanceArgs struct {
	currRec *mecha_record.MechaLanceInstance
	nextRec *mecha_record.MechaLanceInstance
}

func (m *Domain) validateMechaLanceInstanceRecForCreate(rec *mecha_record.MechaLanceInstance) error {
	args := &validateMechaLanceInstanceArgs{nextRec: rec}
	return validateMechaLanceInstanceRec(args, false)
}

func (m *Domain) validateMechaLanceInstanceRecForUpdate(currRec, nextRec *mecha_record.MechaLanceInstance) error {
	args := &validateMechaLanceInstanceArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaLanceInstanceRec(args, true)
}

func validateMechaLanceInstanceRec(args *validateMechaLanceInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceInstanceMechaLanceID, rec.MechaLanceID); err != nil {
		return err
	}

	// GameSubscriptionInstanceID is NULL for computer-opponent lances.
	if rec.GameSubscriptionInstanceID.Valid {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaLanceInstanceGameSubscriptionInstanceID, rec.GameSubscriptionInstanceID.String); err != nil {
			return err
		}
	}

	return nil
}
