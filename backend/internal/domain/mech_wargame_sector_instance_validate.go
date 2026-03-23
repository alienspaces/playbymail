package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameSectorInstanceArgs struct {
	nextRec *mech_wargame_record.MechWargameSectorInstance
}

func (m *Domain) validateMechWargameSectorInstanceRecForCreate(rec *mech_wargame_record.MechWargameSectorInstance) error {
	args := &validateMechWargameSectorInstanceArgs{nextRec: rec}
	return validateMechWargameSectorInstanceRec(args)
}

func validateMechWargameSectorInstanceRec(args *validateMechWargameSectorInstanceArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameSectorInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameSectorInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameSectorInstanceMechWargameSectorID, rec.MechWargameSectorID); err != nil {
		return err
	}

	return nil
}
