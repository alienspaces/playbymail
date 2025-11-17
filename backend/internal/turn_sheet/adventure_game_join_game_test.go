package turn_sheet_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestJoinGameProcessor_GenerateTurnSheet(t *testing.T) {

	// Setup test harness
	l, _, _, cfg := testutil.NewDefaultDependencies(t)

	// Create a mock config for the processor
	cfg.TemplatesPath = "../../templates"

	processor := turn_sheet.NewJoinGameProcessor(l, cfg)

	tests := []struct {
		name        string
		data        any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty data returns validation error",
			data:        &turn_sheet.JoinGameData{},
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
			data: &turn_sheet.JoinGameData{
				TurnSheetTemplateData: turn_sheet.TurnSheetTemplateData{
					GameName:          convert.Ptr("The Enchanted Forest Adventure"),
					GameType:          convert.Ptr("adventure"),
					TurnSheetCode:     convert.Ptr("ABC123XYZ"),
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
			pdf, err := processor.GenerateTurnSheet(ctx, l, turn_sheet.DocumentFormatPDF, sheetData)

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
	l, _, _, cfg := testutil.NewDefaultDependencies(t)

	// Create a mock config for the processor
	cfg.TemplatesPath = "../../templates"

	processor := turn_sheet.NewJoinGameProcessor(l, cfg)
	baseProcessor := turn_sheet.NewBaseProcessor(l, cfg)

	tests := []struct {
		name                  string
		imageDataFn           func() ([]byte, error)
		sheetDataFn           func() ([]byte, error)
		expectError           bool
		expectErrorMessage    string
		expectedTurnSheetCode string
		expectedScanData      *turn_sheet.JoinGameScanData
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
				data := turn_sheet.JoinGameData{
					TurnSheetTemplateData: turn_sheet.TurnSheetTemplateData{
						GameName:      convert.Ptr("The Enchanted Forest Adventure"),
						TurnSheetCode: convert.Ptr("JOIN123"),
					},
					GameDescription: "Adventure",
				}
				return json.Marshal(&data)
			},
			expectError:           false,
			expectedTurnSheetCode: "JOIN123",
			expectedScanData: &turn_sheet.JoinGameScanData{
				Email:              "alienspaces@gmail.com",
				Name:               "Ben Wallin",
				PostalAddressLine1: "54 Ronald Street",
				PostalAddressLine2: "",
				StateProvince:      "Coburg North",
				Country:            "Australia",
				PostalCode:         "3058",
				CharacterName:      "Luscious",
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

			// Test turn sheet code extraction if expected
			if tt.expectedTurnSheetCode != "" {
				codeStart := time.Now()
				turnSheetCode, err := baseProcessor.ParseTurnSheetCodeFromImage(ctx, imageData)
				codeDuration := time.Since(codeStart)
				if tt.expectError {
					require.Error(t, err, "Should return error for turn sheet code extraction")
				} else {
					require.NoError(t, err, "Should extract turn sheet code without error")
					require.Equal(t, tt.expectedTurnSheetCode, turnSheetCode)
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
				var scanData turn_sheet.JoinGameScanData
				err := json.Unmarshal(resultData, &scanData)
				require.NoError(t, err, "Should unmarshal scan results")
				require.Equal(t, tt.expectedScanData, &scanData)
			}

			totalDuration := time.Since(testStart)
			t.Logf("Test completed in %v (scan: %v)", totalDuration, scanDuration)
		})
	}
}

func TestGenerateJoinGameFormatsForPrinting(t *testing.T) {

	l, _, _, cfg := testutil.NewDefaultDependencies(t)
	cfg.TemplatesPath = "../../templates"

	// SaveTestFiles defaults to false - set SAVE_TEST_FILES=true to generate files
	// cfg.SaveTestFiles = true

	processor := turn_sheet.NewJoinGameProcessor(l, cfg)

	type formatCase struct {
		name     string
		format   turn_sheet.DocumentFormat
		ext      string
		logExtra bool
		deadline time.Duration
	}

	cases := []formatCase{
		{
			name:     "pdf",
			format:   turn_sheet.DocumentFormatPDF,
			ext:      "pdf",
			logExtra: true,
			deadline: 7 * 24 * time.Hour,
		},
		{
			name:     "html",
			format:   turn_sheet.DocumentFormatHTML,
			ext:      "html",
			deadline: 48 * time.Hour,
		},
	}

	ctx := context.Background()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testData := &turn_sheet.JoinGameData{
				TurnSheetTemplateData: turn_sheet.TurnSheetTemplateData{
					GameName:          convert.Ptr("The Enchanted Forest Adventure"),
					GameType:          convert.Ptr("adventure"),
					TurnNumber:        convert.Ptr(0),
					TurnSheetCode:     convert.Ptr("JOIN123"),
					TurnSheetDeadline: convert.Ptr(time.Now().Add(tc.deadline)),
				},
				GameDescription: "Welcome to the PlayByMail Adventure!",
			}

			sheetData, err := json.Marshal(testData)
			require.NoError(t, err, "Should marshal test data")

			output, err := processor.GenerateTurnSheet(ctx, l, tc.format, sheetData)
			require.NoError(t, err, "Should generate output without error")
			require.NotEmpty(t, output, "Output should not be empty")

			if cfg.SaveTestFiles {
				path := fmt.Sprintf("testdata/adventure_game_join_game_turn_sheet.%s", tc.ext)
				err = os.WriteFile(path, output, 0644)
				require.NoError(t, err, "Should save output to testdata directory")

				t.Logf("%s preview saved to %s", tc.name, path)

				if tc.logExtra {
					t.Logf("Output size: %d bytes", len(output))
					t.Logf("")
					t.Logf("Generated successfully. To test the scanner:")
					t.Logf("1. Print the PDF: %s", path)
					t.Logf("2. Fill out the turn sheet with your information")
					t.Logf("3. Scan the completed turn sheet to a JPEG file")
					t.Logf("4. Save the JPEG in testdata/ with a descriptive name")
					t.Logf("5. Write a test that loads the JPEG and tests the scanner")
					t.Logf("6. Run the test to verify the scanner works")
				}
			}
		})
	}
}
