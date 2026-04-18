package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameWeapon string = "mecha_game_weapon"
)

const (
	FieldMechaGameWeaponID          string = "id"
	FieldMechaGameWeaponGameID      string = "game_id"
	FieldMechaGameWeaponName        string = "name"
	FieldMechaGameWeaponDescription string = "description"
	FieldMechaGameWeaponDamage      string = "damage"
	FieldMechaGameWeaponHeatCost    string = "heat_cost"
	FieldMechaGameWeaponRangeBand   string = "range_band"
	FieldMechaGameWeaponMountSize   string = "mount_size"
	FieldMechaGameWeaponCreatedAt   string = "created_at"
	FieldMechaGameWeaponUpdatedAt   string = "updated_at"
	FieldMechaGameWeaponDeletedAt   string = "deleted_at"
)

const (
	WeaponRangeBandShort  string = "short"
	WeaponRangeBandMedium string = "medium"
	WeaponRangeBandLong   string = "long"

	WeaponMountSizeSmall  string = "small"
	WeaponMountSizeMedium string = "medium"
	WeaponMountSizeLarge  string = "large"
)

type MechaGameWeapon struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Damage      int    `db:"damage"`
	HeatCost    int    `db:"heat_cost"`
	RangeBand   string `db:"range_band"`
	MountSize   string `db:"mount_size"`
}

func (r *MechaGameWeapon) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameWeaponGameID] = r.GameID
	args[FieldMechaGameWeaponName] = r.Name
	args[FieldMechaGameWeaponDescription] = r.Description
	args[FieldMechaGameWeaponDamage] = r.Damage
	args[FieldMechaGameWeaponHeatCost] = r.HeatCost
	args[FieldMechaGameWeaponRangeBand] = r.RangeBand
	args[FieldMechaGameWeaponMountSize] = r.MountSize
	return args
}
