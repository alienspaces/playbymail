package adventure_game

import (
	"context"
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GenerateTurnSheets generates all turn sheet records for an adventure game turn
func (p *AdventureGame) GenerateTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance) error {
	l := p.Logger.WithFunctionContext("AdventureGame/GenerateTurnSheets")

	l.Info("generating adventure game turn sheets for instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

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
		err := p.generateCharacterTurnSheets(ctx, gameInstanceRec, characterInstanceRec)
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

// generateCharacterTurnSheets processes all turn sheets for a specific character
func (p *AdventureGame) generateCharacterTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstance *adventure_game_record.AdventureGameCharacterInstance) error {
	l := p.Logger.WithFunctionContext("AdventureGame/processCharacterTurnSheets")

	l.Info("generating turn sheets for game instance ID >%s< character instance ID >%s< turn >%d<", gameInstanceRec.ID, characterInstance.ID, gameInstanceRec.CurrentTurn)

	// For each turn sheet type supported by the game generate a turn sheet for this character
	for _, turnSheetType := range adventure_game_record.AdventureGameSheetTypes.ToSlice() {
		// Generate a turn sheet for this character
		turnSheetRec, err := p.generateTurnSheet(ctx, gameInstanceRec, characterInstance, turnSheetType)
		if err != nil {
			l.Warn("failed to generate turn sheet >%s< for character >%s< error >%v<", turnSheetType, characterInstance.ID, err)
			return err
		}

		l.Info("generated turn sheet >%s< for character >%s<", turnSheetRec.ID, characterInstance.ID)
	}

	return nil
}

// generateTurnSheet generates a single turn sheet for a character
func (p *AdventureGame) generateTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstance *adventure_game_record.AdventureGameCharacterInstance, turnSheetType string) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("AdventureGame/generateTurnSheet")

	l.Info("generating turn sheet type >%s< for game instance ID >%s< character instance ID >%s<", turnSheetType, gameInstanceRec.ID, characterInstance.ID)

	// Get the appropriate processor for this sheet type
	processor, exists := p.Processors[turnSheetType]
	if !exists {
		l.Warn("unsupported sheet type >%s< for character >%s<", turnSheetType, characterInstance.ID)
		return nil, fmt.Errorf("unsupported sheet type: %s", turnSheetType)
	}

	// Create next turn sheet using the sheet-specific processor
	return processor.CreateNextTurnSheet(ctx, characterInstance)
}
