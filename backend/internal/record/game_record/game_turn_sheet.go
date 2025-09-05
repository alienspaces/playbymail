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
	FieldGameTurnSheetResultData       string = "result_data"
	FieldGameTurnSheetScannedAt        string = "scanned_at"
	FieldGameTurnSheetScannedBy        string = "scanned_by"
	FieldGameTurnSheetScanQuality      string = "scan_quality"
	FieldGameTurnSheetProcessingStatus string = "processing_status"
	FieldGameTurnSheetErrorMessage     string = "error_message"
	FieldGameTurnSheetCreatedAt        string = "created_at"
	FieldGameTurnSheetUpdatedAt        string = "updated_at"
	FieldGameTurnSheetDeletedAt        string = "deleted_at"
)

// Turn sheet type constants for different game types
const (
	// Adventure game sheet types
	AdventureSheetTypeLocationChoice = "location_choice"
	AdventureSheetTypeCombat         = "combat"
	AdventureSheetTypeInventory      = "inventory"
)

type GameTurnSheet struct {
	record.Record
	GameID           string          `db:"game_id"`
	GameInstanceID   string          `db:"game_instance_id"`
	AccountID        string          `db:"account_id"`
	TurnNumber       int             `db:"turn_number"`
	SheetType        string          `db:"sheet_type"`
	SheetOrder       int             `db:"sheet_order"`
	SheetData        json.RawMessage `db:"sheet_data"`
	IsCompleted      bool            `db:"is_completed"`
	CompletedAt      sql.NullTime    `db:"completed_at"`
	ResultData       json.RawMessage `db:"result_data"`
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
	args[FieldGameTurnSheetResultData] = r.ResultData
	args[FieldGameTurnSheetScannedAt] = r.ScannedAt
	args[FieldGameTurnSheetScannedBy] = r.ScannedBy
	args[FieldGameTurnSheetScanQuality] = r.ScanQuality
	args[FieldGameTurnSheetProcessingStatus] = r.ProcessingStatus
	args[FieldGameTurnSheetErrorMessage] = r.ErrorMessage
	return args
}

// GetAdventureGameSheetTypes returns the sheet types for adventure games
func GetAdventureGameSheetTypes() []string {
	return []string{
		AdventureSheetTypeLocationChoice, // Always required for adventure games
		AdventureSheetTypeCombat,         // Optional - only when combat occurs
		AdventureSheetTypeInventory,      // Optional - only when inventory changes
	}
}
