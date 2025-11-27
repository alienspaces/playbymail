package domain

import (
	"fmt"

	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameCharacterInstanceRecForCreate(rec *adventure_game_record.AdventureGameCharacterInstance) error {
	return validateAdventureGameCharacterInstanceRec(rec, false)
}

func (m *Domain) validateAdventureGameCharacterInstanceRecForUpdate(rec *adventure_game_record.AdventureGameCharacterInstance) error {
	return validateAdventureGameCharacterInstanceRec(rec, true)
}

func validateAdventureGameCharacterInstanceRec(rec *adventure_game_record.AdventureGameCharacterInstance, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceAdventureGameCharacterID, rec.AdventureGameCharacterID); err != nil {
		return err
	}

	if rec.AdventureGameLocationInstanceID != "" {
	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCharacterInstanceAdventureGameLocationInstanceID, rec.AdventureGameLocationInstanceID); err != nil {
		return err
		}
	}

	if rec.Health < 0 {
		return InvalidField(
			adventure_game_record.FieldAdventureGameCharacterInstanceHealth,
			fmt.Sprintf("%d", rec.Health),
			"health must be zero or greater",
		)
	}

	return nil
}
