package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaSectorLinkArgs struct {
	currRec *mecha_record.MechaSectorLink
	nextRec *mecha_record.MechaSectorLink
}

func (m *Domain) validateMechaSectorLinkRecForCreate(rec *mecha_record.MechaSectorLink) error {
	args := &validateMechaSectorLinkArgs{nextRec: rec}
	return validateMechaSectorLinkRec(args, false)
}

func (m *Domain) validateMechaSectorLinkRecForUpdate(currRec, nextRec *mecha_record.MechaSectorLink) error {
	args := &validateMechaSectorLinkArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaSectorLinkRec(args, true)
}

func validateMechaSectorLinkRec(args *validateMechaSectorLinkArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaSectorLinkID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSectorLinkGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSectorLinkFromMechaSectorID, rec.FromMechaSectorID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaSectorLinkToMechaSectorID, rec.ToMechaSectorID); err != nil {
		return err
	}

	if rec.FromMechaSectorID == rec.ToMechaSectorID {
		return InvalidField(mecha_record.FieldMechaSectorLinkFromMechaSectorID, rec.FromMechaSectorID, "sector link cannot point to the same sector")
	}

	return nil
}
