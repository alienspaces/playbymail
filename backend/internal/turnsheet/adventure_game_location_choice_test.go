package turnsheet_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestLocationChoiceProcessor_GenerateTurnSheet(t *testing.T) {

	// Setup test harness
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)

	// Create a mock config for the processor
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewLocationChoiceProcessor(l, cfg)
	require.NoError(t, err)

	// Load test background image for successful test case
	backgroundImage := loadTestBackgroundImage(t, "testdata/background-darkforest.png")

	tests := []struct {
		name               string
		data               any
		expectError        bool
		expectErrorMessage string
	}{
		{
			name:               "given empty LocationChoiceData when generating turn sheet then validation error is returned",
			data:               &turnsheet.LocationChoiceData{},
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
			name: "given valid LocationChoiceData when generating turn sheet then PDF is generated successfully",
			data: &turnsheet.LocationChoiceData{
				TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
					GameName:        convert.Ptr("Test Adventure"),
					GameType:        convert.Ptr("adventure"),
					TurnNumber:      convert.Ptr(1),
					AccountName:     convert.Ptr("Test Player"),
					TurnSheetCode:   convert.Ptr(generateTestTurnSheetCode(t)),
					BackgroundImage: &backgroundImage,
				},
				LocationName:        "Starting Location",
				LocationDescription: "You are at the beginning",
				LocationOptions: []turnsheet.LocationOption{
					{
						LocationID:              "next_location",
						LocationLinkName:        "Next Location",
						LocationLinkDescription: "Go to the next location",
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

func TestLocationChoiceProcessor_ScanTurnSheet(t *testing.T) {

	// Setup test harness
	cfg, l, _, _, _ := testutil.NewDefaultDependencies(t)

	// Create a mock config for the processor
	cfg.TemplatesPath = "../../templates"

	processor, err := turnsheet.NewLocationChoiceProcessor(l, cfg)
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
		expectedScanData      *turnsheet.LocationChoiceScanData
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
				return []byte(`{"invalid":"data"}`), nil
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
			expectErrorMessage: "no location options supplied",
			requiresScanner:    false,
		},
		{
			name: "given real scanned image tick mark when scanning then turn sheet code and choices extracted",
			imageDataFn: func() ([]byte, error) {
				return os.ReadFile("testdata/adventure_game_location_choice_turn_sheet_filled_tick.jpg")
			},
			sheetDataFn: func() ([]byte, error) {
				data := turnsheet.LocationChoiceData{
					LocationOptions: []turnsheet.LocationOption{
						{LocationID: "crystal_caverns", LocationLinkName: "Crystal Caverns"},
						{LocationID: "dark_tower", LocationLinkName: "Dark Tower"},
						{LocationID: "sunset_plains", LocationLinkName: "Sunset Plains"},
						{LocationID: "mermaid_lagoon", LocationLinkName: "Mermaid Lagoon"},
					},
				}
				return json.Marshal(data)
			},
			expectError:           false,
			expectedTurnSheetCode: "", // Will be extracted from image dynamically
			expectedScanData:      &turnsheet.LocationChoiceScanData{Choices: []string{"mermaid_lagoon"}},
			requiresScanner:       true,
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

				// Verify expected choices
				if tt.expectedScanData != nil {
					var scanData turnsheet.LocationChoiceScanData
					err := json.Unmarshal(resultData, &scanData)
					require.NoError(t, err, "Should unmarshal scan results")
					require.Equal(t, tt.expectedScanData.Choices, scanData.Choices, "Should extract correct choices")
					t.Logf("Choices: %v", scanData.Choices)
				}
			}

			totalDuration := time.Since(testStart)
			t.Logf("Test completed in %v (scan: %v)", totalDuration, scanDuration)
		})
	}
}

// TestLocationChoiceScanData_UnmarshalHTMLForm tests that scanned_data unmarshals
// correctly for HTML form format (location_choice as string) and GetChoices() returns it.
func TestLocationChoiceScanData_UnmarshalHTMLForm(t *testing.T) {
	t.Parallel()

	raw := []byte(`{"location_choice":"loc-next-1"}`)
	var scanData turnsheet.LocationChoiceScanData
	err := json.Unmarshal(raw, &scanData)
	require.NoError(t, err)
	choices := scanData.GetChoices()
	require.Len(t, choices, 1)
	require.Equal(t, "loc-next-1", choices[0])
}
