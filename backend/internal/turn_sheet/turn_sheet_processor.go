package turn_sheet

import (
	"context"
	"fmt"
	"maps"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// DocumentScanner defines the interface for scanning completed turn sheet documents
type DocumentScanner interface {
	// ScanTurnSheet scans a turn sheet image and extracts player choices/directions
	// sheetData: JSON-encoded sheet data from the database
	// Returns: JSON-encoded scan results to store in the database
	ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error)
}

// DocumentGenerator defines the interface for generating turn sheet documents
type DocumentGenerator interface {
	// GenerateTurnSheet generates a turn sheet document with the provided data
	// sheetData: JSON-encoded sheet data from the database
	// Returns: PDF bytes for the generated turn sheet
	GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, sheetData []byte) ([]byte, error)
}

// DocumentProcessor defines the interface for processing turn sheet documents (generation + scanning)
type DocumentProcessor interface {
	DocumentScanner
	DocumentGenerator
}

// GetDocumentProcessor returns the document processor for a specific turn sheet type
func GetDocumentProcessor(l logger.Logger, cfg config.Config, sheetType string) (DocumentProcessor, error) {
	// Get turn sheet processor map
	turnSheetProcessorMap, err := getDocumentProcessorMap(l, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create document processor map: %w", err)
	}

	// Get processor for turn sheet type
	processor, exists := turnSheetProcessorMap[sheetType]
	if !exists {
		return nil, fmt.Errorf("no processor registered for turn sheet type: %s", sheetType)
	}

	return processor, nil
}

// getDocumentProcessorMap returns a map of document processors for all turn sheet types
func getDocumentProcessorMap(l logger.Logger, cfg config.Config) (map[string]DocumentProcessor, error) {
	processors := make(map[string]DocumentProcessor)

	adventureProcessors, err := getAdventureGameDocumentProcessorMap(l, cfg)
	if err != nil {
		return nil, err
	}
	maps.Copy(processors, adventureProcessors)

	return processors, nil
}

func getAdventureGameDocumentProcessorMap(l logger.Logger, cfg config.Config) (map[string]DocumentProcessor, error) {
	processors := make(map[string]DocumentProcessor)

	locationChoiceProcessor, err := NewLocationChoiceProcessor(l, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create location choice processor: %w", err)
	}
	processors[adventure_game_record.AdventureSheetTypeLocationChoice] = locationChoiceProcessor

	joinGameProcessor, err := NewJoinGameProcessor(l, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create join game processor: %w", err)
	}
	processors[adventure_game_record.AdventureSheetTypeJoinGame] = joinGameProcessor

	inventoryManagementProcessor, err := NewInventoryManagementProcessor(l, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory management processor: %w", err)
	}
	processors[adventure_game_record.AdventureSheetTypeInventoryManagement] = inventoryManagementProcessor

	return processors, nil
}
