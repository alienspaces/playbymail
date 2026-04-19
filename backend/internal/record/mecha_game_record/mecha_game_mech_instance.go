package mecha_game_record

import (
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameMechInstance string = "mecha_game_mech_instance"
)

const (
	FieldMechaGameMechInstanceID                    string = "id"
	FieldMechaGameMechInstanceGameID                string = "game_id"
	FieldMechaGameMechInstanceGameInstanceID        string = "game_instance_id"
	FieldMechaGameMechInstanceMechaGameSquadInstanceID  string = "mecha_game_squad_instance_id"
	FieldMechaGameMechInstanceMechaGameSectorInstanceID string = "mecha_game_sector_instance_id"
	FieldMechaGameMechInstanceMechaGameChassisID        string = "mecha_game_chassis_id"
	FieldMechaGameMechInstanceCallsign              string = "callsign"
	FieldMechaGameMechInstanceCurrentArmor          string = "current_armor"
	FieldMechaGameMechInstanceCurrentStructure      string = "current_structure"
	FieldMechaGameMechInstanceCurrentHeat           string = "current_heat"
	FieldMechaGameMechInstancePilotSkill            string = "pilot_skill"
	FieldMechaGameMechInstanceExperiencePoints      string = "experience_points"
	FieldMechaGameMechInstanceStatus                string = "status"
	FieldMechaGameMechInstanceWeaponConfig          string = "weapon_config"
	FieldMechaGameMechInstanceEquipmentConfig       string = "equipment_config"
	FieldMechaGameMechInstanceAmmoRemaining         string = "ammo_remaining"
	FieldMechaGameMechInstanceIsRefitting           string = "is_refitting"
	FieldMechaGameMechInstanceCreatedAt             string = "created_at"
	FieldMechaGameMechInstanceUpdatedAt             string = "updated_at"
	FieldMechaGameMechInstanceDeletedAt             string = "deleted_at"
)

const (
	MechInstanceStatusOperational string = "operational"
	MechInstanceStatusDamaged     string = "damaged"
	MechInstanceStatusDestroyed   string = "destroyed"
	MechInstanceStatusShutdown    string = "shutdown"
)

type MechaGameMechInstance struct {
	record.Record
	GameID                string              `db:"game_id"`
	GameInstanceID        string              `db:"game_instance_id"`
	MechaGameSquadInstanceID  string              `db:"mecha_game_squad_instance_id"`
	MechaGameSectorInstanceID string              `db:"mecha_game_sector_instance_id"`
	MechaGameChassisID        string              `db:"mecha_game_chassis_id"`
	Callsign              string              `db:"callsign"`
	CurrentArmor          int                 `db:"current_armor"`
	CurrentStructure      int                 `db:"current_structure"`
	CurrentHeat           int                 `db:"current_heat"`
	PilotSkill            int                 `db:"pilot_skill"`
	ExperiencePoints      int                 `db:"experience_points"`
	Status                string                 `db:"status"`
	WeaponConfig          []WeaponConfigEntry    `db:"-"`
	WeaponConfigJSON      []byte                 `db:"weapon_config"`
	EquipmentConfig       []EquipmentConfigEntry `db:"-"`
	EquipmentConfigJSON   []byte                 `db:"equipment_config"`
	AmmoRemaining         int                    `db:"ammo_remaining"`
	IsRefitting           bool                   `db:"is_refitting"`
}

func (r *MechaGameMechInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameMechInstanceGameID] = r.GameID
	args[FieldMechaGameMechInstanceGameInstanceID] = r.GameInstanceID
	args[FieldMechaGameMechInstanceMechaGameSquadInstanceID] = r.MechaGameSquadInstanceID
	args[FieldMechaGameMechInstanceMechaGameSectorInstanceID] = r.MechaGameSectorInstanceID
	args[FieldMechaGameMechInstanceMechaGameChassisID] = r.MechaGameChassisID
	args[FieldMechaGameMechInstanceCallsign] = r.Callsign
	args[FieldMechaGameMechInstanceCurrentArmor] = r.CurrentArmor
	args[FieldMechaGameMechInstanceCurrentStructure] = r.CurrentStructure
	args[FieldMechaGameMechInstanceCurrentHeat] = r.CurrentHeat
	args[FieldMechaGameMechInstancePilotSkill] = r.PilotSkill
	args[FieldMechaGameMechInstanceExperiencePoints] = r.ExperiencePoints
	args[FieldMechaGameMechInstanceStatus] = r.Status
	weaponJSON, _ := json.Marshal(r.WeaponConfig)
	args[FieldMechaGameMechInstanceWeaponConfig] = string(weaponJSON)
	equipmentJSON, _ := json.Marshal(r.EquipmentConfig)
	args[FieldMechaGameMechInstanceEquipmentConfig] = string(equipmentJSON)
	args[FieldMechaGameMechInstanceAmmoRemaining] = r.AmmoRemaining
	args[FieldMechaGameMechInstanceIsRefitting] = r.IsRefitting
	return args
}
