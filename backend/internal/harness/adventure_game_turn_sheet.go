package harness

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createAdventureGameTurnSheetRec(cfg AdventureGameTurnSheetConfig, gameInstanceRec *game_record.GameInstance) (*adventure_game_record.AdventureGameTurnSheet, error) {
	l := t.Logger("createAdventureGameTurnSheetRec")

	if gameInstanceRec == nil {
		return nil, fmt.Errorf("game instance record is nil for adventure game turn sheet record >%#v<", cfg)
	}

	if cfg.GameCharacterInstanceRef == "" {
		return nil, fmt.Errorf("adventure game turn sheet record must have a GameCharacterInstanceRef")
	}

	turnSheetRec, err := t.createGameTurnSheetRec(cfg.GameTurnSheetConfig, gameInstanceRec)
	if err != nil {
		l.Warn("failed creating game_turn_sheet record >%v<", err)
		return nil, err
	}

	// Get the character instance
	characterInstance, err := t.Data.GetAdventureGameCharacterInstanceRecByRef(cfg.GameCharacterInstanceRef)
	if err != nil {
		l.Error("could not resolve GameCharacterInstanceRef >%s< to a valid character instance ID", cfg.GameCharacterInstanceRef)
		return nil, fmt.Errorf("could not resolve GameCharacterInstanceRef >%s< to a valid character instance ID", cfg.GameCharacterInstanceRef)
	}

	// Create adventure game turn sheet record to link the turn sheet to the character instance
	adventureGameTurnSheet := &adventure_game_record.AdventureGameTurnSheet{
		GameID:                           gameInstanceRec.GameID,
		AdventureGameCharacterInstanceID: characterInstance.ID,
		GameTurnSheetID:                  turnSheetRec.ID,
	}

	l.Debug("creating adventure_game_turn_sheet record >%#v<", adventureGameTurnSheet)

	adventureGameTurnSheetRec, err := t.Domain.(*domain.Domain).CreateAdventureGameTurnSheetRec(adventureGameTurnSheet)
	if err != nil {
		l.Warn("failed creating adventure_game_turn_sheet record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddAdventureGameTurnSheetRec(adventureGameTurnSheetRec)

	// Add to teardown data store
	t.teardownData.AddAdventureGameTurnSheetRec(adventureGameTurnSheetRec)

	l.Debug("created adventure_game_turn_sheet record ID >%s<", adventureGameTurnSheetRec.ID)

	return adventureGameTurnSheetRec, nil
}
