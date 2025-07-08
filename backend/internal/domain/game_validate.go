package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record"
)

type validateGameArgs struct {
	next *record.Game
	curr *record.Game
}

func (m *Domain) populateGameValidateArgs(next, curr *record.Game) (*validateGameArgs, error) {
	args := &validateGameArgs{
		curr: curr,
		next: next,
	}
	return args, nil
}

func (m *Domain) validateGameRecForCreate(rec *record.Game) error {
	args, err := m.populateGameValidateArgs(rec, nil)
	if err != nil {
		return err
	}
	return validateGameRecForCreate(args)
}

func (m *Domain) validateGameRecForUpdate(next, curr *record.Game) error {
	args, err := m.populateGameValidateArgs(next, curr)
	if err != nil {
		return err
	}
	return validateGameRecForUpdate(args)
}

func (m *Domain) validateGameRecForDelete(rec *record.Game) error {
	args, err := m.populateGameValidateArgs(rec, nil)
	if err != nil {
		return err
	}
	return validateGameRecForDelete(args)
}

func validateGameRecForCreate(args *validateGameArgs) error {
	return validateGameRec(args)
}

func validateGameRecForUpdate(args *validateGameArgs) error {
	return validateGameRec(args)
}

func validateGameRec(args *validateGameArgs) error {
	rec := args.next

	if err := domain.ValidateStringField(record.FieldGameName, rec.Name); err != nil {
		return err
	}

	if rec.GameType != record.GameTypeAdventure {
		return InvalidFieldValue("game_type")
	}

	return nil
}

func validateGameRecForDelete(args *validateGameArgs) error {
	rec := args.next

	if err := domain.ValidateUUIDField(record.FieldGameID, rec.ID); err != nil {
		return err
	}

	return nil
}
