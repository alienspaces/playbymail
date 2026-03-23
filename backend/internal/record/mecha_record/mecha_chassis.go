package mecha_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaChassis string = "mecha_chassis"
)

const (
	FieldMechaChassisID             string = "id"
	FieldMechaChassisGameID         string = "game_id"
	FieldMechaChassisName           string = "name"
	FieldMechaChassisDescription    string = "description"
	FieldMechaChassisChassisClass   string = "chassis_class"
	FieldMechaChassisArmorPoints    string = "armor_points"
	FieldMechaChassisStructurePoints string = "structure_points"
	FieldMechaChassisHeatCapacity   string = "heat_capacity"
	FieldMechaChassisSpeed          string = "speed"
	FieldMechaChassisCreatedAt      string = "created_at"
	FieldMechaChassisUpdatedAt      string = "updated_at"
	FieldMechaChassisDeletedAt      string = "deleted_at"
)

const (
	ChassisClassLight   string = "light"
	ChassisClassMedium  string = "medium"
	ChassisClassHeavy   string = "heavy"
	ChassisClassAssault string = "assault"
)

type MechaChassis struct {
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

func (r *MechaChassis) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaChassisGameID] = r.GameID
	args[FieldMechaChassisName] = r.Name
	args[FieldMechaChassisDescription] = r.Description
	args[FieldMechaChassisChassisClass] = r.ChassisClass
	args[FieldMechaChassisArmorPoints] = r.ArmorPoints
	args[FieldMechaChassisStructurePoints] = r.StructurePoints
	args[FieldMechaChassisHeatCapacity] = r.HeatCapacity
	args[FieldMechaChassisSpeed] = r.Speed
	return args
}
