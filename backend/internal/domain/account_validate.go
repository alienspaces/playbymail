package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

type validateAccountArgs struct {
	next *record.Account
	curr *record.Account
}

func (m *Domain) populateAccountValidateArgs(next, curr *record.Account) (*validateAccountArgs, error) {
	args := &validateAccountArgs{
		curr: curr,
		next: next,
	}
	return args, nil
}

func (m *Domain) validateAccountRecForCreate(rec *record.Account) error {
	args, err := m.populateAccountValidateArgs(rec, nil)
	if err != nil {
		return err
	}
	return validateAccountRecForCreate(args)
}

func (m *Domain) validateAccountRecForUpdate(next, curr *record.Account) error {
	args, err := m.populateAccountValidateArgs(next, curr)
	if err != nil {
		return err
	}
	return validateAccountRecForUpdate(args)
}

func (m *Domain) validateAccountRecForDelete(rec *record.Account) error {
	args, err := m.populateAccountValidateArgs(rec, nil)
	if err != nil {
		return err
	}
	return validateAccountRecForDelete(args)
}

func validateAccountRecForCreate(args *validateAccountArgs) error {
	return validateAccountRec(args)
}

func validateAccountRecForUpdate(args *validateAccountArgs) error {
	return validateAccountRec(args)
}

func validateAccountRec(args *validateAccountArgs) error {
	rec := args.next

	if err := domain.ValidateStringField(record.FieldAccountEmail, rec.Email); err != nil {
		return err
	}

	return nil
}

func validateAccountRecForDelete(args *validateAccountArgs) error {
	rec := args.next

	if err := domain.ValidateUUIDField(record.FieldAccountID, rec.ID); err != nil {
		return err
	}

	return nil
}
