package account_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// AccountUser
const (
	TableAccountUser string = "account_user"
)

const (
	FieldAccountUserID                         string = "id"
	FieldAccountUserAccountID                  string = "account_id"
	FieldAccountUserEmail                      string = "email"
	FieldAccountUserVerificationToken          string = "verification_token"
	FieldAccountUserVerificationTokenExpiresAt string = "verification_token_expires_at"
	FieldAccountUserSessionToken               string = "session_token"
	FieldAccountUserSessionTokenExpiresAt      string = "session_token_expires_at"
	FieldAccountUserStatus                     string = "status"
	FieldAccountUserCreatedAt                  string = "created_at"
	FieldAccountUserUpdatedAt                  string = "updated_at"
)

const (
	AccountUserStatusPendingApproval string = "pending_approval"
	AccountUserStatusActive          string = "active"
	AccountUserStatusDisabled        string = "disabled"
)

type AccountUser struct {
	record.Record
	AccountID                  string         `db:"account_id"`
	Email                      string         `db:"email"`
	VerificationToken          sql.NullString `db:"verification_token"`
	VerificationTokenExpiresAt sql.NullTime   `db:"verification_token_expires_at"`
	SessionToken               sql.NullString `db:"session_token"`
	SessionTokenExpiresAt      sql.NullTime   `db:"session_token_expires_at"`
	Status                     string         `db:"status"`
}

func (r *AccountUser) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAccountUserAccountID] = r.AccountID
	args[FieldAccountUserEmail] = r.Email
	args[FieldAccountUserVerificationToken] = r.VerificationToken
	args[FieldAccountUserVerificationTokenExpiresAt] = r.VerificationTokenExpiresAt
	args[FieldAccountUserSessionToken] = r.SessionToken
	args[FieldAccountUserSessionTokenExpiresAt] = r.SessionTokenExpiresAt
	args[FieldAccountUserStatus] = r.Status
	return args
}
