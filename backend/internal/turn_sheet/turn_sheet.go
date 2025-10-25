package turn_sheet

import (
	"fmt"
	"maps"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet/adventure_game/location_choice"
	"gitlab.com/alienspaces/playbymail/internal/turn_sheet/types"
)

// GetDocumentProcessor returns the document processor for a specific turn sheet type
func GetDocumentProcessor(l logger.Logger, sheetType string) (types.DocumentProcessor, error) {

	// Get turn sheet processor map
	turnSheetProcessorMap := getDocumentProcessorMap(l)

	// Get processor for turn sheet type
	processor, exists := turnSheetProcessorMap[sheetType]
	if !exists {
		return nil, fmt.Errorf("no processor registered for turn sheet type: %s", sheetType)
	}

	return processor, nil
}

// getDocumentProcessorMap returns a map of document processors for all turn sheet types
func getDocumentProcessorMap(l logger.Logger) map[string]types.DocumentProcessor {

	processors := make(map[string]types.DocumentProcessor)

	maps.Copy(processors, getAdventureGameDocumentProcessorMap(l))

	return processors
}

func getAdventureGameDocumentProcessorMap(l logger.Logger) map[string]types.DocumentProcessor {
	return map[string]types.DocumentProcessor{
		adventure_game_record.AdventureSheetTypeLocationChoice: location_choice.NewLocationChoiceProcessor(l),
	}
}
