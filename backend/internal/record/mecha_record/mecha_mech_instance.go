package mecha_record

import (
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaMechInstance string = "mecha_mech_instance"
)

const (
	FieldMechaMechInstanceID                    string = "id"
	FieldMechaMechInstanceGameID                string = "game_id"
	FieldMechaMechInstanceGameInstanceID        string = "game_instance_id"
	FieldMechaMechInstanceMechaLanceInstanceID  string = "mecha_lance_instance_id"
	FieldMechaMechInstanceMechaSectorInstanceID string = "mecha_sector_instance_id"
	FieldMechaMechInstanceMechaChassisID        string = "mecha_chassis_id"
	FieldMechaMechInstanceCallsign              string = "callsign"
	FieldMechaMechInstanceCurrentArmor          string = "current_armor"
	FieldMechaMechInstanceCurrentStructure      string = "current_structure"
	FieldMechaMechInstanceCurrentHeat           string = "current_heat"
	FieldMechaMechInstancePilotSkill            string = "pilot_skill"
	FieldMechaMechInstanceStatus                string = "status"
	FieldMechaMechInstanceWeaponConfig          string = "weapon_config"
	FieldMechaMechInstanceIsRefitting           string = "is_refitting"
	FieldMechaMechInstanceCreatedAt             string = "created_at"
	FieldMechaMechInstanceUpdatedAt             string = "updated_at"
	FieldMechaMechInstanceDeletedAt             string = "deleted_at"
)

const (
	MechInstanceStatusOperational string = "operational"
	MechInstanceStatusDamaged     string = "damaged"
	MechInstanceStatusDestroyed   string = "destroyed"
	MechInstanceStatusShutdown    string = "shutdown"
)

type MechaMechInstance struct {
	record.Record
	GameID                string              `db:"game_id"`
	GameInstanceID        string              `db:"game_instance_id"`
	MechaLanceInstanceID  string              `db:"mecha_lance_instance_id"`
	MechaSectorInstanceID string              `db:"mecha_sector_instance_id"`
	MechaChassisID        string              `db:"mecha_chassis_id"`
	Callsign              string              `db:"callsign"`
	CurrentArmor          int                 `db:"current_armor"`
	CurrentStructure      int                 `db:"current_structure"`
	CurrentHeat           int                 `db:"current_heat"`
	PilotSkill            int                 `db:"pilot_skill"`
	Status                string              `db:"status"`
	WeaponConfig          []WeaponConfigEntry `db:"-"`
	WeaponConfigJSON      []byte              `db:"weapon_config"`
	IsRefitting           bool                `db:"is_refitting"`
}

func (r *MechaMechInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaMechInstanceGameID] = r.GameID
	args[FieldMechaMechInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechaMechInstanceMechaLanceInstanceID] = r.MechaLanceInstanceID
	args[FieldMechaMechInstanceMechaSectorInstanceID] = r.MechaSectorInstanceID
	args[FieldMechaMechInstanceMechaChassisID] = r.MechaChassisID
	args[FieldMechaMechInstanceCallsign] = r.Callsign
	args[FieldMechaMechInstanceCurrentArmor] = r.CurrentArmor
	args[FieldMechaMechInstanceCurrentStructure] = r.CurrentStructure
	args[FieldMechaMechInstanceCurrentHeat] = r.CurrentHeat
	args[FieldMechaMechInstancePilotSkill] = r.PilotSkill
	args[FieldMechaMechInstanceStatus] = r.Status
	weaponJSON, _ := json.Marshal(r.WeaponConfig)
	args[FieldMechaMechInstanceWeaponConfig] = string(weaponJSON)
	args[FieldMechaMechInstanceIsRefitting] = r.IsRefitting
	return args
}
