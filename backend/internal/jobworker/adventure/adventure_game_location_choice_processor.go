package adventure

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// AdventureGameLocationChoiceProcessor processes location choice turn sheets for adventure games
type AdventureGameLocationChoiceProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewAdventureGameLocationChoiceProcessor creates a new adventure game location choice processor
func NewAdventureGameLocationChoiceProcessor(l logger.Logger, d *domain.Domain) *AdventureGameLocationChoiceProcessor {
	return &AdventureGameLocationChoiceProcessor{
		Logger: l,
		Domain: d,
	}
}

// ProcessLocationChoice processes a single location choice turn sheet
func (p *AdventureGameLocationChoiceProcessor) ProcessLocationChoice(ctx context.Context, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("AdventureGameLocationChoiceProcessor/ProcessLocationChoice")

	l.Info("processing location choice for turn sheet >%s<", turnSheet.ID)

	// Verify this is a location choice sheet
	if turnSheet.SheetType != adventure_game_record.AdventureSheetTypeLocationChoice {
		l.Warn("expected location choice sheet type, got >%s<", turnSheet.SheetType)
		return fmt.Errorf("invalid sheet type: expected %s, got %s", adventure_game_record.AdventureSheetTypeLocationChoice, turnSheet.SheetType)
	}

	// TODO: Implement actual location choice processing logic
	// This will involve:
	// 1. Parse the player's location choice from turnSheet.SheetData
	// 2. Validate the choice is valid for the character's current location
	// 3. Update character's location in the game state
	// 4. Generate any location-specific events or encounters

	return nil
}
