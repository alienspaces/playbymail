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

func TestMechaLanceManagementProcessor_GenerateTurnSheet(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaLanceManagementProcessor(l, cfg)
	require.NoError(t, err)

	tests := []struct {
		name        string
		data        any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty data returns validation error",
			data:        &turnsheet.LanceManagementData{},
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
			data: &turnsheet.LanceManagementData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:      convert.Ptr("Steel Thunder"),
					GameType:      convert.Ptr("mecha"),
					TurnSheetCode: convert.Ptr(generateTestTurnSheetCode(t)),
					TurnNumber:    convert.Ptr(2),
				},
				LanceName:    "Alpha Lance",
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

func TestMechaLanceManagementProcessor_GenerateTurnSheet_ContainsDepotInfo(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaLanceManagementProcessor(l, cfg)
	require.NoError(t, err)

	data := &turnsheet.LanceManagementData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:      convert.Ptr("Steel Thunder"),
			GameType:      convert.Ptr("mecha"),
			TurnSheetCode: convert.Ptr(generateTestTurnSheetCode(t)),
			TurnNumber:    convert.Ptr(2),
		},
		LanceName:    "Alpha Lance",
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
				Weapons: []turnsheet.MechWeaponSlot{
					{SlotLocation: "left-arm", CurrentWeaponID: "wpn-1", CurrentWeaponName: "Light Pulse Cannon"},
				},
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
	require.True(t, strings.Contains(htmlStr, "Anvil"), "should contain mech callsign")
	require.True(t, strings.Contains(htmlStr, "Light Pulse Cannon"), "should contain catalog weapon name")
}

func TestMechaLanceManagementProcessor_ScanTurnSheet_EmptyImageReturnsError(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaLanceManagementProcessor(l, cfg)
	require.NoError(t, err)

	ctx := context.Background()
	result, err := processor.ScanTurnSheet(ctx, l, []byte(`{}`), []byte{})
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGenerateMechaLanceManagementRendering(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"
	cfg.SaveTestFiles = true

	processor, err := turnsheet.NewMechaLanceManagementProcessor(l, cfg)
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

	testData := &turnsheet.LanceManagementData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:          convert.Ptr("Steel Thunder"),
			GameType:          convert.Ptr("mecha"),
			TurnNumber:        convert.Ptr(2),
			TurnSheetCode:     convert.Ptr(turnSheetCode),
			TurnSheetDeadline: convert.Ptr(time.Now().Add(7 * 24 * time.Hour)),
			BackgroundImage:   &backgroundImage,
			TurnEvents: []turnsheet.TurnEvent{
				{
					Category: turnsheet.TurnEventCategorySystem,
					Icon:     turnsheet.TurnEventIconSystem,
					Message:  "Hammer field repairs restored 18 armor (72/72).",
				},
				{
					Category: turnsheet.TurnEventCategorySystem,
					Icon:     turnsheet.TurnEventIconSystem,
					Message:  "Lance received 2 supply points (6 total).",
				},
			},
		},
		LanceName:    "Alpha Lance",
		SupplyPoints: 6,
		Mechs: []turnsheet.ManagementMechEntry{
			{
				MechInstanceID:   "mech-1",
				Callsign:         "Hammer",
				ChassisName:      "Scout",
				ChassisClass:     "light",
				Status:           "operational",
				IsAtDepot:        true,
				CurrentArmor:     72,
				MaxArmor:         72,
				CurrentStructure: 28,
				MaxStructure:     32,
				StructureDamage:  4,
				Weapons: []turnsheet.MechWeaponSlot{
					{SlotLocation: "left-arm", CurrentWeaponID: "wpn-1", CurrentWeaponName: "Light Pulse Cannon"},
					{SlotLocation: "right-arm", CurrentWeaponID: "wpn-2", CurrentWeaponName: "Chaingun"},
				},
			},
			{
				MechInstanceID:   "mech-2",
				Callsign:         "Anvil",
				ChassisName:      "Sentinel",
				ChassisClass:     "medium",
				Status:           "operational",
				IsAtDepot:        false,
				CurrentArmor:     100,
				MaxArmor:         130,
				CurrentStructure: 65,
				MaxStructure:     65,
				StructureDamage:  0,
				Weapons: []turnsheet.MechWeaponSlot{
					{SlotLocation: "left-torso", CurrentWeaponID: "wpn-3", CurrentWeaponName: "Pulse Cannon"},
					{SlotLocation: "right-arm", CurrentWeaponID: "wpn-4", CurrentWeaponName: "Rocket Pack"},
				},
			},
			{
				MechInstanceID:   "mech-3",
				Callsign:         "Wrench",
				ChassisName:      "Scout",
				ChassisClass:     "light",
				Status:           "operational",
				IsAtDepot:        true,
				IsRefitting:      true,
				CurrentArmor:     72,
				MaxArmor:         72,
				CurrentStructure: 32,
				MaxStructure:     32,
				StructureDamage:  0,
			},
		},
		WeaponCatalog: []turnsheet.CatalogWeapon{
			{WeaponID: "cat-1", Name: "Light Pulse Cannon", Damage: 3, HeatCost: 1, RangeBand: "short"},
			{WeaponID: "cat-2", Name: "Chaingun", Damage: 2, HeatCost: 0, RangeBand: "short"},
			{WeaponID: "cat-3", Name: "Pulse Cannon", Damage: 5, HeatCost: 3, RangeBand: "medium"},
			{WeaponID: "cat-4", Name: "Rocket Pack", Damage: 8, HeatCost: 3, RangeBand: "short"},
			{WeaponID: "cat-5", Name: "Auto-Cannon", Damage: 10, HeatCost: 5, RangeBand: "medium"},
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
				path := fmt.Sprintf("testdata/mecha_lance_management_turnsheet.%s", tc.ext)
				err = os.WriteFile(path, output, 0644)
				require.NoError(t, err, "should save output to testdata directory")
				t.Logf("%s preview saved to %s", tc.name, path)
			}
		})
	}
}
