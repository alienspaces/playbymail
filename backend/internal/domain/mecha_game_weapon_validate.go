package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

type validateMechaGameWeaponArgs struct {
	currRec *mecha_game_record.MechaGameWeapon
	nextRec *mecha_game_record.MechaGameWeapon
}

func (m *Domain) validateMechaGameWeaponRecForCreate(rec *mecha_game_record.MechaGameWeapon) error {
	args := &validateMechaGameWeaponArgs{nextRec: rec}
	return validateMechaGameWeaponRec(args, false)
}

func (m *Domain) validateMechaGameWeaponRecForUpdate(currRec, nextRec *mecha_game_record.MechaGameWeapon) error {
	args := &validateMechaGameWeaponArgs{currRec: currRec, nextRec: nextRec}
	return validateMechaGameWeaponRec(args, true)
}

func validateMechaGameWeaponRec(args *validateMechaGameWeaponArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameWeaponID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(mecha_game_record.FieldMechaGameWeaponGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(mecha_game_record.FieldMechaGameWeaponName, rec.Name); err != nil {
		return err
	}

	validRangeBands := map[string]bool{
		mecha_game_record.WeaponRangeBandShort:  true,
		mecha_game_record.WeaponRangeBandMedium: true,
		mecha_game_record.WeaponRangeBandLong:   true,
	}
	if rec.RangeBand == "" {
		rec.RangeBand = mecha_game_record.WeaponRangeBandMedium
	}
	if !validRangeBands[rec.RangeBand] {
		return InvalidField(mecha_game_record.FieldMechaGameWeaponRangeBand, rec.RangeBand, "must be one of: short, medium, long")
	}

	validMountSizes := map[string]bool{
		mecha_game_record.WeaponMountSizeSmall:  true,
		mecha_game_record.WeaponMountSizeMedium: true,
		mecha_game_record.WeaponMountSizeLarge:  true,
	}
	if rec.MountSize == "" {
		rec.MountSize = mecha_game_record.WeaponMountSizeMedium
	}
	if !validMountSizes[rec.MountSize] {
		return InvalidField(mecha_game_record.FieldMechaGameWeaponMountSize, rec.MountSize, "must be one of: small, medium, large")
	}

	if rec.Damage <= 0 || rec.Damage > maxWeaponDamage {
		return InvalidField(mecha_game_record.FieldMechaGameWeaponDamage, "", "damage must be between 1 and 20")
	}

	if rec.HeatCost < 0 || rec.HeatCost > maxWeaponHeatCost {
		return InvalidField(mecha_game_record.FieldMechaGameWeaponHeatCost, "", "heat_cost must be between 0 and 20")
	}

	if rec.AmmoCapacity < 0 || rec.AmmoCapacity > maxWeaponAmmoCapacity {
		return InvalidField(mecha_game_record.FieldMechaGameWeaponAmmoCapacity, "", "ammo_capacity must be between 0 and 200 (0 = no ammo tracking)")
	}

	return nil
}

// Upper bounds on weapon stats. These keep designer input within playable
// ranges relative to chassis armour, structure, and heat capacity (see
// mecha_game_chassis_validate.go).
const (
	maxWeaponDamage       = 20
	maxWeaponHeatCost     = 20
	maxWeaponAmmoCapacity = 200
)
