package domain

import (
	"database/sql"
	"time"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GetGameTurnSheetRec retrieves a game turn sheet record by ID
func (m *Domain) GetGameTurnSheetRec(recID string, lock *coresql.Lock) (*game_record.GameTurnSheet, error) {
	r := m.GameTurnSheetRepository()
	rec, err := r.GetOne(recID, lock)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// CreateGameTurnSheetRec creates a new game turn sheet record
func (m *Domain) CreateGameTurnSheetRec(rec *game_record.GameTurnSheet) (*game_record.GameTurnSheet, error) {
	r := m.GameTurnSheetRepository()
	rec, err := r.CreateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

// UpdateGameTurnSheetRec updates an existing game turn sheet record
func (m *Domain) UpdateGameTurnSheetRec(rec *game_record.GameTurnSheet) (*game_record.GameTurnSheet, error) {
	r := m.GameTurnSheetRepository()
	rec, err := r.UpdateOne(rec)
	if err != nil {
		return rec, err
	}
	return rec, nil
}

// DeleteGameTurnSheetRec deletes a game turn sheet record
func (m *Domain) DeleteGameTurnSheetRec(recID string) error {
	r := m.GameTurnSheetRepository()
	if err := r.DeleteOne(recID); err != nil {
		return err
	}
	return nil
}

// RemoveGameTurnSheetRec removes a game turn sheet record
func (m *Domain) RemoveGameTurnSheetRec(recID string) error {
	r := m.GameTurnSheetRepository()
	if err := r.RemoveOne(recID); err != nil {
		return err
	}
	return nil
}

// GetGameTurnSheetRecsByGameInstance retrieves all turn sheets for a game instance
func (m *Domain) GetGameTurnSheetRecsByGameInstance(gameInstanceID string, turnNumber int) ([]*game_record.GameTurnSheet, error) {

	r := m.GameTurnSheetRepository()
	recs, err := r.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: game_record.FieldGameTurnSheetGameInstanceID,
				Val: gameInstanceID,
			},
			{
				Col: game_record.FieldGameTurnSheetTurnNumber,
				Val: turnNumber,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// GetGameTurnSheetRecsByAccount retrieves all turn sheets for an account
func (m *Domain) GetGameTurnSheetRecsByAccount(accountID string) ([]*game_record.GameTurnSheet, error) {

	r := m.GameTurnSheetRepository()
	recs, err := r.GetMany(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: game_record.FieldGameTurnSheetAccountID,
				Val: accountID,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return recs, nil
}

// MarkGameTurnSheetAsScanned marks a turn sheet as scanned with quality score
func (m *Domain) MarkGameTurnSheetAsScanned(turnSheetID string, scanQuality float64, scannedBy string) error {

	rec, err := m.GetGameTurnSheetRec(turnSheetID, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	rec.ScannedAt = sql.NullTime{Time: now, Valid: true}
	rec.ScanQuality = sql.NullFloat64{Float64: scanQuality, Valid: true}
	rec.ScannedBy = sql.NullString{String: scannedBy, Valid: true}
	rec.ProcessingStatus = "scanned"

	_, err = m.UpdateGameTurnSheetRec(rec)
	return err
}

// MarkGameTurnSheetAsCompleted marks a turn sheet as completed with result data
func (m *Domain) MarkGameTurnSheetAsCompleted(turnSheetID string, ScannedData []byte) error {
	rec, err := m.GetGameTurnSheetRec(turnSheetID, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	rec.IsCompleted = true
	rec.CompletedAt = sql.NullTime{Time: now, Valid: true}
	rec.ScannedData = ScannedData
	rec.ProcessingStatus = "completed"

	_, err = m.UpdateGameTurnSheetRec(rec)
	return err
}

// MarkGameTurnSheetAsError marks a turn sheet as having an error
func (m *Domain) MarkGameTurnSheetAsError(turnSheetID string, errorMessage string) error {
	rec, err := m.GetGameTurnSheetRec(turnSheetID, nil)
	if err != nil {
		return err
	}

	rec.ProcessingStatus = "error"
	rec.ErrorMessage = sql.NullString{String: errorMessage, Valid: true}

	_, err = m.UpdateGameTurnSheetRec(rec)
	return err
}
