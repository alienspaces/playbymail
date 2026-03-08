package game_record

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	TableAccountGameView = "account_game_view"
)

const (
	FieldAccountGameViewID               = "id"
	FieldAccountGameViewAccountID        = "account_id"
	FieldAccountGameViewAccountName      = "account_name"
	FieldAccountGameViewGameID           = "game_id"
	FieldAccountGameViewGameName         = "game_name"
	FieldAccountGameViewGameType         = "game_type"
	FieldAccountGameViewDescription      = "description"
	FieldAccountGameViewTurnDurationHrs  = "turn_duration_hours"
	FieldAccountGameViewGameStatus       = "game_status"
	FieldAccountGameViewIsDesigner       = "is_designer"
	FieldAccountGameViewIsManager        = "is_manager"
	FieldAccountGameViewCanManage        = "can_manage"
	FieldAccountGameViewCreatedAt        = "created_at"
	FieldAccountGameViewUpdatedAt        = "updated_at"
	FieldAccountGameViewDeletedAt        = "deleted_at"
)

type AccountGameView struct {
	ID               string       `db:"id"`
	AccountID        string       `db:"account_id"`
	AccountName      string       `db:"account_name"`
	GameID           string       `db:"game_id"`
	GameName         string       `db:"game_name"`
	GameType         string       `db:"game_type"`
	Description      string       `db:"description"`
	TurnDurationHours int         `db:"turn_duration_hours"`
	GameStatus       string       `db:"game_status"`
	IsDesigner       bool         `db:"is_designer"`
	IsManager        bool         `db:"is_manager"`
	CanManage        bool         `db:"can_manage"`
	CreatedAt        time.Time    `db:"created_at"`
	UpdatedAt        sql.NullTime `db:"updated_at"`
	DeletedAt        sql.NullTime `db:"deleted_at"`
}

func (r *AccountGameView) ToNamedArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		FieldAccountGameViewID:              r.ID,
		FieldAccountGameViewAccountID:       r.AccountID,
		FieldAccountGameViewAccountName:     r.AccountName,
		FieldAccountGameViewGameID:          r.GameID,
		FieldAccountGameViewGameName:        r.GameName,
		FieldAccountGameViewGameType:        r.GameType,
		FieldAccountGameViewDescription:     r.Description,
		FieldAccountGameViewTurnDurationHrs: r.TurnDurationHours,
		FieldAccountGameViewGameStatus:      r.GameStatus,
		FieldAccountGameViewIsDesigner:      r.IsDesigner,
		FieldAccountGameViewIsManager:       r.IsManager,
		FieldAccountGameViewCanManage:       r.CanManage,
		FieldAccountGameViewCreatedAt:       r.CreatedAt,
		FieldAccountGameViewUpdatedAt:       r.UpdatedAt,
		FieldAccountGameViewDeletedAt:       r.DeletedAt,
	}
}
