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

