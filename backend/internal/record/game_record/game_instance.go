package game_record

import (
	"time"

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
	FieldGameInstanceMaxTurns            string = "max_turns"
	FieldGameInstanceTurnDeadlineHours   string = "turn_deadline_hours"
	FieldGameInstanceLastTurnProcessedAt string = "last_turn_processed_at"
	FieldGameInstanceNextTurnDeadline    string = "next_turn_deadline"
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
	GameID              string     `db:"game_id"`
	Status              string     `db:"status"`
	CurrentTurn         int        `db:"current_turn"`
	MaxTurns            *int       `db:"max_turns"`
	TurnDeadlineHours   int        `db:"turn_deadline_hours"`
	LastTurnProcessedAt *time.Time `db:"last_turn_processed_at"`
	NextTurnDeadline    *time.Time `db:"next_turn_deadline"`
	StartedAt           *time.Time `db:"started_at"`
	CompletedAt         *time.Time `db:"completed_at"`
}

func (r *GameInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameInstanceGameID] = r.GameID
	args[FieldGameInstanceStatus] = r.Status
	args[FieldGameInstanceCurrentTurn] = r.CurrentTurn
	args[FieldGameInstanceMaxTurns] = r.MaxTurns
	args[FieldGameInstanceTurnDeadlineHours] = r.TurnDeadlineHours
	args[FieldGameInstanceLastTurnProcessedAt] = r.LastTurnProcessedAt
	args[FieldGameInstanceNextTurnDeadline] = r.NextTurnDeadline
	args[FieldGameInstanceStartedAt] = r.StartedAt
	args[FieldGameInstanceCompletedAt] = r.CompletedAt
	return args
} 