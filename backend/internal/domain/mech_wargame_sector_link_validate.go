package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameSectorLinkArgs struct {
	currRec *mech_wargame_record.MechWargameSectorLink
	nextRec *mech_wargame_record.MechWargameSectorLink
}

func (m *Domain) validateMechWargameSectorLinkRecForCreate(rec *mech_wargame_record.MechWargameSectorLink) error {
	args := &validateMechWargameSectorLinkArgs{nextRec: rec}
	return validateMechWargameSectorLinkRec(args, false)
}

func (m *Domain) validateMechWargameSectorLinkRecForUpdate(currRec, nextRec *mech_wargame_record.MechWargameSectorLink) error {
	args := &validateMechWargameSectorLinkArgs{currRec: currRec, nextRec: nextRec}
	return validateMechWargameSectorLinkRec(args, true)
}

func validateMechWargameSectorLinkRec(args *validateMechWargameSectorLinkArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameSectorLinkID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameSectorLinkGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameSectorLinkFromMechWargameSectorID, rec.FromMechWargameSectorID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameSectorLinkToMechWargameSectorID, rec.ToMechWargameSectorID); err != nil {
		return err
	}

	if rec.FromMechWargameSectorID == rec.ToMechWargameSectorID {
		return InvalidField(mech_wargame_record.FieldMechWargameSectorLinkFromMechWargameSectorID, rec.FromMechWargameSectorID, "sector link cannot point to the same sector")
	}

	return nil
}
