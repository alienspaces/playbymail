package turn_sheet_processor

import (
	"fmt"
	"math/rand"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// applyItemEffectsForAction queries all item effects matching (itemDefinitionID, actionType),
// filters against required conditions, then applies them atomically to the character.
// Returns the combined result descriptions and whether the character record was mutated.
// weapon_damage and armor_defense effects are passive stats consumed by ResolveEquipmentStats
// and are intentionally skipped here.
func applyItemEffectsForAction(
	l logger.Logger,
	d *domain.Domain,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	itemInstanceRec *adventure_game_record.AdventureGameItemInstance,
	actionType string,
) (resultDescriptions []string, characterMutated bool, err error) {
	l = l.WithFunctionContext("applyItemEffectsForAction")

	effects, err := d.GetManyAdventureGameItemEffectRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemEffectAdventureGameItemID, Val: itemInstanceRec.AdventureGameItemID},
			{Col: adventure_game_record.FieldAdventureGameItemEffectActionType, Val: actionType},
		},
	})
	if err != nil {
		return nil, false, fmt.Errorf("failed to query item effects for item >%s< action >%s<: %w", itemInstanceRec.AdventureGameItemID, actionType, err)
	}

	// Filter by required conditions.
	var matchingEffects []*adventure_game_record.AdventureGameItemEffect
	for _, effect := range effects {
		if !itemEffectRequiredConditionsMet(l, d, characterInstanceRec, effect) {
			continue
		}
		matchingEffects = append(matchingEffects, effect)
	}

	if len(matchingEffects) == 0 {
		l.Info("no matching effects for item >%s< action >%s<", itemInstanceRec.AdventureGameItemID, actionType)
		return nil, false, nil
	}

	for _, effect := range matchingEffects {
		desc, mutated, applyErr := applyItemEffect(l, d, gameInstanceRec, characterInstanceRec, effect)
		if applyErr != nil {
			l.Warn("failed to apply item effect >%s< >%v<", effect.ID, applyErr)
			return resultDescriptions, characterMutated, fmt.Errorf("failed to apply item effect: %w", applyErr)
		}
		if mutated {
			characterMutated = true
		}
		if desc != "" {
			resultDescriptions = append(resultDescriptions, desc)
		}
	}

	return resultDescriptions, characterMutated, nil
}

// itemEffectRequiredConditionsMet returns true if all required conditions for the effect are satisfied.
func itemEffectRequiredConditionsMet(
	l logger.Logger,
	d *domain.Domain,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	effect *adventure_game_record.AdventureGameItemEffect,
) bool {
	// Check required item in inventory.
	if effect.RequiredAdventureGameItemID.Valid && effect.RequiredAdventureGameItemID.String != "" {
		itemInstances, err := d.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceRec.ID},
				{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: effect.RequiredAdventureGameItemID.String},
			},
		})
		if err != nil {
			l.Warn("failed to check required item for effect >%s< >%v<", effect.ID, err)
			return false
		}
		hasItem := false
		for _, inst := range itemInstances {
			if !inst.IsUsed {
				hasItem = true
				break
			}
		}
		if !hasItem {
			l.Info("required item >%s< not in inventory — skipping effect >%s<", effect.RequiredAdventureGameItemID.String, effect.ID)
			return false
		}
	}

	// Check required location.
	if effect.RequiredAdventureGameLocationID.Valid && effect.RequiredAdventureGameLocationID.String != "" {
		locationInstances, err := d.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameLocationInstanceAdventureGameLocationID, Val: effect.RequiredAdventureGameLocationID.String},
			},
		})
		if err != nil {
			l.Warn("failed to check required location for effect >%s< >%v<", effect.ID, err)
			return false
		}
		atLocation := false
		for _, inst := range locationInstances {
			if inst.ID == characterInstanceRec.AdventureGameLocationInstanceID {
				atLocation = true
				break
			}
		}
		if !atLocation {
			l.Info("required location not current location — skipping effect >%s<", effect.ID)
			return false
		}
	}

	return true
}

// applyItemEffect dispatches a single item effect and returns (resultDescription, characterMutated, error).
// weapon_damage and armor_defense effects are passive stats — they are skipped here and consumed
// by ResolveEquipmentStats when computing combat stats for equipped items.
func applyItemEffect(
	l logger.Logger,
	d *domain.Domain,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	effect *adventure_game_record.AdventureGameItemEffect,
) (string, bool, error) {
	l = l.WithFunctionContext("applyItemEffect")
	l.Info("applying item effect >%s< type >%s<", effect.ID, effect.EffectType)

	switch effect.EffectType {
	case adventure_game_record.AdventureGameItemEffectEffectTypeInfo,
		adventure_game_record.AdventureGameItemEffectEffectTypeNothing:
		return effect.ResultDescription, false, nil

	case adventure_game_record.AdventureGameItemEffectEffectTypeWeaponDamage,
		adventure_game_record.AdventureGameItemEffectEffectTypeArmorDefense:
		// Passive equipped stat — not applied as a one-shot effect.
		return "", false, nil

	case adventure_game_record.AdventureGameItemEffectEffectTypeDamageTarget,
		adventure_game_record.AdventureGameItemEffectEffectTypeDamageWielder,
		adventure_game_record.AdventureGameItemEffectEffectTypeHealTarget,
		adventure_game_record.AdventureGameItemEffectEffectTypeHealWielder:
		applyItemHealthEffect(l, characterInstanceRec, effect)
		return effect.ResultDescription, true, nil

	case adventure_game_record.AdventureGameItemEffectEffectTypeGiveItem,
		adventure_game_record.AdventureGameItemEffectEffectTypeRemoveItem:
		if err := applyItemInventoryEffect(l, d, gameInstanceRec, characterInstanceRec, effect); err != nil {
			return "", false, err
		}
		return effect.ResultDescription, false, nil

	case adventure_game_record.AdventureGameItemEffectEffectTypeOpenLink,
		adventure_game_record.AdventureGameItemEffectEffectTypeCloseLink:
		if err := applyItemLinkEffect(l, d, effect); err != nil {
			return "", false, err
		}
		return effect.ResultDescription, false, nil

	case adventure_game_record.AdventureGameItemEffectEffectTypeTeleport,
		adventure_game_record.AdventureGameItemEffectEffectTypeSummonCreature:
		mutated, err := applyItemWorldEffect(l, d, gameInstanceRec, characterInstanceRec, effect)
		if err != nil {
			return "", false, err
		}
		return effect.ResultDescription, mutated, nil
	}

	l.Warn("unhandled item effect type >%s< for effect >%s<", effect.EffectType, effect.ID)
	return "", false, nil
}

// applyItemHealthEffect applies damage or heal health effects to the character.
// damage_target and heal_target treat the character as the target in the item-use context
// (there is no explicit target selection during inventory management).
func applyItemHealthEffect(
	l logger.Logger,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	effect *adventure_game_record.AdventureGameItemEffect,
) {
	if !effect.ResultValueMin.Valid || !effect.ResultValueMax.Valid {
		return
	}
	minVal := int(effect.ResultValueMin.Int32)
	maxVal := int(effect.ResultValueMax.Int32)
	amount := minVal
	if maxVal > minVal {
		amount = minVal + rand.Intn(maxVal-minVal+1)
	}

	switch effect.EffectType {
	case adventure_game_record.AdventureGameItemEffectEffectTypeDamageTarget,
		adventure_game_record.AdventureGameItemEffectEffectTypeDamageWielder:
		characterInstanceRec.Health -= amount
		if characterInstanceRec.Health < 0 {
			characterInstanceRec.Health = 0
		}
		l.Info("item dealt %d damage to character >%s< (health now %d)", amount, characterInstanceRec.ID, characterInstanceRec.Health)

	case adventure_game_record.AdventureGameItemEffectEffectTypeHealTarget,
		adventure_game_record.AdventureGameItemEffectEffectTypeHealWielder:
		characterInstanceRec.Health += amount
		if characterInstanceRec.Health > 100 {
			characterInstanceRec.Health = 100
		}
		l.Info("item healed character >%s< for %d (health now %d)", characterInstanceRec.ID, amount, characterInstanceRec.Health)
	}
}

// applyItemInventoryEffect applies give_item and remove_item effects.
func applyItemInventoryEffect(
	l logger.Logger,
	d *domain.Domain,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	effect *adventure_game_record.AdventureGameItemEffect,
) error {
	switch effect.EffectType {
	case adventure_game_record.AdventureGameItemEffectEffectTypeGiveItem:
		if !effect.ResultAdventureGameItemID.Valid || effect.ResultAdventureGameItemID.String == "" {
			return nil
		}
		itemInstance := &adventure_game_record.AdventureGameItemInstance{
			GameID:                           gameInstanceRec.GameID,
			GameInstanceID:                   gameInstanceRec.ID,
			AdventureGameItemID:              effect.ResultAdventureGameItemID.String,
			AdventureGameCharacterInstanceID: nullstring.FromString(characterInstanceRec.ID),
		}
		if _, err := d.CreateAdventureGameItemInstanceRec(itemInstance); err != nil {
			return fmt.Errorf("failed to give item: %w", err)
		}
		l.Info("gave item >%s< to character >%s<", effect.ResultAdventureGameItemID.String, characterInstanceRec.ID)

	case adventure_game_record.AdventureGameItemEffectEffectTypeRemoveItem:
		if !effect.ResultAdventureGameItemID.Valid || effect.ResultAdventureGameItemID.String == "" {
			return nil
		}
		itemInstances, err := d.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceRec.ID},
				{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: effect.ResultAdventureGameItemID.String},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to query item instances for removal: %w", err)
		}
		for _, inst := range itemInstances {
			if !inst.IsUsed {
				if err := d.DeleteAdventureGameItemInstanceRec(inst.ID); err != nil {
					l.Warn("failed to remove item instance >%v<", err)
				}
				l.Info("removed item >%s< from character >%s<", effect.ResultAdventureGameItemID.String, characterInstanceRec.ID)
				break
			}
		}
	}
	return nil
}

// applyItemLinkEffect applies open_link and close_link effects.
func applyItemLinkEffect(
	l logger.Logger,
	d *domain.Domain,
	effect *adventure_game_record.AdventureGameItemEffect,
) error {
	switch effect.EffectType {
	case adventure_game_record.AdventureGameItemEffectEffectTypeOpenLink:
		if !effect.ResultAdventureGameLocationLinkID.Valid || effect.ResultAdventureGameLocationLinkID.String == "" {
			return nil
		}
		requirements, err := d.GetManyAdventureGameLocationLinkRequirementRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID, Val: effect.ResultAdventureGameLocationLinkID.String},
				{Col: adventure_game_record.FieldAdventureGameLocationLinkRequirementPurpose, Val: adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get link requirements: %w", err)
		}
		for _, req := range requirements {
			if err := d.DeleteAdventureGameLocationLinkRequirementRec(req.ID); err != nil {
				l.Warn("failed to remove link requirement >%v<", err)
			}
		}
		l.Info("opened link >%s<", effect.ResultAdventureGameLocationLinkID.String)

	case adventure_game_record.AdventureGameItemEffectEffectTypeCloseLink:
		if !effect.ResultAdventureGameLocationLinkID.Valid || effect.ResultAdventureGameLocationLinkID.String == "" {
			return nil
		}
		linkID := effect.ResultAdventureGameLocationLinkID.String
		linkRec, err := d.GetAdventureGameLocationLinkRec(linkID, nil)
		if err != nil {
			return fmt.Errorf("failed to get link record for close_link: %w", err)
		}

		req := &adventure_game_record.AdventureGameLocationLinkRequirement{
			GameID:                      linkRec.GameID,
			AdventureGameLocationLinkID: linkID,
			Purpose:                     adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse,
			Quantity:                    1,
		}

		hasItem := effect.ResultAdventureGameItemID.Valid && effect.ResultAdventureGameItemID.String != ""
		hasCreature := effect.ResultAdventureGameCreatureID.Valid && effect.ResultAdventureGameCreatureID.String != ""

		switch {
		case hasItem:
			req.AdventureGameItemID = effect.ResultAdventureGameItemID
			req.Condition = adventure_game_record.AdventureGameLocationLinkRequirementConditionInInventory
		case hasCreature:
			req.AdventureGameCreatureID = effect.ResultAdventureGameCreatureID
			req.Condition = adventure_game_record.AdventureGameLocationLinkRequirementConditionNoneAliveInGame
		default:
			l.Warn("close_link effect >%s< has no result item or creature — link cannot be locked", effect.ID)
			return nil
		}

		if _, err := d.CreateAdventureGameLocationLinkRequirementRec(req); err != nil {
			return fmt.Errorf("failed to create link requirement for close_link: %w", err)
		}
		l.Info("closed link >%s< (item=%v creature=%v)", linkID, hasItem, hasCreature)
	}
	return nil
}

// applyItemWorldEffect applies teleport and summon_creature effects.
// Returns whether the character instance record was mutated.
func applyItemWorldEffect(
	l logger.Logger,
	d *domain.Domain,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	effect *adventure_game_record.AdventureGameItemEffect,
) (bool, error) {
	switch effect.EffectType {
	case adventure_game_record.AdventureGameItemEffectEffectTypeTeleport:
		if !effect.ResultAdventureGameLocationID.Valid || effect.ResultAdventureGameLocationID.String == "" {
			return false, nil
		}
		locationInstances, err := d.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, Val: gameInstanceRec.ID},
				{Col: adventure_game_record.FieldAdventureGameLocationInstanceAdventureGameLocationID, Val: effect.ResultAdventureGameLocationID.String},
			},
		})
		if err != nil {
			return false, fmt.Errorf("failed to get destination location instance: %w", err)
		}
		if len(locationInstances) == 0 {
			return false, fmt.Errorf("no location instance found for location >%s<", effect.ResultAdventureGameLocationID.String)
		}
		destLocationInstance := locationInstances[0]
		characterInstanceRec.AdventureGameLocationInstanceID = destLocationInstance.ID
		updatedChar, err := d.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec)
		if err != nil {
			return false, fmt.Errorf("failed to teleport character: %w", err)
		}
		*characterInstanceRec = *updatedChar
		l.Info("teleported character >%s< to location instance >%s<", characterInstanceRec.ID, destLocationInstance.ID)
		return true, nil

	case adventure_game_record.AdventureGameItemEffectEffectTypeSummonCreature:
		if !effect.ResultAdventureGameCreatureID.Valid || effect.ResultAdventureGameCreatureID.String == "" {
			return false, nil
		}
		creatureRec, err := d.GetAdventureGameCreatureRec(effect.ResultAdventureGameCreatureID.String, nil)
		if err != nil {
			return false, fmt.Errorf("failed to get creature definition for summon: %w", err)
		}
		creatureInstance := &adventure_game_record.AdventureGameCreatureInstance{
			GameID:                          gameInstanceRec.GameID,
			GameInstanceID:                  gameInstanceRec.ID,
			AdventureGameCreatureID:         effect.ResultAdventureGameCreatureID.String,
			AdventureGameLocationInstanceID: characterInstanceRec.AdventureGameLocationInstanceID,
			Health:                          creatureRec.MaxHealth,
		}
		if _, err := d.CreateAdventureGameCreatureInstanceRec(creatureInstance); err != nil {
			return false, fmt.Errorf("failed to summon creature: %w", err)
		}
		l.Info("summoned creature >%s< at location instance >%s<", effect.ResultAdventureGameCreatureID.String, characterInstanceRec.AdventureGameLocationInstanceID)
	}

	return false, nil
}

// itemHasUseEffects returns true if an item has at least one use-action effect,
// which determines whether the item can be used by the player.
func itemHasUseEffects(l logger.Logger, d *domain.Domain, itemDefinitionID string) bool {
	effects, err := d.GetManyAdventureGameItemEffectRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemEffectAdventureGameItemID, Val: itemDefinitionID},
			{Col: adventure_game_record.FieldAdventureGameItemEffectActionType, Val: adventure_game_record.AdventureGameItemEffectActionTypeUse},
		},
	})
	if err != nil {
		l.Warn("failed to query use effects for item >%s< >%v<", itemDefinitionID, err)
		return false
	}
	return len(effects) > 0
}

// buildItemEffectTurnEvents creates turn events from item effect result descriptions and effect types.
func buildItemEffectTurnEvents(
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	descriptions []string,
	effectTypes []string,
) {
	if len(descriptions) == 0 {
		return
	}

	// Determine icon based on effect types applied.
	icon := turnsheet.TurnEventIconInventory
	for _, et := range effectTypes {
		switch et {
		case adventure_game_record.AdventureGameItemEffectEffectTypeDamageTarget,
			adventure_game_record.AdventureGameItemEffectEffectTypeDamageWielder:
			icon = turnsheet.TurnEventIconCombat
		case adventure_game_record.AdventureGameItemEffectEffectTypeHealTarget,
			adventure_game_record.AdventureGameItemEffectEffectTypeHealWielder:
			icon = turnsheet.TurnEventIconHeal
		case adventure_game_record.AdventureGameItemEffectEffectTypeTeleport:
			icon = turnsheet.TurnEventIconMovement
		}
	}

	_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
		Category: turnsheet.TurnEventCategoryInventory,
		Icon:     icon,
		Message:  strings.Join(descriptions, " "),
	})
}
