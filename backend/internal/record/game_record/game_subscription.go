package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// Table and field constants
const (
	TableGameSubscription = "game_subscription"
)

const (
	FieldGameSubscriptionID                      = "id"
	FieldGameSubscriptionGameID                  = "game_id"
	FieldGameSubscriptionAccountID               = "account_id"
	FieldGameSubscriptionAccountContactID        = "account_contact_id"
	FieldGameSubscriptionSubscriptionType        = "subscription_type"
	FieldGameSubscriptionCreatedAt               = "created_at"
	FieldGameSubscriptionStatus                  = "status"
	FieldGameSubscriptionGameInstanceID          = "game_instance_id"
	FieldGameSubscriptionTurnSheetToken          = "turn_sheet_token"
	FieldGameSubscriptionTurnSheetTokenExpiresAt = "turn_sheet_token_expires_at"
)

const (
	GameSubscriptionTypePlayer   = "Player"
	GameSubscriptionTypeManager  = "Manager"
	GameSubscriptionTypeDesigner = "Designer"
)

const (
	GameSubscriptionStatusPendingApproval = "pending_approval"
	GameSubscriptionStatusActive          = "active"
	GameSubscriptionStatusRevoked         = "revoked"
)

// GameSubscription represents a subscription to a game (Player, Manager, Designer)
type GameSubscription struct {
	record.Record
	GameID                  string         `db:"game_id"`
	AccountID               string         `db:"account_id"`
	AccountContactID        sql.NullString `db:"account_contact_id"`
	SubscriptionType        string         `db:"subscription_type"`
	Status                  string         `db:"status"`
	GameInstanceID          sql.NullString `db:"game_instance_id"`
	TurnSheetToken          sql.NullString `db:"turn_sheet_token"`
	TurnSheetTokenExpiresAt sql.NullTime   `db:"turn_sheet_token_expires_at"`
}

func (r *GameSubscription) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameSubscriptionGameID] = r.GameID
	args[FieldGameSubscriptionAccountID] = r.AccountID
	args[FieldGameSubscriptionAccountContactID] = r.AccountContactID
	args[FieldGameSubscriptionSubscriptionType] = r.SubscriptionType
	args[FieldGameSubscriptionStatus] = r.Status
	args[FieldGameSubscriptionGameInstanceID] = r.GameInstanceID
	args[FieldGameSubscriptionTurnSheetToken] = r.TurnSheetToken
	args[FieldGameSubscriptionTurnSheetTokenExpiresAt] = r.TurnSheetTokenExpiresAt
	return args
}
