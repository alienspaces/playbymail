package domain

import (
	"fmt"

	"gitlab.com/alienspaces/playbymail/core/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

type validateAdventureGameLocationLinkRequirementArgs struct {
	nextRec *adventure_game_record.AdventureGameLocationLinkRequirement
	currRec *adventure_game_record.AdventureGameLocationLinkRequirement
}

func (m *Domain) populateAdventureGameLocationLinkRequirementValidateArgs(currRec, nextRec *adventure_game_record.AdventureGameLocationLinkRequirement) (*validateAdventureGameLocationLinkRequirementArgs, error) {
	args := &validateAdventureGameLocationLinkRequirementArgs{
		currRec: currRec,
		nextRec: nextRec,
	}
	return args, nil
}

func (m *Domain) validateAdventureGameLocationLinkRequirementRecForCreate(rec *adventure_game_record.AdventureGameLocationLinkRequirement) error {
	args, err := m.populateAdventureGameLocationLinkRequirementValidateArgs(nil, rec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationLinkRequirementRecForCreate(args)
}

func (m *Domain) validateAdventureGameLocationLinkRequirementRecForUpdate(currRec, nextRec *adventure_game_record.AdventureGameLocationLinkRequirement) error {
	args, err := m.populateAdventureGameLocationLinkRequirementValidateArgs(currRec, nextRec)
	if err != nil {
		return err
	}
	return validateAdventureGameLocationLinkRequirementRecForUpdate(args)
}

func validateAdventureGameLocationLinkRequirementRecForCreate(args *validateAdventureGameLocationLinkRequirementArgs) error {
	return validateAdventureGameLocationLinkRequirementRec(args, false)
}

func validateAdventureGameLocationLinkRequirementRecForUpdate(args *validateAdventureGameLocationLinkRequirementArgs) error {
	return validateAdventureGameLocationLinkRequirementRec(args, true)
}

func validateAdventureGameLocationLinkRequirementRec(args *validateAdventureGameLocationLinkRequirementArgs, requireID bool) error {
	rec := args.nextRec

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
