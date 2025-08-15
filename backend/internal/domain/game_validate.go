package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type validateGameArgs struct {
	next *game_record.Game
	curr *game_record.Game
}

func (m *Domain) populateGameValidateArgs(next, curr *game_record.Game) (*validateGameArgs, error) {
	args := &validateGameArgs{
		curr: curr,
		next: next,
	}
	return args, nil
}

func (m *Domain) validateGameRecForCreate(rec *game_record.Game) error {
	args, err := m.populateGameValidateArgs(rec, nil)
	if err != nil {
		return err
	}
	return validateGameRecForCreate(args)
}

func (m *Domain) validateGameRecForUpdate(next, curr *game_record.Game) error {
	args, err := m.populateGameValidateArgs(next, curr)
	if err != nil {
		return err
	}
	return validateGameRecForUpdate(args)
}

func (m *Domain) validateGameRecForDelete(rec *game_record.Game) error {
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

	if err := domain.ValidateStringField(game_record.FieldGameName, rec.Name); err != nil {
		return err
	}

	if rec.GameType != game_record.GameTypeAdventure {
		return InvalidFieldValue("game_type")
	}

	if rec.TurnDurationHours <= 0 {
		return InvalidFieldValue("turn_duration_hours")
	}

	return nil
}

func validateGameRecForDelete(args *validateGameArgs) error {
	rec := args.next

	if err := domain.ValidateUUIDField(game_record.FieldGameID, rec.ID); err != nil {
		return err
	}

	return nil
}
