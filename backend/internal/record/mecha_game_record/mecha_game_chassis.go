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
	FieldMechaGameChassisSmallSlots     string = "small_slots"
	FieldMechaGameChassisMediumSlots    string = "medium_slots"
	FieldMechaGameChassisLargeSlots     string = "large_slots"
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

// DefaultSlotsForChassisClass returns the default small/medium/large slot
// counts to pre-fill for a newly created chassis of the given class. These
// match the per-class backfill applied by the mecha_game_chassis_slots
// migration so a chassis created via any path ends up with the same shape as
// seeded data of the same class. Values intentionally leave headroom for the
// kinds of loadouts that class is expected to field.
func DefaultSlotsForChassisClass(class string) (small, medium, large int) {
	switch class {
	case ChassisClassLight:
		return 2, 1, 0
	case ChassisClassHeavy:
		return 2, 2, 2
	case ChassisClassAssault:
		return 2, 3, 3
	case ChassisClassMedium:
		fallthrough
	default:
		return 2, 2, 1
	}
}

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
	SmallSlots      int    `db:"small_slots"`
	MediumSlots     int    `db:"medium_slots"`
	LargeSlots      int    `db:"large_slots"`
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
	args[FieldMechaGameChassisSmallSlots] = r.SmallSlots
	args[FieldMechaGameChassisMediumSlots] = r.MediumSlots
	args[FieldMechaGameChassisLargeSlots] = r.LargeSlots
	return args
}
