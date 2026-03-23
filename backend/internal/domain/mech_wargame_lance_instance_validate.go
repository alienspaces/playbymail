package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameLanceInstanceArgs struct {
	nextRec *mech_wargame_record.MechWargameLanceInstance
}

func (m *Domain) validateMechWargameLanceInstanceRecForCreate(rec *mech_wargame_record.MechWargameLanceInstance) error {
	args := &validateMechWargameLanceInstanceArgs{nextRec: rec}
	return validateMechWargameLanceInstanceRec(args)
}

func validateMechWargameLanceInstanceRec(args *validateMechWargameLanceInstanceArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceInstanceMechWargameLanceID, rec.MechWargameLanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameLanceInstanceGameSubscriptionInstanceID, rec.GameSubscriptionInstanceID); err != nil {
		return err
	}

	return nil
}
