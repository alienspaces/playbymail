package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameChassis string = "mecha_game_chassis"
)

const (
	FieldMechaGameChassisID             string = "id"
	FieldMechaGameChassisGameID         string = "game_id"
	FieldMechaGameChassisName           string = "name"
	FieldMechaGameChassisDescription    string = "description"
	FieldMechaGameChassisChassisClass   string = "chassis_class"
	FieldMechaGameChassisArmorPoints    string = "armor_points"
	FieldMechaGameChassisStructurePoints string = "structure_points"
	FieldMechaGameChassisHeatCapacity   string = "heat_capacity"
	FieldMechaGameChassisSpeed          string = "speed"
	FieldMechaGameChassisCreatedAt      string = "created_at"
	FieldMechaGameChassisUpdatedAt      string = "updated_at"
	FieldMechaGameChassisDeletedAt      string = "deleted_at"
)

const (
	ChassisClassLight   string = "light"
	ChassisClassMedium  string = "medium"
	ChassisClassHeavy   string = "heavy"
	ChassisClassAssault string = "assault"
)

type MechaGameChassis struct {
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

func (r *MechaGameChassis) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameChassisGameID] = r.GameID
	args[FieldMechaGameChassisName] = r.Name
	args[FieldMechaGameChassisDescription] = r.Description
	args[FieldMechaGameChassisChassisClass] = r.ChassisClass
	args[FieldMechaGameChassisArmorPoints] = r.ArmorPoints
	args[FieldMechaGameChassisStructurePoints] = r.StructurePoints
	args[FieldMechaGameChassisHeatCapacity] = r.HeatCapacity
	args[FieldMechaGameChassisSpeed] = r.Speed
	return args
}
