package scanner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/harness"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/deps"
)

func TestNewScanner(t *testing.T) {
	tests := []struct {
		name   string
		logger logger.Logger
		domain *domain.Domain
	}{
		{
			name: "creates scanner with valid dependencies",
			logger: func() logger.Logger {
				cfg, _ := config.Parse()
				l, _, _, _ := deps.NewDefaultDependencies(cfg)
				return l
			}(),
			domain: &domain.Domain{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.logger, tt.domain)

			require.NotNil(t, scanner, "Scanner should not be nil")
			require.Equal(t, tt.logger, scanner.logger, "Logger should match")
			require.Equal(t, tt.domain, scanner.domain, "Domain should match")
		})
	}
}

func TestScanner_ParseTurnSheetCodeFromImage(t *testing.T) {
	tests := []struct {
		name        string
		imageData   []byte
		expectError bool
		validate    func(t *testing.T, code string, err error)
	}{
		{
			name:        "returns error for unimplemented OCR",
			imageData:   []byte("mock image data"),
			expectError: true,
			validate: func(t *testing.T, code string, err error) {
				require.Empty(t, code, "Code should be empty on error")
				require.Error(t, err, "Should return error for unimplemented OCR")
				require.Contains(t, err.Error(), "OCR not implemented", "Error should mention OCR not implemented")
			},
		},
		{
			name:        "handles empty image data",
			imageData:   []byte{},
			expectError: true,
			validate: func(t *testing.T, code string, err error) {
				require.Empty(t, code, "Code should be empty on error")
				require.Error(t, err, "Should return error for empty image data")
				require.Contains(t, err.Error(), "OCR not implemented", "Error should mention OCR not implemented")
			},
		},
		{
			name:        "handles nil image data",
			imageData:   nil,
			expectError: true,
			validate: func(t *testing.T, code string, err error) {
				require.Empty(t, code, "Code should be empty on error")
				require.Error(t, err, "Should return error for nil image data")
				require.Contains(t, err.Error(), "OCR not implemented", "Error should mention OCR not implemented")
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

			// Create scanner
			scanner := NewScanner(l, d)

			// Execute test
			code, err := scanner.ParseTurnSheetCodeFromImage(context.Background(), tt.imageData)

			// Verify results
			tt.validate(t, code, err)
		})
	}
}
