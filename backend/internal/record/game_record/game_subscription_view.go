package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// Table and field constants
const (
	TableGameSubscriptionView = "game_subscription_view"
)

const (
	FieldGameSubscriptionViewID                   = "id"
	FieldGameSubscriptionViewGameID               = "game_id"
	FieldGameSubscriptionViewAccountID            = "account_id"
	FieldGameSubscriptionViewAccountUserContactID = "account_contact_id"
	FieldGameSubscriptionViewSubscriptionType     = "subscription_type"
	FieldGameSubscriptionViewStatus               = "status"
	FieldGameSubscriptionViewInstanceLimit        = "instance_limit"
	FieldGameSubscriptionViewGameInstanceIDs      = "game_instance_ids"
	FieldGameSubscriptionViewCreatedAt            = "created_at"
	FieldGameSubscriptionViewUpdatedAt            = "updated_at"
	FieldGameSubscriptionViewDeletedAt            = "deleted_at"
)

type GameSubscriptionView struct {
	record.Record
	GameID               string         `db:"game_id"`
	AccountID            string         `db:"account_id"`
	AccountUserContactID sql.NullString `db:"account_contact_id"`
	SubscriptionType     string         `db:"subscription_type"`
	Status               string         `db:"status"`
	InstanceLimit        sql.NullInt32  `db:"instance_limit"`
	GameInstanceIDs      []string       `db:"game_instance_ids"`
}

func (r *GameSubscriptionView) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameSubscriptionViewGameID] = r.GameID
	args[FieldGameSubscriptionViewAccountID] = r.AccountID
	args[FieldGameSubscriptionViewAccountUserContactID] = r.AccountUserContactID
	args[FieldGameSubscriptionViewSubscriptionType] = r.SubscriptionType
	args[FieldGameSubscriptionViewStatus] = r.Status
	args[FieldGameSubscriptionViewInstanceLimit] = r.InstanceLimit
	args[FieldGameSubscriptionViewGameInstanceIDs] = r.GameInstanceIDs
	return args
}
