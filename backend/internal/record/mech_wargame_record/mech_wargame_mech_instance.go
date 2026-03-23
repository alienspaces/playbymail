package mech_wargame_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameMechInstance string = "mech_wargame_mech_instance"
)

const (
	FieldMechWargameMechInstanceID                         string = "id"
	FieldMechWargameMechInstanceGameID                     string = "game_id"
	FieldMechWargameMechInstanceGameInstanceID             string = "game_instance_id"
	FieldMechWargameMechInstanceMechWargameLanceInstanceID string = "mech_wargame_lance_instance_id"
	FieldMechWargameMechInstanceMechWargameSectorInstanceID string = "mech_wargame_sector_instance_id"
	FieldMechWargameMechInstanceMechWargameChassisID       string = "mech_wargame_chassis_id"
	FieldMechWargameMechInstanceCallsign                   string = "callsign"
	FieldMechWargameMechInstanceCurrentArmor               string = "current_armor"
	FieldMechWargameMechInstanceCurrentStructure           string = "current_structure"
	FieldMechWargameMechInstanceCurrentHeat                string = "current_heat"
	FieldMechWargameMechInstancePilotSkill                 string = "pilot_skill"
	FieldMechWargameMechInstanceStatus                     string = "status"
	FieldMechWargameMechInstanceCreatedAt                  string = "created_at"
	FieldMechWargameMechInstanceUpdatedAt                  string = "updated_at"
	FieldMechWargameMechInstanceDeletedAt                  string = "deleted_at"
)

const (
	MechInstanceStatusOperational string = "operational"
	MechInstanceStatusDamaged     string = "damaged"
	MechInstanceStatusDestroyed   string = "destroyed"
	MechInstanceStatusShutdown    string = "shutdown"
)

type MechWargameMechInstance struct {
	record.Record
	GameID                       string `db:"game_id"`
	GameInstanceID               string `db:"game_instance_id"`
	MechWargameLanceInstanceID   string `db:"mech_wargame_lance_instance_id"`
	MechWargameSectorInstanceID  string `db:"mech_wargame_sector_instance_id"`
	MechWargameChassisID         string `db:"mech_wargame_chassis_id"`
	Callsign                     string `db:"callsign"`
	CurrentArmor                 int    `db:"current_armor"`
	CurrentStructure             int    `db:"current_structure"`
	CurrentHeat                  int    `db:"current_heat"`
	PilotSkill                   int    `db:"pilot_skill"`
	Status                       string `db:"status"`
}

func (r *MechWargameMechInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameMechInstanceGameID] = r.GameID
	args[FieldMechWargameMechInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechWargameMechInstanceMechWargameLanceInstanceID] = r.MechWargameLanceInstanceID
	args[FieldMechWargameMechInstanceMechWargameSectorInstanceID] = r.MechWargameSectorInstanceID
	args[FieldMechWargameMechInstanceMechWargameChassisID] = r.MechWargameChassisID
	args[FieldMechWargameMechInstanceCallsign] = r.Callsign
	args[FieldMechWargameMechInstanceCurrentArmor] = r.CurrentArmor
	args[FieldMechWargameMechInstanceCurrentStructure] = r.CurrentStructure
	args[FieldMechWargameMechInstanceCurrentHeat] = r.CurrentHeat
	args[FieldMechWargameMechInstancePilotSkill] = r.PilotSkill
	args[FieldMechWargameMechInstanceStatus] = r.Status
	return args
}
