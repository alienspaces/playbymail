package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
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

	if err := validatePostalAddressGroup(rec); err != nil {
		return err
	}

	return nil
}

// validatePostalAddressGroup checks that if any postal address field is set, all are required.
func validatePostalAddressGroup(rec *account_record.AccountUserContact) error {
	p1 := nullstring.ToString(rec.PostalAddressLine1)
	sp := nullstring.ToString(rec.StateProvince)
	co := nullstring.ToString(rec.Country)
	pc := nullstring.ToString(rec.PostalCode)
	hasAny := p1 != "" || sp != "" || co != "" || pc != ""
	if !hasAny {
		return nil
	}
	if p1 == "" {
		return coreerror.NewInvalidDataError("postal_address_line1 is required when other postal fields are set")
	}
	if sp == "" {
		return coreerror.NewInvalidDataError("state_province is required when other postal fields are set")
	}
	if co == "" {
		return coreerror.NewInvalidDataError("country is required when other postal fields are set")
	}
	if pc == "" {
		return coreerror.NewInvalidDataError("postal_code is required when other postal fields are set")
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

	if err := validatePostalAddressGroup(nextRec); err != nil {
		return err
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
