package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

func (m *Domain) validateAccountContactRecForCreate(rec *account_record.AccountContact) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField("account_id", rec.AccountID); err != nil {
		return err
	}

	if rec.Name == "" {
		return coreerror.NewInvalidDataError("name is required")
	}

	if rec.PostalAddressLine1 == "" {
		return coreerror.NewInvalidDataError("postal_address_line1 is required")
	}

	if rec.StateProvince == "" {
		return coreerror.NewInvalidDataError("state_province is required")
	}

	if rec.Country == "" {
		return coreerror.NewInvalidDataError("country is required")
	}

	if rec.PostalCode == "" {
		return coreerror.NewInvalidDataError("postal_code is required")
	}

	return nil
}

func (m *Domain) validateAccountContactRecForUpdate(nextRec, currRec *account_record.AccountContact) error {
	if nextRec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if nextRec.AccountID != currRec.AccountID {
		return coreerror.NewInvalidDataError("account_id cannot be updated")
	}

	if nextRec.Name == "" {
		return coreerror.NewInvalidDataError("name is required")
	}

	if nextRec.PostalAddressLine1 == "" {
		return coreerror.NewInvalidDataError("postal_address_line1 is required")
	}

	if nextRec.StateProvince == "" {
		return coreerror.NewInvalidDataError("state_province is required")
	}

	if nextRec.Country == "" {
		return coreerror.NewInvalidDataError("country is required")
	}

	if nextRec.PostalCode == "" {
		return coreerror.NewInvalidDataError("postal_code is required")
	}

	return nil
}

func (m *Domain) validateAccountContactRecForDelete(rec *account_record.AccountContact) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	return nil
}
