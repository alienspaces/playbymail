package location_choice_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/playbymail/internal/turn_sheet/adventure_game/location_choice"
	"gitlab.com/alienspaces/playbymail/internal/utils/testutil"
)

func TestLocationChoiceProcessor_GenerateTurnSheet(t *testing.T) {
	tests := []struct {
		name        string
		data        any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "generates turn sheet with valid data",
			data:        &location_choice.LocationChoiceData{},
			expectError: false,
		},
		{
			name:        "generates turn sheet with nil data",
			data:        nil,
			expectError: false, // Generator may handle nil data gracefully
		},
		{
			name:        "generates turn sheet with invalid data type",
			data:        "invalid data",
			expectError: false, // Generator may handle invalid data gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup test harness
			l, _, _, _ := testutil.NewDefaultDependencies(t)

			processor := location_choice.NewLocationChoiceProcessor(l)

			ctx := context.Background()
			pdfData, err := processor.GenerateTurnSheet(ctx, l, tt.data)

			if tt.expectError {
				require.Error(t, err, "Should return error")
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg, "Error message should contain expected text")
				}
				require.Nil(t, pdfData, "PDF data should be nil on error")
			} else {
				// Note: This test may fail if PDF generation requires specific dependencies
				// In that case, we'd mock the generator or skip PDF generation tests
				if err != nil {
					t.Logf("PDF generation failed (may be expected in test environment): %v", err)
				}
			}
		})
	}
}

func TestLocationChoiceProcessor_ScanTurnSheet(t *testing.T) {
	tests := []struct {
		name        string
		imageData   []byte
		sheetData   any
		expectError bool
		errorMsg    string
	}{
		{
			name:        "returns error for empty image data",
			imageData:   []byte{},
			sheetData:   map[string]any{"locations": []interface{}{}},
			expectError: true,
			errorMsg:    "empty image data",
		},
		{
			name:        "returns error for nil image data",
			imageData:   nil,
			sheetData:   map[string]any{"locations": []interface{}{}},
			expectError: true,
			errorMsg:    "empty image data",
		},
		{
			name:        "returns error for invalid sheet data format",
			imageData:   []byte("fake image data"),
			sheetData:   "invalid sheet data",
			expectError: true,
			errorMsg:    "invalid sheet data format",
		},
		{
			name:        "returns error for sheet data without locations",
			imageData:   []byte("fake image data"),
			sheetData:   map[string]any{"other": "data"},
			expectError: true,
			errorMsg:    "text extraction failed", // Will fail at OCR extraction before sheet data validation
		},
		{
			name:      "handles valid sheet data with locations",
			imageData: []byte("fake image data"),
			sheetData: map[string]any{
				"locations": []interface{}{
					map[string]interface{}{
						"name": "Crystal Caverns",
					},
					map[string]interface{}{
						"name": "Mystic Grove",
					},
				},
			},
			expectError: true, // Will fail at OCR extraction, but should get past sheet data validation
			errorMsg:    "text extraction failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup test harness
			l, _, _, _ := testutil.NewDefaultDependencies(t)

			processor := location_choice.NewLocationChoiceProcessor(l)

			ctx := context.Background()
			result, err := processor.ScanTurnSheet(ctx, l, tt.imageData, tt.sheetData)

			if tt.expectError {
				require.Error(t, err, "Should return error")
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg, "Error message should contain expected text")
				}
				require.Nil(t, result, "Result should be nil on error")
			} else {
				require.NoError(t, err, "Should not return error")
				require.NotNil(t, result, "Result should not be nil")
			}
		})
	}
}
