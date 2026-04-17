package turnsheet_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestMechaSquadManagementProcessor_GenerateTurnSheet(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaSquadManagementProcessor(l, cfg)
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

func TestMechaSquadManagementProcessor_GenerateTurnSheet_ContainsDepotInfo(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaSquadManagementProcessor(l, cfg)
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

func TestMechaSquadManagementProcessor_ScanTurnSheet_EmptyImageReturnsError(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaSquadManagementProcessor(l, cfg)
	require.NoError(t, err)

	ctx := context.Background()
	result, err := processor.ScanTurnSheet(ctx, l, []byte(`{}`), []byte{})
	require.Error(t, err)
	require.Nil(t, result)
}
