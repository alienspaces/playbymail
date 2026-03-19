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

var validItemConditions = map[string]bool{
	adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory: true,
	adventure_game_record.AdventureGameLocationLinkRequirementConditionEquipped:    true,
}

var validCreatureConditions = map[string]bool{
	adventure_game_record.AdventureGameLocationLinkRequirementConditionDeadAtLocation:      true,
	adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveAtLocation: true,
	adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveInGame:     true,
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

	// Exactly one of item or creature must be set
	hasItem := rec.AdventureGameItemID.Valid && rec.AdventureGameItemID.String != ""
	hasCreature := rec.AdventureGameCreatureID.Valid && rec.AdventureGameCreatureID.String != ""

	if hasItem == hasCreature {
		return InvalidField(
			"adventure_game_item_id / adventure_game_creature_id",
			"",
			"exactly one of adventure_game_item_id or adventure_game_creature_id must be set",
		)
	}

	if hasItem {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkRequirementAdventureGameItemID, rec.AdventureGameItemID.String); err != nil {
			return err
		}
		if !validItemConditions[rec.Condition] {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationLinkRequirementCondition,
				rec.Condition,
				fmt.Sprintf("item requirements must use one of: in_inventory, equipped; got %q", rec.Condition),
			)
		}
	}

	if hasCreature {
		if err := domain.ValidateUUIDField(adventure_game_record.FieldAdventureGameLocationLinkRequirementAdventureGameCreatureID, rec.AdventureGameCreatureID.String); err != nil {
			return err
		}
		if !validCreatureConditions[rec.Condition] {
			return InvalidField(
				adventure_game_record.FieldAdventureGameLocationLinkRequirementCondition,
				rec.Condition,
				fmt.Sprintf("creature requirements must use one of: dead_at_location, none_alive_at_location, none_alive_in_game; got %q", rec.Condition),
			)
		}
	}

	switch rec.Purpose {
	case adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
		adventure_game_record.AdventureGameLocationLinkRequirementPurposeVisible:
	default:
		return InvalidField(
			adventure_game_record.FieldAdventureGameLocationLinkRequirementPurpose,
			rec.Purpose,
			fmt.Sprintf("purpose must be one of: traverse, visible; got %q", rec.Purpose),
		)
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
