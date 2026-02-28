package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

type validateAccountArgs struct {
	nextRec *account_record.AccountUser
	currRec *account_record.AccountUser
}

func (m *Domain) populateAccountValidateArgs(currRec, nextRec *account_record.AccountUser) (*validateAccountArgs, error) {
	args := &validateAccountArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAccountRecForCreate(rec *account_record.AccountUser) error {
	args, err := m.populateAccountValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAccountRecForCreate(args)
}

func (m *Domain) validateAccountRecForUpdate(currRec, nextRec *account_record.AccountUser) error {
	args, err := m.populateAccountValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAccountRecForUpdate(args)
}

func (m *Domain) validateAccountRecForDelete(rec *account_record.AccountUser) error {
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

	if rec.Email == "" {
		return coreerror.NewInvalidDataError("email is required")
	}

	if rec.Status == "" {
		rec.Status = account_record.AccountUserStatusActive
	}

	if err := validateAccountStatus(rec.Status); err != nil {
		return err
	}

	return nil
}

func validateAccountRecForUpdate(args *validateAccountArgs) error {
	nextRec := args.nextRec
	currRec := args.currRec

	if nextRec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if nextRec.Email != currRec.Email {
		return coreerror.NewInvalidDataError("email cannot be updated")
	}

	if nextRec.Status == "" {
		nextRec.Status = currRec.Status
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
	case account_record.AccountUserStatusPendingApproval,
		account_record.AccountUserStatusActive,
		account_record.AccountUserStatusDisabled:
		return nil
	default:
		return coreerror.NewInvalidDataError("invalid account status >%s<", status)
	}
}
