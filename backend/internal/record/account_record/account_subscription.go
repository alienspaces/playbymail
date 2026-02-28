package account_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// AccountSubscription
const (
	TableAccountSubscription string = "account_subscription"
)

const (
	FieldAccountSubscriptionID                 string = "id"
	FieldAccountSubscriptionAccountID          string = "account_id"
	FieldAccountSubscriptionAccountUserID      string = "account_user_id"
	FieldAccountSubscriptionSubscriptionType   string = "subscription_type"
	FieldAccountSubscriptionSubscriptionPeriod string = "subscription_period"
	FieldAccountSubscriptionStatus             string = "status"
	FieldAccountSubscriptionAutoRenew          string = "auto_renew"
	FieldAccountSubscriptionExpiresAt          string = "expires_at"
	FieldAccountSubscriptionCreatedAt          string = "created_at"
	FieldAccountSubscriptionUpdatedAt          string = "updated_at"
	FieldAccountSubscriptionDeletedAt          string = "deleted_at"
)

const (
	// Game designer subscriptions
	// The basic game designer subscription is free and allows the account to create a limited number of games.
	// The professional game designer subscription is paid and allows the account to create unlimited games.
	AccountSubscriptionTypeBasicGameDesigner        string = "basic_game_designer"
	AccountSubscriptionTypeProfessionalGameDesigner string = "professional_game_designer"

	// Manager subscriptions
	// The basic manager subscription is free and allows the account to manage a limited number of games.
	// The professional manager subscription is paid and allows the account to manage unlimited games.
	AccountSubscriptionTypeBasicManager        string = "basic_manager"
	AccountSubscriptionTypeProfessionalManager string = "professional_manager"

	// Player subscriptions
	// The basic player subscription is free and allows the account to play a limited number of games.
	// The professional player subscription is paid and allows the account to play unlimited games.
	AccountSubscriptionTypeBasicPlayer        string = "basic_player"
	AccountSubscriptionTypeProfessionalPlayer string = "professional_player"
)

const (
	AccountSubscriptionStatusActive  string = "active"
	AccountSubscriptionStatusExpired string = "expired"
)

const (
	AccountSubscriptionPeriodMonth   string = "month"
	AccountSubscriptionPeriodYear    string = "year"
	AccountSubscriptionPeriodEternal string = "eternal"
)

type AccountSubscription struct {
	record.Record
	AccountID          sql.NullString `db:"account_id"`
	AccountUserID      sql.NullString `db:"account_user_id"`
	SubscriptionType   string         `db:"subscription_type"`
	SubscriptionPeriod string         `db:"subscription_period"`
	Status             string         `db:"status"`
	AutoRenew          bool           `db:"auto_renew"`
	ExpiresAt          sql.NullTime   `db:"expires_at"`
}

func (r *AccountSubscription) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAccountSubscriptionAccountID] = r.AccountID
	args[FieldAccountSubscriptionAccountUserID] = r.AccountUserID
	args[FieldAccountSubscriptionSubscriptionType] = r.SubscriptionType
	args[FieldAccountSubscriptionSubscriptionPeriod] = r.SubscriptionPeriod
	args[FieldAccountSubscriptionStatus] = r.Status
	args[FieldAccountSubscriptionAutoRenew] = r.AutoRenew
	args[FieldAccountSubscriptionExpiresAt] = r.ExpiresAt
	return args
}
