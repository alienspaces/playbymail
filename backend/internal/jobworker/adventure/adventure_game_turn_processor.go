package adventure

import (
	"context"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
)

// AdventureGameTurnProcessor coordinates turn processing for adventure games
type AdventureGameTurnProcessor struct {
	Logger                  logger.Logger
	Domain                  *domain.Domain
	LocationChoiceProcessor *AdventureGameLocationChoiceProcessor
}

// NewAdventureGameTurnProcessor creates a new adventure game turn processor
func NewAdventureGameTurnProcessor(l logger.Logger, d *domain.Domain) *AdventureGameTurnProcessor {
	return &AdventureGameTurnProcessor{
		Logger:                  l,
		Domain:                  d,
		LocationChoiceProcessor: NewAdventureGameLocationChoiceProcessor(l, d),
	}
}

// ProcessTurn processes all turn sheets for an adventure game turn
func (p *AdventureGameTurnProcessor) ProcessTurn(ctx context.Context, gameInstanceID string, turnNumber int) error {
	l := p.Logger.WithFunctionContext("AdventureGameTurnProcessor/ProcessTurn")

	l.Info("processing adventure game turn for instance >%s< turn >%d<", gameInstanceID, turnNumber)

	// TODO: Get all turn sheets for this game instance and turn
	// TODO: Group turn sheets by type
	// TODO: Process each sheet type in order
	// TODO: Update game state based on results

	// For now, just process location choice sheets (which are always present)
	// This is the only sheet type we're certain about

	return nil
}
