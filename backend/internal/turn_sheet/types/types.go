package types

import (
	"context"
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

// Turn sheet types
const (
	TurnSheetTypeLocationChoice string = "location_choice"
)

// TurnSheetData represents the data for a turn sheet
//
// All turn sheet types use this same data structure
type TurnSheetTemplateData struct {
	// Game data
	GameName *string
	GameType *string

	// Game instance data
	TurnNumber *int

	// Account data
	AccountName *string

	// Background images
	BackgroundTop    *string
	BackgroundMiddle *string
	BackgroundBottom *string

	// Turn sheet
	TurnSheetDeadline *time.Time
	TurnSheetCode     *string
}

// DocumentScanner defines the interface for scanning completed turn sheet documents
type DocumentScanner interface {
	// ScanTurnSheet scans a turn sheet image and extracts player choices/directions
	ScanTurnSheet(ctx context.Context, l logger.Logger, imageData []byte, sheetData any) (any, error)
}

// DocumentGenerator defines the interface for generating turn sheet documents
type DocumentGenerator interface {
	// GenerateTurnSheet generates a turn sheet document with the provided data
	GenerateTurnSheet(ctx context.Context, l logger.Logger, data any) ([]byte, error)
}

// DocumentProcessor defines the interface for processing turn sheet documents (generation + scanning)
type DocumentProcessor interface {
	DocumentScanner
	DocumentGenerator
}
