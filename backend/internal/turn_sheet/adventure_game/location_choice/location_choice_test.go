package location_choice_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/generator"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/scanner"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet/adventure_game/location_choice"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

// TestGenerateLocationChoiceHTML tests HTML generation for location choice
// turn sheets
func TestGenerateLocationChoiceHTML(t *testing.T) {
	tests := []struct {
		name        string
		data        *location_choice.LocationChoiceData
		expectError bool
		validate    func(t *testing.T, html string)
	}{
		{
			name: "generate with single location option",
			data: &location_choice.LocationChoiceData{
				LocationName:        "Crystal Caverns",
				LocationDescription: "A vast network of glittering caves.",
				LocationOptions: []location_choice.LocationOption{
					{
						LocationID:              "loc-2",
						LocationLinkName:        "Floating Islands",
						LocationLinkDescription: "Sky islands connected by bridges.",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, html string) {
				require.Contains(t, html, "Crystal Caverns")
				require.Contains(t, html, "Floating Islands")
				require.Contains(t, html, "glittering caves")
			},
		},
		{
			name: "generate with multiple location options",
			data: &location_choice.LocationChoiceData{
				LocationName:        "Mystic Grove",
				LocationDescription: "An ancient forest filled with magic.",
				LocationOptions: []location_choice.LocationOption{
					{
						LocationID:              "loc-2",
						LocationLinkName:        "Crystal Caverns",
						LocationLinkDescription: "Glittering underground caves.",
					},
					{
						LocationID:              "loc-3",
						LocationLinkName:        "Floating Islands",
						LocationLinkDescription: "Sky islands connected by bridges.",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, html string) {
				require.Contains(t, html, "Mystic Grove")
				require.Contains(t, html, "Crystal Caverns")
				require.Contains(t, html, "Floating Islands")
			},
		},
		{
			name: "generate with no location options",
			data: &location_choice.LocationChoiceData{
				LocationName:        "Dead End",
				LocationDescription: "A place with nowhere to go.",
				LocationOptions:     []location_choice.LocationOption{},
			},
			expectError: false,
			validate: func(t *testing.T, html string) {
				require.Contains(t, html, "Dead End")
				require.Contains(t, html, "nowhere to go")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// setup test harness
			h := deps.NewHarness(t)
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()

			// setup test directories
			templateDir := getTemplateDir(t)
			outputDir := t.TempDir()

			// create generator
			gen := generator.NewPDFGenerator(h.Log.(*log.Log), templateDir, outputDir)

			// prepare template data
			templateData := createTestTemplateData(t, tt.data)

			// generate HTML
			html, err := gen.GenerateHTML(ctx, "adventure_game/location_choice/template/content.template", templateData)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, html)

			if tt.validate != nil {
				tt.validate(t, html)
			}
		})
	}
}

// TestGenerateLocationChoicePDF tests PDF generation for location choice
// turn sheets
func TestGenerateLocationChoicePDF(t *testing.T) {
	// skip if running in CI without Chrome
	if os.Getenv("CI") == "true" && os.Getenv("GOOGLE_CHROME_SHIM") == "" {
		t.Skip("skipping PDF generation test in CI without Chrome")
	}

	tests := []struct {
		name        string
		data        *location_choice.LocationChoiceData
		expectError bool
		validate    func(t *testing.T, pdfData []byte)
	}{
		{
			name: "generate PDF with location options",
			data: &location_choice.LocationChoiceData{
				LocationName:        "Mystic Grove",
				LocationDescription: "A peaceful forest clearing with ancient trees and magical energy.",
				LocationOptions: []location_choice.LocationOption{
					{
						LocationID:              "crystal_caverns",
						LocationLinkName:        "Crystal Caverns",
						LocationLinkDescription: "Deep underground caves filled with glowing crystals.",
					},
					{
						LocationID:              "floating_islands",
						LocationLinkName:        "Floating Islands",
						LocationLinkDescription: "Mysterious islands suspended in the sky above the clouds.",
					},
					{
						LocationID:              "shadow_valley",
						LocationLinkName:        "Shadow Valley",
						LocationLinkDescription: "A dark valley shrouded in perpetual twilight.",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, pdfData []byte) {
				require.NotEmpty(t, pdfData)
				// basic PDF validation
				require.True(t, len(pdfData) > 100, "PDF should be larger than 100 bytes")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// enable test mode for mock PDF generation
			t.Setenv("TESTING", "true")

			ctx := context.Background()

			// setup test directories
			templateDir := getTemplateDir(t)
			outputDir := t.TempDir()

			// create generator with simple logger (no harness)
			cfg := config.Config{}
			logger, err := log.NewLogger(cfg)
			require.NoError(t, err)
			gen := generator.NewPDFGenerator(logger, templateDir, outputDir)

			// prepare template data
			templateData := createTestTemplateData(t, tt.data)

			// generate PDF
			pdfData, err := gen.GeneratePDF(ctx, "adventure_game/location_choice/template/content.template", templateData)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Save PDF for physical testing when flag is enabled
			if os.Getenv("SAVE_PDF_FOR_TESTING") == "true" {
				testDataDir := filepath.Join("testdata")
				err = os.MkdirAll(testDataDir, 0755)
				require.NoError(t, err)

				pdfPath := filepath.Join(testDataDir, "location_choice_test_sheet.pdf")
				err = os.WriteFile(pdfPath, pdfData, 0644)
				require.NoError(t, err)

				t.Logf("Saved PDF for physical testing: %s", pdfPath)
				t.Logf("PDF size: %d bytes", len(pdfData))
				t.Logf("Instructions:")
				t.Logf("1. Print the PDF")
				t.Logf("2. Fill out the turn sheet by hand")
				t.Logf("3. Take a photo of the completed sheet")
				t.Logf("4. Save the photo as 'location_choice_filled.jpg' in testdata/")
				t.Logf("5. Use the photo for OCR testing")
			}

			if tt.validate != nil {
				tt.validate(t, pdfData)
			}
		})
	}
}

// TestScanLocationChoice tests scanning completed location choice turn sheets
func TestLocationChoiceScan(t *testing.T) {
	tests := []struct {
		name         string
		imageFile    string
		expectedData *location_choice.LocationChoiceScanData
		expectError  bool
	}{
		{
			name:      "scan not implemented yet",
			imageFile: "valid_location_choice.png",
			expectedData: &location_choice.LocationChoiceScanData{
				Choices: []string{"loc-1"},
			},
			expectError: true, // OCR not implemented yet
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// setup test harness
			h := deps.NewHarness(t)
			_, err := h.Setup()
			require.NoError(t, err)
			defer func() {
				err = h.Teardown()
				require.NoError(t, err)
			}()

			// setup scanner
			s := scanner.NewScanner(h.Logger("TestScanLocationChoice"), h.Domain.(*domain.Domain))

			// load test image
			imageData := loadTestImage(t, tt.imageFile)

			// extract text from image
			_, err = s.ExtractTextFromImage(ctx, imageData)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// TODO: When OCR is implemented, add validation:
			// 1. Parse scanned text into LocationChoiceScanData
			// 2. Validate against expected data
			// 3. Verify turn sheet code extraction
		})
	}
}

// TestLocationChoiceIntegration tests the complete workflow of generating
// and scanning location choice turn sheets
func TestLocationChoiceIntegration(t *testing.T) {
	h := deps.NewHarness(t)
	_, err := h.Setup()
	require.NoError(t, err)
	defer func() {
		err = h.Teardown()
		require.NoError(t, err)
	}()

	ctx := context.Background()

	// create test game turn sheet record
	gameTurnSheetRec := &game_record.GameTurnSheet{
		GameID:           h.Data.GameRecs[0].Record.ID,
		GameInstanceID:   h.Data.GameInstanceRecs[0].Record.ID,
		AccountID:        h.Data.AccountRecs[0].Record.ID,
		TurnNumber:       1,
		SheetType:        "location_choice",
		SheetOrder:       1,
		ProcessingStatus: "pending",
	}

	// prepare sheet data
	sheetData := &location_choice.LocationChoiceData{
		LocationName:        "Crystal Caverns",
		LocationDescription: "A vast network of glittering caves.",
		LocationOptions: []location_choice.LocationOption{
			{
				LocationID:              "loc-2",
				LocationLinkName:        "Floating Islands",
				LocationLinkDescription: "Sky islands connected by bridges.",
			},
		},
	}

	sheetDataJSON, err := json.Marshal(sheetData)
	require.NoError(t, err)
	gameTurnSheetRec.SheetData = sheetDataJSON

	// create turn sheet record
	createdRec, err := h.Domain.(*domain.Domain).GameTurnSheetRepository().CreateOne(gameTurnSheetRec)
	require.NoError(t, err)
	require.NotNil(t, createdRec)

	// generate PDF from sheet_data
	templateDir := getTemplateDir(t)
	outputDir := t.TempDir()
	gen := generator.NewPDFGenerator(h.Log.(*log.Log), templateDir, outputDir)

	templateData := generator.TemplateData{
		AccountRec:      h.Data.AccountRecs[0],
		GameRec:         h.Data.GameRecs[0],
		GameInstanceRec: h.Data.GameInstanceRecs[0],
		TurnSheetData:   sheetData,
		TurnSheetCode:   "TEST-CODE-123",
	}

	pdfData, err := gen.GeneratePDF(ctx, "adventure_game/location_choice/template/content.template", templateData)
	require.NoError(t, err)
	require.NotEmpty(t, pdfData)

	// Save PDF to testdata directory for physical testing
	testDataDir := filepath.Join("testdata")
	err = os.MkdirAll(testDataDir, 0755)
	require.NoError(t, err)

	pdfPath := filepath.Join(testDataDir, "location_choice_test_sheet.pdf")
	err = os.WriteFile(pdfPath, pdfData, 0644)
	require.NoError(t, err)

	t.Logf("Generated test PDF: %s", pdfPath)
	t.Logf("PDF size: %d bytes", len(pdfData))
	t.Logf("Instructions:")
	t.Logf("1. Print the PDF")
	t.Logf("2. Fill out the turn sheet by hand")
	t.Logf("3. Take a photo of the completed sheet")
	t.Logf("4. Save the photo as 'location_choice_filled.jpg' in testdata/")
	t.Logf("5. Use the photo for OCR testing")
}

// Helper functions

// createTestTemplateData creates template data for testing
func createTestTemplateData(t *testing.T, turnSheetData *location_choice.LocationChoiceData) generator.TemplateData {
	t.Helper()

	accountRec := &account_record.Account{
		Name: "Test Player",
	}
	accountRec.Record.ID = "test-account-id"

	gameRec := &game_record.Game{
		Name: "Test Adventure Game",
	}
	gameRec.Record.ID = "test-game-id"

	gameInstanceRec := &game_record.GameInstance{
		CurrentTurn: 1,
	}
	gameInstanceRec.Record.ID = "test-instance-id"

	return generator.TemplateData{
		AccountRec:      accountRec,
		GameRec:         gameRec,
		GameInstanceRec: gameInstanceRec,
		Header: map[string]any{
			"title": "Location Choice",
		},
		Content: map[string]any{
			"instructions": "Choose your next destination",
		},
		Footer: map[string]any{
			"deadline": "Submit by turn deadline",
		},
		TurnSheetData: turnSheetData,
		TurnSheetCode: "TEST-LOC-CHOICE-123",
	}
}

// getTemplateDir returns the path to the turn sheet template directory
func getTemplateDir(t *testing.T) string {
	t.Helper()

	// get the absolute path to the turn_sheet directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// navigate up to internal/turn_sheet
	templateDir := filepath.Join(cwd, "..", "..")
	absPath, err := filepath.Abs(templateDir)
	require.NoError(t, err)

	return absPath
}

// loadTestImage loads a test image from testdata directory
func loadTestImage(t *testing.T, filename string) []byte {
	t.Helper()

	// for now, return mock image data since we don't have real test images yet
	// TODO: Add real test images to testdata/ when OCR is implemented
	return []byte("mock-image-data-for-testing")
}
