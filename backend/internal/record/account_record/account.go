package account_record

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
	FieldAccountVerificationToken          string = "verification_token"
	FieldAccountVerificationTokenExpiresAt string = "verification_token_expires_at"
	FieldAccountSessionToken               string = "session_token"
	FieldAccountSessionTokenExpiresAt      string = "session_token_expires_at"
	FieldAccountStatus                     string = "status"
)

const (
	AccountStatusPendingApproval string = "pending_approval"
	AccountStatusActive          string = "active"
	AccountStatusDisabled        string = "disabled"
)

type Account struct {
	record.Record
	Email                      string         `db:"email"`
	VerificationToken          sql.NullString `db:"verification_token"`
	VerificationTokenExpiresAt sql.NullTime   `db:"verification_token_expires_at"`
	SessionToken               sql.NullString `db:"session_token"`
	SessionTokenExpiresAt      sql.NullTime   `db:"session_token_expires_at"`
	Status                     string         `db:"status"`
}

func (r *Account) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAccountEmail] = r.Email
	args[FieldAccountVerificationToken] = r.VerificationToken
	args[FieldAccountVerificationTokenExpiresAt] = r.VerificationTokenExpiresAt
	args[FieldAccountSessionToken] = r.SessionToken
	args[FieldAccountSessionTokenExpiresAt] = r.SessionTokenExpiresAt
	args[FieldAccountStatus] = r.Status
	return args
}
