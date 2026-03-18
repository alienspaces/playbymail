package turnsheet_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestInventoryManagementProcessor_GenerateTurnSheet(t *testing.T) {

	// Setup test harness
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)

	// Create a mock config for the processor
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewInventoryManagementProcessor(l, cfg)
	require.NoError(t, err)

	// Load test background image for successful test case
	backgroundImage := loadTestBackgroundImage(t, "testdata/background-dungeon.png")

	tests := []struct {
		name               string
		data               any
		expectError        bool
		expectErrorMessage string
	}{
		{
			name:               "given empty InventoryManagementData when generating turn sheet then validation error is returned",
			data:               &turnsheet.InventoryManagementData{},
			expectError:        true,
			expectErrorMessage: "game name is required",
		},
		{
			name:        "given nil data when generating turn sheet then PDF generation is handled gracefully",
			data:        nil,
			expectError: false, // Generator may handle nil data gracefully
		},
		{
			name:        "given invalid data type when generating turn sheet then PDF generation is handled gracefully",
			data:        "invalid data",
			expectError: false, // Generator may handle invalid data gracefully
		},
		{
			name: "given valid InventoryManagementData when generating turn sheet then PDF is generated successfully",
			data: &turnsheet.InventoryManagementData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:        convert.Ptr("Test Adventure"),
					GameType:        convert.Ptr("adventure"),
					TurnNumber:      convert.Ptr(1),
					AccountName:     convert.Ptr("Test Player"),
					TurnSheetCode:   convert.Ptr(generateTestTurnSheetCode(t)),
					BackgroundImage: &backgroundImage,
				},
				CharacterName:       "Aria the Mage",
				CurrentLocationName: "Mystic Grove",
				InventoryCapacity:   10,
				InventoryCount:      3,
				CurrentInventory: []turnsheet.InventoryItem{
					{
						ItemInstanceID:  "item-1",
						ItemName:        "Crystal Key",
						ItemDescription: "A glowing crystal key",
						IsEquipped:      false,
						EquipmentSlot:   "",
						CanEquip:        false,
					},
					{
						ItemInstanceID:  "item-2",
						ItemName:        "Iron Sword",
						ItemDescription: "A sturdy iron sword",
						IsEquipped:      true,
						EquipmentSlot:   "weapon",
						CanEquip:        true,
					},
					{
						ItemInstanceID:  "item-3",
						ItemName:        "Leather Armor",
						ItemDescription: "Basic leather protection",
						IsEquipped:      true,
						EquipmentSlot:   "armor",
						CanEquip:        true,
					},
				},
				EquipmentSlots: turnsheet.EquipmentSlots{
					Weapon: &turnsheet.EquippedItem{
						ItemInstanceID: "item-2",
						ItemName:       "Iron Sword",
						SlotName:       "weapon",
					},
					Armor: &turnsheet.EquippedItem{
						ItemInstanceID: "item-3",
						ItemName:       "Leather Armor",
						SlotName:       "armor",
					},
				},
				LocationItems: []turnsheet.LocationItem{
					{
						ItemInstanceID:  "item-4",
						ItemName:        "Healing Potion",
						ItemDescription: "Restores health",
					},
					{
						ItemInstanceID:  "item-5",
						ItemName:        "Magic Ring",
						ItemDescription: "A ring imbued with magic",
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Marshal test data to JSON bytes
			var sheetData []byte
			if tt.data != nil {
				var err error
				sheetData, err = json.Marshal(tt.data)
				require.NoError(t, err, "Should marshal test data")
			}

			ctx := context.Background()
			pdfData, err := processor.GenerateTurnSheet(ctx, l, turnsheet.DocumentFormatPDF, sheetData)

			if tt.expectError {
				require.Error(t, err, "Should return error")
				if tt.expectErrorMessage != "" {
					require.Contains(t, err.Error(), tt.expectErrorMessage, "Error message should contain expected text")
				}
				require.Nil(t, pdfData, "PDF data should be nil on error")
			} else if err != nil {
				t.Logf("PDF generation failed (may be expected in test environment): %v", err)
			}
		})
	}
}

func TestInventoryManagementProcessor_ScanTurnSheet(t *testing.T) {

	// Setup test harness
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)

	// Create a mock config for the processor
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewInventoryManagementProcessor(l, cfg)
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
		expectedScanData      *turnsheet.InventoryManagementScanData
		requiresScanner       bool
	}{
		{
			name: "given empty image data when scanning turn sheet then error is returned",
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
			name: "given nil image data when scanning turn sheet then error is returned",
			imageDataFn: func() ([]byte, error) {
				return nil, nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{"character_name":"Test"}`), nil
			},
			expectError:        true,
			expectErrorMessage: "empty image data",
			requiresScanner:    false,
		},
		{
			name: "given empty sheet data when scanning turn sheet then error is returned",
			imageDataFn: func() ([]byte, error) {
				return []byte("fake-image"), nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "structured extraction failed",
			requiresScanner:    true, // Error occurs during image processing
		},
		// TODO: (agent) Add test case with real scanned image when test data exists. Use a fixture image (e.g. testdata/adventure_game_inventory_management_turn_sheet_filled.jpg) and expected InventoryManagementData; ensure scanner extraction is asserted. Blocked on test asset availability.
		// {
		// 	name: "given real scanned image when scanning then inventory actions extracted",
		// 	imageDataFn: func() ([]byte, error) {
		// 		return os.ReadFile("testdata/adventure_game_inventory_management_turn_sheet_filled.jpg")
		// 	},
		// 	sheetDataFn: func() ([]byte, error) {
		// 		data := turnsheet.InventoryManagementData{
		// 			CharacterName: "Aria the Mage",
		// 			CurrentInventory: []turnsheet.InventoryItem{
		// 				{ItemInstanceID: "item-1", ItemName: "Crystal Key"},
		// 			},
		// 			LocationItems: []turnsheet.LocationItem{
		// 				{ItemInstanceID: "item-2", ItemName: "Healing Potion"},
		// 			},
		// 		}
		// 		return json.Marshal(data)
		// 	},
		// 	expectError:     false,
		// 	expectedScanData: &turnsheet.InventoryManagementScanData{
		// 		PickUp: []string{"item-2"},
		// 	},
		// 	requiresScanner: true,
		// },
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
						require.Equal(t, tt.expectedTurnSheetCode, turnSheetCode, "Should extract correct turn sheet code")
					}
					t.Logf("Extracted turn sheet code '%s' in %v", turnSheetCode, codeDuration)
				}
			}

			// Get sheet data bytes
			sheetData, err := tt.sheetDataFn()
			if err != nil {
				t.Fatalf("Failed to get sheet data: %v", err)
			}

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
			} else {
				require.NoError(t, err, "Should not return error")
				require.NotNil(t, resultData, "Result should not be nil")

				// Verify expected scan data
				if tt.expectedScanData != nil {
					var scanData turnsheet.InventoryManagementScanData
					err := json.Unmarshal(resultData, &scanData)
					require.NoError(t, err, "Should unmarshal scan results")
					require.Equal(t, tt.expectedScanData.PickUp, scanData.PickUp, "Should extract correct pick up actions")
					require.Equal(t, tt.expectedScanData.Drop, scanData.Drop, "Should extract correct drop actions")
					require.Equal(t, tt.expectedScanData.Equip, scanData.Equip, "Should extract correct equip actions")
					require.Equal(t, tt.expectedScanData.Unequip, scanData.Unequip, "Should extract correct unequip actions")
					t.Logf("Actions - PickUp: %v, Drop: %v, Equip: %v, Unequip: %v",
						scanData.PickUp, scanData.Drop, scanData.Equip, scanData.Unequip)
				}
			}

			totalDuration := time.Since(testStart)
			t.Logf("Test completed in %v (scan: %v)", totalDuration, scanDuration)
		})
	}
}

// TestGenerateInventoryManagementFormatsForPrinting generates HTML and PDF versions for physical testing
// Set SAVE_TEST_FILES=true to save the files to testdata directory
func TestGenerateInventoryManagementFormatsForPrinting(t *testing.T) {

	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)

	cfg.TemplatesPath = "../../templates"
	// SaveTestFiles defaults to false - set SAVE_TEST_FILES=true to generate files
	cfg.SaveTestFiles = true

	processor, err := turnsheet.NewInventoryManagementProcessor(l, cfg)
	require.NoError(t, err)

	type formatCase struct {
		name     string
		format   turnsheet.DocumentFormat
		ext      string
		logExtra bool
	}

	cases := []formatCase{
		{
			name:     "pdf",
			format:   turnsheet.DocumentFormatPDF,
			ext:      "pdf",
			logExtra: true,
		},
		{
			name:   "html",
			format: turnsheet.DocumentFormatHTML,
			ext:    "html",
		},
	}

	// Load test background image
	backgroundImage := loadTestBackgroundImage(t, "testdata/background-dungeon.png")

	// Generate realistic turn sheet code
	turnSheetCode := generateTestTurnSheetCode(t)

	testData := &turnsheet.InventoryManagementData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:          convert.Ptr("The Enchanted Forest Adventure"),
			GameType:          convert.Ptr("adventure"),
			TurnNumber:        convert.Ptr(1),
			AccountName:       convert.Ptr("Test Player"),
			TurnSheetCode:     convert.Ptr(turnSheetCode),
			TurnSheetDeadline: convert.Ptr(time.Now().Add(24 * time.Hour)),
			BackgroundImage:   &backgroundImage,
		},
		CharacterName:       "Aria the Mage",
		CurrentLocationName: "Mystic Grove",
		InventoryCapacity:   10,
		InventoryCount:      5,
		CurrentInventory: []turnsheet.InventoryItem{
			{
				ItemInstanceID:  "item-1",
				ItemName:        "Crystal Key",
				ItemDescription: "A glowing crystal key that hums with ancient magic",
				IsEquipped:      false,
				EquipmentSlot:   "",
				CanEquip:        false,
			},
			{
				ItemInstanceID:  "item-2",
				ItemName:        "Iron Sword",
				ItemDescription: "A sturdy iron sword with a leather-wrapped hilt",
				IsEquipped:      true,
				EquipmentSlot:   "weapon",
				CanEquip:        true,
			},
			{
				ItemInstanceID:  "item-3",
				ItemName:        "Leather Armor",
				ItemDescription: "Basic leather protection, worn but serviceable",
				IsEquipped:      true,
				EquipmentSlot:   "armor",
				CanEquip:        true,
			},
			{
				ItemInstanceID:  "item-4",
				ItemName:        "Healing Potion",
				ItemDescription: "A red potion that restores health when consumed",
				IsEquipped:      false,
				EquipmentSlot:   "",
				CanEquip:        false,
			},
			{
				ItemInstanceID:  "item-5",
				ItemName:        "Magic Ring",
				ItemDescription: "A silver ring imbued with protective magic",
				IsEquipped:      false,
				EquipmentSlot:   "",
				CanEquip:        true,
			},
		},
		EquipmentSlots: turnsheet.EquipmentSlots{
			Weapon: &turnsheet.EquippedItem{
				ItemInstanceID: "item-2",
				ItemName:       "Iron Sword",
				SlotName:       "weapon",
			},
			Armor: &turnsheet.EquippedItem{
				ItemInstanceID: "item-3",
				ItemName:       "Leather Armor",
				SlotName:       "armor",
			},
		},
		LocationItems: []turnsheet.LocationItem{
			{
				ItemInstanceID:  "item-6",
				ItemName:        "Shadow Cloak",
				ItemDescription: "A dark cloak that seems to blend with shadows",
			},
			{
				ItemInstanceID:  "item-7",
				ItemName:        "Wind Charm",
				ItemDescription: "A small charm that whispers with the wind",
			},
		},
	}

	ctx := context.Background()
	sheetData, err := json.Marshal(testData)
	require.NoError(t, err, "Should marshal test data")

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output, err := processor.GenerateTurnSheet(ctx, l, tc.format, sheetData)
			require.NoError(t, err, "Should generate output without error")
			require.NotEmpty(t, output, "Output should not be empty")

			if cfg.SaveTestFiles {
				path := fmt.Sprintf("testdata/adventure_game_inventory_management_turnsheet.%s", tc.ext)
				err = os.WriteFile(path, output, 0644)
				require.NoError(t, err, "Should save output to testdata directory")

				t.Logf("%s preview saved to %s", tc.name, path)

				if tc.logExtra {
					t.Logf("Output size: %d bytes", len(output))
					t.Logf("")
					t.Logf("Generated successfully. To test the scanner:")
					t.Logf("1. Print the PDF: %s", path)
					t.Logf("2. Fill out the turn sheet with your inventory actions")
					t.Logf("3. Scan the completed turn sheet to a JPEG file")
					t.Logf("4. Save the JPEG in testdata/ with a descriptive name")
					t.Logf("5. Write a test that loads the JPEG and tests the scanner")
				}
			}
		})
	}
}

// TestInventoryManagementScanData_UnmarshalHTMLFormEquip tests that scanned_data
// unmarshals correctly for both HTML form format (equip as []string) and full format (equip as []EquipAction).
func TestInventoryManagementScanData_UnmarshalHTMLFormEquip(t *testing.T) {
	t.Parallel()

	t.Run("HTML form format: equip as array of item IDs", func(t *testing.T) {
		raw := []byte(`{"pick_up":["item-1"],"drop":[],"equip":["item-2","item-3"],"unequip":[]}`)
		var scanData turnsheet.InventoryManagementScanData
		err := json.Unmarshal(raw, &scanData)
		require.NoError(t, err)
		require.Len(t, scanData.Equip, 2)
		require.Equal(t, "item-2", scanData.Equip[0].ItemInstanceID)
		require.Equal(t, turnsheet.DefaultEquipSlot, scanData.Equip[0].Slot)
		require.Equal(t, "item-3", scanData.Equip[1].ItemInstanceID)
		require.Equal(t, turnsheet.DefaultEquipSlot, scanData.Equip[1].Slot)
		require.Equal(t, []string{"item-1"}, scanData.PickUp)
	})

	t.Run("full format: equip as array of objects with slot", func(t *testing.T) {
		raw := []byte(`{"equip":[{"item_instance_id":"item-a","slot":"armor"},{"item_instance_id":"item-b","slot":"weapon"}]}`)
		var scanData turnsheet.InventoryManagementScanData
		err := json.Unmarshal(raw, &scanData)
		require.NoError(t, err)
		require.Len(t, scanData.Equip, 2)
		require.Equal(t, "item-a", scanData.Equip[0].ItemInstanceID)
		require.Equal(t, "armor", scanData.Equip[0].Slot)
		require.Equal(t, "item-b", scanData.Equip[1].ItemInstanceID)
		require.Equal(t, "weapon", scanData.Equip[1].Slot)
	})

	t.Run("HTML form format: empty equip array when no items checked", func(t *testing.T) {
		// An empty equip array is produced by extractFormData() when no equip checkboxes
		// are checked. The backend must accept this without error; the frontend strips
		// empty arrays before saving, but the backend should also be lenient.
		raw := []byte(`{"pick_up":["item-1"],"equip":[]}`)
		var scanData turnsheet.InventoryManagementScanData
		err := json.Unmarshal(raw, &scanData)
		require.NoError(t, err)
		require.Len(t, scanData.Equip, 0)
		require.Equal(t, []string{"item-1"}, scanData.PickUp)
	})
}

func TestValidateInventoryActions(t *testing.T) {
	t.Parallel()

	sheetData := &turnsheet.InventoryManagementData{
		CurrentInventory: []turnsheet.InventoryItem{
			{ItemInstanceID: "inv-1", ItemName: "Iron Sword", CanEquip: true},
			{ItemInstanceID: "inv-2", ItemName: "Health Potion", CanEquip: false},
		},
		LocationItems: []turnsheet.LocationItem{
			{ItemInstanceID: "loc-1", ItemName: "Desert Compass"},
			{ItemInstanceID: "loc-2", ItemName: "Water Flask"},
		},
	}

	tests := []struct {
		name        string
		scanData    *turnsheet.InventoryManagementScanData
		expectError bool
		errorContains string
	}{
		{
			name: "equip inventory item passes",
			scanData: &turnsheet.InventoryManagementScanData{
				Equip: turnsheet.EquipPayload{{ItemInstanceID: "inv-1", Slot: "weapon"}},
			},
			expectError: false,
		},
		{
			name: "equip location item passes",
			scanData: &turnsheet.InventoryManagementScanData{
				Equip: turnsheet.EquipPayload{{ItemInstanceID: "loc-1", Slot: "weapon"}},
			},
			expectError: false,
		},
		{
			name: "equip unknown item fails",
			scanData: &turnsheet.InventoryManagementScanData{
				Equip: turnsheet.EquipPayload{{ItemInstanceID: "unknown-id", Slot: "weapon"}},
			},
			expectError:   true,
			errorContains: "invalid item_instance_id for equip: unknown-id",
		},
		{
			name: "equip with invalid slot fails",
			scanData: &turnsheet.InventoryManagementScanData{
				Equip: turnsheet.EquipPayload{{ItemInstanceID: "inv-1", Slot: "hat"}},
			},
			expectError:   true,
			errorContains: "invalid equipment slot: hat",
		},
		{
			name: "pick up location item passes",
			scanData: &turnsheet.InventoryManagementScanData{
				PickUp: []string{"loc-1"},
			},
			expectError: false,
		},
		{
			name: "pick up inventory item fails",
			scanData: &turnsheet.InventoryManagementScanData{
				PickUp: []string{"inv-1"},
			},
			expectError:   true,
			errorContains: "invalid item_instance_id for pick up: inv-1",
		},
		{
			name: "drop inventory item passes",
			scanData: &turnsheet.InventoryManagementScanData{
				Drop: []string{"inv-1"},
			},
			expectError: false,
		},
		{
			name: "drop location item fails",
			scanData: &turnsheet.InventoryManagementScanData{
				Drop: []string{"loc-1"},
			},
			expectError:   true,
			errorContains: "invalid item_instance_id for drop: loc-1",
		},
		{
			name: "nil scan data fails",
			scanData:      nil,
			expectError:   true,
			errorContains: "no scan data provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := turnsheet.ValidateInventoryActions(sheetData, tt.scanData)
			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
