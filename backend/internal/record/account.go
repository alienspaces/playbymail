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
	FieldAccountID                         string = "id"
	FieldAccountEmail                      string = "email"
	FieldAccountName                       string = "name"
	FieldAccountVerificationToken          string = "verification_token"
	FieldAccountVerificationTokenExpiresAt string = "verification_token_expires_at"
	FieldAccountSessionToken               string = "session_token"
	FieldAccountSessionTokenExpiresAt      string = "session_token_expires_at"
)

type Account struct {
	record.Record
	Email                      string         `db:"email"`
	Name                       string         `db:"name"`
	VerificationToken          sql.NullString `db:"verification_token"`
	VerificationTokenExpiresAt sql.NullTime   `db:"verification_token_expires_at"`
	SessionToken               sql.NullString `db:"session_token"`
	SessionTokenExpiresAt      sql.NullTime   `db:"session_token_expires_at"`
}

func (r *Account) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAccountEmail] = r.Email
	args[FieldAccountName] = r.Name
	args[FieldAccountVerificationToken] = r.VerificationToken
	args[FieldAccountVerificationTokenExpiresAt] = r.VerificationTokenExpiresAt
	args[FieldAccountSessionToken] = r.SessionToken
	args[FieldAccountSessionTokenExpiresAt] = r.SessionTokenExpiresAt
	return args
}
