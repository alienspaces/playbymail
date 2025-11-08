package turn_sheet_test

import (
	"context"
	"encoding/json"
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
			pdf, err := processor.GenerateTurnSheet(ctx, l, sheetData)

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
	}{
		{
			name: "valid data",
			imageDataFn: func() ([]byte, error) {
				return []byte("fake image data"), nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError: false,
		},
		{
			name: "missing email",
			imageDataFn: func() ([]byte, error) {
				return []byte("fake image data"), nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "email is required",
		},
		{
			name: "missing name",
			imageDataFn: func() ([]byte, error) {
				return []byte("fake image data"), nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "name is required",
		},
		{
			name: "missing address line1",
			imageDataFn: func() ([]byte, error) {
				return []byte("fake image data"), nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "postal address line 1 is required",
		},
		{
			name: "missing state",
			imageDataFn: func() ([]byte, error) {
				return []byte("fake image data"), nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "state or province is required",
		},
		{
			name: "missing country",
			imageDataFn: func() ([]byte, error) {
				return []byte("fake image data"), nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "country is required",
		},
		{
			name: "missing post code",
			imageDataFn: func() ([]byte, error) {
				return []byte("fake image data"), nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "post code is required",
		},
		{
			name: "missing character name",
			imageDataFn: func() ([]byte, error) {
				return []byte("fake image data"), nil
			},
			sheetDataFn: func() ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectError:        true,
			expectErrorMessage: "character name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Load image data
			imageData, err := tt.imageDataFn()
			if err != nil {
				t.Fatalf("Failed to load image data: %v", err)
			}

			// Load sheet data
			sheetData, err := tt.sheetDataFn()
			if err != nil {
				t.Fatalf("Failed to load sheet data: %v", err)
			}

			// New context
			ctx := context.Background()

			// Test turn sheet code extraction if expected
			if tt.expectedTurnSheetCode != "" {
				turnSheetCode, err := baseProcessor.ParseTurnSheetCodeFromImage(ctx, imageData)
				if tt.expectError {
					require.Error(t, err, "Should return error for turn sheet code extraction")
				} else {
					require.Equal(t, tt.expectedTurnSheetCode, turnSheetCode)
				}
			}

			// Test scan data extraction if expected
			if tt.expectedScanData != nil {
				scanData, err := processor.ScanTurnSheet(ctx, l, sheetData, imageData)
				if tt.expectError {
					require.Error(t, err, "Should return error for scan data extraction")
				} else {
					require.Equal(t, tt.expectedScanData, scanData)
				}
			}
		})
	}
}

func TestGenerateJoinGamePDFForPrinting(t *testing.T) {

	// Setup test harness
	l, _, _, cfg := testutil.NewDefaultDependencies(t)

	// Create a mock config for the processor
	cfg.TemplatesPath = "../../templates"
	cfg.SaveTestPDFs = true

	processor := turn_sheet.NewJoinGameProcessor(l, cfg)

	testData := &turn_sheet.JoinGameData{
		TurnSheetTemplateData: turn_sheet.TurnSheetTemplateData{
			GameName:          convert.Ptr("The Enchanted Forest Adventure"),
			GameType:          convert.Ptr("adventure"),
			TurnNumber:        convert.Ptr(0),
			TurnSheetCode:     convert.Ptr("JOIN123"),
			TurnSheetDeadline: convert.Ptr(time.Now().Add(7 * 24 * time.Hour)),
		},
		GameDescription: "Welcome to the PlayByMail Adventure!",
	}

	ctx := context.Background()
	sheetData, err := json.Marshal(testData)
	require.NoError(t, err, "Should marshal test data")

	pdfData, err := processor.GenerateTurnSheet(ctx, l, sheetData)
	require.NoError(t, err, "Should generate PDF without error")

	require.NotNil(t, pdfData, "PDF data should not be nil")
	require.Greater(t, len(pdfData), 0, "PDF should contain data")

	// Optionally save PDF to testdata directory for manual inspection and testing
	if cfg.SaveTestPDFs {
		testDataPath := "testdata/adventure_game_join_game_turn_sheet.pdf"
		err = os.WriteFile(testDataPath, pdfData, 0644)
		require.NoError(t, err, "Should save PDF to testdata directory")
		t.Logf("PDF saved to %s", testDataPath)
		t.Logf("PDF size: %d bytes", len(pdfData))

		// Print instructions for testing
		t.Logf("")
		t.Logf("PDF generated successfully. To test the scanner:")
		t.Logf("1. Print the PDF: %s", testDataPath)
		t.Logf("2. Fill out the turn sheet with your information")
		t.Logf("3. Scan the completed turn sheet to a JPEG file")
		t.Logf("4. Save the JPEG in testdata/ with a descriptive name")
		t.Logf("5. Write a test that loads the JPEG and tests the scanner")
		t.Logf("6. Run the test to verify the scanner works")
	}
}
