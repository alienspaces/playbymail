package domain

import (
	coreerror "gitlab.com/alienspaces/playbymail/core/error"
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameTurnSheetRecForCreate(rec *adventure_game_record.AdventureGameTurnSheet) error {
	return validateAdventureGameTurnSheetRec(rec, false)
}

func (m *Domain) validateAdventureGameTurnSheetRecForUpdate(rec *adventure_game_record.AdventureGameTurnSheet) error {
	return validateAdventureGameTurnSheetRec(rec, true)
}

func validateAdventureGameTurnSheetRec(rec *adventure_game_record.AdventureGameTurnSheet, requireID bool) error {
	if rec == nil {
		return coreerror.NewInvalidDataError("record is nil")
	}

	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameTurnSheetID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameTurnSheetGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID, rec.AdventureGameCharacterInstanceID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameTurnSheetGameTurnSheetID, rec.GameTurnSheetID); err != nil {
		return err
	}

	return nil
}
