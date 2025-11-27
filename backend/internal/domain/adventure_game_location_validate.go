package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameLocationRecForCreate(rec *adventure_game_record.AdventureGameLocation) error {
	return validateAdventureGameLocationRec(rec, false)
}

func (m *Domain) validateAdventureGameLocationRecForUpdate(rec *adventure_game_record.AdventureGameLocation) error {
	return validateAdventureGameLocationRec(rec, true)
}

func validateAdventureGameLocationRec(rec *adventure_game_record.AdventureGameLocation, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameLocationName, rec.Name); err != nil {
		return err
	}

	if len(rec.Name) > 255 {
		return InvalidField(adventure_game_record.FieldAdventureGameLocationName, rec.Name, "name exceeds maximum length of 255 characters")
	}

	return nil
}
