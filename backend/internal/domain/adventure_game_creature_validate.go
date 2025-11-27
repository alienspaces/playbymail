package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameCreatureRecForCreate(rec *adventure_game_record.AdventureGameCreature) error {
	return validateAdventureGameCreatureRec(rec, false)
}

func (m *Domain) validateAdventureGameCreatureRecForUpdate(rec *adventure_game_record.AdventureGameCreature) error {
	return validateAdventureGameCreatureRec(rec, true)
}

func validateAdventureGameCreatureRec(rec *adventure_game_record.AdventureGameCreature, requireID bool) error {
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
