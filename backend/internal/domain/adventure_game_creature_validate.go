package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameCreatureArgs struct {
	nextRec *adventure_game_record.AdventureGameCreature
	currRec *adventure_game_record.AdventureGameCreature
}

func (m *Domain) populateAdventureGameCreatureValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameCreature) (*validateAdventureGameCreatureArgs, error) {
	args := &validateAdventureGameCreatureArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameCreatureRecForCreate(rec *adventure_game_record.AdventureGameCreature) error {
	args, err := m.populateAdventureGameCreatureValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameCreatureRecForCreate(args)
}

func (m *Domain) validateAdventureGameCreatureRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameCreature) error {
	args, err := m.populateAdventureGameCreatureValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameCreatureRecForUpdate(args)
}

func validateAdventureGameCreatureRecForCreate(args *validateAdventureGameCreatureArgs) error {
	return validateAdventureGameCreatureRec(args, false)
}

func validateAdventureGameCreatureRecForUpdate(args *validateAdventureGameCreatureArgs) error {
	return validateAdventureGameCreatureRec(args, true)
}

func validateAdventureGameCreatureRec(args *validateAdventureGameCreatureArgs, requireID bool) error {
	rec := args.nextRec

	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameCreatureName, rec.Name); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameCreatureDescription, rec.Description); err != nil {
		return err
	}

	return nil
}
