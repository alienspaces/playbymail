package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
)

func (m *Domain) validateAccountRecForCreate(rec *account_record.Account) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if rec.Email == "" {
		return coreerror.NewInvalidDataError("email is required")
	}

	if rec.Status == "" {
		rec.Status = account_record.AccountStatusActive
	}

	if err := validateAccountStatus(rec.Status); err != nil {
		return err
	}

	return nil
}

func (m *Domain) validateAccountRecForUpdate(nextRec, currRec *account_record.Account) error {
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

func (m *Domain) validateAccountRecForDelete(rec *account_record.Account) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	return nil
}

func validateAccountStatus(status string) error {
	switch status {
	case account_record.AccountStatusPendingApproval,
		account_record.AccountStatusActive,
		account_record.AccountStatusDisabled:
		return nil
	default:
		return coreerror.NewInvalidDataError("invalid account status >%s<", status)
	}
}
