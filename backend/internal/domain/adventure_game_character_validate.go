package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameCharacterRecForCreate(rec *adventure_game_record.AdventureGameCharacter) error {
	return validateAdventureGameCharacterRec(rec, false)
}

func (m *Domain) validateAdventureGameCharacterRecForUpdate(rec *adventure_game_record.AdventureGameCharacter) error {
	return validateAdventureGameCharacterRec(rec, true)
}

func validateAdventureGameCharacterRec(rec *adventure_game_record.AdventureGameCharacter, requireID bool) error {
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
