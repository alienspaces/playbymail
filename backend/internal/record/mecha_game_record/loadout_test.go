package mecha_game_record

import (
	"encoding/json"
	"testing"
)

// TestToNamedArgs_PreservesRawLoadoutOnReadModifyWrite guards against the
// regression where a read-modify-write cycle silently wiped the persisted
// weapon_config / equipment_config. pgx populates the `db:"weapon_config"`
// []byte field during a read but leaves the `db:"-"` struct slice nil.
// Naively re-marshaling the struct slice in ToNamedArgs would then write
// "null" back to jsonb and strip every mech of its loadout.
func TestToNamedArgs_PreservesRawLoadoutOnReadModifyWrite(t *testing.T) {
	rawWeapons := []byte(`[{"weapon_id":"w-1","slot_location":"left-arm"}]`)
	rawEquipment := []byte(`[{"equipment_id":"e-1","slot_location":"head"}]`)

	t.Run("mech_instance uses raw JSON when struct is empty", func(t *testing.T) {
		rec := &MechaGameMechInstance{
			WeaponConfigJSON:    rawWeapons,
			EquipmentConfigJSON: rawEquipment,
		}

		args := rec.ToNamedArgs()

		if got := args[FieldMechaGameMechInstanceWeaponConfig]; got != string(rawWeapons) {
			t.Fatalf("weapon_config: want %q, got %q", rawWeapons, got)
		}
		if got := args[FieldMechaGameMechInstanceEquipmentConfig]; got != string(rawEquipment) {
			t.Fatalf("equipment_config: want %q, got %q", rawEquipment, got)
		}
	})

	t.Run("mech_instance uses struct when populated", func(t *testing.T) {
		rec := &MechaGameMechInstance{
			WeaponConfig: []WeaponConfigEntry{
				{WeaponID: "w-2", SlotLocation: "right-arm"},
			},
			WeaponConfigJSON: rawWeapons,
		}

		args := rec.ToNamedArgs()

		encoded, ok := args[FieldMechaGameMechInstanceWeaponConfig].(string)
		if !ok {
			t.Fatalf("weapon_config: expected string, got %T", args[FieldMechaGameMechInstanceWeaponConfig])
		}
		var out []WeaponConfigEntry
		if err := json.Unmarshal([]byte(encoded), &out); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if len(out) != 1 || out[0].WeaponID != "w-2" {
			t.Fatalf("weapon_config: want struct payload, got %q", encoded)
		}
	})

	t.Run("mech_instance writes empty array when both sources empty", func(t *testing.T) {
		rec := &MechaGameMechInstance{}
		args := rec.ToNamedArgs()
		if got := args[FieldMechaGameMechInstanceWeaponConfig]; got != "[]" {
			t.Fatalf("weapon_config: want %q, got %q", "[]", got)
		}
	})

	t.Run("squad_mech uses raw JSON when struct is empty", func(t *testing.T) {
		rec := &MechaGameSquadMech{
			WeaponConfigJSON:    rawWeapons,
			EquipmentConfigJSON: rawEquipment,
		}

		args := rec.ToNamedArgs()

		if got := args[FieldMechaGameSquadMechWeaponConfig]; got != string(rawWeapons) {
			t.Fatalf("weapon_config: want %q, got %q", rawWeapons, got)
		}
		if got := args[FieldMechaGameSquadMechEquipmentConfig]; got != string(rawEquipment) {
			t.Fatalf("equipment_config: want %q, got %q", rawEquipment, got)
		}
	})
}
