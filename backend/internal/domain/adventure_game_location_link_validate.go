package domain

import (
	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameLocationLinkRecForCreate(rec *adventure_game_record.AdventureGameLocationLink) error {
	return validateAdventureGameLocationLinkRec(rec, false)
}

func (m *Domain) validateAdventureGameLocationLinkRecForUpdate(rec *adventure_game_record.AdventureGameLocationLink) error {
	return validateAdventureGameLocationLinkRec(rec, true)
}

func validateAdventureGameLocationLinkRec(rec *adventure_game_record.AdventureGameLocationLink, requireID bool) error {
	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkFromAdventureGameLocationID, rec.FromAdventureGameLocationID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkToAdventureGameLocationID, rec.ToAdventureGameLocationID); err != nil {
		return err
	}

	if err := domain.ValidateStringField(adventure_game_record.FieldAdventureGameLocationLinkName, rec.Name); err != nil {
		return err
	}

	if len(rec.Name) > 255 {
		return InvalidField(adventure_game_record.FieldAdventureGameLocationLinkName, rec.Name, "name exceeds maximum length of 255 characters")
	}

	return nil
}
