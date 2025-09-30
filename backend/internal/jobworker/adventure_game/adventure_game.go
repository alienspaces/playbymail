package adventure_game

import (
	"context"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// AdventureGame coordinates turn processing for adventure games
type AdventureGame struct {
	Logger     logger.Logger
	Domain     *domain.Domain
	Processors map[string]TurnSheetProcessor
}

// TurnSheetProcessor defines the interface for processing individual turn sheets in adventure games
type TurnSheetProcessor interface {
	// GetSheetType returns the sheet type this processor handles
	GetSheetType() string

	// ProcessTurnSheet processes a single turn sheet of a specific type
	ProcessTurnSheet(ctx context.Context, turnSheet *game_record.GameTurnSheet) error

	// GenerateTurnSheet generates a turn sheet of a specific type for a character
	GenerateTurnSheet(ctx context.Context, characterInstance *adventure_game_record.AdventureGameCharacterInstance) (*game_record.GameTurnSheet, error)
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

// getCharacterInstancesForGameInstance retrieves all character instances for a game instance
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

// initializeTurnSheetProcessors creates and registers all available adventure game turn sheet processors
func (p *AdventureGame) initializeTurnSheetProcessors() (map[string]TurnSheetProcessor, error) {
	l := p.Logger.WithFunctionContext("AdventureGame/initializeTurnSheetProcessors")

	processors := make(map[string]TurnSheetProcessor)

	// Register location choice processor
	locationChoiceProcessor, err := NewAdventureGameLocationChoiceProcessor(l, p.Domain)
	if err != nil {
		l.Warn("failed to initialize location choice processor >%v<", err)
		return nil, err
	}
	processors[adventure_game_record.AdventureSheetTypeLocationChoice] = locationChoiceProcessor

	// Future turn sheet types can be registered here
	// processors[adventure_game_record.AdventureSheetTypeCombat] = combatProcessor
	// processors[adventure_game_record.AdventureSheetTypeInventory] = inventoryProcessor
	// processors[adventure_game_record.AdventureSheetTypeDialogue] = dialogueProcessor

	return processors, nil
}
