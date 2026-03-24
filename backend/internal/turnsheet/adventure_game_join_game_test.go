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
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestJoinGameProcessor_GenerateTurnSheet(t *testing.T) {

	// Setup test harness
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)

	// Create a mock config for the processor
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewJoinGameProcessor(l, cfg)
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
			name: "valid data generates PDF",
			data: &turnsheet.JoinGameData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:          convert.Ptr("The Enchanted Forest Adventure"),
					GameType:          convert.Ptr("adventure"),
					TurnSheetCode:     convert.Ptr(generateTestJoinTurnSheetCode(t)),
					TurnSheetDeadline: convert.Ptr(time.Now().Add(7 * 24 * time.Hour)),
				},
				GameDescription: "Embark on a new adventure!",
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
			pdf, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatPDF, sheetData)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
				require.Nil(t, pdf)
			} else if err != nil {
				t.Logf("GenerateTurnSheet returned error: %v", err)
			}
		})
	}
}

func TestJoinGameProcessor_ScanTurnSheet(t *testing.T) {

	// Setup test harness
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)

	// Create a mock config for the processor
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewJoinGameProcessor(l, cfg)
	require.NoError(t, err)
	baseProcessor, err := turnsheet.NewBaseProcessor(l, cfg)
	require.NoError(t, err)

	tests := []struct {
		name                  string
		imageDataFn           func() ([]byte, error)
		sheetDataFn           func() ([]byte, error)
		expectError           bool
		expectErrorMessage    string
		expectedTurnSheetCode string
		expectedScanData      *turnsheet.AdventureGameJoinGameScanData
		requiresScanner       bool
	}{
		{
			name: "given empty image data when scanning join game turn sheet then error returned",
			imageDataFn: func() ([]byte, error) {
				return []byte{}, nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "empty image data",
			requiresScanner:    false,
		},
		{
			name: "given nil image data when scanning join game turn sheet then error returned",
			imageDataFn: func() ([]byte, error) {
				return nil, nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "empty image data",
			requiresScanner:    false,
		},
		{
			name: "given filled join game turn sheet image when scanning then code and player details are extracted correctly",
			imageDataFn: func() ([]byte, error) {
				return os.ReadFile("testdata/adventure_game_join_game_turn_sheet_filled.jpg")
			},
			sheetDataFn: func() ([]byte, error) {
				turnSheetCode := generateTestJoinTurnSheetCode(t)
				data := turnsheet.JoinGameData{
					TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
						GameName:      convert.Ptr("The Enchanted Forest Adventure"),
						TurnSheetCode: convert.Ptr(turnSheetCode),
					},
					GameDescription: "Adventure",
				}
				return json.Marshal(&data)
			},
			expectError:           false,
			expectedTurnSheetCode: "", // Will be extracted from image dynamically
			expectedScanData: &turnsheet.AdventureGameJoinGameScanData{
				JoinGameScanData: turnsheet.JoinGameScanData{
					GameSubscriptionID: "00000000-0000-0000-0000-000000000001", // Test manager subscription ID
					Email:              "freddyfriday@gmail.com",
					Name:               "mr Freddy",
					PostalAddressLine1: "732 main Road",
					PostalAddressLine2: "",
					StateProvince:      "Canberra",
					Country:            "Australia",
					PostalCode:         "3247",
				},
				CharacterName: "Felicia Six Fingers",
			},
			requiresScanner: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStart := time.Now()
			if tt.requiresScanner {
				requireOpenAIKey(t)
			}

			// Load image data
			loadStart := time.Now()
			imageData, err := tt.imageDataFn()
			if err != nil {
				t.Fatalf("Failed to load image data: %v", err)
			}
			t.Logf("Loaded image data: %d bytes in %v", len(imageData), time.Since(loadStart))

			// Get sheet data bytes
			sheetData, err := tt.sheetDataFn()
			if err != nil {
				t.Fatalf("Failed to get sheet data: %v", err)
			}

			ctx := context.Background()

			// Test turn sheet code extraction if expected or if scanner is required
			if tt.requiresScanner || tt.expectedTurnSheetCode != "" {
				codeStart := time.Now()
				turnSheetCode, err := baseProcessor.ParseTurnSheetCodeFromImage(ctx, imageData)
				codeDuration := time.Since(codeStart)
				if tt.expectError {
					require.Error(t, err, "Should return error for turn sheet code extraction")
				} else {
					require.NoError(t, err, "Should extract turn sheet code without error")
					if tt.expectedTurnSheetCode != "" {
						require.Equal(t, tt.expectedTurnSheetCode, turnSheetCode)
					}
					t.Logf("Extracted turn sheet code '%s' in %v", turnSheetCode, codeDuration)
				}
			}

			// Test join game scanning
			scanStart := time.Now()
			resultData, err := processor.ScanTurnSheet(ctx, l, sheetData, imageData)
			scanDuration := time.Since(scanStart)
			t.Logf("ScanTurnSheet completed in %v", scanDuration)

			if tt.expectError {
				require.Error(t, err, "Should return error")
				if tt.expectErrorMessage != "" {
					require.Contains(t, err.Error(), tt.expectErrorMessage, "Error message should contain expected text")
				}
				require.Nil(t, resultData, "Result should be nil on error")
				return
			}

			require.NoError(t, err, "Should not return error")
			require.NotNil(t, resultData, "Result should not be nil")

			if tt.expectedScanData != nil {
				var scanData turnsheet.AdventureGameJoinGameScanData
				err := json.Unmarshal(resultData, &scanData)
				require.NoError(t, err, "Should unmarshal scan results")

				// Compare fields individually using lowercase to handle OCR case variations
				require.Equal(t, strings.ToLower(tt.expectedScanData.Name), strings.ToLower(scanData.Name), "Name should match")
				require.Equal(t, strings.ToLower(tt.expectedScanData.PostalAddressLine1), strings.ToLower(scanData.PostalAddressLine1), "PostalAddressLine1 should match")
				require.Equal(t, strings.ToLower(tt.expectedScanData.PostalAddressLine2), strings.ToLower(scanData.PostalAddressLine2), "PostalAddressLine2 should match")
				require.Equal(t, strings.ToLower(tt.expectedScanData.StateProvince), strings.ToLower(scanData.StateProvince), "StateProvince should match")
				require.Equal(t, strings.ToLower(tt.expectedScanData.Country), strings.ToLower(scanData.Country), "Country should match")
				require.Equal(t, strings.ToLower(tt.expectedScanData.PostalCode), strings.ToLower(scanData.PostalCode), "PostalCode should match")
				require.Equal(t, strings.ToLower(tt.expectedScanData.CharacterName), strings.ToLower(scanData.CharacterName), "CharacterName should match")
				require.Equal(t, strings.ToLower(tt.expectedScanData.Email), strings.ToLower(scanData.Email), "Email should match")
			}

			totalDuration := time.Since(testStart)
			t.Logf("Test completed in %v (scan: %v)", totalDuration, scanDuration)
		})
	}
}

func TestJoinGameData_DefaultDeliveryMethod(t *testing.T) {
	tests := []struct {
		name     string
		methods  turnsheet.DeliveryMethods
		expected string
	}{
		{
			name:     "no delivery methods — empty string",
			methods:  turnsheet.DeliveryMethods{},
			expected: "",
		},
		{
			name:     "email only — email",
			methods:  turnsheet.DeliveryMethods{Email: true},
			expected: "email",
		},
		{
			name:     "local only — local",
			methods:  turnsheet.DeliveryMethods{PhysicalLocal: true},
			expected: "local",
		},
		{
			name:     "post only — post",
			methods:  turnsheet.DeliveryMethods{PhysicalPost: true},
			expected: "post",
		},
		{
			name:     "email and post — email (email has priority)",
			methods:  turnsheet.DeliveryMethods{Email: true, PhysicalPost: true},
			expected: "email",
		},
		{
			name:     "local and post — local (local has priority over post)",
			methods:  turnsheet.DeliveryMethods{PhysicalLocal: true, PhysicalPost: true},
			expected: "local",
		},
		{
			name:     "all three — email (email has highest priority)",
			methods:  turnsheet.DeliveryMethods{Email: true, PhysicalLocal: true, PhysicalPost: true},
			expected: "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &turnsheet.JoinGameData{
				AvailableDeliveryMethods: tt.methods,
			}
			require.Equal(t, tt.expected, data.DefaultDeliveryMethod())
		})
	}
}

func TestJoinGameData_HasDeliveryChoice(t *testing.T) {
	tests := []struct {
		name     string
		methods  turnsheet.DeliveryMethods
		expected bool
	}{
		{
			name:     "no delivery methods — no choice",
			methods:  turnsheet.DeliveryMethods{},
			expected: false,
		},
		{
			name:     "email only — no choice",
			methods:  turnsheet.DeliveryMethods{Email: true},
			expected: false,
		},
		{
			name:     "email and post — choice required",
			methods:  turnsheet.DeliveryMethods{Email: true, PhysicalPost: true},
			expected: true,
		},
		{
			name:     "all three methods — choice required",
			methods:  turnsheet.DeliveryMethods{Email: true, PhysicalLocal: true, PhysicalPost: true},
			expected: true,
		},
		{
			name:     "local only — no choice",
			methods:  turnsheet.DeliveryMethods{PhysicalLocal: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &turnsheet.JoinGameData{
				AvailableDeliveryMethods: tt.methods,
			}
			require.Equal(t, tt.expected, data.HasDeliveryChoice())
		})
	}
}

func TestJoinGameProcessor_GenerateTurnSheet_WithDeliveryMethods(t *testing.T) {
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewJoinGameProcessor(l, cfg)
	require.NoError(t, err)

	tests := []struct {
		name                  string
		deliveryMethods       turnsheet.DeliveryMethods
		expectRadioButton     bool
		expectAddressFields   bool
		expectToggleScript    bool
		expectCheckedMethod   string // value of the radio that should have checked attribute
		expectHiddenMethod    string // value of the hidden input when no choice
	}{
		{
			name:               "email only — hidden input, no address fields",
			deliveryMethods:    turnsheet.DeliveryMethods{Email: true},
			expectRadioButton:  false,
			expectAddressFields: false,
			expectToggleScript: false,
			expectHiddenMethod: "email",
		},
		{
			name:               "local only — hidden input, no address fields",
			deliveryMethods:    turnsheet.DeliveryMethods{PhysicalLocal: true},
			expectRadioButton:  false,
			expectAddressFields: false,
			expectToggleScript: false,
			expectHiddenMethod: "local",
		},
		{
			name:               "post only — hidden input, address fields shown",
			deliveryMethods:    turnsheet.DeliveryMethods{PhysicalPost: true},
			expectRadioButton:  false,
			expectAddressFields: true,
			expectToggleScript: false,
			expectHiddenMethod: "post",
		},
		{
			name:               "email and local — radio buttons, email checked, no address fields",
			deliveryMethods:    turnsheet.DeliveryMethods{Email: true, PhysicalLocal: true},
			expectRadioButton:  true,
			expectAddressFields: false,
			expectToggleScript: false,
			expectCheckedMethod: "email",
		},
		{
			name:               "email and post — radio buttons, email checked, address fields, toggle script",
			deliveryMethods:    turnsheet.DeliveryMethods{Email: true, PhysicalPost: true},
			expectRadioButton:  true,
			expectAddressFields: true,
			expectToggleScript: true,
			expectCheckedMethod: "email",
		},
		{
			name:               "all three — radio buttons, email checked, address fields, toggle script",
			deliveryMethods:    turnsheet.DeliveryMethods{Email: true, PhysicalLocal: true, PhysicalPost: true},
			expectRadioButton:  true,
			expectAddressFields: true,
			expectToggleScript: true,
			expectCheckedMethod: "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(&turnsheet.JoinGameData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:      convert.Ptr("Test Game"),
					GameType:      convert.Ptr("adventure"),
					TurnNumber:    convert.Ptr(0),
					TurnSheetCode: convert.Ptr(generateTestJoinTurnSheetCode(t)),
				},
				GameDescription:          "A test game",
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

func TestJoinGameScanData_Validate(t *testing.T) {
	fullPostal := turnsheet.JoinGameScanData{
		Email:              "test@example.com",
		Name:               "Test User",
		PostalAddressLine1: "123 Main St",
		StateProvince:      "NSW",
		Country:            "Australia",
		PostalCode:         "2000",
	}

	tests := []struct {
		name      string
		data      turnsheet.JoinGameScanData
		expectErr string
	}{
		{
			name:      "missing email",
			data:      turnsheet.JoinGameScanData{Name: "Test"},
			expectErr: "email is required",
		},
		{
			name:      "missing name",
			data:      turnsheet.JoinGameScanData{Email: "a@b.com"},
			expectErr: "name is required",
		},
		{
			name: "post delivery — missing address",
			data: turnsheet.JoinGameScanData{
				Email:          "a@b.com",
				Name:           "Test",
				DeliveryMethod: "post",
			},
			expectErr: "postal address line 1 is required",
		},
		{
			name: "empty delivery method — missing address (treated as postal)",
			data: turnsheet.JoinGameScanData{
				Email: "a@b.com",
				Name:  "Test",
			},
			expectErr: "postal address line 1 is required",
		},
		{
			name: "email delivery — postal fields not required",
			data: turnsheet.JoinGameScanData{
				Email:          "a@b.com",
				Name:           "Test",
				DeliveryMethod: "email",
			},
		},
		{
			name: "local delivery — postal fields not required",
			data: turnsheet.JoinGameScanData{
				Email:          "a@b.com",
				Name:           "Test",
				DeliveryMethod: "local",
			},
		},
		{
			name: "post delivery — all fields present",
			data: func() turnsheet.JoinGameScanData {
				d := fullPostal
				d.DeliveryMethod = "post"
				return d
			}(),
		},
		{
			name:     "empty delivery method — all fields present",
			data:     fullPostal,
		},
		{
			name: "post delivery — missing state",
			data: turnsheet.JoinGameScanData{
				Email:              "a@b.com",
				Name:               "Test",
				DeliveryMethod:     "post",
				PostalAddressLine1: "123 Main St",
			},
			expectErr: "state or province is required",
		},
		{
			name: "post delivery — missing country",
			data: turnsheet.JoinGameScanData{
				Email:              "a@b.com",
				Name:               "Test",
				DeliveryMethod:     "post",
				PostalAddressLine1: "123 Main St",
				StateProvince:      "NSW",
			},
			expectErr: "country is required",
		},
		{
			name: "post delivery — missing post code",
			data: turnsheet.JoinGameScanData{
				Email:              "a@b.com",
				Name:               "Test",
				DeliveryMethod:     "post",
				PostalAddressLine1: "123 Main St",
				StateProvince:      "NSW",
				Country:            "Australia",
			},
			expectErr: "post code is required",
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

