package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaMechInstanceArgs struct {
	currRec *mecha_record.MechaMechInstance
	nextRec *mecha_record.MechaMechInstance
}

func (m *Domain) validateMechaMechInstanceRecForCreate(rec *mecha_record.MechaMechInstance) error {
	args := &validateMechaMechInstanceArgs{nextRec: rec}
	return validateMechaMechInstanceRec(args, false)
}

func (m *Domain) validateMechaMechInstanceRecForUpdate(currRec, nextRec *mecha_record.MechaMechInstance) error {
	args := &validateMechaMechInstanceArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaMechInstanceRec(args, true)
}

func validateMechaMechInstanceRec(args *validateMechaMechInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaMechInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaMechInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaMechInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaMechInstanceMechaLanceInstanceID, rec.MechaLanceInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaMechInstanceMechaSectorInstanceID, rec.MechaSectorInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaMechInstanceMechaChassisID, rec.MechaChassisID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaMechInstanceCallsign, rec.Callsign); err != nil {
		return err
	}

	validStatuses := map[string]bool{
		mecha_record.MechInstanceStatusOperational: true,
		mecha_record.MechInstanceStatusDamaged:     true,
		mecha_record.MechInstanceStatusDestroyed:   true,
		mecha_record.MechInstanceStatusShutdown:    true,
	}
	if rec.Status == "" {
		rec.Status = mecha_record.MechInstanceStatusOperational
	}
	if !validStatuses[rec.Status] {
		return InvalidField(mecha_record.FieldMechaMechInstanceStatus, rec.Status, "must be one of: operational, damaged, destroyed, shutdown")
	}

	return nil
}
