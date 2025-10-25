package scanner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestNewImageScanner(t *testing.T) {
	tests := []struct {
		name   string
		logger logger.Logger
	}{
		{
			name: "creates image scanner with valid dependencies",
			logger: func() logger.Logger {
				cfg, _ := config.Parse()
				l, _, _, _ := deps.NewDefaultDependencies(cfg)
				return l
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewImageScanner(tt.logger)

			require.NotNil(t, scanner, "Scanner should not be nil")
			require.Equal(t, tt.logger, scanner.logger, "Logger should match")
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

			// Create scanner
			scanner := NewImageScanner(l)

			// Execute test
			text, err := scanner.ExtractTextFromImage(context.Background(), tt.imageData)

			// Verify results
			tt.validate(t, text, err)
		})
	}
}
