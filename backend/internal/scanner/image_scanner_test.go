package scanner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/log"
)

func TestNewImageScanner(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "creates image scanner with valid dependencies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.NewDefaultLogger()
			scanner := NewImageScanner(logger)

			require.NotNil(t, scanner, "Scanner should not be nil")
			require.Equal(t, logger, scanner.logger, "Logger should match")
		})
	}
}

func TestImageScanner_ExtractTextFromImage(t *testing.T) {
	tests := []struct {
		name        string
		imageData   []byte
		expectError bool
		validate    func(t *testing.T, text string, err error)
	}{
		{
			name:        "returns error for empty image data",
			imageData:   []byte{},
			expectError: true,
			validate: func(t *testing.T, text string, err error) {
				require.Empty(t, text, "Text should be empty on error")
				require.Error(t, err, "Should return error for empty image data")
				require.Contains(t, err.Error(), "empty image data", "Error should mention empty image data")
			},
		},
		{
			name:        "returns error for nil image data",
			imageData:   nil,
			expectError: true,
			validate: func(t *testing.T, text string, err error) {
				require.Empty(t, text, "Text should be empty on error")
				require.Error(t, err, "Should return error for nil image data")
				require.Contains(t, err.Error(), "empty image data", "Error should mention empty image data")
			},
		},
		{
			name:        "returns error for image data too small",
			imageData:   []byte("small"),
			expectError: true,
			validate: func(t *testing.T, text string, err error) {
				require.Empty(t, text, "Text should be empty on error")
				require.Error(t, err, "Should return error for small image data")
				require.Contains(t, err.Error(), "too small", "Error should mention image data too small")
			},
		},
		{
			name:        "handles invalid image data gracefully",
			imageData:   make([]byte, 200), // Valid size but invalid image data
			expectError: true,
			validate: func(t *testing.T, text string, err error) {
				require.Empty(t, text, "Text should be empty on error")
				require.Error(t, err, "Should return error for invalid image data")
				// OCR may fail with various error messages for invalid data
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger
			logger := log.NewDefaultLogger()

			// Create scanner
			scanner := NewImageScanner(logger)

			// Execute test
			text, err := scanner.ExtractTextFromImage(context.Background(), tt.imageData)

			// Verify results
			tt.validate(t, text, err)
		})
	}
}
