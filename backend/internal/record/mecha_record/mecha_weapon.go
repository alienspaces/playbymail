package mecha_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaWeapon string = "mecha_weapon"
)

const (
	FieldMechaWeaponID          string = "id"
	FieldMechaWeaponGameID      string = "game_id"
	FieldMechaWeaponName        string = "name"
	FieldMechaWeaponDescription string = "description"
	FieldMechaWeaponDamage      string = "damage"
	FieldMechaWeaponHeatCost    string = "heat_cost"
	FieldMechaWeaponRangeBand   string = "range_band"
	FieldMechaWeaponMountSize   string = "mount_size"
	FieldMechaWeaponCreatedAt   string = "created_at"
	FieldMechaWeaponUpdatedAt   string = "updated_at"
	FieldMechaWeaponDeletedAt   string = "deleted_at"
)

const (
	WeaponRangeBandShort  string = "short"
	WeaponRangeBandMedium string = "medium"
	WeaponRangeBandLong   string = "long"

	WeaponMountSizeSmall  string = "small"
	WeaponMountSizeMedium string = "medium"
	WeaponMountSizeLarge  string = "large"
)

type MechaWeapon struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Damage      int    `db:"damage"`
	HeatCost    int    `db:"heat_cost"`
	RangeBand   string `db:"range_band"`
	MountSize   string `db:"mount_size"`
}

func (r *MechaWeapon) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaWeaponGameID] = r.GameID
	args[FieldMechaWeaponName] = r.Name
	args[FieldMechaWeaponDescription] = r.Description
	args[FieldMechaWeaponDamage] = r.Damage
	args[FieldMechaWeaponHeatCost] = r.HeatCost
	args[FieldMechaWeaponRangeBand] = r.RangeBand
	args[FieldMechaWeaponMountSize] = r.MountSize
	return args
}
