package adventuregamescanner

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestNewAdventureGameLocationScanner(t *testing.T) {
	tests := []struct {
		name        string
		logger      logger.Logger
		domain      *domain.Domain
		expectError bool
	}{
		{
			name: "creates scanner with valid dependencies",
			logger: func() logger.Logger {
				cfg, _ := config.Parse()
				l, _, _, _ := deps.NewDefaultDependencies(cfg)
				return l
			}(),
			domain:      &domain.Domain{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewAdventureGameLocationScanner(tt.logger, tt.domain)

			if tt.expectError {
				require.Nil(t, scanner, "Expected nil scanner")
			} else {
				require.NotNil(t, scanner, "Expected non-nil scanner")
			}
		})
	}
}

func TestAdventureGameLocationScanner_ScanLocationChoiceSheet(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T, d *domain.Domain) (*game_record.GameTurnSheet, []byte)
		expectError bool
		validate    func(t *testing.T, result *game_record.GameTurnSheet)
	}{
		{
			name: "scans location choice turn sheet successfully",
			setup: func(t *testing.T, d *domain.Domain) (*game_record.GameTurnSheet, []byte) {
				// Create sheet data
				sheetData := map[string]interface{}{
					"type":     "location_choice",
					"options":  []string{"north", "south", "east", "west"},
					"selected": "",
				}
				sheetDataBytes, _ := json.Marshal(sheetData)

				// Generate UUIDs for the turn sheet
				gameID, _ := uuid.NewRandom()
				instanceID, _ := uuid.NewRandom()
				accountID, _ := uuid.NewRandom()
				turnSheetID, _ := uuid.NewRandom()

				turnSheet := &game_record.GameTurnSheet{
					GameID:         gameID.String(),
					GameInstanceID: instanceID.String(),
					AccountID:      accountID.String(),
					SheetType:      "location_choice",
					SheetData:      sheetDataBytes,
				}
				turnSheet.ID = turnSheetID.String()

				return turnSheet, []byte("mock image data")
			},
			expectError: false,
			validate: func(t *testing.T, result *game_record.GameTurnSheet) {
				require.NotNil(t, result, "Result should not be nil")
				require.NotNil(t, result.ScannedData, "ScannedData should be set")

				// Verify the result data contains player choices
				var playerChoices map[string]interface{}
				err := json.Unmarshal(result.ScannedData, &playerChoices)
				require.NoError(t, err, "Should unmarshal ScannedData")
				require.Contains(t, playerChoices, "a", "Should contain choices")

				// Verify the choices are properly extracted
				choices, ok := playerChoices["a"].([]interface{})
				require.True(t, ok, "Choices should be an array")
				require.Len(t, choices, 1, "Should have one choice")
				require.Equal(t, "The dark alley", choices[0], "Should have mock extracted choice")
			},
		},
		{
			name: "handles empty image data",
			setup: func(t *testing.T, d *domain.Domain) (*game_record.GameTurnSheet, []byte) {
				// Create sheet data
				sheetData := map[string]interface{}{
					"type":     "location_choice",
					"options":  []string{"north", "south", "east", "west"},
					"selected": "",
				}
				sheetDataBytes, _ := json.Marshal(sheetData)

				// Generate UUIDs for the turn sheet
				gameID, _ := uuid.NewRandom()
				instanceID, _ := uuid.NewRandom()
				accountID, _ := uuid.NewRandom()
				turnSheetID, _ := uuid.NewRandom()

				turnSheet := &game_record.GameTurnSheet{
					GameID:         gameID.String(),
					GameInstanceID: instanceID.String(),
					AccountID:      accountID.String(),
					SheetType:      "location_choice",
					SheetData:      sheetDataBytes,
				}
				turnSheet.ID = turnSheetID.String()

				return turnSheet, []byte{}
			},
			expectError: false,
			validate: func(t *testing.T, result *game_record.GameTurnSheet) {
				require.NotNil(t, result, "Result should not be nil")
				require.NotNil(t, result.ScannedData, "ScannedData should be set even with empty image")

				// Verify the result data contains player choices
				var playerChoices map[string]interface{}
				err := json.Unmarshal(result.ScannedData, &playerChoices)
				require.NoError(t, err, "Should unmarshal ScannedData")
				require.Contains(t, playerChoices, "a", "Should contain choices")
			},
		},
		{
			name: "handles nil turn sheet record",
			setup: func(t *testing.T, d *domain.Domain) (*game_record.GameTurnSheet, []byte) {
				return nil, []byte("mock image data")
			},
			expectError: true,
			validate: func(t *testing.T, result *game_record.GameTurnSheet) {
				require.Nil(t, result, "Result should be nil on error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup harness
			dcfg := harness.DefaultDataConfig()
			cfg, err := config.Parse()
			require.NoError(t, err, "Parse returns without error")

			l, s, j, err := deps.NewDefaultDependencies(cfg)
			require.NoError(t, err, "Default dependencies returns without error")

			h, err := harness.NewTesting(l, s, j, cfg, dcfg)
			require.NoError(t, err, "NewTesting returns without error")

			_, err = h.Setup()
			require.NoError(t, err, "Setup returns without error")
			defer func() {
				err = h.Teardown()
				require.NoError(t, err, "Teardown returns without error")
			}()

			// Get domain
			d := h.Domain.(*domain.Domain)

			// Setup test data
			turnSheet, imageData := tt.setup(t, d)

			// Create scanner
			scanner := NewAdventureGameLocationScanner(func() logger.Logger {
				cfg, _ := config.Parse()
				l, _, _, _ := deps.NewDefaultDependencies(cfg)
				return l
			}(), d)

			// Execute test
			result, err := scanner.ScanLocationChoiceSheet(context.Background(), turnSheet, imageData)

			// Verify results
			if tt.expectError {
				require.Error(t, err, "Expected error but got none")
				require.Nil(t, result, "Expected nil result on error")
			} else {
				require.NoError(t, err, "Unexpected error: %v", err)
				require.NotNil(t, result, "Expected non-nil result")
			}

			tt.validate(t, result)
		})
	}
}

func TestLocationChoiceScanData_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		scanData *LocationChoiceScanData
		validate func(t *testing.T, jsonData []byte)
	}{
		{
			name: "serializes valid scan data",
			scanData: &LocationChoiceScanData{
				Choices: []string{"The dark alley", "The bright path"},
			},
			validate: func(t *testing.T, jsonData []byte) {
				require.NotEmpty(t, jsonData, "JSON data should not be empty")

				// Verify we can unmarshal it back
				var unmarshaled LocationChoiceScanData
				err := json.Unmarshal(jsonData, &unmarshaled)
				require.NoError(t, err, "Should unmarshal without error")
				require.Len(t, unmarshaled.Choices, 2, "Should have two choices")
				require.Contains(t, unmarshaled.Choices, "The dark alley", "Should contain first choice")
				require.Contains(t, unmarshaled.Choices, "The bright path", "Should contain second choice")
			},
		},
		{
			name: "serializes empty scan data",
			scanData: &LocationChoiceScanData{
				Choices: []string{},
			},
			validate: func(t *testing.T, jsonData []byte) {
				require.NotEmpty(t, jsonData, "JSON data should not be empty")

				// Verify we can unmarshal it back
				var unmarshaled LocationChoiceScanData
				err := json.Unmarshal(jsonData, &unmarshaled)
				require.NoError(t, err, "Should unmarshal without error")
				require.Len(t, unmarshaled.Choices, 0, "Should have no choices")
			},
		},
		{
			name: "serializes nil choices",
			scanData: &LocationChoiceScanData{
				Choices: nil,
			},
			validate: func(t *testing.T, jsonData []byte) {
				require.NotEmpty(t, jsonData, "JSON data should not be empty")

				// Verify we can unmarshal it back
				var unmarshaled LocationChoiceScanData
				err := json.Unmarshal(jsonData, &unmarshaled)
				require.NoError(t, err, "Should unmarshal without error")
				require.Nil(t, unmarshaled.Choices, "Choices should be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.scanData)
			require.NoError(t, err, "Should marshal without error")

			tt.validate(t, jsonData)
		})
	}
}
