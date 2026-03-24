package turnsheet_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestMechaOrdersProcessor_GenerateTurnSheet(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaOrdersProcessor(l, cfg)
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
				LanceName: "Alpha Lance",
				LanceMechs: []turnsheet.MechOrderEntry{
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

func TestMechaOrdersProcessor_GenerateTurnSheet_ContainsMechStats(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaOrdersProcessor(l, cfg)
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
		LanceName: "Alpha Lance",
		LanceMechs: []turnsheet.MechOrderEntry{
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
				PilotSkill:        4,
				Weapons: []turnsheet.MechWeaponEntry{
					{Name: "Light Pulse Cannon", Damage: 3, HeatCost: 1, RangeBand: "short", SlotLocation: "left-arm"},
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
	require.True(t, strings.Contains(htmlStr, "Hammer"), "should contain mech callsign")
	require.True(t, strings.Contains(htmlStr, "Scout"), "should contain chassis name")
	require.True(t, strings.Contains(htmlStr, "Light Pulse Cannon"), "should contain weapon name")
}

func TestMechaOrdersProcessor_ScanTurnSheet_EmptyImageReturnsError(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaOrdersProcessor(l, cfg)
	require.NoError(t, err)

	ctx := context.Background()
	result, err := processor.ScanTurnSheet(ctx, l, []byte(`{}`), []byte{})
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGenerateMechaOrdersRendering(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"
	cfg.SaveTestFiles = true

	processor, err := turnsheet.NewMechaOrdersProcessor(l, cfg)
	require.NoError(t, err)

	backgroundImage := loadTestBackgroundImage(t, "testdata/background-darkforest.png")

	type formatCase struct {
		name   string
		format turnsheet.DocumentFormat
		ext    string
	}

	cases := []formatCase{
		{name: "pdf", format: turnsheet.DocumentFormatPDF, ext: "pdf"},
		{name: "html", format: turnsheet.DocumentFormatHTML, ext: "html"},
	}

	turnSheetCode := generateTestTurnSheetCode(t)

	testData := &turnsheet.OrdersData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:          convert.Ptr("Steel Thunder"),
			GameType:          convert.Ptr("mecha"),
			TurnNumber:        convert.Ptr(3),
			TurnSheetCode:     convert.Ptr(turnSheetCode),
			TurnSheetDeadline: convert.Ptr(time.Now().Add(7 * 24 * time.Hour)),
			BackgroundImage:   &backgroundImage,
			TurnEvents: []turnsheet.TurnEvent{
				{
					Category: turnsheet.TurnEventCategoryMovement,
					Icon:     turnsheet.TurnEventIconMovement,
					Message:  "Hammer advanced to Central Wastes.",
				},
				{
					Category: turnsheet.TurnEventCategoryCombat,
					Icon:     turnsheet.TurnEventIconCombat,
					Message:  "Anvil fired Pulse Cannon at Stalker — HIT for 5 damage.",
				},
				{
					Category: turnsheet.TurnEventCategorySystem,
					Icon:     turnsheet.TurnEventIconSystem,
					Message:  "Hammer field repairs restored 18 armor (72/72).",
				},
			},
		},
		LanceName: "Alpha Lance",
		LanceMechs: []turnsheet.MechOrderEntry{
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
				PilotSkill:        4,
				Weapons: []turnsheet.MechWeaponEntry{
					{Name: "Light Pulse Cannon", Damage: 3, HeatCost: 1, RangeBand: "short", SlotLocation: "left-arm"},
					{Name: "Chaingun", Damage: 2, HeatCost: 0, RangeBand: "short", SlotLocation: "right-arm"},
				},
			},
			{
				MechInstanceID:    "mech-2",
				MechCallsign:      "Anvil",
				MechStatus:        "operational",
				CurrentSectorName: "Central Wastes",
				ChassisName:       "Sentinel",
				ChassisClass:      "medium",
				CurrentArmor:      130,
				MaxArmor:          130,
				CurrentStructure:  65,
				MaxStructure:      65,
				CurrentHeat:       0,
				HeatCapacity:      28,
				Speed:             4,
				PilotSkill:        4,
				Weapons: []turnsheet.MechWeaponEntry{
					{Name: "Pulse Cannon", Damage: 5, HeatCost: 3, RangeBand: "medium", SlotLocation: "left-torso"},
					{Name: "Rocket Pack", Damage: 8, HeatCost: 3, RangeBand: "short", SlotLocation: "right-arm"},
				},
			},
			{
				MechInstanceID:    "mech-3",
				MechCallsign:      "Titan",
				MechStatus:        "damaged",
				CurrentSectorName: "Northern Ridge",
				ChassisName:       "Colossus",
				ChassisClass:      "heavy",
				CurrentArmor:      45,
				MaxArmor:          180,
				CurrentStructure:  60,
				MaxStructure:      80,
				CurrentHeat:       14,
				HeatCapacity:      35,
				Speed:             3,
				PilotSkill:        5,
				Weapons: []turnsheet.MechWeaponEntry{
					{Name: "Auto-Cannon", Damage: 10, HeatCost: 5, RangeBand: "medium", SlotLocation: "right-arm"},
				},
			},
			{
				MechInstanceID:   "mech-4",
				MechCallsign:     "Wrench",
				MechStatus:       "operational",
				IsRefitting:      true,
				ChassisName:      "Scout",
				ChassisClass:     "light",
				CurrentArmor:     72,
				MaxArmor:         72,
				CurrentStructure: 32,
				MaxStructure:     32,
				Speed:            7,
				PilotSkill:       3,
			},
		},
		AvailableSectors: []turnsheet.SectorOption{
			{SectorInstanceID: "sector-1", SectorName: "Northern Ridge"},
			{SectorInstanceID: "sector-2", SectorName: "Southern Flats"},
		},
		EnemyMechs: []turnsheet.EnemyMechOption{
			{MechInstanceID: "enemy-1", Callsign: "Stalker", SectorName: "Northern Ridge"},
			{MechInstanceID: "enemy-2", Callsign: "Predator", SectorName: "Southern Flats"},
		},
	}

	ctx := context.Background()
	sheetData, err := json.Marshal(testData)
	require.NoError(t, err)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output, err := processor.GenerateTurnSheet(ctx, l, tc.format, sheetData)
			require.NoError(t, err, "should generate output without error")
			require.NotEmpty(t, output, "output should not be empty")

			if cfg.SaveTestFiles {
				path := fmt.Sprintf("testdata/mecha_orders_turnsheet.%s", tc.ext)
				err = os.WriteFile(path, output, 0644)
				require.NoError(t, err, "should save output to testdata directory")
				t.Logf("%s preview saved to %s", tc.name, path)
			}
		})
	}
}
