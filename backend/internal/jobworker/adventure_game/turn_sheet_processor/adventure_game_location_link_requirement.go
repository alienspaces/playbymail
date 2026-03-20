package turn_sheet_processor

import (
	"fmt"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
)

// evaluateLinkRequirements returns (isVisible, isTraversable, error) for a location link.
// isVisible=false means the link must not appear on the sheet at all.
// isVisible=true, isTraversable=false means it appears locked.
// isVisible=true, isTraversable=true means it appears with a radio button.
func evaluateLinkRequirements(
	l logger.Logger,
	d *domain.Domain,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	fromLocationInstanceRec *adventure_game_record.AdventureGameLocationInstance,
	linkRec *adventure_game_record.AdventureGameLocationLink,
) (isVisible bool, isTraversable bool, err error) {
	requirements, err := d.GetManyAdventureGameLocationLinkRequirementRecs(&coresql.Options{
		Params: []coresql.Param{
			{
				Col: adventure_game_record.FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID,
				Val: linkRec.ID,
			},
		},
	})
	if err != nil {
		return false, false, fmt.Errorf("failed to get link requirements: %w", err)
	}

	// No requirements — link is always visible and traversable.
	if len(requirements) == 0 {
		return true, true, nil
	}

	// Evaluate all visible requirements first (AND logic — all must pass).
	for _, req := range requirements {
		if req.Purpose != adventure_game_record.AdventureGameLocationLinkRequirementPurposeVisible {
			continue
		}
		met, err := evaluateSingleRequirement(l, d, gameInstanceRec, characterInstanceRec, fromLocationInstanceRec, req)
		if err != nil {
			return false, false, err
		}
		if !met {
			return false, false, nil // hidden
		}
	}

	// Evaluate all traverse requirements (AND logic).
	for _, req := range requirements {
		if req.Purpose != adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse {
			continue
		}
		met, err := evaluateSingleRequirement(l, d, gameInstanceRec, characterInstanceRec, fromLocationInstanceRec, req)
		if err != nil {
			return true, false, err
		}
		if !met {
			return true, false, nil // visible but locked
		}
	}

	return true, true, nil
}

// evaluateSingleRequirement returns whether a single link requirement condition is satisfied.
func evaluateSingleRequirement(
	l logger.Logger,
	d *domain.Domain,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	fromLocationInstanceRec *adventure_game_record.AdventureGameLocationInstance,
	req *adventure_game_record.AdventureGameLocationLinkRequirement,
) (bool, error) {
	switch req.Condition {

	// Item conditions.
	case adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory:
		return characterHasItemInInventory(d, characterInstanceRec.ID, req.AdventureGameItemID.String, req.Quantity)

	case adventure_game_record.AdventureGameLocationLinkRequirementConditionEquipped:
		return characterHasItemEquipped(d, characterInstanceRec.ID, req.AdventureGameItemID.String)

	// Creature conditions.
	case adventure_game_record.AdventureGameLocationLinkRequirementConditionDeadAtLocation:
		return creatureDeadAtLocation(d, gameInstanceRec.ID, fromLocationInstanceRec.ID, req.AdventureGameCreatureID.String, req.Quantity)

	case adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveAtLocation:
		return noCreaturesAliveAtLocation(d, gameInstanceRec.ID, fromLocationInstanceRec.ID, req.AdventureGameCreatureID.String)

	case adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveInGame:
		return noCreaturesAliveInGame(d, gameInstanceRec.ID, req.AdventureGameCreatureID.String)

	default:
		l.Warn("unknown requirement condition >%s<", req.Condition)
		return false, fmt.Errorf("unknown requirement condition: %s", req.Condition)
	}
}

// characterHasItemInInventory returns true if the character holds at least quantity
// unused instances of the given item definition.
func characterHasItemInInventory(d *domain.Domain, characterInstanceID, itemID string, quantity int) (bool, error) {
	itemInstances, err := d.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceID},
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: itemID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query item instances: %w", err)
	}
	count := 0
	for _, inst := range itemInstances {
		if !inst.IsUsed {
			count++
		}
	}
	return count >= quantity, nil
}

// characterHasItemEquipped returns true if the character has at least one instance
// of the given item definition currently equipped.
func characterHasItemEquipped(d *domain.Domain, characterInstanceID, itemID string) (bool, error) {
	itemInstances, err := d.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceID},
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: itemID},
			{Col: adventure_game_record.FieldAdventureGameItemInstanceIsEquipped, Val: true},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query equipped item instances: %w", err)
	}
	return len(itemInstances) > 0, nil
}

// creatureDeadAtLocation returns true if at least quantity instances of the given
// creature definition are dead at the specified location.
func creatureDeadAtLocation(d *domain.Domain, gameInstanceID, fromLocationInstanceID, creatureID string, quantity int) (bool, error) {
	allInstances, err := d.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameCreatureID, Val: creatureID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, Val: fromLocationInstanceID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query creature instances: %w", err)
	}
	deadCount := 0
	for _, inst := range allInstances {
		if inst.Health <= 0 {
			deadCount++
		}
	}
	return deadCount >= quantity, nil
}

// noCreaturesAliveAtLocation returns true if no alive instances of the given creature
// definition exist at the specified location.
func noCreaturesAliveAtLocation(d *domain.Domain, gameInstanceID, fromLocationInstanceID, creatureID string) (bool, error) {
	allInstances, err := d.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameCreatureID, Val: creatureID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, Val: fromLocationInstanceID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query creature instances: %w", err)
	}
	for _, inst := range allInstances {
		if inst.Health > 0 {
			return false, nil
		}
	}
	return true, nil
}

// noCreaturesAliveInGame returns true if no alive instances of the given creature
// definition exist anywhere in the game instance.
func noCreaturesAliveInGame(d *domain.Domain, gameInstanceID, creatureID string) (bool, error) {
	allInstances, err := d.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameCreatureID, Val: creatureID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to query creature instances: %w", err)
	}
	for _, inst := range allInstances {
		if inst.Health > 0 {
			return false, nil
		}
	}
	return true, nil
}
