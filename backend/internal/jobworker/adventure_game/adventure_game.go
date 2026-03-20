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

// AdventureGame is the turn sheet processor for adventure games.
//   - Only existing character instances are processed during turn processing.
//   - New player onboarding (join game) is handled by AdventureGameJoinGameProcessor
//     in process_subscription.go.
//
// # Extension guide for new game types
//
// File layout (mirror this package's structure):
//   - One file per turn sheet processor in turn_sheet_processor/ (e.g. adventure_game_*_processor.go)
//   - One file per effect subsystem (e.g. adventure_game_*_effect_processor.go)
//   - Shared query helpers go in turn_sheet_processor/helpers.go
//   - Package-level constants go in turn_sheet_processor/constants.go and adventure_game/constants.go
//
// Registering a new processor:
//   - Implement the TurnSheetProcessor interface (ProcessTurnSheetResponse + CreateNextTurnSheet)
//   - Register the new processor in initializeTurnSheetProcessors() below
//
// Null type construction:
//   - Use core/nullstring, core/nullint64 etc. for all null-type construction
//   - Avoid raw sql.Null* struct literals with explicit field assignments; use the core helpers
//
// Function/method argument order (enforced throughout this package):
//  1. context.Context  — only on interface method boundaries
//  2. logger.Logger    — all package-level helpers and unexported sub-functions
//  3. *domain.Domain   — all package-level helpers and unexported sub-functions
//  4. remaining domain arguments

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
	processors[adventure_game_record.AdventureGameTurnSheetTypeLocationChoice] = locationChoiceProcessor

	// Register inventory management processor
	inventoryManagementProcessor, err := turn_sheet_processor.NewAdventureGameInventoryManagementProcessor(l, p.Domain)
	if err != nil {
		l.Warn("failed to initialize inventory management processor >%v<", err)
		return nil, err
	}
	processors[adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement] = inventoryManagementProcessor

	// Register monster encounter processor
	creatureEncounterProcessor, err := turn_sheet_processor.NewAdventureGameCreatureEncounterProcessor(l, p.Domain)
	if err != nil {
		l.Warn("failed to initialize creature encounter processor >%v<", err)
		return nil, err
	}
	processors[adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter] = creatureEncounterProcessor

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
