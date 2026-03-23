package mech_wargame_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameChassis string = "mech_wargame_chassis"
)

const (
	FieldMechWargameChassisID             string = "id"
	FieldMechWargameChassisGameID         string = "game_id"
	FieldMechWargameChassisName           string = "name"
	FieldMechWargameChassisDescription    string = "description"
	FieldMechWargameChassisChassisClass   string = "chassis_class"
	FieldMechWargameChassisArmorPoints    string = "armor_points"
	FieldMechWargameChassisStructurePoints string = "structure_points"
	FieldMechWargameChassisHeatCapacity   string = "heat_capacity"
	FieldMechWargameChassisSpeed          string = "speed"
	FieldMechWargameChassisCreatedAt      string = "created_at"
	FieldMechWargameChassisUpdatedAt      string = "updated_at"
	FieldMechWargameChassisDeletedAt      string = "deleted_at"
)

const (
	ChassisClassLight   string = "light"
	ChassisClassMedium  string = "medium"
	ChassisClassHeavy   string = "heavy"
	ChassisClassAssault string = "assault"
)

type MechWargameChassis struct {
	record.Record
	GameID          string `db:"game_id"`
	Name            string `db:"name"`
	Description     string `db:"description"`
	ChassisClass    string `db:"chassis_class"`
	ArmorPoints     int    `db:"armor_points"`
	StructurePoints int    `db:"structure_points"`
	HeatCapacity    int    `db:"heat_capacity"`
	Speed           int    `db:"speed"`
}

func (r *MechWargameChassis) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameChassisGameID] = r.GameID
	args[FieldMechWargameChassisName] = r.Name
	args[FieldMechWargameChassisDescription] = r.Description
	args[FieldMechWargameChassisChassisClass] = r.ChassisClass
	args[FieldMechWargameChassisArmorPoints] = r.ArmorPoints
	args[FieldMechWargameChassisStructurePoints] = r.StructurePoints
	args[FieldMechWargameChassisHeatCapacity] = r.HeatCapacity
	args[FieldMechWargameChassisSpeed] = r.Speed
	return args
}
