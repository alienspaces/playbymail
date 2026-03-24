package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

type validateMechaWeaponArgs struct {
	currRec *mecha_record.MechaWeapon
	nextRec *mecha_record.MechaWeapon
}

func (m *Domain) validateMechaWeaponRecForCreate(rec *mecha_record.MechaWeapon) error {
	args := &validateMechaWeaponArgs{nextRec: rec}
	return validateMechaWeaponRec(args, false)
}

func (m *Domain) validateMechaWeaponRecForUpdate(currRec, nextRec *mecha_record.MechaWeapon) error {
	args := &validateMechaWeaponArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaWeaponRec(args, true)
}

func validateMechaWeaponRec(args *validateMechaWeaponArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_record.FieldMechaWeaponID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_record.FieldMechaWeaponGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_record.FieldMechaWeaponName, rec.Name); err != nil {
		return err
	}

	validRangeBands := map[string]bool{
		mecha_record.WeaponRangeBandShort:  true,
		mecha_record.WeaponRangeBandMedium: true,
		mecha_record.WeaponRangeBandLong:   true,
	}
	if rec.RangeBand == "" {
		rec.RangeBand = mecha_record.WeaponRangeBandMedium
	}
	if !validRangeBands[rec.RangeBand] {
		return InvalidField(mecha_record.FieldMechaWeaponRangeBand, rec.RangeBand, "must be one of: short, medium, long")
	}

	validMountSizes := map[string]bool{
		mecha_record.WeaponMountSizeSmall:  true,
		mecha_record.WeaponMountSizeMedium: true,
		mecha_record.WeaponMountSizeLarge:  true,
	}
	if rec.MountSize == "" {
		rec.MountSize = mecha_record.WeaponMountSizeMedium
	}
	if !validMountSizes[rec.MountSize] {
		return InvalidField(mecha_record.FieldMechaWeaponMountSize, rec.MountSize, "must be one of: small, medium, large")
	}

	if rec.Damage <= 0 {
		return InvalidField(mecha_record.FieldMechaWeaponDamage, "", "damage must be greater than 0")
	}

	if rec.HeatCost < 0 {
		return InvalidField(mecha_record.FieldMechaWeaponHeatCost, "", "heat_cost must be >= 0")
	}

	return nil
}
