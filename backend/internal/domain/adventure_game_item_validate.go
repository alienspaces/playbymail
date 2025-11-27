package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameItemRecForCreate(rec *adventure_game_record.AdventureGameItem) error {
	return validateAdventureGameItemRec(rec, false)
}

func (m *Domain) validateAdventureGameItemRecForUpdate(rec *adventure_game_record.AdventureGameItem) error {
	return validateAdventureGameItemRec(rec, true)
}

func validateAdventureGameItemRec(rec *adventure_game_record.AdventureGameItem, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameItemGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameItemName, rec.Name); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameItemDescription, rec.Description); err != nil {
		return err
	}

	return nil
}
