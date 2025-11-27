package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameCreatureInstanceRecForCreate(rec *adventure_game_record.AdventureGameCreatureInstance) error {
	return validateAdventureGameCreatureInstanceRec(rec, false)
}

func (m *Domain) validateAdventureGameCreatureInstanceRecForUpdate(rec *adventure_game_record.AdventureGameCreatureInstance) error {
	return validateAdventureGameCreatureInstanceRec(rec, true)
}

func validateAdventureGameCreatureInstanceRec(rec *adventure_game_record.AdventureGameCreatureInstance, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameCreatureID, rec.AdventureGameCreatureID); err != nil {
		return err
	}

	if rec.AdventureGameLocationInstanceID != "" {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, rec.AdventureGameLocationInstanceID); err != nil {
			return err
		}
	}

	if rec.Health < 0 {
		return InvalidField(
			adventure_game_record.FieldAdventureGameCreatureInstanceHealth,
			fmt.Sprintf("%d", rec.Health),
			"health must be zero or greater",
		)
	}

	return nil
}
