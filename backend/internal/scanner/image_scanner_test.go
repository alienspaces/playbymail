package scanner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
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
			cfg := config.Config{}
			scanner := NewImageScanner(logger, cfg)

			require.NotNil(t, scanner, "Scanner should not be nil")
			require.Equal(t, logger, scanner.logger, "Logger should match")
			require.Equal(t, cfg, scanner.cfg, "Config should match")
		})
	}
}

func TestImageScanner_ExtractTextFromImage(t *testing.T) {
	tests := []struct {
		name        string
		imageData   []byte
		mock        TextExtractorFunc
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
			name:      "delegates to text extractor",
			imageData: []byte("pretend-image"),
			mock: func(_ context.Context, _ []byte) (string, error) {
				return "mock text", nil
			},
			validate: func(t *testing.T, text string, err error) {
				require.NoError(t, err)
				require.Equal(t, "mock text", text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger
			logger := log.NewDefaultLogger()
			cfg := config.Config{}

			// Create scanner
			scanner := NewImageScanner(logger, cfg)
			if tt.mock != nil {
				scanner.SetTextExtractor(tt.mock)
			}

			// Execute test
			text, err := scanner.ExtractTextFromImage(context.Background(), tt.imageData)

			// Verify results
			tt.validate(t, text, err)
		})
	}
}

func TestImageScanner_ExtractStructuredData(t *testing.T) {
	logger := log.NewDefaultLogger()
	cfg := config.Config{}
	scanner := NewImageScanner(logger, cfg)

	mockResponse := []byte(`{"email":"alienspaces@gmail.com"}`)
	scanner.SetStructuredExtractor(func(ctx context.Context, req StructuredScanRequest) ([]byte, error) {
		require.NotNil(t, ctx)
		require.NotEmpty(t, req.FilledImage)
		require.NotNil(t, req.ExpectedJSONSchema)
		return mockResponse, nil
	})

	result, err := scanner.ExtractStructuredData(context.Background(), StructuredScanRequest{
		Instructions:       "test",
		FilledImage:        []byte("image"),
		ExpectedJSONSchema: map[string]any{"email": ""},
	})

	require.NoError(t, err)
	require.Equal(t, mockResponse, result)

	_, err = scanner.ExtractStructuredData(context.Background(), StructuredScanRequest{
		Instructions:       "test",
		ExpectedJSONSchema: map[string]any{"email": ""},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "filled image")
}
