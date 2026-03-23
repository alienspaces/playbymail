package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mech_wargame_record"
)

type validateMechWargameWeaponArgs struct {
	currRec *mech_wargame_record.MechWargameWeapon
	nextRec *mech_wargame_record.MechWargameWeapon
}

func (m *Domain) validateMechWargameWeaponRecForCreate(rec *mech_wargame_record.MechWargameWeapon) error {
	args := &validateMechWargameWeaponArgs{nextRec: rec}
	return validateMechWargameWeaponRec(args, false)
}

func (m *Domain) validateMechWargameWeaponRecForUpdate(currRec, nextRec *mech_wargame_record.MechWargameWeapon) error {
	args := &validateMechWargameWeaponArgs{currRec: currRec, nextRec: nextRec}
	return validateMechWargameWeaponRec(args, true)
}

func validateMechWargameWeaponRec(args *validateMechWargameWeaponArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameWeaponID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mech_wargame_record.FieldMechWargameWeaponGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mech_wargame_record.FieldMechWargameWeaponName, rec.Name); err != nil {
		return err
	}

	validRangeBands := map[string]bool{
		mech_wargame_record.WeaponRangeBandShort:  true,
		mech_wargame_record.WeaponRangeBandMedium: true,
		mech_wargame_record.WeaponRangeBandLong:   true,
	}
	if rec.RangeBand == "" {
		rec.RangeBand = mech_wargame_record.WeaponRangeBandMedium
	}
	if !validRangeBands[rec.RangeBand] {
		return InvalidField(mech_wargame_record.FieldMechWargameWeaponRangeBand, rec.RangeBand, "must be one of: short, medium, long")
	}

	validMountSizes := map[string]bool{
		mech_wargame_record.WeaponMountSizeSmall:  true,
		mech_wargame_record.WeaponMountSizeMedium: true,
		mech_wargame_record.WeaponMountSizeLarge:  true,
	}
	if rec.MountSize == "" {
		rec.MountSize = mech_wargame_record.WeaponMountSizeMedium
	}
	if !validMountSizes[rec.MountSize] {
		return InvalidField(mech_wargame_record.FieldMechWargameWeaponMountSize, rec.MountSize, "must be one of: small, medium, large")
	}

	if rec.Damage <= 0 {
		return InvalidField(mech_wargame_record.FieldMechWargameWeaponDamage, "", "damage must be greater than 0")
	}

	if rec.HeatCost < 0 {
		return InvalidField(mech_wargame_record.FieldMechWargameWeaponHeatCost, "", "heat_cost must be >= 0")
	}

	return nil
}
