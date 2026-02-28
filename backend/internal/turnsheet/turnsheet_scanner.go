package turnsheet

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// TurnSheetScanner defines the interface for scanning turn sheets
type TurnSheetScanner interface {
	// GetTurnSheetCodeFromImage extracts the turn sheet code from image data
	GetTurnSheetCodeFromImage(ctx context.Context, l logger.Logger, imageData []byte) (string, error)
	// GetTurnSheetScanData scans a turn sheet image and returns the extracted data
	GetTurnSheetScanData(ctx context.Context, l logger.Logger, sheetType string, sheetData []byte, imageData []byte) ([]byte, error)
}

// turnSheetScanner is the default implementation of TurnSheetScanner
type turnSheetScanner struct {
	cfg config.Config
}

// NewScanner creates a new turn sheet scanner instance
func NewScanner(cfg config.Config) (TurnSheetScanner, error) {
	return &turnSheetScanner{
		cfg: cfg,
	}, nil
}

// GetTurnSheetCodeFromImage extracts the turn sheet code from image data
func (s *turnSheetScanner) GetTurnSheetCodeFromImage(ctx context.Context, l logger.Logger, imageData []byte) (string, error) {
	baseProcessor, err := NewBaseProcessor(l, s.cfg)
	if err != nil {
		return "", fmt.Errorf("failed to create base processor: %w", err)
	}

	return baseProcessor.ParseTurnSheetCodeFromImage(ctx, imageData)
}

// GetTurnSheetScanData scans a turn sheet image and returns the extracted data
func (s *turnSheetScanner) GetTurnSheetScanData(ctx context.Context, l logger.Logger, sheetType string, sheetData []byte, imageData []byte) ([]byte, error) {
	processor, err := GetDocumentProcessor(l, s.cfg, sheetType)
	if err != nil {
		return nil, fmt.Errorf("failed to get processor for sheet type %s: %w", sheetType, err)
	}

	return processor.ScanTurnSheet(ctx, l, sheetData, imageData)
}
