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
	FieldGameInstanceID                                string = "id"
	FieldGameInstanceGameID                            string = "game_id"
	FieldGameInstanceGameSubscriptionID                string = "game_subscription_id"
	FieldGameInstanceStatus                            string = "status"
	FieldGameInstanceCurrentTurn                       string = "current_turn"
	FieldGameInstanceDeliveryPhysicalPost              string = "delivery_physical_post"
	FieldGameInstanceDeliveryPhysicalLocal             string = "delivery_physical_local"
	FieldGameInstanceDeliveryEmail                     string = "delivery_email"
	FieldGameInstanceRequiredPlayerCount               string = "required_player_count"
	FieldGameInstanceIsClosedTesting                   string = "is_closed_testing"
	FieldGameInstanceClosedTestingJoinGameKey          string = "closed_testing_join_game_key"
	FieldGameInstanceClosedTestingJoinGameKeyExpiresAt string = "closed_testing_join_game_key_expires_at"
	FieldGameInstanceStartedAt                         string = "started_at"
	FieldGameInstanceCompletedAt                       string = "completed_at"
	FieldGameInstanceLastTurnProcessedAt               string = "last_turn_processed_at"
	FieldGameInstanceNextTurnDueAt                     string = "next_turn_due_at"
	FieldGameInstanceCreatedAt                         string = "created_at"
	FieldGameInstanceUpdatedAt                         string = "updated_at"
	FieldGameInstanceDeletedAt                         string = "deleted_at"
)

// Game instance status constants
const (
	GameInstanceStatusCreated   = "created"
	GameInstanceStatusStarted   = "started"
	GameInstanceStatusPaused    = "paused"
	GameInstanceStatusCompleted = "completed"
	GameInstanceStatusCancelled = "cancelled"
)

type GameInstance struct {
	record.Record
	GameID                            string         `db:"game_id"`
	GameSubscriptionID                string         `db:"game_subscription_id"`
	Status                            string         `db:"status"`
	CurrentTurn                       int            `db:"current_turn"`
	LastTurnProcessedAt               sql.NullTime   `db:"last_turn_processed_at"`
	NextTurnDueAt                     sql.NullTime   `db:"next_turn_due_at"`
	StartedAt                         sql.NullTime   `db:"started_at"`
	CompletedAt                       sql.NullTime   `db:"completed_at"`
	DeliveryPhysicalPost              bool           `db:"delivery_physical_post"`
	DeliveryPhysicalLocal             bool           `db:"delivery_physical_local"`
	DeliveryEmail                     bool           `db:"delivery_email"`
	RequiredPlayerCount               int            `db:"required_player_count"`
	IsClosedTesting                   bool           `db:"is_closed_testing"`
	ClosedTestingJoinGameKey          sql.NullString `db:"closed_testing_join_game_key"`
	ClosedTestingJoinGameKeyExpiresAt sql.NullTime   `db:"closed_testing_join_game_key_expires_at"`
}

func (r *GameInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameInstanceGameID] = r.GameID
	args[FieldGameInstanceGameSubscriptionID] = r.GameSubscriptionID
	args[FieldGameInstanceStatus] = r.Status
	args[FieldGameInstanceCurrentTurn] = r.CurrentTurn
	args[FieldGameInstanceLastTurnProcessedAt] = r.LastTurnProcessedAt
	args[FieldGameInstanceNextTurnDueAt] = r.NextTurnDueAt
	args[FieldGameInstanceStartedAt] = r.StartedAt
	args[FieldGameInstanceCompletedAt] = r.CompletedAt
	args[FieldGameInstanceDeliveryPhysicalPost] = r.DeliveryPhysicalPost
	args[FieldGameInstanceDeliveryPhysicalLocal] = r.DeliveryPhysicalLocal
	args[FieldGameInstanceDeliveryEmail] = r.DeliveryEmail
	args[FieldGameInstanceRequiredPlayerCount] = r.RequiredPlayerCount
	args[FieldGameInstanceIsClosedTesting] = r.IsClosedTesting
	args[FieldGameInstanceClosedTestingJoinGameKey] = r.ClosedTestingJoinGameKey
	args[FieldGameInstanceClosedTestingJoinGameKeyExpiresAt] = r.ClosedTestingJoinGameKeyExpiresAt
	return args
}
