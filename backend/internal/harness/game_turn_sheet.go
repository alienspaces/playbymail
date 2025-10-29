package harness

import (
	"encoding/json"
	"fmt"
	"time"

	"database/sql"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

func (t *Testing) createGameTurnSheetRec(cfg GameTurnSheetConfig, gameInstanceRec *game_record.GameInstance) (*game_record.GameTurnSheet, error) {
	l := t.Logger("createGameTurnSheetRec")

	if gameInstanceRec == nil {
		return nil, fmt.Errorf("game instance record is nil for game_turn_sheet record >%#v<", cfg)
	}

	if cfg.GameCharacterInstanceRef == "" {
		return nil, fmt.Errorf("game_turn_sheet record must have a GameCharacterInstanceRef")
	}

	// Get the character instance
	characterInstance, err := t.Data.GetGameCharacterInstanceRecByRef(cfg.GameCharacterInstanceRef)
	if err != nil {
		l.Error("could not resolve GameCharacterInstanceRef >%s< to a valid character instance ID", cfg.GameCharacterInstanceRef)
		return nil, fmt.Errorf("could not resolve GameCharacterInstanceRef >%s< to a valid character instance ID", cfg.GameCharacterInstanceRef)
	}

	// Get the account ID from the character
	characterRec, err := t.Data.GetGameCharacterRecByID(characterInstance.AdventureGameCharacterID)
	if err != nil {
		l.Error("could not get character record >%s<", characterInstance.AdventureGameCharacterID)
		return nil, fmt.Errorf("could not get character record >%s<", characterInstance.AdventureGameCharacterID)
	}

	var rec *game_record.GameTurnSheet
	if cfg.Record != nil {
		recCopy := *cfg.Record
		rec = &recCopy
	} else {
		rec = &game_record.GameTurnSheet{}
	}

	rec = t.applyGameTurnSheetRecDefaultValues(rec)

	rec.GameID = gameInstanceRec.GameID
	rec.GameInstanceID = gameInstanceRec.ID
	rec.AccountID = characterRec.AccountID
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
	createdRec, err := t.Domain.(*domain.Domain).CreateGameTurnSheetRec(rec)
	if err != nil {
		l.Warn("failed creating game_turn_sheet record >%v<", err)
		return nil, err
	}

	// Create adventure game turn sheet record to link the turn sheet to the character instance
	adventureGameTurnSheet := &adventure_game_record.AdventureGameTurnSheet{
		GameID:                           gameInstanceRec.GameID,
		AdventureGameCharacterInstanceID: characterInstance.ID,
		GameTurnSheetID:                  createdRec.ID,
	}

	l.Debug("creating adventure_game_turn_sheet record >%#v<", adventureGameTurnSheet)

	createdAdventureTurnSheet, err := t.Domain.(*domain.Domain).CreateAdventureGameTurnSheetRec(adventureGameTurnSheet)
	if err != nil {
		l.Warn("failed creating adventure_game_turn_sheet record >%v<", err)
		return nil, err
	}

	l.Debug("created adventure_game_turn_sheet record ID >%s<", createdAdventureTurnSheet.ID)

	// Add to data store
	t.Data.AddGameTurnSheetRec(createdRec)

	// Add to teardown data store
	t.teardownData.AddGameTurnSheetRec(createdRec)

	// Add to references store
	if cfg.Reference != "" {
		t.Data.Refs.GameTurnSheetRefs[cfg.Reference] = createdRec.ID
	}

	return createdRec, nil
}

func (t *Testing) applyGameTurnSheetRecDefaultValues(rec *game_record.GameTurnSheet) *game_record.GameTurnSheet {
	if rec == nil {
		rec = &game_record.GameTurnSheet{}
	}

	// Set default processing status if not already set
	if rec.ProcessingStatus == "" {
		rec.ProcessingStatus = "pending"
	}

	// Set default sheet data if not already set
	if len(rec.SheetData) == 0 {
		rec.SheetData = json.RawMessage("{}")
	}

	return rec
}
