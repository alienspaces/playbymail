package account_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// Account (Tenant)
const (
	TableAccount string = "account"
)

const (
	FieldAccountName     string = "name"
	FieldAccountStatus   string = "status"
	FieldAccountTimezone string = "timezone"
)

const (
	AccountStatusActive   string = "active"
	AccountStatusDisabled string = "disabled"
)

type Account struct {
	record.Record
	Name     string         `db:"name"`
	Status   string         `db:"status"`
	Timezone sql.NullString `db:"timezone"`
}

// ToNamedArgs converts the struct to a map of named arguments
// Note: Record.ToNamedArgs() handles common fields (ID, CreatedAt, etc.)
func (r *Account) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAccountName] = r.Name
	args[FieldAccountStatus] = r.Status
	args[FieldAccountTimezone] = r.Timezone
	return args
}
