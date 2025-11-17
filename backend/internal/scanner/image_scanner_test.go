package scanner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/internal/agent"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// mockMultiModalAgent is a mock implementation of agent.MultiModalAgent for testing
type mockMultiModalAgent struct {
	extractTextFunc           func(ctx context.Context, req agent.TextExtractionRequest) (string, error)
	extractStructuredDataFunc func(ctx context.Context, req agent.StructuredExtractionRequest) ([]byte, error)
	generateContentFunc       func(ctx context.Context, req agent.ContentGenerationRequest) (string, error)
	analyzeTextFunc           func(ctx context.Context, req agent.TextAnalysisRequest) (agent.TextAnalysisResult, error)
}

func (m *mockMultiModalAgent) ExtractText(ctx context.Context, req agent.TextExtractionRequest) (string, error) {
	if m.extractTextFunc != nil {
		return m.extractTextFunc(ctx, req)
	}
	return "", nil
}

func (m *mockMultiModalAgent) ExtractStructuredData(ctx context.Context, req agent.StructuredExtractionRequest) ([]byte, error) {
	if m.extractStructuredDataFunc != nil {
		return m.extractStructuredDataFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockMultiModalAgent) GenerateContent(ctx context.Context, req agent.ContentGenerationRequest) (string, error) {
	if m.generateContentFunc != nil {
		return m.generateContentFunc(ctx, req)
	}
	return "", nil
}

func (m *mockMultiModalAgent) AnalyzeText(ctx context.Context, req agent.TextAnalysisRequest) (agent.TextAnalysisResult, error) {
	if m.analyzeTextFunc != nil {
		return m.analyzeTextFunc(ctx, req)
	}
	return agent.TextAnalysisResult{}, nil
}

func TestNewImageScanner(t *testing.T) {
	t.Run("creates image scanner with valid dependencies", func(t *testing.T) {
		logger := log.NewDefaultLogger()
		cfg := config.Config{}
		scanner, err := NewImageScanner(logger, cfg)

		require.NoError(t, err, "Should not return error")
		require.NotNil(t, scanner, "Scanner should not be nil")
		require.Equal(t, logger, scanner.logger, "Logger should match")
		require.Equal(t, cfg, scanner.cfg, "Config should match")
		require.NotNil(t, scanner.agent, "Agent should be configured")
	})
}

func TestImageScanner_ExtractTextFromImage(t *testing.T) {
	tests := []struct {
		name        string
		imageData   []byte
		mock        *mockMultiModalAgent
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
			name:      "delegates to agent",
			imageData: []byte("pretend-image"),
			mock: &mockMultiModalAgent{
				extractTextFunc: func(_ context.Context, req agent.TextExtractionRequest) (string, error) {
					require.NotEmpty(t, req.ImageData, "Image data should be provided")
					return "mock text", nil
				},
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
			scanner, err := NewImageScanner(logger, cfg)
			require.NoError(t, err)

			if tt.mock != nil {
				scanner.SetAgent(tt.mock)
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
	scanner, err := NewImageScanner(logger, cfg)
	require.NoError(t, err)

	mockResponse := []byte(`{"email":"alienspaces@gmail.com"}`)
	mockAgent := &mockMultiModalAgent{
		extractStructuredDataFunc: func(ctx context.Context, req agent.StructuredExtractionRequest) ([]byte, error) {
			require.NotNil(t, ctx)
			require.NotEmpty(t, req.FilledImage.Data, "Filled image data should be provided")
			require.NotNil(t, req.ExpectedSchema, "Expected schema should be provided")
			return mockResponse, nil
		},
	}
	scanner.SetAgent(mockAgent)

	result, err := scanner.ExtractStructuredData(context.Background(), StructuredScanRequest{
		Instructions:       "test",
		FilledImage:        []byte("image"),
		FilledImageMIME:    "image/png",
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
