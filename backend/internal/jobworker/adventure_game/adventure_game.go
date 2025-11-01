package adventure_game

import (
	"context"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/jobworker/adventure_game/turn_sheet_processor"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// AdventureGame coordinates turn processing for adventure games
// NOTE: Assumes only existing character instances are present in the game instance
// New player onboarding (join game) should be handled separately through API endpoints
type AdventureGame struct {
	Logger     logger.Logger
	Domain     *domain.Domain
	Processors map[string]TurnSheetProcessor
}

// TurnSheetProcessor defines the interface for processing turn sheet business logic in adventure games
type TurnSheetProcessor interface {
	// GetSheetType returns the sheet type this processor handles
	GetSheetType() string

	// ProcessTurnSheetResponse processes a single turn sheet response and updates game state
	ProcessTurnSheetResponse(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstance *adventure_game_record.AdventureGameCharacterInstance, turnSheet *game_record.GameTurnSheet) error

	// CreateNextTurnSheet creates a new turn sheet record for the next turn
	CreateNextTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstance *adventure_game_record.AdventureGameCharacterInstance) (*game_record.GameTurnSheet, error)
}

// NewAdventureGame creates a new adventure game turn processor
func NewAdventureGame(l logger.Logger, d *domain.Domain) (*AdventureGame, error) {
	l = l.WithFunctionContext("NewAdventureGame")

	g := &AdventureGame{
		Logger: l,
		Domain: d,
	}

	p, err := g.initializeTurnSheetProcessors()
	if err != nil {
		l.Warn("failed to initialize turn sheet processors >%v<", err)
		return nil, err
	}
	g.Processors = p

	return g, nil
}

// initializeTurnSheetProcessors creates and registers all available adventure game turn sheet business processors
// To add new turn sheet types: 1) Create processor in internal/jobworker/adventure_game/turn_sheet_processor/ 2) Register here
func (p *AdventureGame) initializeTurnSheetProcessors() (map[string]TurnSheetProcessor, error) {
	l := p.Logger.WithFunctionContext("AdventureGame/initializeTurnSheetProcessors")

	processors := make(map[string]TurnSheetProcessor)

	// Register location choice processor
	locationChoiceProcessor, err := turn_sheet_processor.NewAdventureGameLocationChoiceProcessor(l, p.Domain)
	if err != nil {
		l.Warn("failed to initialize location choice processor >%v<", err)
		return nil, err
	}
	processors[adventure_game_record.AdventureSheetTypeLocationChoice] = locationChoiceProcessor

	// TODO: Add new adventure game turn sheet processors here
	// Example: processors[adventure_game_record.AdventureSheetTypeCombat] = combatProcessor
	// Example: processors[adventure_game_record.AdventureSheetTypeInventory] = inventoryProcessor
	// Example: processors[adventure_game_record.AdventureSheetTypeDialogue] = dialogueProcessor

	return processors, nil
}

// getCharacterInstancesForGameInstance retrieves all character instances for a game instance
// NOTE: This assumes only existing character instances - new players are not added during turn processing
func (p *AdventureGame) getCharacterInstancesForGameInstance(_ context.Context, gameInstanceRec *game_record.GameInstance) ([]*adventure_game_record.AdventureGameCharacterInstance, error) {
	l := p.Logger.WithFunctionContext("AdventureGame/getCharacterInstancesForGameInstance")

	characterInstanceRecs, err := p.Domain.GetManyAdventureGameCharacterInstanceRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: adventure_game_record.FieldAdventureGameCharacterInstanceGameInstanceID,
					Val: gameInstanceRec.ID,
				},
			},
		},
	)
	if err != nil {
		l.Error("failed to get character instances error >%v<", err)
		return nil, err
	}

	return characterInstanceRecs, nil
}
