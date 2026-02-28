package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// Table and field constants
const (
	TableGameSubscriptionInstance = "game_subscription_instance"
)

const (
	FieldGameSubscriptionInstanceID                      = "id"
	FieldGameSubscriptionInstanceAccountID               = "account_id"
	FieldGameSubscriptionInstanceGameSubscriptionID      = "game_subscription_id"
	FieldGameSubscriptionInstanceGameInstanceID          = "game_instance_id"
	FieldGameSubscriptionInstanceTurnSheetToken          = "turn_sheet_token"
	FieldGameSubscriptionInstanceTurnSheetTokenExpiresAt = "turn_sheet_token_expires_at"
	FieldGameSubscriptionInstanceCreatedAt               = "created_at"
	FieldGameSubscriptionInstanceUpdatedAt               = "updated_at"
	FieldGameSubscriptionInstanceDeletedAt               = "deleted_at"
)

// GameSubscriptionInstance represents a link between a game subscription and a game instance
type GameSubscriptionInstance struct {
	record.Record
	AccountID               string         `db:"account_id"`
	GameSubscriptionID      string         `db:"game_subscription_id"`
	GameInstanceID          string         `db:"game_instance_id"`
	TurnSheetToken          sql.NullString `db:"turn_sheet_token"`
	TurnSheetTokenExpiresAt sql.NullTime   `db:"turn_sheet_token_expires_at"`
}

func (r *GameSubscriptionInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameSubscriptionInstanceAccountID] = r.AccountID
	args[FieldGameSubscriptionInstanceGameSubscriptionID] = r.GameSubscriptionID
	args[FieldGameSubscriptionInstanceGameInstanceID] = r.GameInstanceID
	args[FieldGameSubscriptionInstanceTurnSheetToken] = r.TurnSheetToken
	args[FieldGameSubscriptionInstanceTurnSheetTokenExpiresAt] = r.TurnSheetTokenExpiresAt
	return args
}
