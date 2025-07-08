package record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// Account
const (
	TableAccount string = "account"
)

const (
	FieldAccountID             string = "id"
	FieldAccountEmail          string = "email"
	FieldAccountName           string = "name"
	FieldAccountToken          string = "token"
	FieldAccountTokenExpiresAt string = "token_expires_at"
)

type Account struct {
	record.Record
	Email          string         `db:"email"`
	Name           string         `db:"name"`
	Token          sql.NullString `db:"token"`
	TokenExpiresAt sql.NullTime   `db:"token_expires_at"`
}

func (r *Account) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAccountEmail] = r.Email
	args[FieldAccountName] = r.Name
	args[FieldAccountToken] = r.Token
	args[FieldAccountTokenExpiresAt] = r.TokenExpiresAt
	return args
}
