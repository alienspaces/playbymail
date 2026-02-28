package account_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableAccountUserContact string = "account_user_contact"
)

const (
	FieldAccountUserContactID                 string = "id"
	FieldAccountUserContactAccountUserID      string = "account_user_id"
	FieldAccountUserContactName               string = "name"
	FieldAccountUserContactPostalAddressLine1 string = "postal_address_line1"
	FieldAccountUserContactPostalAddressLine2 string = "postal_address_line2"
	FieldAccountUserContactStateProvince      string = "state_province"
	FieldAccountUserContactCountry            string = "country"
	FieldAccountUserContactPostalCode         string = "postal_code"
	FieldAccountUserContactCreatedAt          string = "created_at"
	FieldAccountUserContactUpdatedAt          string = "updated_at"
	FieldAccountUserContactDeletedAt          string = "deleted_at"
)

type AccountUserContact struct {
	record.Record
	AccountUserID      string         `db:"account_user_id"`
	Name               string         `db:"name"`
	PostalAddressLine1 string         `db:"postal_address_line1"`
	PostalAddressLine2 sql.NullString `db:"postal_address_line2"`
	StateProvince      string         `db:"state_province"`
	Country            string         `db:"country"`
	PostalCode         string         `db:"postal_code"`
}

func (r *AccountUserContact) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAccountUserContactAccountUserID] = r.AccountUserID
	args[FieldAccountUserContactName] = r.Name
	args[FieldAccountUserContactPostalAddressLine1] = r.PostalAddressLine1
	args[FieldAccountUserContactPostalAddressLine2] = r.PostalAddressLine2
	args[FieldAccountUserContactStateProvince] = r.StateProvince
	args[FieldAccountUserContactCountry] = r.Country
	args[FieldAccountUserContactPostalCode] = r.PostalCode
	return args
}
