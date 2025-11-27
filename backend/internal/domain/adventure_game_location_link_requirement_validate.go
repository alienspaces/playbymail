package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func (m *Domain) validateAdventureGameLocationLinkRequirementRecForCreate(rec *adventure_game_record.AdventureGameLocationLinkRequirement) error {
	return validateAdventureGameLocationLinkRequirementRec(rec, false)
}

func (m *Domain) validateAdventureGameLocationLinkRequirementRecForUpdate(rec *adventure_game_record.AdventureGameLocationLinkRequirement) error {
	return validateAdventureGameLocationLinkRequirementRec(rec, true)
}

func validateAdventureGameLocationLinkRequirementRec(rec *adventure_game_record.AdventureGameLocationLinkRequirement, requireID bool) error {
	if requireID {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkRequirementID, rec.ID); err != nil {
			return err
		}
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkRequirementGameID, rec.GameID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID, rec.AdventureGameLocationLinkID); err != nil {
		return err
	}

	if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkRequirementAdventureGameItemID, rec.AdventureGameItemID); err != nil {
		return err
	}

	if rec.Quantity <= 0 {
		return InvalidField(
			adventure_game_record.FieldAdventureGameLocationLinkRequirementQuantity,
			fmt.Sprintf("%d", rec.Quantity),
			"quantity must be greater than 0",
		)
	}

	return nil
}
