package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// WeaponConfigEntry represents a weapon assigned to a mech at a specific slot.
type WeaponConfigEntry struct {
	WeaponID     string `json:"weapon_id"`
	SlotLocation string `json:"slot_location"`
}

// EquipmentConfigEntry represents an equipment item assigned to a mech at a
// specific slot. Equipment and weapons share the same chassis slot budget via
// the Mountable abstraction.
type EquipmentConfigEntry struct {
	EquipmentID  string `json:"equipment_id"`
	SlotLocation string `json:"slot_location"`
}

const (
	TableMechaGameSquadMech string = "mecha_game_squad_mech"
)

const (
	FieldMechaGameSquadMechID             string = "id"
	FieldMechaGameSquadMechGameID         string = "game_id"
	FieldMechaGameSquadMechMechaGameSquadID   string = "mecha_game_squad_id"
	FieldMechaGameSquadMechMechaGameChassisID string = "mecha_game_chassis_id"
	FieldMechaGameSquadMechCallsign       string = "callsign"
	FieldMechaGameSquadMechWeaponConfig   string = "weapon_config"
	FieldMechaGameSquadMechEquipmentConfig string = "equipment_config"
	FieldMechaGameSquadMechCreatedAt      string = "created_at"
	FieldMechaGameSquadMechUpdatedAt      string = "updated_at"
	FieldMechaGameSquadMechDeletedAt      string = "deleted_at"
)

type MechaGameSquadMech struct {
	record.Record
	GameID           string              `db:"game_id"`
	MechaGameSquadID     string              `db:"mecha_game_squad_id"`
	MechaGameChassisID   string              `db:"mecha_game_chassis_id"`
	Callsign            string                 `db:"callsign"`
	WeaponConfig        []WeaponConfigEntry    `db:"-"`
	WeaponConfigJSON    []byte                 `db:"weapon_config"`
	EquipmentConfig     []EquipmentConfigEntry `db:"-"`
	EquipmentConfigJSON []byte                 `db:"equipment_config"`
}

func (r *MechaGameSquadMech) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameSquadMechGameID] = r.GameID
	args[FieldMechaGameSquadMechMechaGameSquadID] = r.MechaGameSquadID
	args[FieldMechaGameSquadMechMechaGameChassisID] = r.MechaGameChassisID
	args[FieldMechaGameSquadMechCallsign] = r.Callsign
	// See MechaGameMechInstance.ToNamedArgs: falling back to the raw JSON
	// bytes when the decoded struct is empty prevents read-modify-write
	// cycles from nulling out the persisted loadout.
	args[FieldMechaGameSquadMechWeaponConfig] = marshalConfigForWrite(r.WeaponConfig, r.WeaponConfigJSON)
	args[FieldMechaGameSquadMechEquipmentConfig] = marshalConfigForWrite(r.EquipmentConfig, r.EquipmentConfigJSON)
	return args
}
