package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameMechInstanceArgs struct {
	currRec *mech_wargame_record.MechWargameMechInstance
	nextRec *mech_wargame_record.MechWargameMechInstance
}

func (m *Domain) validateMechWargameMechInstanceRecForCreate(rec *mech_wargame_record.MechWargameMechInstance) error {
	args := &validateMechWargameMechInstanceArgs{nextRec: rec}
	return validateMechWargameMechInstanceRec(args, false)
}

func (m *Domain) validateMechWargameMechInstanceRecForUpdate(currRec, nextRec *mech_wargame_record.MechWargameMechInstance) error {
	args := &validateMechWargameMechInstanceArgs{currRec: currRec, nextRec: nextRec}
	return validateMechWargameMechInstanceRec(args, true)
}

func validateMechWargameMechInstanceRec(args *validateMechWargameMechInstanceArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameMechInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameMechInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameMechInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameMechInstanceMechWargameLanceInstanceID, rec.MechWargameLanceInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameMechInstanceMechWargameSectorInstanceID, rec.MechWargameSectorInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameMechInstanceMechWargameChassisID, rec.MechWargameChassisID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mech_wargame_record.FieldMechWargameMechInstanceCallsign, rec.Callsign); err != nil {
		return err
	}

	validStatuses := map[string]bool{
		mech_wargame_record.MechInstanceStatusOperational: true,
		mech_wargame_record.MechInstanceStatusDamaged:     true,
		mech_wargame_record.MechInstanceStatusDestroyed:   true,
		mech_wargame_record.MechInstanceStatusShutdown:    true,
	}
	if rec.Status == "" {
		rec.Status = mech_wargame_record.MechInstanceStatusOperational
	}
	if !validStatuses[rec.Status] {
		return InvalidField(mech_wargame_record.FieldMechWargameMechInstanceStatus, rec.Status, "must be one of: operational, damaged, destroyed, shutdown")
	}

	return nil
}
