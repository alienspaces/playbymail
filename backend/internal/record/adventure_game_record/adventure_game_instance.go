package adventure_game_record

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameInstance = "game_instance"

const (
	FieldGameInstanceID                           = "id"
	FieldAdventureGameInstanceGameID              = "game_id"
	FieldAdventureGameInstanceCreatedAt           = "created_at"
	FieldAdventureGameInstanceUpdatedAt           = "updated_at"
	FieldAdventureGameInstanceDeletedAt           = "deleted_at"
	FieldAdventureGameInstanceStatus              = "status"
	FieldAdventureGameInstanceCurrentTurn         = "current_turn"
	FieldAdventureGameInstanceLastTurnProcessedAt = "last_turn_processed_at"
	FieldAdventureGameInstanceNextTurnDueAt       = "next_turn_due_at"
	FieldAdventureGameInstanceStartedAt           = "started_at"
	FieldAdventureGameInstanceCompletedAt         = "completed_at"
	FieldAdventureGameInstanceGameConfig          = "game_config"
)

// Game instance status constants
const (
	GameInstanceStatusCreated   = "created"
	GameInstanceStatusStarted   = "started"
	GameInstanceStatusPaused    = "paused"
	GameInstanceStatusCompleted = "completed"
	GameInstanceStatusCancelled = "cancelled"
)

type AdventureGameInstance struct {
	record.Record
	GameID              string          `db:"game_id"`
	Status              string          `db:"status"`
	CurrentTurn         int             `db:"current_turn"`
	LastTurnProcessedAt *time.Time      `db:"last_turn_processed_at"`
	NextTurnDueAt       *time.Time      `db:"next_turn_due_at"`
	StartedAt           *time.Time      `db:"started_at"`
	CompletedAt         *time.Time      `db:"completed_at"`
	GameConfig          json.RawMessage `db:"game_config"`
}

func (r *AdventureGameInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameInstanceGameID] = r.GameID
	args[FieldAdventureGameInstanceStatus] = r.Status
	args[FieldAdventureGameInstanceCurrentTurn] = r.CurrentTurn
	args[FieldAdventureGameInstanceLastTurnProcessedAt] = r.LastTurnProcessedAt
	args[FieldAdventureGameInstanceNextTurnDueAt] = r.NextTurnDueAt
	args[FieldAdventureGameInstanceStartedAt] = r.StartedAt
	args[FieldAdventureGameInstanceCompletedAt] = r.CompletedAt
	args[FieldAdventureGameInstanceGameConfig] = r.GameConfig
	return args
}
