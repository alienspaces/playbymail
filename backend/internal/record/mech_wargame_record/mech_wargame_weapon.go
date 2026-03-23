package mech_wargame_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameWeapon string = "mech_wargame_weapon"
)

const (
	FieldMechWargameWeaponID          string = "id"
	FieldMechWargameWeaponGameID      string = "game_id"
	FieldMechWargameWeaponName        string = "name"
	FieldMechWargameWeaponDescription string = "description"
	FieldMechWargameWeaponDamage      string = "damage"
	FieldMechWargameWeaponHeatCost    string = "heat_cost"
	FieldMechWargameWeaponRangeBand   string = "range_band"
	FieldMechWargameWeaponMountSize   string = "mount_size"
	FieldMechWargameWeaponCreatedAt   string = "created_at"
	FieldMechWargameWeaponUpdatedAt   string = "updated_at"
	FieldMechWargameWeaponDeletedAt   string = "deleted_at"
)

const (
	WeaponRangeBandShort  string = "short"
	WeaponRangeBandMedium string = "medium"
	WeaponRangeBandLong   string = "long"

	WeaponMountSizeSmall  string = "small"
	WeaponMountSizeMedium string = "medium"
	WeaponMountSizeLarge  string = "large"
)

type MechWargameWeapon struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Damage      int    `db:"damage"`
	HeatCost    int    `db:"heat_cost"`
	RangeBand   string `db:"range_band"`
	MountSize   string `db:"mount_size"`
}

func (r *MechWargameWeapon) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameWeaponGameID] = r.GameID
	args[FieldMechWargameWeaponName] = r.Name
	args[FieldMechWargameWeaponDescription] = r.Description
	args[FieldMechWargameWeaponDamage] = r.Damage
	args[FieldMechWargameWeaponHeatCost] = r.HeatCost
	args[FieldMechWargameWeaponRangeBand] = r.RangeBand
	args[FieldMechWargameWeaponMountSize] = r.MountSize
	return args
}
