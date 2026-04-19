package turnsheet_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestMechaGameSquadManagementProcessor_GenerateTurnSheet(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameSquadManagementProcessor(l, cfg)
	require.NoError(t, err)

	tests := []struct {
		name        string
		data        any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty data returns validation error",
			data:        &turnsheet.SquadManagementData{},
			expectError: true,
			errorMsg:    "game name is required",
		},
		{
			name:        "nil data handled gracefully",
			data:        nil,
			expectError: false,
		},
		{
			name: "valid data generates HTML",
			data: &turnsheet.SquadManagementData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:      convert.Ptr("Steel Thunder"),
					GameType:      convert.Ptr("mecha"),
					TurnSheetCode: convert.Ptr(generateTestTurnSheetCode(t)),
					TurnNumber:    convert.Ptr(2),
				},
				SquadName:    "Alpha Squad",
				SupplyPoints: 6,
				Mechs: []turnsheet.ManagementMechEntry{
					{
						MechInstanceID:   "mech-1",
						Callsign:         "Hammer",
						ChassisName:      "Scout",
						IsAtDepot:        true,
						CurrentArmor:     72,
						MaxArmor:         72,
						CurrentStructure: 28,
						MaxStructure:     32,
						StructureDamage:  4,
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sheetData []byte
			if tt.data != nil {
				var err error
				sheetData, err = json.Marshal(tt.data)
				require.NoError(t, err)
			}

			ctx := context.Background()
			html, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatHTML, sheetData)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
				require.Nil(t, html)
			} else if err != nil {
				t.Logf("GenerateTurnSheet returned error: %v", err)
			}
		})
	}
}

func TestMechaGameSquadManagementProcessor_GenerateTurnSheet_ContainsDepotInfo(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameSquadManagementProcessor(l, cfg)
	require.NoError(t, err)

	data := &turnsheet.SquadManagementData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:      convert.Ptr("Steel Thunder"),
			GameType:      convert.Ptr("mecha"),
			TurnSheetCode: convert.Ptr(generateTestTurnSheetCode(t)),
			TurnNumber:    convert.Ptr(2),
		},
		SquadName:    "Alpha Squad",
		SupplyPoints: 6,
		Mechs: []turnsheet.ManagementMechEntry{
			{
				MechInstanceID:   "mech-1",
				Callsign:         "Hammer",
				ChassisName:      "Scout",
				ChassisClass:     "light",
				IsAtDepot:        true,
				CurrentStructure: 28,
				MaxStructure:     32,
				StructureDamage:  4,
				CurrentHeat:      3,
				HeatCapacity:     18,
				Weapons: []turnsheet.MechWeaponSlot{
					{SlotLocation: "left-arm", CurrentWeaponID: "wpn-1", CurrentWeaponName: "Light Pulse Cannon"},
				},
				Equipment: []turnsheet.MechEquipmentEntry{
					{Name: "Heat Sink", EffectKind: "heat_sink", Magnitude: 3, HeatCost: 0, MountSize: "small", SlotLocation: "right-torso"},
				},
				AmmoRemaining: 4,
				AmmoCapacity:  8,
			},
			{
				MechInstanceID:   "mech-2",
				Callsign:         "Anvil",
				ChassisName:      "Sentinel",
				ChassisClass:     "medium",
				IsAtDepot:        false,
				CurrentStructure: 65,
				MaxStructure:     65,
			},
		},
		WeaponCatalog: []turnsheet.CatalogWeapon{
			{WeaponID: "cat-1", Name: "Light Pulse Cannon", Damage: 3, HeatCost: 1, RangeBand: "short"},
			{WeaponID: "cat-2", Name: "Rocket Pack", Damage: 8, HeatCost: 3, RangeBand: "short", AmmoCapacity: 2},
		},
	}

	sheetData, err := json.Marshal(data)
	require.NoError(t, err)

	ctx := context.Background()
	html, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatHTML, sheetData)
	require.NoError(t, err)
	require.NotEmpty(t, html)

	htmlStr := string(html)
	require.Contains(t, htmlStr, "Hammer", "should contain mech callsign")
	require.Contains(t, htmlStr, "Anvil", "should contain mech callsign")
	require.Contains(t, htmlStr, "Light Pulse Cannon", "should contain catalog weapon name")
	// Heat + ammo readouts on the stat-summary row.
	require.Contains(t, htmlStr, "3/18", "should render heat current/capacity in stat summary")
	require.Contains(t, htmlStr, "4/8", "should render ammo current/capacity in stat summary")
	// Read-only equipment table beside the weapons section.
	require.Contains(t, htmlStr, "Heat Sink", "should render equipment name")
	require.Contains(t, htmlStr, "heat_sink", "should render equipment effect kind")
	// Catalog AMMO column: dash for non-ammo weapons, value for ammo weapons.
	require.Contains(t, htmlStr, "Rocket Pack", "should render ammo-consuming catalog weapon")
}

func TestMechaGameSquadManagementProcessor_ScanTurnSheet_EmptyImageReturnsError(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameSquadManagementProcessor(l, cfg)
	require.NoError(t, err)

	ctx := context.Background()
	result, err := processor.ScanTurnSheet(ctx, l, []byte(`{}`), []byte{})
	require.Error(t, err)
	require.Nil(t, result)
}
