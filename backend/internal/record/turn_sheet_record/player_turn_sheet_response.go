package turn_sheet_record

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// PlayerTurnSheetResponse
const (
	TablePlayerTurnSheetResponse string = "player_turn_sheet_response"
)

const (
	FieldPlayerTurnSheetResponseID           string = "id"
	FieldPlayerTurnSheetResponseSheetID      string = "sheet_id"
	FieldPlayerTurnSheetResponseResponseData string = "response_data"
	FieldPlayerTurnSheetResponseCreatedAt    string = "created_at"
	FieldPlayerTurnSheetResponseUpdatedAt    string = "updated_at"
	FieldPlayerTurnSheetResponseDeletedAt    string = "deleted_at"
)

type PlayerTurnSheetResponse struct {
	record.Record
	SheetID      string          `db:"sheet_id"`
	ResponseData json.RawMessage `db:"response_data"`
	CreatedAt    sql.NullTime    `db:"created_at"`
	UpdatedAt    sql.NullTime    `db:"updated_at"`
	DeletedAt    sql.NullTime    `db:"deleted_at"`
}

func (r *PlayerTurnSheetResponse) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldPlayerTurnSheetResponseSheetID] = r.SheetID
	args[FieldPlayerTurnSheetResponseResponseData] = r.ResponseData
	args[FieldPlayerTurnSheetResponseCreatedAt] = r.CreatedAt
	args[FieldPlayerTurnSheetResponseUpdatedAt] = r.UpdatedAt
	args[FieldPlayerTurnSheetResponseDeletedAt] = r.DeletedAt
	return args
}
