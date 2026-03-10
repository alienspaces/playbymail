package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

type validateAccountArgs struct {
	nextRec *account_record.Account
	currRec *account_record.Account
}

func (m *Domain) populateAccountValidateArgs(currRec, nextRec *account_record.Account) (*validateAccountArgs, error) {
	args := &validateAccountArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAccountRecForCreate(rec *account_record.Account) error {
	args, err := m.populateAccountValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAccountRecForCreate(args)
}

func (m *Domain) validateAccountRecForUpdate(currRec, nextRec *account_record.Account) error {
	args, err := m.populateAccountValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAccountRecForUpdate(args)
}

func (m *Domain) validateAccountRecForDelete(rec *account_record.Account) error {
	args, err := m.populateAccountValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAccountRecForDelete(args)
}

func validateAccountRecForCreate(args *validateAccountArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if rec.Status == "" {
		return coreerror.NewInvalidDataError("status is required")
	}

	if err := validateAccountStatus(rec.Status); err != nil {
		return err
	}

	return nil
}

func validateAccountRecForUpdate(args *validateAccountArgs) error {
	nextRec := args.nextRec

	if nextRec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if nextRec.Status == "" {
		return coreerror.NewInvalidDataError("status is required")
	}

	if err := validateAccountStatus(nextRec.Status); err != nil {
		return err
	}

	return nil
}

func validateAccountRecForDelete(args *validateAccountArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	return nil
}

func validateAccountStatus(status string) error {
	switch status {
	case account_record.AccountStatusActive,
		account_record.AccountStatusDisabled:
		return nil
	default:
		return coreerror.NewInvalidDataError("invalid account status >%s<", status)
	}
}
