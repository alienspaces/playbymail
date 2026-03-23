package mech_wargame_record

import (
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameLanceMech string = "mech_wargame_lance_mech"
)

const (
	FieldMechWargameLanceMechID                   string = "id"
	FieldMechWargameLanceMechGameID               string = "game_id"
	FieldMechWargameLanceMechMechWargameLanceID   string = "mech_wargame_lance_id"
	FieldMechWargameLanceMechMechWargameChassisID string = "mech_wargame_chassis_id"
	FieldMechWargameLanceMechCallsign             string = "callsign"
	FieldMechWargameLanceMechWeaponConfig         string = "weapon_config"
	FieldMechWargameLanceMechCreatedAt            string = "created_at"
	FieldMechWargameLanceMechUpdatedAt            string = "updated_at"
	FieldMechWargameLanceMechDeletedAt            string = "deleted_at"
)

// WeaponConfigEntry represents a weapon assigned to a mech at a specific slot.
type WeaponConfigEntry struct {
	WeaponID     string `json:"weapon_id"`
	SlotLocation string `json:"slot_location"`
}

type MechWargameLanceMech struct {
	record.Record
	GameID                  string              `db:"game_id"`
	MechWargameLanceID      string              `db:"mech_wargame_lance_id"`
	MechWargameChassisID    string              `db:"mech_wargame_chassis_id"`
	Callsign                string              `db:"callsign"`
	WeaponConfig            []WeaponConfigEntry `db:"-"`
	WeaponConfigJSON        []byte              `db:"weapon_config"`
}

func (r *MechWargameLanceMech) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameLanceMechGameID] = r.GameID
	args[FieldMechWargameLanceMechMechWargameLanceID] = r.MechWargameLanceID
	args[FieldMechWargameLanceMechMechWargameChassisID] = r.MechWargameChassisID
	args[FieldMechWargameLanceMechCallsign] = r.Callsign

	weaponJSON, _ := json.Marshal(r.WeaponConfig)
	args[FieldMechWargameLanceMechWeaponConfig] = string(weaponJSON)
	return args
}
