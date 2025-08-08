package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameInstance
const (
	TableGameInstance string = "game_instance"
)

const (
	FieldGameInstanceID                  string = "id"
	FieldGameInstanceGameID              string = "game_id"
	FieldGameInstanceStatus              string = "status"
	FieldGameInstanceCurrentTurn         string = "current_turn"
	FieldGameInstanceLastTurnProcessedAt string = "last_turn_processed_at"
	FieldGameInstanceNextTurnDueAt       string = "next_turn_due_at"
	FieldGameInstanceStartedAt           string = "started_at"
	FieldGameInstanceCompletedAt         string = "completed_at"
	FieldGameInstanceCreatedAt           string = "created_at"
	FieldGameInstanceUpdatedAt           string = "updated_at"
	FieldGameInstanceDeletedAt           string = "deleted_at"
)

// Game instance status constants
const (
	GameInstanceStatusCreated   = "created"
	GameInstanceStatusStarting  = "starting"
	GameInstanceStatusRunning   = "running"
	GameInstanceStatusPaused    = "paused"
	GameInstanceStatusCompleted = "completed"
	GameInstanceStatusCancelled = "cancelled"
)

type GameInstance struct {
	record.Record
	GameID              string       `db:"game_id"`
	Status              string       `db:"status"`
	CurrentTurn         int          `db:"current_turn"`
	LastTurnProcessedAt sql.NullTime `db:"last_turn_processed_at"`
	NextTurnDueAt       sql.NullTime `db:"next_turn_due_at"`
	StartedAt           sql.NullTime `db:"started_at"`
	CompletedAt         sql.NullTime `db:"completed_at"`
}

func (r *GameInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameInstanceGameID] = r.GameID
	args[FieldGameInstanceStatus] = r.Status
	args[FieldGameInstanceCurrentTurn] = r.CurrentTurn
	args[FieldGameInstanceLastTurnProcessedAt] = r.LastTurnProcessedAt
	args[FieldGameInstanceNextTurnDueAt] = r.NextTurnDueAt
	args[FieldGameInstanceStartedAt] = r.StartedAt
	args[FieldGameInstanceCompletedAt] = r.CompletedAt
	return args
}
