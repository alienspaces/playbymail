package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameCharacterArgs struct {
	nextRec *adventure_game_record.AdventureGameCharacter
	currRec *adventure_game_record.AdventureGameCharacter
}

func (m *Domain) populateAdventureGameCharacterValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameCharacter) (*validateAdventureGameCharacterArgs, error) {
	args := &validateAdventureGameCharacterArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameCharacterRecForCreate(rec *adventure_game_record.AdventureGameCharacter) error {
	args, err := m.populateAdventureGameCharacterValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameCharacterRecForCreate(args)
}

func (m *Domain) validateAdventureGameCharacterRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameCharacter) error {
	args, err := m.populateAdventureGameCharacterValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameCharacterRecForUpdate(args)
}

func validateAdventureGameCharacterRecForCreate(args *validateAdventureGameCharacterArgs) error {
	return validateAdventureGameCharacterRec(args, false)
}

func validateAdventureGameCharacterRecForUpdate(args *validateAdventureGameCharacterArgs) error {
	return validateAdventureGameCharacterRec(args, true)
}

func validateAdventureGameCharacterRec(args *validateAdventureGameCharacterArgs, requireID bool) error {
	rec := args.nextRec

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameCharacterName, rec.Name); err != nil {
		return err
	}

	if len(rec.Name) > 128 {
		return InvalidField(adventure_game_record.FieldAdventureGameCharacterName, rec.Name, "name is too long")
	}

	return nil
}
