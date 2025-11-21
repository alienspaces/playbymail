package account_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableAccountContact string = "account_contact"
)

const (
	FieldAccountContactID                 string = "id"
	FieldAccountContactAccountID          string = "account_id"
	FieldAccountContactName               string = "name"
	FieldAccountContactPostalAddressLine1 string = "postal_address_line1"
	FieldAccountContactPostalAddressLine2 string = "postal_address_line2"
	FieldAccountContactStateProvince      string = "state_province"
	FieldAccountContactCountry            string = "country"
	FieldAccountContactPostalCode         string = "postal_code"
	FieldAccountContactCreatedAt          string = "created_at"
	FieldAccountContactUpdatedAt          string = "updated_at"
	FieldAccountContactDeletedAt          string = "deleted_at"
)

type AccountContact struct {
	record.Record
	AccountID          string         `db:"account_id"`
	Name               string         `db:"name"`
	PostalAddressLine1 string         `db:"postal_address_line1"`
	PostalAddressLine2 sql.NullString `db:"postal_address_line2"`
	StateProvince      string         `db:"state_province"`
	Country            string         `db:"country"`
	PostalCode         string         `db:"postal_code"`
}

func (r *AccountContact) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAccountContactAccountID] = r.AccountID
	args[FieldAccountContactName] = r.Name
	args[FieldAccountContactPostalAddressLine1] = r.PostalAddressLine1
	args[FieldAccountContactPostalAddressLine2] = r.PostalAddressLine2
	args[FieldAccountContactStateProvince] = r.StateProvince
	args[FieldAccountContactCountry] = r.Country
	args[FieldAccountContactPostalCode] = r.PostalCode
	return args
}
