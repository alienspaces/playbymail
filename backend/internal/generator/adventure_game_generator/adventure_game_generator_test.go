package adventuregamegenerator

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/generator"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// createValidTemplateData creates a valid TemplateData struct for testing
func createValidTemplateData() generator.TemplateData {
	accountRec := &account_record.Account{
		Record: record.Record{ID: "account-1"},
		Name:   "Aria the Mage",
		Email:  "aria@example.com",
	}

	gameRec := &game_record.Game{
		Record:   record.Record{ID: "game-1"},
		Name:     "The Enchanted Forest Adventure",
		GameType: game_record.GameTypeAdventure,
	}

	gameInstanceRec := &game_record.GameInstance{
		Record:      record.Record{ID: "instance-1"},
		GameID:      gameRec.ID,
		CurrentTurn: 1,
		Status:      game_record.GameInstanceStatusStarted,
	}

	return generator.TemplateData{
		AccountRec:       accountRec,
		GameInstanceRec:  gameInstanceRec,
		GameRec:          gameRec,
		BackgroundTop:    "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iIzAwN2YwMCIvPjwvc3ZnPg==",
		BackgroundMiddle: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iIzAwN2YwMCIvPjwvc3ZnPg==",
		BackgroundBottom: "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iIzAwN2YwMCIvPjwvc3ZnPg==",
		TurnSheetData: &AdventureGameLocationTurnSheetData{
			CurrentLocationRec: &adventure_game_record.AdventureGameLocation{
				Record:      record.Record{ID: "location-1"},
				GameID:      gameRec.ID,
				Name:        "Mystic Grove",
				Description: "A peaceful clearing with ancient trees and magical creatures.",
			},
			AvailableLocationLinkRecs: []*adventure_game_record.AdventureGameLocationLink{
				{
					Record:                    record.Record{ID: "link-1"},
					GameID:                    gameRec.ID,
					Name:                      "Path to Crystal Caverns",
					Description:               "A winding path leads down into the depths where glowing crystals illuminate the underground passages.",
					ToAdventureGameLocationID: "location-2",
				},
			},
			NextTurnDeadline: "Next Friday at 5:00 PM",
		},
	}
}

// getExpectedHTMLContains returns the expected strings that should be found in generated HTML
func getExpectedHTMLContains() []string {
	return []string{
		"The Enchanted Forest Adventure",
		"Turn 1",
		"Aria the Mage",
		"Mystic Grove",
		"Path to Crystal Caverns",
		"Next Friday at 5:00 PM",
	}
}

func TestAdventureGameLocationGenerator(t *testing.T) {
	// Use actual templates directory
	templateDir := "."

	// Create logger for testing
	cfg := config.Config{
		LogLevel:    "debug",
		LogIsPretty: true,
	}
	logger, err := log.NewLogger(cfg)
	require.NoError(t, err, "Failed to create logger")

	// Create domain for testing
	domain := &domain.Domain{}

	// Create generator with temporary output directory
	tempOutputDir := t.TempDir()
	gen := NewAdventureGameLocationGenerator(logger, domain, templateDir, tempOutputDir)

	t.Run("GenerateLocationChoiceTurnSheet", func(t *testing.T) {
		tests := []struct {
			name            string
			data            generator.TemplateData
			expectError     bool
			expectedMinSize int
		}{
			{
				name:            "valid location choice data",
				data:            createValidTemplateData(),
				expectError:     false,
				expectedMinSize: 100000, // PDF should be at least 100KB
			},
			{
				name: "invalid turn sheet data type",
				data: func() generator.TemplateData {
					data := createValidTemplateData()
					data.TurnSheetData = "invalid data type"
					return data
				}(),
				expectError:     true,
				expectedMinSize: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()
				turnSheetCode := "test-turn-sheet-code-1"
				pdfData, err := gen.GenerateLocationChoiceTurnSheet(ctx, tt.data, turnSheetCode)

				if tt.expectError {
					require.Error(t, err, "Expected error but got none")
					require.Nil(t, pdfData, "Expected nil PDF data on error")
				} else {
					require.NoError(t, err, "Unexpected error: %v", err)
					require.NotNil(t, pdfData, "Expected PDF data but got nil")
					require.GreaterOrEqual(t, len(pdfData), tt.expectedMinSize, "PDF size too small")
				}
			})
		}
	})

	t.Run("GenerateLocationChoiceTurnSheetToFile", func(t *testing.T) {

		tests := []struct {
			name        string
			data        generator.TemplateData
			filename    string
			expectError bool
		}{
			{
				name:        "valid data with file output",
				data:        createValidTemplateData(),
				filename:    "test_location_choice.pdf",
				expectError: false,
			},
			{
				name: "invalid data type",
				data: func() generator.TemplateData {
					data := createValidTemplateData()
					data.TurnSheetData = "invalid data type"
					return data
				}(),
				filename:    "test_invalid.pdf",
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()
				turnSheetCode := "test-turn-sheet-code-1"
				err := gen.GenerateLocationChoiceTurnSheetToFile(ctx, tt.data, turnSheetCode, tt.filename)

				if tt.expectError {
					require.Error(t, err, "Expected error but got none")
					// Verify file was not created
					_, fileErr := os.Stat(filepath.Join(tempOutputDir, tt.filename))
					require.Error(t, fileErr, "File should not exist on error")
				} else {
					require.NoError(t, err, "Unexpected error: %v", err)
					// Verify file was created and has content
					filePath := filepath.Join(tempOutputDir, tt.filename)
					fileInfo, err := os.Stat(filePath)
					require.NoError(t, err, "File should exist at %s", filePath)
					require.Greater(t, fileInfo.Size(), int64(100000), "File should be substantial size")
				}
			})
		}
	})

	t.Run("GenerateLocationChoiceTurnSheetHTML", func(t *testing.T) {
		tests := []struct {
			name             string
			data             generator.TemplateData
			expectError      bool
			expectedContains []string
		}{
			{
				name:             "valid location choice data",
				data:             createValidTemplateData(),
				expectError:      false,
				expectedContains: getExpectedHTMLContains(),
			},
			{
				name: "invalid turn sheet data type",
				data: func() generator.TemplateData {
					data := createValidTemplateData()
					data.TurnSheetData = "invalid data type"
					return data
				}(),
				expectError:      true,
				expectedContains: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()
				turnSheetCode := "test-turn-sheet-code-1"
				html, err := gen.GenerateLocationChoiceTurnSheetHTML(ctx, tt.data, turnSheetCode)

				if tt.expectError {
					require.Error(t, err, "Expected error but got none")
					require.Empty(t, html, "Expected empty HTML on error")
				} else {
					require.NoError(t, err, "Unexpected error: %v", err)
					require.NotEmpty(t, html, "Expected HTML content")

					// Verify HTML contains expected content
					for _, expected := range tt.expectedContains {
						require.Contains(t, html, expected, "HTML should contain: %s", expected)
					}
				}
			})
		}
	})
}
