package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

type validateGameArgs struct {
	nextRec *game_record.Game
	currRec *game_record.Game
}

func (m *Domain) populateGameValidateArgs(nextRec, currRec *game_record.Game) (*validateGameArgs, error) {
	args := &validateGameArgs{
		currRec: currRec,
		nextRec: nextRec,
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

func (m *Domain) validateGameRecForUpdate(nextRec, currRec *game_record.Game) error {
	args, err := m.populateGameValidateArgs(nextRec, currRec)
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
	rec := args.nextRec

	if err := domain.ValidateStringField(game_record.FieldGameName, rec.Name); err != nil {
		return err
	}

	if rec.GameType != game_record.GameTypeAdventure {
		return InvalidField(game_record.FieldGameType, rec.GameType, "game type is not valid")
	}

	if rec.TurnDurationHours <= 0 {
		return InvalidField(game_record.FieldGameTurnDurationHours, fmt.Sprintf("%d", rec.TurnDurationHours), "turn duration hours must be greater than 0")
	}

	if err := domain.ValidateStringField(game_record.FieldGameDescription, rec.Description); err != nil {
		return err
	}

	return nil
}

func validateGameRecForDelete(args *validateGameArgs) error {
	rec := args.nextRec

	if err := domain.ValidateUUIDField(game_record.FieldGameID, rec.ID); err != nil {
		return err
	}

	return nil
}
