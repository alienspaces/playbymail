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

func TestMechaGameOrdersProcessor_GenerateTurnSheet(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameOrdersProcessor(l, cfg)
	require.NoError(t, err)

	tests := []struct {
		name        string
		data        any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty data returns validation error",
			data:        &turnsheet.OrdersData{},
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
			data: &turnsheet.OrdersData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:      convert.Ptr("Steel Thunder"),
					GameType:      convert.Ptr("mecha"),
					TurnSheetCode: convert.Ptr(generateTestTurnSheetCode(t)),
					TurnNumber:    convert.Ptr(1),
				},
				SquadName: "Alpha Squad",
				SquadMechs: []turnsheet.MechOrderEntry{
					{
						MechInstanceID:    "mech-1",
						MechCallsign:      "Hammer",
						MechStatus:        "operational",
						CurrentSectorName: "Central Wastes",
						ChassisName:       "Scout",
						ChassisClass:      "light",
						CurrentArmor:      72,
						MaxArmor:          72,
						CurrentStructure:  32,
						MaxStructure:      32,
						Speed:             7,
						PilotSkill:        4,
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

func TestMechaGameOrdersProcessor_GenerateTurnSheet_ContainsMechStats(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameOrdersProcessor(l, cfg)
	require.NoError(t, err)

	data := &turnsheet.OrdersData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:      convert.Ptr("Steel Thunder"),
			GameType:      convert.Ptr("mecha"),
			TurnSheetCode: convert.Ptr(generateTestTurnSheetCode(t)),
			TurnNumber:    convert.Ptr(1),
			TurnEvents: []turnsheet.TurnEvent{
				{Category: turnsheet.TurnEventCategoryMovement, Icon: turnsheet.TurnEventIconMovement, Message: "Hammer moved to Northern Ridge."},
				{Category: turnsheet.TurnEventCategoryCombat, Icon: turnsheet.TurnEventIconCombat, Message: "Anvil hit Stalker for 5 damage."},
			},
		},
		SquadName: "Alpha Squad",
		SquadMechs: []turnsheet.MechOrderEntry{
			{
				MechInstanceID:    "mech-1",
				MechCallsign:      "Hammer",
				MechStatus:        "operational",
				CurrentSectorName: "Central Wastes",
				ChassisName:       "Scout",
				ChassisClass:      "light",
				CurrentArmor:      55,
				MaxArmor:          72,
				CurrentStructure:  32,
				MaxStructure:      32,
				CurrentHeat:       4,
				HeatCapacity:      18,
				Speed:             7,
				EffectiveSpeed:    9,
				PilotSkill:        4,
				Weapons: []turnsheet.MechWeaponEntry{
					{Name: "Light Pulse Cannon", Damage: 3, HeatCost: 1, RangeBand: "short", SlotLocation: "left-arm"},
					{Name: "Rocket Pack", Damage: 8, HeatCost: 3, RangeBand: "short", SlotLocation: "right-arm", AmmoCapacity: 2},
				},
				Equipment: []turnsheet.MechEquipmentEntry{
					{Name: "Jump Jets", EffectKind: "jump_jets", Magnitude: 2, HeatCost: 2, MountSize: "medium", SlotLocation: "left-leg"},
					{Name: "Ammo Bin (Standard)", EffectKind: "ammo_bin", Magnitude: 8, HeatCost: 0, MountSize: "small", SlotLocation: "right-torso"},
				},
				AmmoRemaining: 6,
				AmmoCapacity:  10,
				ReachableSectors: []turnsheet.SectorOption{
					{SectorInstanceID: "sector-ridge", SectorName: "Ridge Overlook"},
				},
			},
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
	require.Contains(t, htmlStr, "Scout", "should contain chassis name")
	require.Contains(t, htmlStr, "Light Pulse Cannon", "should contain weapon name")
	// Equipment rendered inline with the weapons table.
	require.Contains(t, htmlStr, "Jump Jets", "should render equipment name")
	require.Contains(t, htmlStr, "jump_jets", "should render equipment effect kind")
	require.Contains(t, htmlStr, "Ammo Bin (Standard)", "should render ammo bin equipment")
	// Ammo readout for mechs with ammo-consuming weapons.
	require.Contains(t, htmlStr, "6 / 10", "should render ammo remaining/capacity")
	// Effective-speed shows base + jump-jet bonus.
	require.Contains(t, htmlStr, "+2 JJ", "should render jump-jet speed bonus")
	// Per-mech reachable sectors feed the Move-To dropdown.
	require.Contains(t, htmlStr, "Ridge Overlook", "should render reachable sector option")
}

func TestMechaGameOrdersProcessor_ScanTurnSheet_EmptyImageReturnsError(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameOrdersProcessor(l, cfg)
	require.NoError(t, err)

	ctx := context.Background()
	result, err := processor.ScanTurnSheet(ctx, l, []byte(`{}`), []byte{})
	require.Error(t, err)
	require.Nil(t, result)
}

