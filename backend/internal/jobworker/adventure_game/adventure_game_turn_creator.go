package adventure_game

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// CreateTurnSheets creates all turn sheet records for the current turn of an adventure game instance
func (p *AdventureGame) CreateTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance) error {
	l := p.Logger.WithFunctionContext("AdventureGame/CreateTurnSheets")

	l.Info("creating adventure game turn sheets for instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	// Get all character instances for this game instance
	characterInstanceRecs, err := p.getCharacterInstancesForGameInstance(ctx, gameInstanceRec)
	if err != nil {
		l.Error("failed to get character instances for game instance >%s< error >%v<", gameInstanceRec.ID, err)
		return err
	}

	l.Info("found >%d< character instances for game instance >%s<", len(characterInstanceRecs), gameInstanceRec.ID)

	if len(characterInstanceRecs) == 0 {
		l.Info("no character instances found for game instance >%s<", gameInstanceRec.ID)
		return nil
	}

	// TODO: Wrap this all in a new database transaction so we can roll back
	// changes for a single character if anything fails.

	// Process turn sheets for each character
	var errs []error
	for _, characterInstanceRec := range characterInstanceRecs {
		err := p.createCharacterTurnSheets(ctx, gameInstanceRec, characterInstanceRec)
		if err != nil {
			l.Warn("failed to process turn sheets for character >%s< error >%v<", characterInstanceRec.ID, err)
			// Continue processing other characters even if one fails
			errs = append(errs, err)
			continue
		}
	}

	// If there were any errors we cannot continue processing the turn
	// until the errors are resolved.
	if len(errs) > 0 {
		l.Warn("failed to process turn sheets for some characters error >%v<", errs)
		return fmt.Errorf("failed to process turn sheets for some characters error >%v<", errs)
	}

	return nil
}

// createCharacterTurnSheets creates all of the current game turn's turn sheets for a character
func (p *AdventureGame) createCharacterTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstance *adventure_game_record.AdventureGameCharacterInstance) error {
	l := p.Logger.WithFunctionContext("AdventureGame/processCharacterTurnSheets")

	l.Info("creating turn sheets for game instance ID >%s< character instance ID >%s< turn number >%d<", gameInstanceRec.ID, characterInstance.ID, gameInstanceRec.CurrentTurn)

	// For each turn sheet type supported by the game create a turn sheet for this character
	for _, turnSheetType := range adventure_game_record.AdventureGameSheetTypes.ToSlice() {
		// Create a turn sheet for this character
		turnSheetRec, err := p.createTurnSheet(ctx, gameInstanceRec, characterInstance, turnSheetType)
		if err != nil {
			l.Warn("failed to create turn sheet >%s< for character >%s< error >%v<", turnSheetType, characterInstance.ID, err)
			return err
		}

		l.Info("created turn sheet >%s< for character instance ID >%s< turn sheet type >%s< turn number >%d<", turnSheetRec.ID, characterInstance.ID, turnSheetType, gameInstanceRec.CurrentTurn)
	}

	return nil
}

// createTurnSheet creates a single turn sheet for a character
func (p *AdventureGame) createTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstance *adventure_game_record.AdventureGameCharacterInstance, turnSheetType string) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("AdventureGame/createTurnSheet")

	l.Info("creating turn sheet type >%s< for game instance ID >%s< character instance ID >%s<", turnSheetType, gameInstanceRec.ID, characterInstance.ID)

	// Get the appropriate processor for this sheet type
	processor, exists := p.Processors[turnSheetType]
	if !exists {
		l.Warn("unsupported sheet type >%s< for character >%s<", turnSheetType, characterInstance.ID)
		return nil, fmt.Errorf("unsupported sheet type: %s", turnSheetType)
	}

	// Create next turn sheet using the sheet-specific processor
	return processor.CreateNextTurnSheet(ctx, gameInstanceRec, characterInstance)
}
