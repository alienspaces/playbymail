package mecha_record

import (
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaLanceMech string = "mecha_lance_mech"
)

const (
	FieldMechaLanceMechID                   string = "id"
	FieldMechaLanceMechGameID               string = "game_id"
	FieldMechaLanceMechMechaLanceID   string = "mecha_lance_id"
	FieldMechaLanceMechMechaChassisID string = "mecha_chassis_id"
	FieldMechaLanceMechCallsign             string = "callsign"
	FieldMechaLanceMechWeaponConfig         string = "weapon_config"
	FieldMechaLanceMechCreatedAt            string = "created_at"
	FieldMechaLanceMechUpdatedAt            string = "updated_at"
	FieldMechaLanceMechDeletedAt            string = "deleted_at"
)

// WeaponConfigEntry represents a weapon assigned to a mech at a specific slot.
type WeaponConfigEntry struct {
	WeaponID     string `json:"weapon_id"`
	SlotLocation string `json:"slot_location"`
}

type MechaLanceMech struct {
	record.Record
	GameID                  string              `db:"game_id"`
	MechaLanceID      string              `db:"mecha_lance_id"`
	MechaChassisID    string              `db:"mecha_chassis_id"`
	Callsign                string              `db:"callsign"`
	WeaponConfig            []WeaponConfigEntry `db:"-"`
	WeaponConfigJSON        []byte              `db:"weapon_config"`
}

func (r *MechaLanceMech) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaLanceMechGameID] = r.GameID
	args[FieldMechaLanceMechMechaLanceID] = r.MechaLanceID
	args[FieldMechaLanceMechMechaChassisID] = r.MechaChassisID
	args[FieldMechaLanceMechCallsign] = r.Callsign

	weaponJSON, _ := json.Marshal(r.WeaponConfig)
	args[FieldMechaLanceMechWeaponConfig] = string(weaponJSON)
	return args
}
