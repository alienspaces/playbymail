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
	args["email"] = r.Email
	args["name"] = r.Name
	args["token"] = r.Token
	args["token_expires_at"] = r.TokenExpiresAt
	return args
}
