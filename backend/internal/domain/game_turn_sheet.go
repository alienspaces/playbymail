package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// GetGameTurnSheetRec retrieves a game turn sheet record by ID
func (m *Domain) GetGameTurnSheetRec(recID string, lock *coresql.Lock) (*game_record.GameTurnSheet, error) {
	l := m.Logger("GetGameTurnSheetRec")

	l.Debug("getting game_turn_sheet record ID >%s<", recID)

	if err := domain.ValidateUUIDField("id", recID); err != nil {
		return nil, err
	}

	r := m.GameTurnSheetRepository()

	rec, err := r.GetOne(recID, lock)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, coreerror.NewNotFoundError(game_record.TableGameTurnSheet, recID)
	} else if err != nil {
		return nil, databaseError(err)
	}

	return rec, nil
}

// CreateGameTurnSheetRec creates a new game turn sheet record
func (m *Domain) CreateGameTurnSheetRec(rec *game_record.GameTurnSheet) (*game_record.GameTurnSheet, error) {
	l := m.Logger("CreateGameTurnSheetRec")

	l.Debug("creating game_turn_sheet record >%#v<", rec)

	if err := m.validateGameTurnSheetRecForCreate(rec); err != nil {
		l.Warn("failed to validate game_turn_sheet record >%v<", err)
		return rec, err
	}

	r := m.GameTurnSheetRepository()

	rec, err := r.CreateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// UpdateGameTurnSheetRec updates an existing game turn sheet record
func (m *Domain) UpdateGameTurnSheetRec(rec *game_record.GameTurnSheet) (*game_record.GameTurnSheet, error) {
	l := m.Logger("UpdateGameTurnSheetRec")

	l.Debug("updating game_turn_sheet record >%#v<", rec)

	currRec, err := m.GetGameTurnSheetRec(rec.ID, coresql.ForUpdateNoWait)
	if err != nil {
		return rec, err
	}

	l.Debug("updating game_turn_sheet record >%#v<", rec)

	if err := m.validateGameTurnSheetRecForUpdate(currRec, rec); err != nil {
		l.Warn("failed to validate game_turn_sheet record >%v<", err)
		return rec, err
	}

	r := m.GameTurnSheetRepository()

	rec, err = r.UpdateOne(rec)
	if err != nil {
		return rec, databaseError(err)
	}

	return rec, nil
}

// DeleteGameTurnSheetRec deletes a game turn sheet record
func (m *Domain) DeleteGameTurnSheetRec(recID string) error {
	l := m.Logger("DeleteGameTurnSheetRec")

	l.Debug("deleting game_turn_sheet record ID >%s<", recID)

	_, err := m.GetGameTurnSheetRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameTurnSheetRepository()

	if err := r.DeleteOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// RemoveGameTurnSheetRec removes a game turn sheet record
func (m *Domain) RemoveGameTurnSheetRec(recID string) error {
	l := m.Logger("RemoveGameTurnSheetRec")

	l.Debug("removing game_turn_sheet record ID >%s<", recID)

	_, err := m.GetGameTurnSheetRec(recID, coresql.ForUpdateNoWait)
	if err != nil {
		return err
	}

	r := m.GameTurnSheetRepository()

	if err := r.RemoveOne(recID); err != nil {
		return databaseError(err)
	}

	return nil
}

// GetGameTurnSheetRecsByGameInstance retrieves all turn sheets for a game instance
func (m *Domain) GetGameTurnSheetRecsByGameInstance(gameInstanceID string, turnNumber int) ([]*game_record.GameTurnSheet, error) {
	l := m.Logger("GetGameTurnSheetRecsByGameInstance")

	l.Debug("getting game_turn_sheet records for game_instance_id >%s< turn_number >%d<", gameInstanceID, turnNumber)

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
		return nil, databaseError(err)
	}

	return recs, nil
}

// GetGameTurnSheetRecsByAccount retrieves all turn sheets for an account
func (m *Domain) GetGameTurnSheetRecsByAccount(accountID string) ([]*game_record.GameTurnSheet, error) {
	l := m.Logger("GetGameTurnSheetRecsByAccount")

	l.Debug("getting game_turn_sheet records for account_id >%s<", accountID)

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
		return nil, databaseError(err)
	}

	return recs, nil
}

// MarkGameTurnSheetAsScanned marks a turn sheet as scanned
func (m *Domain) MarkGameTurnSheetAsScanned(turnSheetID string, scannedBy string) error {
	l := m.Logger("MarkGameTurnSheetAsScanned")

	l.Debug("marking game_turn_sheet record ID >%s< as scanned", turnSheetID)

	rec, err := m.GetGameTurnSheetRec(turnSheetID, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	rec.ScannedAt = sql.NullTime{Time: now, Valid: true}
	rec.ScannedBy = sql.NullString{String: scannedBy, Valid: true}
	rec.ProcessingStatus = "scanned"

	_, err = m.UpdateGameTurnSheetRec(rec)
	return err
}

// MarkGameTurnSheetAsCompleted marks a turn sheet as completed with result data
func (m *Domain) MarkGameTurnSheetAsCompleted(turnSheetID string, scannedData []byte) error {
	l := m.Logger("MarkGameTurnSheetAsCompleted")

	l.Debug("marking game_turn_sheet record ID >%s< as completed", turnSheetID)

	rec, err := m.GetGameTurnSheetRec(turnSheetID, nil)
	if err != nil {
		return err
	}

	now := time.Now()
	rec.IsCompleted = true
	rec.CompletedAt = sql.NullTime{Time: now, Valid: true}
	rec.ScannedData = scannedData
	rec.ProcessingStatus = "completed"

	_, err = m.UpdateGameTurnSheetRec(rec)
	return err
}

// MarkGameTurnSheetAsError marks a turn sheet as having an error
func (m *Domain) MarkGameTurnSheetAsError(turnSheetID string, errorMessage string) error {
	l := m.Logger("MarkGameTurnSheetAsError")

	l.Debug("marking game_turn_sheet record ID >%s< as error", turnSheetID)

	rec, err := m.GetGameTurnSheetRec(turnSheetID, nil)
	if err != nil {
		return err
	}

	rec.ProcessingStatus = "error"
	rec.ErrorMessage = sql.NullString{String: errorMessage, Valid: true}

	_, err = m.UpdateGameTurnSheetRec(rec)
	return err
}
