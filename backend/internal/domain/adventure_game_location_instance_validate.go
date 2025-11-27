package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameLocationInstanceRecForCreate(rec *adventure_game_record.AdventureGameLocationInstance) error {
	return validateAdventureGameLocationInstanceRec(rec, false)
}

func (m *Domain) validateAdventureGameLocationInstanceRecForUpdate(rec *adventure_game_record.AdventureGameLocationInstance) error {
	return validateAdventureGameLocationInstanceRec(rec, true)
}

func validateAdventureGameLocationInstanceRec(rec *adventure_game_record.AdventureGameLocationInstance, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationInstanceID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationInstanceGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, rec.GameInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationInstanceAdventureGameLocationID, rec.AdventureGameLocationID); err != nil {
		return err
	}

	return nil
}

