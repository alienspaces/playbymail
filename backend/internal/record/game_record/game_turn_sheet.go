package game_record

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameTurnSheet
const (
	TableGameTurnSheet string = "game_turn_sheet"
)

const (
	FieldGameTurnSheetID               string = "id"
	FieldGameTurnSheetGameID           string = "game_id"
	FieldGameTurnSheetGameInstanceID   string = "game_instance_id"
	FieldGameTurnSheetAccountID        string = "account_id"
	FieldGameTurnSheetTurnNumber       string = "turn_number"
	FieldGameTurnSheetSheetType        string = "sheet_type"
	FieldGameTurnSheetSheetOrder       string = "sheet_order"
	FieldGameTurnSheetSheetData        string = "sheet_data"
	FieldGameTurnSheetIsCompleted      string = "is_completed"
	FieldGameTurnSheetCompletedAt      string = "completed_at"
	FieldGameTurnSheetScannedData      string = "scanned_data"
	FieldGameTurnSheetScannedAt        string = "scanned_at"
	FieldGameTurnSheetScannedBy        string = "scanned_by"
	FieldGameTurnSheetScanQuality      string = "scan_quality"
	FieldGameTurnSheetProcessingStatus string = "processing_status"
	FieldGameTurnSheetErrorMessage     string = "error_message"
	FieldGameTurnSheetCreatedAt        string = "created_at"
	FieldGameTurnSheetUpdatedAt        string = "updated_at"
	FieldGameTurnSheetDeletedAt        string = "deleted_at"
)

// Turn sheet processing status constants
// - pending: The turn sheet has not been processed yet
// - processed: The turn sheet has been processed successfully
// - error: The turn sheet has an error
const (
	TurnSheetProcessingStatusPending   string = "pending"
	TurnSheetProcessingStatusProcessed string = "processed"
	TurnSheetProcessingStatusError     string = "error"
)

type GameTurnSheet struct {
	record.Record
	GameID           string          `db:"game_id"`
	GameInstanceID   sql.NullString  `db:"game_instance_id"`
	AccountID        string          `db:"account_id"`
	TurnNumber       int             `db:"turn_number"`
	SheetType        string          `db:"sheet_type"`
	SheetOrder       int             `db:"sheet_order"`
	SheetData        json.RawMessage `db:"sheet_data"`
	IsCompleted      bool            `db:"is_completed"`
	CompletedAt      sql.NullTime    `db:"completed_at"`
	ScannedData      json.RawMessage `db:"scanned_data"`
	ScannedAt        sql.NullTime    `db:"scanned_at"`
	ScannedBy        sql.NullString  `db:"scanned_by"`
	ScanQuality      sql.NullFloat64 `db:"scan_quality"`
	ProcessingStatus string          `db:"processing_status"`
	ErrorMessage     sql.NullString  `db:"error_message"`
}

func (r *GameTurnSheet) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameTurnSheetGameID] = r.GameID
	args[FieldGameTurnSheetGameInstanceID] = r.GameInstanceID
	args[FieldGameTurnSheetAccountID] = r.AccountID
	args[FieldGameTurnSheetTurnNumber] = r.TurnNumber
	args[FieldGameTurnSheetSheetType] = r.SheetType
	args[FieldGameTurnSheetSheetOrder] = r.SheetOrder
	args[FieldGameTurnSheetSheetData] = r.SheetData
	args[FieldGameTurnSheetIsCompleted] = r.IsCompleted
	args[FieldGameTurnSheetCompletedAt] = r.CompletedAt
	args[FieldGameTurnSheetScannedData] = r.ScannedData
	args[FieldGameTurnSheetScannedAt] = r.ScannedAt
	args[FieldGameTurnSheetScannedBy] = r.ScannedBy
	args[FieldGameTurnSheetScanQuality] = r.ScanQuality
	args[FieldGameTurnSheetProcessingStatus] = r.ProcessingStatus
	args[FieldGameTurnSheetErrorMessage] = r.ErrorMessage
	return args
}
