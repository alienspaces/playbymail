package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

type validateAccountUserContactArgs struct {
	nextRec *account_record.AccountUserContact
	currRec *account_record.AccountUserContact
}

func (m *Domain) populateAccountUserContactValidateArgs(currRec, nextRec *account_record.AccountUserContact) (*validateAccountUserContactArgs, error) {
	args := &validateAccountUserContactArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAccountUserContactRecForCreate(rec *account_record.AccountUserContact) error {
	args, err := m.populateAccountUserContactValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAccountUserContactRecForCreate(args)
}

func (m *Domain) validateAccountUserContactRecForUpdate(currRec, nextRec *account_record.AccountUserContact) error {
	args, err := m.populateAccountUserContactValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAccountUserContactRecForUpdate(args)
}

func (m *Domain) validateAccountUserContactRecForDelete(rec *account_record.AccountUserContact) error {
	args, err := m.populateAccountUserContactValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAccountUserContactRecForDelete(args)
}

func validateAccountUserContactRecForCreate(args *validateAccountUserContactArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField("account_user_id", rec.AccountUserID); err != nil {
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

func validateAccountUserContactRecForUpdate(args *validateAccountUserContactArgs) error {
	nextRec := args.nextRec
	currRec := args.currRec

	if nextRec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if nextRec.AccountUserID != currRec.AccountUserID {
		return coreerror.NewInvalidDataError("account_user_id cannot be updated")
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

func validateAccountUserContactRecForDelete(args *validateAccountUserContactArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	return nil
}
