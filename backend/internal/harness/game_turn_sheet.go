package harness

import (
	"fmt"
	"time"

	"database/sql"

	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

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
