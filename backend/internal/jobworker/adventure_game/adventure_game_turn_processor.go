package adventure_game

import (
	"context"
	"fmt"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// ProcessTurn processes all turn sheets for an adventure game turn
// NOTE: Assumes only existing character instances are present - new players not added during processing
func (p *AdventureGame) ProcessTurnSheets(ctx context.Context, gameInstanceRec *game_record.GameInstance) error {
	l := p.Logger.WithFunctionContext("AdventureGame/ProcessTurn")

	l.Info("processing adventure game turn for instance >%s< turn >%d<", gameInstanceRec.ID, gameInstanceRec.CurrentTurn)

	// Get all character instances for this game instance
	// NOTE: This only retrieves existing character instances - new players are not added during turn processing
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
		err := p.processCharacterTurnSheets(ctx, characterInstanceRec, gameInstanceRec.CurrentTurn)
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

// processCharacterTurnSheets processes all turn sheets for a specific character
func (p *AdventureGame) processCharacterTurnSheets(ctx context.Context, characterInstance *adventure_game_record.AdventureGameCharacterInstance, turnNumber int) error {
	l := p.Logger.WithFunctionContext("AdventureGame/processCharacterTurnSheets")

	l.Info("processing turn sheets for character >%s< turn >%d<", characterInstance.ID, turnNumber)

	// Get turn sheets for this character and turn
	turnSheetRecs, err := p.getTurnSheetsForCharacter(characterInstance, turnNumber)
	if err != nil {
		l.Error("failed to get turn sheets for character >%s< turn >%d< error >%v<", characterInstance.ID, turnNumber, err)
		return err
	}

	l.Info("found >%d< turn sheets for character >%s< turn >%d<", len(turnSheetRecs), characterInstance.ID, turnNumber)

	if len(turnSheetRecs) == 0 {
		l.Info("no turn sheets found for character >%s< turn >%d<", characterInstance.ID, turnNumber)
		return nil
	}

	// Process each turn sheet for this character
	for _, turnSheet := range turnSheetRecs {
		err := p.processTurnSheet(ctx, characterInstance, turnSheet)
		if err != nil {
			l.Warn("failed to process turn sheet >%s< for character >%s< error >%v<", turnSheet.ID, characterInstance.ID, err)
			return err
		}
	}

	return nil
}

// processTurnSheet processes a single turn sheet for a character
func (p *AdventureGame) processTurnSheet(ctx context.Context, characterInstance *adventure_game_record.AdventureGameCharacterInstance, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("AdventureGame/processTurnSheet")

	l.Info("processing turn sheet >%s< type >%s< for character >%s<", turnSheet.ID, turnSheet.SheetType, characterInstance.ID)

	// Get the appropriate processor for this sheet type
	processor, exists := p.Processors[turnSheet.SheetType]
	if !exists {
		l.Warn("unsupported sheet type >%s< for turn sheet >%s<", turnSheet.SheetType, turnSheet.ID)
		return fmt.Errorf("unsupported sheet type: %s", turnSheet.SheetType)
	}

	// Process turn sheet using the sheet-specific processor
	return processor.ProcessTurnSheet(ctx, turnSheet)
}

// getTurnSheetsForCharacter retrieves turn sheets for a specific character and turn
func (p *AdventureGame) getTurnSheetsForCharacter(characterInstance *adventure_game_record.AdventureGameCharacterInstance, turnNumber int) ([]*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("AdventureGame/getTurnSheetsForCharacter")

	adventureGameTurnSheetRecs, err := p.Domain.GetManyAdventureGameTurnSheetRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: adventure_game_record.FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID,
					Val: characterInstance.ID,
				},
			},
		},
	)
	if err != nil {
		l.Error("failed to get adventure game turn sheets for character >%s< turn >%d< error >%v<", characterInstance.ID, turnNumber, err)
		return nil, err
	}

	var turnSheetRecs []*game_record.GameTurnSheet
	for _, adventureGameTurnSheetRec := range adventureGameTurnSheetRecs {
		turnSheetRec, err := p.Domain.GetGameTurnSheetRec(adventureGameTurnSheetRec.GameTurnSheetID, nil)
		if err != nil {
			l.Error("failed to get turn sheet >%s< for adventure game turn sheet >%s< error >%v<", adventureGameTurnSheetRec.GameTurnSheetID, adventureGameTurnSheetRec.ID, err)
			return nil, err
		}
		turnSheetRecs = append(turnSheetRecs, turnSheetRec)
	}

	return turnSheetRecs, nil
}
