package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

type validateAccountUserArgs struct {
	nextRec *account_record.AccountUser
	currRec *account_record.AccountUser
}

func (m *Domain) populateAccountUserValidateArgs(currRec, nextRec *account_record.AccountUser) (*validateAccountUserArgs, error) {
	args := &validateAccountUserArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAccountUserRecForCreate(rec *account_record.AccountUser) error {
	args, err := m.populateAccountUserValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAccountUserRecForCreate(args)
}

func (m *Domain) validateAccountUserRecForUpdate(currRec, nextRec *account_record.AccountUser) error {
	args, err := m.populateAccountUserValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAccountUserRecForUpdate(args)
}

func (m *Domain) validateAccountUserRecForDelete(rec *account_record.AccountUser) error {
	args, err := m.populateAccountUserValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAccountUserRecForDelete(args)
}

func validateAccountUserRecForCreate(args *validateAccountUserArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if err := domain.ValidateUUIDField("account_id", rec.AccountID); err != nil {
		return err
	}

	if rec.Email == "" {
		return coreerror.NewInvalidDataError("email is required")
	}

	if rec.Status == "" {
		rec.Status = account_record.AccountUserStatusPendingApproval
	}

	if err := validateAccountUserStatus(rec.Status); err != nil {
		return err
	}

	return nil
}

func validateAccountUserRecForUpdate(args *validateAccountUserArgs) error {
	nextRec := args.nextRec
	currRec := args.currRec

	if nextRec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if nextRec.AccountID != currRec.AccountID {
		return coreerror.NewInvalidDataError("account_id cannot be updated")
	}

	if nextRec.Email != currRec.Email {
		return coreerror.NewInvalidDataError("email cannot be updated")
	}

	if nextRec.Status == "" {
		nextRec.Status = currRec.Status
	}

	if err := validateAccountUserStatus(nextRec.Status); err != nil {
		return err
	}

	return nil
}

func validateAccountUserRecForDelete(args *validateAccountUserArgs) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	return nil
}

func validateAccountUserStatus(status string) error {
	switch status {
	case account_record.AccountUserStatusPendingApproval,
		account_record.AccountUserStatusActive,
		account_record.AccountUserStatusDisabled:
		return nil
	default:
		return coreerror.NewInvalidDataError("invalid account user status >%s<", status)
	}
}
