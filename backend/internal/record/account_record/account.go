package account_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// Account (Tenant)
const (
	TableAccount string = "account"
)

const (
	FieldAccountName   string = "name"
	FieldAccountStatus string = "status"
)

type Account struct {
	record.Record
	Name   string `db:"name"`
	Status string `db:"status"`
}

// ToNamedArgs converts the struct to a map of named arguments
// Note: Record.ToNamedArgs() handles common fields (ID, CreatedAt, etc.)
func (r *Account) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAccountName] = r.Name
	args[FieldAccountStatus] = r.Status
	return args
}
