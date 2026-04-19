package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/record/mecha_game_record"
)

// newTestChassisSlots builds a chassis record with just enough data for the
// slot-fit validator. Stats other than the slot counts are left at zero so
// the helper is obviously only suitable for loadout tests.
func newTestChassisSlots(small, medium, large int) *mecha_game_record.MechaGameChassis {
	return &mecha_game_record.MechaGameChassis{
		SmallSlots:  small,
		MediumSlots: medium,
		LargeSlots:  large,
	}
}

func TestSquadMechLoadoutFitResult(t *testing.T) {
	// Shared item catalog: weapons and equipment keyed by the same ids used
	// in WeaponConfigEntry / EquipmentConfigEntry below.
	smallLaser := &mecha_game_record.MechaGameWeapon{Name: "Small Laser", MountSize: mecha_game_record.WeaponMountSizeSmall}
	plasma := &mecha_game_record.MechaGameWeapon{Name: "Plasma", MountSize: mecha_game_record.WeaponMountSizeLarge}
	heatSink := &mecha_game_record.MechaGameEquipment{Name: "Heat Sink", MountSize: mecha_game_record.EquipmentMountSizeSmall}
	giantECM := &mecha_game_record.MechaGameEquipment{Name: "Giant ECM Pod", MountSize: mecha_game_record.EquipmentMountSizeLarge}

	weapons := map[string]*mecha_game_record.MechaGameWeapon{
		"small_laser": smallLaser,
		"plasma":      plasma,
	}
	equipment := map[string]*mecha_game_record.MechaGameEquipment{
		"heat_sink": heatSink,
		"giant_ecm": giantECM,
	}

	cases := []struct {
		name          string
		chassis       *mecha_game_record.MechaGameChassis
		weaponCfg     []mecha_game_record.WeaponConfigEntry
		equipmentCfg  []mecha_game_record.EquipmentConfigEntry
		wantErr       bool
		errField      string // field name expected in the InvalidField error
		errContains   string // substring (typically an item label) expected in the error message
	}{
		{
			name:    "empty loadout on empty chassis passes",
			chassis: newTestChassisSlots(0, 0, 0),
		},
		{
			name:    "nil chassis is a no-op even with entries present",
			chassis: nil,
			weaponCfg: []mecha_game_record.WeaponConfigEntry{
				{WeaponID: "small_laser"},
			},
			equipmentCfg: []mecha_game_record.EquipmentConfigEntry{
				{EquipmentID: "heat_sink"},
			},
		},
		{
			name:    "weapon and equipment share a single small slot — overflow points at equipment_config",
			chassis: newTestChassisSlots(1, 0, 0),
			weaponCfg: []mecha_game_record.WeaponConfigEntry{
				{WeaponID: "small_laser"},
			},
			equipmentCfg: []mecha_game_record.EquipmentConfigEntry{
				{EquipmentID: "heat_sink"},
			},
			wantErr:     true,
			errField:    mecha_game_record.FieldMechaGameSquadMechEquipmentConfig,
			errContains: "Heat Sink",
		},
		{
			name:    "weapon-only overflow points at weapon_config and names the weapon",
			chassis: newTestChassisSlots(0, 0, 0),
			weaponCfg: []mecha_game_record.WeaponConfigEntry{
				{WeaponID: "plasma"},
			},
			wantErr:     true,
			errField:    mecha_game_record.FieldMechaGameSquadMechWeaponConfig,
			errContains: "Plasma",
		},
		{
			name:    "equipment-only overflow points at equipment_config and names the equipment",
			chassis: newTestChassisSlots(0, 0, 0),
			equipmentCfg: []mecha_game_record.EquipmentConfigEntry{
				{EquipmentID: "giant_ecm"},
			},
			wantErr:     true,
			errField:    mecha_game_record.FieldMechaGameSquadMechEquipmentConfig,
			errContains: "Giant ECM Pod",
		},
		{
			name:    "small weapon and small equipment spill up into available large slots",
			chassis: newTestChassisSlots(0, 0, 2),
			weaponCfg: []mecha_game_record.WeaponConfigEntry{
				{WeaponID: "small_laser"},
			},
			equipmentCfg: []mecha_game_record.EquipmentConfigEntry{
				{EquipmentID: "heat_sink"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := squadMechLoadoutFitResult(tc.chassis, tc.weaponCfg, weapons, tc.equipmentCfg, equipment)
			if !tc.wantErr {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.errField)
			require.Contains(t, err.Error(), tc.errContains)
		})
	}
}
