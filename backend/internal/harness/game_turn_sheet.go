package harness

import (
	"fmt"
	"time"

	"database/sql"

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

// TODO: Make the following function generic for all game types but add an adventure game specific function for adventure games

func (t *Testing) createGameTurnSheetRec(cfg GameTurnSheetConfig, gameInstanceRec *game_record.GameInstance) (*game_record.GameTurnSheet, error) {
	l := t.Logger("createGameTurnSheetRec")

	if gameInstanceRec == nil {
		return nil, fmt.Errorf("game instance record is nil for game_turn_sheet record >%#v<", cfg)
	}

	if cfg.AccountRef == "" {
		return nil, fmt.Errorf("game_turn_sheet record must have a AccountRef")
	}

	var rec *game_record.GameTurnSheet
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameTurnSheet{}
	}

	rec = t.applyGameTurnSheetRecDefaultValues(rec)

	// Get the account record
	accountRec, err := t.Data.GetAccountRecByRef(cfg.AccountRef)
	if err != nil {
		l.Warn("failed resolving account ref >%s<: %v", cfg.AccountRef, err)
		return nil, err
	}

	rec.AccountID = accountRec.ID
	rec.GameID = gameInstanceRec.GameID
	rec.GameInstanceID = gameInstanceRec.ID
	rec.TurnNumber = cfg.TurnNumber
	rec.SheetType = cfg.SheetType
	rec.SheetOrder = cfg.SheetOrder
	rec.ProcessingStatus = cfg.ProcessingStatus
	rec.IsCompleted = cfg.IsCompleted

	// Set sheet data if provided
	if cfg.SheetData != "" {
		rec.SheetData = []byte(cfg.SheetData)
	}

	// Set scanned data if provided
	if cfg.ScannedData != "" {
		rec.ScannedData = []byte(cfg.ScannedData)
		if cfg.IsCompleted {
			now := time.Now()
			rec.ScannedAt = sql.NullTime{Time: now, Valid: true}
		}
	}

	l.Debug("creating game_turn_sheet record >%#v<", rec)

	// Create turn sheet record
	turnSheetRec, err := t.Domain.(*domain.Domain).CreateGameTurnSheetRec(rec)
	if err != nil {
		l.Warn("failed creating game_turn_sheet record >%v<", err)
		return nil, err
	}

	// Add to data store
	t.Data.AddGameTurnSheetRec(turnSheetRec)

	// Add to teardown data store
	t.teardownData.AddGameTurnSheetRec(turnSheetRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.GameTurnSheetRefs[cfg.Reference] = turnSheetRec.ID
	}

	return turnSheetRec, nil
}

func (t *Testing) applyGameTurnSheetRecDefaultValues(rec *game_record.GameTurnSheet) *game_record.GameTurnSheet {
	if rec == nil {
		rec = &game_record.GameTurnSheet{}
	}

	// Set default processing status if not already set
	if rec.ProcessingStatus == "" {
		rec.ProcessingStatus = game_record.TurnSheetProcessingStatusPending
	}

	return rec
}
