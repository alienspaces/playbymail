package turnsheet_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestMechaGameJoinGameProcessor_GenerateTurnSheet(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameJoinGameProcessor(l, cfg)
	require.NoError(t, err)

	tests := []struct {
		name        string
		data        any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty data returns validation error",
			data:        &turnsheet.JoinGameData{},
			expectError: true,
			errorMsg:    "game name is required",
		},
		{
			name:        "nil data handled gracefully",
			data:        nil,
			expectError: false,
		},
		{
			name:        "invalid data type handled gracefully",
			data:        "invalid",
			expectError: false,
		},
		{
			name: "valid data generates HTML",
			data: &turnsheet.JoinGameData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:      convert.Ptr("Steel Thunder"),
					GameType:      convert.Ptr("mecha"),
					TurnSheetCode: convert.Ptr(generateTestJoinTurnSheetCode(t)),
					TurnNumber:    convert.Ptr(0),
				},
				GameDescription: "Command a squad of war mechs!",
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

func TestMechaGameJoinGameProcessor_GenerateTurnSheet_HTMLContainsInputElements(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameJoinGameProcessor(l, cfg)
	require.NoError(t, err)

	data, err := json.Marshal(&turnsheet.JoinGameData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:      convert.Ptr("Steel Thunder"),
			GameType:      convert.Ptr("mecha"),
			TurnNumber:    convert.Ptr(0),
			TurnSheetCode: convert.Ptr(generateTestJoinTurnSheetCode(t)),
		},
		GameDescription:          "Command a squad of war mechs!",
		AvailableDeliveryMethods: turnsheet.DeliveryMethods{Email: true},
	})
	require.NoError(t, err)

	ctx := context.Background()
	html, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatHTML, data)
	require.NoError(t, err)
	require.NotEmpty(t, html)

	htmlStr := string(html)

	require.True(t, strings.Contains(htmlStr, `id="email"`),
		"should contain email input element")
	require.True(t, strings.Contains(htmlStr, `id="name"`),
		"should contain name input element")
	require.True(t, strings.Contains(htmlStr, `id="commander_name"`),
		"should contain commander_name input element")
	require.True(t, strings.Contains(htmlStr, `type="email"`),
		"email field should use type=email")
	require.True(t, strings.Contains(htmlStr, "Your Mecha Commander"),
		"should contain mecha commander section title")
	require.False(t, strings.Contains(htmlStr, "form-input-line"),
		"should not contain non-interactive form-input-line divs")
}

func TestMechaGameJoinGameProcessor_GenerateTurnSheet_PreFillsAccountEmail(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameJoinGameProcessor(l, cfg)
	require.NoError(t, err)

	data, err := json.Marshal(&turnsheet.JoinGameData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:      convert.Ptr("Steel Thunder"),
			GameType:      convert.Ptr("mecha"),
			TurnNumber:    convert.Ptr(0),
			TurnSheetCode: convert.Ptr(generateTestJoinTurnSheetCode(t)),
		},
		GameDescription:          "Command a squad of war mechs!",
		AvailableDeliveryMethods: turnsheet.DeliveryMethods{Email: true},
		AccountEmail:             "commander@example.com",
	})
	require.NoError(t, err)

	ctx := context.Background()
	html, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatHTML, data)
	require.NoError(t, err)

	htmlStr := string(html)
	require.True(t, strings.Contains(htmlStr, `value="commander@example.com"`),
		"should pre-fill the email field with the account email")
	require.True(t, strings.Contains(htmlStr, "readonly"),
		"pre-filled email should be readonly")
}

func TestMechaGameJoinGameProcessor_GenerateTurnSheet_WithDeliveryMethods(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewMechaGameJoinGameProcessor(l, cfg)
	require.NoError(t, err)

	tests := []struct {
		name                string
		deliveryMethods     turnsheet.DeliveryMethods
		expectRadioButton   bool
		expectAddressFields bool
		expectToggleScript  bool
		expectCheckedMethod string
		expectHiddenMethod  string
	}{
		{
			name:                "email only — hidden input, no address fields",
			deliveryMethods:     turnsheet.DeliveryMethods{Email: true},
			expectRadioButton:   false,
			expectAddressFields: false,
			expectToggleScript:  false,
			expectHiddenMethod:  "email",
		},
		{
			name:                "local only — hidden input, no address fields",
			deliveryMethods:     turnsheet.DeliveryMethods{PhysicalLocal: true},
			expectRadioButton:   false,
			expectAddressFields: false,
			expectToggleScript:  false,
			expectHiddenMethod:  "local",
		},
		{
			name:                "post only — hidden input, address fields shown",
			deliveryMethods:     turnsheet.DeliveryMethods{PhysicalPost: true},
			expectRadioButton:   false,
			expectAddressFields: true,
			expectToggleScript:  false,
			expectHiddenMethod:  "post",
		},
		{
			name:                "email and local — radio buttons, email checked, no address fields",
			deliveryMethods:     turnsheet.DeliveryMethods{Email: true, PhysicalLocal: true},
			expectRadioButton:   true,
			expectAddressFields: false,
			expectToggleScript:  false,
			expectCheckedMethod: "email",
		},
		{
			name:                "email and post — radio buttons, email checked, address fields, toggle script",
			deliveryMethods:     turnsheet.DeliveryMethods{Email: true, PhysicalPost: true},
			expectRadioButton:   true,
			expectAddressFields: true,
			expectToggleScript:  true,
			expectCheckedMethod: "email",
		},
		{
			name:                "all three — radio buttons, email checked, address fields, toggle script",
			deliveryMethods:     turnsheet.DeliveryMethods{Email: true, PhysicalLocal: true, PhysicalPost: true},
			expectRadioButton:   true,
			expectAddressFields: true,
			expectToggleScript:  true,
			expectCheckedMethod: "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&turnsheet.JoinGameData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:      convert.Ptr("Steel Thunder"),
					GameType:      convert.Ptr("mecha"),
					TurnNumber:    convert.Ptr(0),
					TurnSheetCode: convert.Ptr(generateTestJoinTurnSheetCode(t)),
				},
				GameDescription:          "Command a squad of war mechs!",
				AvailableDeliveryMethods: tt.deliveryMethods,
			})
			require.NoError(t, err)

			ctx := context.Background()
			html, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatHTML, data)
			require.NoError(t, err)
			require.NotEmpty(t, html)

			htmlStr := string(html)

			hasRadio := strings.Contains(htmlStr, `type="radio" name="delivery_method"`)
			require.Equal(t, tt.expectRadioButton, hasRadio,
				"expected radio button presence to be %v", tt.expectRadioButton)

			if tt.expectCheckedMethod != "" {
				checkedRadio := fmt.Sprintf("value=%q checked", tt.expectCheckedMethod)
				require.True(t, strings.Contains(htmlStr, checkedRadio),
					"expected radio value=%q to have checked attribute", tt.expectCheckedMethod)
			}

			if tt.expectHiddenMethod != "" {
				hiddenInput := fmt.Sprintf(`type="hidden" name="delivery_method" value=%q`, tt.expectHiddenMethod)
				require.True(t, strings.Contains(htmlStr, hiddenInput),
					"expected hidden input with value=%q", tt.expectHiddenMethod)
			}

			hasAddress := strings.Contains(htmlStr, `id="postal-address-fields"`)
			require.Equal(t, tt.expectAddressFields, hasAddress,
				"expected address fields presence to be %v", tt.expectAddressFields)

			hasScript := strings.Contains(htmlStr, `postal-address-fields`) && strings.Contains(htmlStr, `toggle(this.value === 'post')`)
			require.Equal(t, tt.expectToggleScript, hasScript,
				"expected toggle script presence to be %v", tt.expectToggleScript)
		})
	}
}

func TestMechaGameJoinGameScanData_Validate(t *testing.T) {
	tests := []struct {
		name      string
		data      turnsheet.MechaGameJoinGameScanData
		expectErr string
	}{
		{
			name: "missing email",
			data: turnsheet.MechaGameJoinGameScanData{
				JoinGameScanData: turnsheet.JoinGameScanData{Name: "Test", DeliveryMethod: "email"},
				CommanderName:    "Commander X",
			},
			expectErr: "email is required",
		},
		{
			name: "missing name",
			data: turnsheet.MechaGameJoinGameScanData{
				JoinGameScanData: turnsheet.JoinGameScanData{Email: "a@b.com", DeliveryMethod: "email"},
				CommanderName:    "Commander X",
			},
			expectErr: "name is required",
		},
		{
			name: "missing commander name",
			data: turnsheet.MechaGameJoinGameScanData{
				JoinGameScanData: turnsheet.JoinGameScanData{Email: "a@b.com", Name: "Test", DeliveryMethod: "email"},
			},
			expectErr: "commander name is required",
		},
		{
			name: "valid email delivery",
			data: turnsheet.MechaGameJoinGameScanData{
				JoinGameScanData: turnsheet.JoinGameScanData{Email: "a@b.com", Name: "Test", DeliveryMethod: "email"},
				CommanderName:    "Commander X",
			},
		},
		{
			name: "valid post delivery",
			data: turnsheet.MechaGameJoinGameScanData{
				JoinGameScanData: turnsheet.JoinGameScanData{
					Email:              "a@b.com",
					Name:               "Test",
					DeliveryMethod:     "post",
					PostalAddressLine1: "123 Main St",
					StateProvince:      "NSW",
					Country:            "Australia",
					PostalCode:         "2000",
				},
				CommanderName: "Commander X",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if tt.expectErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
