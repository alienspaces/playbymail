package turn_sheet_record

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// PlayerTurnSheetResult
const (
	TablePlayerTurnSheetResult string = "player_turn_sheet_result"
)

const (
	FieldPlayerTurnSheetResultID               string = "id"
	FieldPlayerTurnSheetResultTurnSheetID      string = "turn_sheet_id"
	FieldPlayerTurnSheetResultResultData       string = "result_data"
	FieldPlayerTurnSheetResultScannedAt        string = "scanned_at"
	FieldPlayerTurnSheetResultScannedBy        string = "scanned_by"
	FieldPlayerTurnSheetResultScanQuality      string = "scan_quality"
	FieldPlayerTurnSheetResultProcessingStatus string = "processing_status"
	FieldPlayerTurnSheetResultErrorMessage     string = "error_message"
	FieldPlayerTurnSheetResultProcessedAt      string = "processed_at"
	FieldPlayerTurnSheetResultCreatedAt        string = "created_at"
	FieldPlayerTurnSheetResultUpdatedAt        string = "updated_at"
	FieldPlayerTurnSheetResultDeletedAt        string = "deleted_at"
)

type PlayerTurnSheetResult struct {
	record.Record
	TurnSheetID      string          `db:"turn_sheet_id"`
	ResultData       json.RawMessage `db:"result_data"`
	ScannedAt        sql.NullTime    `db:"scanned_at"`
	ScannedBy        sql.NullString  `db:"scanned_by"`
	ScanQuality      sql.NullFloat64 `db:"scan_quality"`
	ProcessingStatus string          `db:"processing_status"`
	ErrorMessage     sql.NullString  `db:"error_message"`
	ProcessedAt      sql.NullTime    `db:"processed_at"`
	CreatedAt        sql.NullTime    `db:"created_at"`
	UpdatedAt        sql.NullTime    `db:"updated_at"`
	DeletedAt        sql.NullTime    `db:"deleted_at"`
}

func (r *PlayerTurnSheetResult) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldPlayerTurnSheetResultTurnSheetID] = r.TurnSheetID
	args[FieldPlayerTurnSheetResultResultData] = r.ResultData
	args[FieldPlayerTurnSheetResultScannedAt] = r.ScannedAt
	args[FieldPlayerTurnSheetResultScannedBy] = r.ScannedBy
	args[FieldPlayerTurnSheetResultScanQuality] = r.ScanQuality
	args[FieldPlayerTurnSheetResultProcessingStatus] = r.ProcessingStatus
	args[FieldPlayerTurnSheetResultErrorMessage] = r.ErrorMessage
	args[FieldPlayerTurnSheetResultProcessedAt] = r.ProcessedAt
	args[FieldPlayerTurnSheetResultCreatedAt] = r.CreatedAt
	args[FieldPlayerTurnSheetResultUpdatedAt] = r.UpdatedAt
	args[FieldPlayerTurnSheetResultDeletedAt] = r.DeletedAt
	return args
}
