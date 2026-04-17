package mecha_record

import (
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// WeaponConfigEntry represents a weapon assigned to a mech at a specific slot.
type WeaponConfigEntry struct {
	WeaponID     string `json:"weapon_id"`
	SlotLocation string `json:"slot_location"`
}

const (
	TableMechaSquadMech string = "mecha_squad_mech"
)

const (
	FieldMechaSquadMechID             string = "id"
	FieldMechaSquadMechGameID         string = "game_id"
	FieldMechaSquadMechMechaSquadID   string = "mecha_squad_id"
	FieldMechaSquadMechMechaChassisID string = "mecha_chassis_id"
	FieldMechaSquadMechCallsign       string = "callsign"
	FieldMechaSquadMechWeaponConfig   string = "weapon_config"
	FieldMechaSquadMechCreatedAt      string = "created_at"
	FieldMechaSquadMechUpdatedAt      string = "updated_at"
	FieldMechaSquadMechDeletedAt      string = "deleted_at"
)

type MechaSquadMech struct {
	record.Record
	GameID           string              `db:"game_id"`
	MechaSquadID     string              `db:"mecha_squad_id"`
	MechaChassisID   string              `db:"mecha_chassis_id"`
	Callsign         string              `db:"callsign"`
	WeaponConfig     []WeaponConfigEntry `db:"-"`
	WeaponConfigJSON []byte              `db:"weapon_config"`
}

func (r *MechaSquadMech) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaSquadMechGameID] = r.GameID
	args[FieldMechaSquadMechMechaSquadID] = r.MechaSquadID
	args[FieldMechaSquadMechMechaChassisID] = r.MechaChassisID
	args[FieldMechaSquadMechCallsign] = r.Callsign
	weaponJSON, _ := json.Marshal(r.WeaponConfig)
	args[FieldMechaSquadMechWeaponConfig] = string(weaponJSON)
	return args
}
