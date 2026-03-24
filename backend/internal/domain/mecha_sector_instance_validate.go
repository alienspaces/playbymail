package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaSectorInstanceArgs struct {
	nextRec *mecha_record.MechaSectorInstance
}

func (m *Domain) validateMechaSectorInstanceRecForCreate(rec *mecha_record.MechaSectorInstance) error {
	args := &validateMechaSectorInstanceArgs{nextRec: rec}
	return validateMechaSectorInstanceRec(args)
}

func validateMechaSectorInstanceRec(args *validateMechaSectorInstanceArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSectorInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSectorInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSectorInstanceMechaSectorID, rec.MechaSectorID); err != nil {
		return err
	}

	return nil
}
