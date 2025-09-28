package adventuregamescanner

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// AdventureGameLocationScanner handles scanning of location choice turn sheets
type AdventureGameLocationScanner struct {
	logger logger.Logger
	domain *domain.Domain
}

// NewAdventureGameLocationScanner creates a new location choice scanner
func NewAdventureGameLocationScanner(l logger.Logger, d *domain.Domain) *AdventureGameLocationScanner {
	return &AdventureGameLocationScanner{
		logger: l,
		domain: d,
	}
}

// scanLocationChoiceSheet processes a location choice turn sheet
func (s *AdventureGameLocationScanner) ScanLocationChoiceSheet(ctx context.Context, turnSheetRec *game_record.GameTurnSheet, imageData []byte) (*game_record.GameTurnSheet, error) {
	l := s.logger.WithFunctionContext("AdventureGameLocationScanner/ScanLocationChoiceSheet")

	if turnSheetRec == nil {
		return nil, fmt.Errorf("turn sheet record cannot be nil")
	}

	l.Info("scanning location choice sheet >%s<", turnSheetRec.ID)

	// TODO: Implement actual OCR/image processing here
	// For now, we'll mock the OCR extraction process
	extractedChoice := "The dark alley" // Mock extracted choice from OCR

	scanData := &LocationChoiceScanData{
		// OCR extracted data
		Choices: []string{extractedChoice},
	}

	// Convert scan data to generic result format
	playerChoices := map[string]interface{}{
		"a": scanData.Choices,
	}

	ScannedData, err := json.Marshal(playerChoices)
	if err != nil {
		l.Warn("failed to marshal player choices >%v<", err)
		return nil, err
	}

	turnSheetRec.ScannedData = ScannedData

	return turnSheetRec, nil
}
