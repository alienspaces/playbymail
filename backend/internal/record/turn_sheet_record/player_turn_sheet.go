package turn_sheet_record

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// PlayerTurnSheet
const (
	TablePlayerTurnSheet string = "player_turn_sheet"
)

const (
	FieldPlayerTurnSheetID             string = "id"
	FieldPlayerTurnSheetGameInstanceID string = "game_instance_id"
	FieldPlayerTurnSheetPlayerID       string = "player_id"
	FieldPlayerTurnSheetTurnNumber     string = "turn_number"
	FieldPlayerTurnSheetSheetType      string = "sheet_type"
	FieldPlayerTurnSheetSheetOrder     string = "sheet_order"
	FieldPlayerTurnSheetSheetData      string = "sheet_data"
	FieldPlayerTurnSheetIsCompleted    string = "is_completed"
	FieldPlayerTurnSheetCompletedAt    string = "completed_at"
	FieldPlayerTurnSheetCreatedAt      string = "created_at"
	FieldPlayerTurnSheetUpdatedAt      string = "updated_at"
	FieldPlayerTurnSheetDeletedAt      string = "deleted_at"
)

type PlayerTurnSheet struct {
	record.Record
	GameInstanceID string          `db:"game_instance_id"`
	PlayerID       string          `db:"player_id"`
	TurnNumber     int             `db:"turn_number"`
	SheetType      string          `db:"sheet_type"`
	SheetOrder     int             `db:"sheet_order"`
	SheetData      json.RawMessage `db:"sheet_data"`
	IsCompleted    bool            `db:"is_completed"`
	CompletedAt    sql.NullTime    `db:"completed_at"`
	CreatedAt      sql.NullTime    `db:"created_at"`
	UpdatedAt      sql.NullTime    `db:"updated_at"`
	DeletedAt      sql.NullTime    `db:"deleted_at"`
}

func (r *PlayerTurnSheet) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldPlayerTurnSheetGameInstanceID] = r.GameInstanceID
	args[FieldPlayerTurnSheetPlayerID] = r.PlayerID
	args[FieldPlayerTurnSheetTurnNumber] = r.TurnNumber
	args[FieldPlayerTurnSheetSheetType] = r.SheetType
	args[FieldPlayerTurnSheetSheetOrder] = r.SheetOrder
	args[FieldPlayerTurnSheetSheetData] = r.SheetData
	args[FieldPlayerTurnSheetIsCompleted] = r.IsCompleted
	args[FieldPlayerTurnSheetCompletedAt] = r.CompletedAt
	args[FieldPlayerTurnSheetCreatedAt] = r.CreatedAt
	args[FieldPlayerTurnSheetUpdatedAt] = r.UpdatedAt
	args[FieldPlayerTurnSheetDeletedAt] = r.DeletedAt
	return args
}
