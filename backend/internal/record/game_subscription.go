package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// Table and field constants
const (
	TableGameSubscription = "game_subscription"
)

const (
	FieldGameSubscriptionID               = "id"
	FieldGameSubscriptionGameID           = "game_id"
	FieldGameSubscriptionAccountID        = "account_id"
	FieldGameSubscriptionSubscriptionType = "subscription_type"
	FieldGameSubscriptionCreatedAt        = "created_at"
)

// GameSubscription represents a subscription to a game (Player, Manager, Collaborator)
type GameSubscription struct {
	record.Record
	GameID           string `db:"game_id"`
	AccountID        string `db:"account_id"`
	SubscriptionType string `db:"subscription_type"`
}

func (r *GameSubscription) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameSubscriptionGameID] = r.GameID
	args[FieldGameSubscriptionAccountID] = r.AccountID
	args[FieldGameSubscriptionSubscriptionType] = r.SubscriptionType
	return args
}
