package turn_sheet_processor

import (
	"fmt"
	"math/rand"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// buildLocationObjects queries visible object instances at a location and
// builds the LocationObjectOption slice.
func (p *AdventureGameLocationChoiceProcessor) buildLocationObjects(
	l logger.Logger,
	gameInstanceID, characterInstanceID, locationInstanceID string,
) ([]turnsheet.LocationObjectOption, error) {
	l = l.WithFunctionContext("buildLocationObjects")

	objectInstances, err := p.Domain.GetManyAdventureGameLocationObjectInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceAdventureGameLocationInstanceID, Val: locationInstanceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object instances: %w", err)
	}

	var result []turnsheet.LocationObjectOption
	for _, inst := range objectInstances {
		if !inst.IsVisible {
			continue
		}

		objectDef, err := p.Domain.GetAdventureGameLocationObjectRec(inst.AdventureGameLocationObjectID, nil)
		if err != nil {
			l.Warn("failed to get object definition >%s< >%v<", inst.AdventureGameLocationObjectID, err)
			continue
		}

		// Load effects that match current state
		// (required_state = current_state OR required_state IS NULL).
		allEffects, err := p.Domain.GetManyAdventureGameLocationObjectEffectRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameLocationObjectEffectAdventureGameLocationObjectID, Val: objectDef.ID},
			},
		})
		if err != nil {
			l.Warn("failed to get object effects >%v<", err)
			continue
		}

		// Collect unique action types available in the current state.
		actionMap := map[string]*turnsheet.LocationObjectActionOption{}
		for _, effect := range allEffects {
			if effect.RequiredAdventureGameLocationObjectStateID.Valid && effect.RequiredAdventureGameLocationObjectStateID.String != inst.CurrentAdventureGameLocationObjectStateID {
				continue
			}
			at := effect.ActionType
			if _, exists := actionMap[at]; exists {
				continue
			}

			action := &turnsheet.LocationObjectActionOption{
				ActionType:      at,
				IsAvailable:     true,
				HasRequiredItem: true,
			}

			if effect.RequiredAdventureGameItemID.Valid && effect.RequiredAdventureGameItemID.String != "" {
				hasItem, itemName, err := p.characterHasItemForEffect(characterInstanceID, effect.RequiredAdventureGameItemID.String)
				if err != nil {
					l.Warn("failed to check required item >%v<", err)
				}
				action.RequiredItemName = itemName
				action.HasRequiredItem = hasItem
				action.IsAvailable = hasItem
			}

			actionMap[at] = action
		}

		// Convert map to ordered slice.
		var actions []turnsheet.LocationObjectActionOption
		for _, action := range actionMap {
			actions = append(actions, *action)
		}

		// Resolve state ID to display name.
		currentStateName := inst.CurrentAdventureGameLocationObjectStateID
		if inst.CurrentAdventureGameLocationObjectStateID != "" {
			if stateRec, err := p.Domain.GetAdventureGameLocationObjectStateRec(inst.CurrentAdventureGameLocationObjectStateID, nil); err == nil {
				currentStateName = stateRec.Name
			} else {
				l.Warn("failed to resolve state ID >%s< to name: %v", inst.CurrentAdventureGameLocationObjectStateID, err)
			}
		}

		result = append(result, turnsheet.LocationObjectOption{
			ObjectInstanceID: inst.ID,
			Name:             objectDef.Name,
			Description:      objectDef.Description,
			CurrentState:     currentStateName,
			Actions:          actions,
		})
	}

	return result, nil
}

// characterHasItemForEffect returns (hasItem, itemName, error).
func (p *AdventureGameLocationChoiceProcessor) characterHasItemForEffect(characterInstanceID, itemID string) (bool, string, error) {
	itemDef, err := p.Domain.GetAdventureGameItemRec(itemID, nil)
	if err != nil {
		return false, "", fmt.Errorf("failed to get item definition: %w", err)
	}

	itemInstances, err := p.Domain.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceID},
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: itemID},
		},
	})
	if err != nil {
		return false, itemDef.Name, fmt.Errorf("failed to query item instances: %w", err)
	}

	for _, inst := range itemInstances {
		if !inst.IsUsed {
			return true, itemDef.Name, nil
		}
	}
	return false, itemDef.Name, nil
}

// applyObjectEffect applies a single effect and returns its result_description.
func (p *AdventureGameLocationChoiceProcessor) applyObjectEffect(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	objectInstance *adventure_game_record.AdventureGameLocationObjectInstance,
	effect *adventure_game_record.AdventureGameLocationObjectEffect,
) (string, error) {
	l = l.WithFunctionContext("applyObjectEffect")
	l.Info("applying effect >%s< type >%s<", effect.ID, effect.EffectType)

	switch effect.EffectType {
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeInfo,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeNothing:
		// no state change — just return the description

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeObjectState,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRevealObject,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHideObject:
		if err := p.applyObjectStateEffect(l, gameInstanceRec, objectInstance, effect); err != nil {
			return "", err
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveItem,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeOpenLink,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeCloseLink:
		if err := p.applyObjectInventoryEffect(l, gameInstanceRec, characterInstanceRec, effect); err != nil {
			return "", err
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeDamage,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHeal:
		p.applyObjectHealthEffect(l, characterInstanceRec, effect)

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeSummonCreature,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeTeleport,
		adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveObject:
		if err := p.applyObjectWorldEffect(l, gameInstanceRec, characterInstanceRec, objectInstance, effect); err != nil {
			return "", err
		}
	}

	return effect.ResultDescription, nil
}

// applyObjectStateEffect handles change_state, change_object_state, reveal_object,
// and hide_object effects.
func (p *AdventureGameLocationChoiceProcessor) applyObjectStateEffect(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	objectInstance *adventure_game_record.AdventureGameLocationObjectInstance,
	effect *adventure_game_record.AdventureGameLocationObjectEffect,
) error {
	switch effect.EffectType {
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeState:
		if !effect.ResultAdventureGameLocationObjectStateID.Valid || effect.ResultAdventureGameLocationObjectStateID.String == "" {
			return nil
		}
		objectInstance.CurrentAdventureGameLocationObjectStateID = effect.ResultAdventureGameLocationObjectStateID.String
		updatedInst, err := p.Domain.UpdateAdventureGameLocationObjectInstanceRec(objectInstance)
		if err != nil {
			return fmt.Errorf("failed to update object instance state: %w", err)
		}
		*objectInstance = *updatedInst
		l.Info("object instance >%s< state changed to >%s<", objectInstance.ID, objectInstance.CurrentAdventureGameLocationObjectStateID)

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeChangeObjectState:
		if !effect.ResultAdventureGameLocationObjectID.Valid || !effect.ResultAdventureGameLocationObjectStateID.Valid {
			return nil
		}
		targets, err := p.getTargetObjectInstances(gameInstanceRec.ID, effect.ResultAdventureGameLocationObjectID.String)
		if err != nil {
			return fmt.Errorf("failed to get target object instances: %w", err)
		}
		for _, t := range targets {
			t.CurrentAdventureGameLocationObjectStateID = effect.ResultAdventureGameLocationObjectStateID.String
			if _, err := p.Domain.UpdateAdventureGameLocationObjectInstanceRec(t); err != nil {
				l.Warn("failed to update target object instance state >%v<", err)
			}
			l.Info("target object instance >%s< state changed to >%s<", t.ID, t.CurrentAdventureGameLocationObjectStateID)
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRevealObject:
		if !effect.ResultAdventureGameLocationObjectID.Valid {
			return nil
		}
		targets, err := p.getTargetObjectInstances(gameInstanceRec.ID, effect.ResultAdventureGameLocationObjectID.String)
		if err != nil {
			return fmt.Errorf("failed to get target object instances: %w", err)
		}
		for _, t := range targets {
			t.IsVisible = true
			if _, err := p.Domain.UpdateAdventureGameLocationObjectInstanceRec(t); err != nil {
				l.Warn("failed to reveal target object instance >%v<", err)
			}
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeHideObject:
		if !effect.ResultAdventureGameLocationObjectID.Valid {
			return nil
		}
		targets, err := p.getTargetObjectInstances(gameInstanceRec.ID, effect.ResultAdventureGameLocationObjectID.String)
		if err != nil {
			return fmt.Errorf("failed to get target object instances: %w", err)
		}
		for _, t := range targets {
			t.IsVisible = false
			if _, err := p.Domain.UpdateAdventureGameLocationObjectInstanceRec(t); err != nil {
				l.Warn("failed to hide target object instance >%v<", err)
			}
		}
	}
	return nil
}

// applyObjectInventoryEffect handles give_item, remove_item, open_link, and
// close_link effects.
func (p *AdventureGameLocationChoiceProcessor) applyObjectInventoryEffect(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	effect *adventure_game_record.AdventureGameLocationObjectEffect,
) error {
	switch effect.EffectType {
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeGiveItem:
		if !effect.ResultAdventureGameItemID.Valid || effect.ResultAdventureGameItemID.String == "" {
			return nil
		}
		itemInstance := &adventure_game_record.AdventureGameItemInstance{
			GameID:                           gameInstanceRec.GameID,
			GameInstanceID:                   gameInstanceRec.ID,
			AdventureGameItemID:              effect.ResultAdventureGameItemID.String,
			AdventureGameCharacterInstanceID: nullstring.FromString(characterInstanceRec.ID),
		}
		if _, err := p.Domain.CreateAdventureGameItemInstanceRec(itemInstance); err != nil {
			return fmt.Errorf("failed to give item: %w", err)
		}
		l.Info("gave item >%s< to character >%s<", effect.ResultAdventureGameItemID.String, characterInstanceRec.ID)

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveItem:
		if !effect.ResultAdventureGameItemID.Valid || effect.ResultAdventureGameItemID.String == "" {
			return nil
		}
		itemInstances, err := p.Domain.GetManyAdventureGameItemInstanceRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceRec.ID},
				{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameItemID, Val: effect.ResultAdventureGameItemID.String},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to query item instances: %w", err)
		}
		for _, inst := range itemInstances {
			if !inst.IsUsed {
				if err := p.Domain.DeleteAdventureGameItemInstanceRec(inst.ID); err != nil {
					l.Warn("failed to remove item instance >%v<", err)
				}
				break
			}
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeOpenLink:
		if !effect.ResultAdventureGameLocationLinkID.Valid || effect.ResultAdventureGameLocationLinkID.String == "" {
			return nil
		}
		requirements, err := p.Domain.GetManyAdventureGameLocationLinkRequirementRecs(&coresql.Options{
			Params: []coresql.Param{
				{Col: adventure_game_record.FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID, Val: effect.ResultAdventureGameLocationLinkID.String},
				{Col: adventure_game_record.FieldAdventureGameLocationLinkRequirementPurpose, Val: adventure_game_record.AdventureGameLocationLinkRequirementPurposeTraverse},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get link requirements: %w", err)
		}
		for _, req := range requirements {
			if err := p.Domain.DeleteAdventureGameLocationLinkRequirementRec(req.ID); err != nil {
				l.Warn("failed to remove link requirement >%v<", err)
			}
		}

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeCloseLink:
		if !effect.ResultAdventureGameLocationLinkID.Valid || effect.ResultAdventureGameLocationLinkID.String == "" {
			return nil
		}
		linkID := effect.ResultAdventureGameLocationLinkID.String

		// Retrieve the link record to get its game_id.
		linkRec, err := p.Domain.GetAdventureGameLocationLinkRec(linkID, nil)
		if err != nil {
			return fmt.Errorf("failed to get link record for close_link: %w", err)
		}

		// Build the new requirement based on whichever result reference is specified.
		// If an item is provided, the link now requires it to be in-inventory to traverse.
		// If a creature is provided, the link now requires the creature to be dead to traverse.
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
			// No item or creature specified — cannot create a valid requirement; skip.
			l.Warn("close_link effect >%s< has no result item or creature — link cannot be locked", effect.ID)
			return nil
		}

		if _, err := p.Domain.CreateAdventureGameLocationLinkRequirementRec(req); err != nil {
			return fmt.Errorf("failed to create link requirement for close_link: %w", err)
		}
		l.Info("link >%s< closed via requirement (item=%v creature=%v)", linkID, hasItem, hasCreature)
	}
	return nil
}

// applyObjectHealthEffect handles damage and heal effects, applying a random amount
// within the configured min/max range.
func (p *AdventureGameLocationChoiceProcessor) applyObjectHealthEffect(
	l logger.Logger,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	effect *adventure_game_record.AdventureGameLocationObjectEffect,
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
	if effect.EffectType == adventure_game_record.AdventureGameLocationObjectEffectEffectTypeDamage {
		characterInstanceRec.Health -= amount
		if characterInstanceRec.Health < 0 {
			characterInstanceRec.Health = 0
		}
		l.Info("object dealt %d damage to character >%s< (health now %d)", amount, characterInstanceRec.ID, characterInstanceRec.Health)
	} else {
		characterInstanceRec.Health += amount
		if characterInstanceRec.Health > 100 {
			characterInstanceRec.Health = 100
		}
		l.Info("object healed %d for character >%s< (health now %d)", amount, characterInstanceRec.ID, characterInstanceRec.Health)
	}
}

// applyObjectWorldEffect handles summon_creature, teleport, and remove_object effects.
func (p *AdventureGameLocationChoiceProcessor) applyObjectWorldEffect(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	objectInstance *adventure_game_record.AdventureGameLocationObjectInstance,
	effect *adventure_game_record.AdventureGameLocationObjectEffect,
) error {
	switch effect.EffectType {
	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeSummonCreature:
		if !effect.ResultAdventureGameCreatureID.Valid || effect.ResultAdventureGameCreatureID.String == "" {
			return nil
		}
		creatureRec, err := p.Domain.GetAdventureGameCreatureRec(effect.ResultAdventureGameCreatureID.String, nil)
		if err != nil {
			return fmt.Errorf("failed to get creature definition for summon: %w", err)
		}
		creatureInstance := &adventure_game_record.AdventureGameCreatureInstance{
			GameID:                          gameInstanceRec.GameID,
			GameInstanceID:                  gameInstanceRec.ID,
			AdventureGameCreatureID:         effect.ResultAdventureGameCreatureID.String,
			AdventureGameLocationInstanceID: objectInstance.AdventureGameLocationInstanceID,
			Health:                          creatureRec.MaxHealth,
		}
		if _, err := p.Domain.CreateAdventureGameCreatureInstanceRec(creatureInstance); err != nil {
			return fmt.Errorf("failed to summon creature: %w", err)
		}
		l.Info("summoned creature >%s< at location >%s<", effect.ResultAdventureGameCreatureID.String, objectInstance.AdventureGameLocationInstanceID)

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeTeleport:
		if !effect.ResultAdventureGameLocationID.Valid || effect.ResultAdventureGameLocationID.String == "" {
			return nil
		}
		destLocInst, err := p.getLocationInstanceForLocation(gameInstanceRec.ID, effect.ResultAdventureGameLocationID.String)
		if err != nil {
			return fmt.Errorf("failed to get destination location instance: %w", err)
		}
		characterInstanceRec.AdventureGameLocationInstanceID = destLocInst.ID
		updatedChar, err := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec)
		if err != nil {
			return fmt.Errorf("failed to teleport character: %w", err)
		}
		*characterInstanceRec = *updatedChar
		l.Info("teleported character >%s< to location instance >%s<", characterInstanceRec.ID, destLocInst.ID)

	case adventure_game_record.AdventureGameLocationObjectEffectEffectTypeRemoveObject:
		if err := p.Domain.DeleteAdventureGameLocationObjectInstanceRec(objectInstance.ID); err != nil {
			return fmt.Errorf("failed to remove object instance: %w", err)
		}
		l.Info("removed object instance >%s<", objectInstance.ID)
	}
	return nil
}

// getTargetObjectInstances returns all object instances for a given object definition
// within a game instance.
func (p *AdventureGameLocationChoiceProcessor) getTargetObjectInstances(gameInstanceID, objectDefID string) ([]*adventure_game_record.AdventureGameLocationObjectInstance, error) {
	return p.Domain.GetManyAdventureGameLocationObjectInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameLocationObjectInstanceAdventureGameLocationObjectID, Val: objectDefID},
		},
	})
}

// getLocationInstanceForLocation finds the location instance for a given game and
// location definition.
func (p *AdventureGameLocationChoiceProcessor) getLocationInstanceForLocation(gameInstanceID, locationID string) (*adventure_game_record.AdventureGameLocationInstance, error) {
	locationInstances, err := p.Domain.GetManyAdventureGameLocationInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameLocationInstanceAdventureGameLocationID, Val: locationID},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(locationInstances) == 0 {
		return nil, fmt.Errorf("no location instance found for location >%s<", locationID)
	}
	return locationInstances[0], nil
}

// processObjectChoice parses "{instance_id}:{action_type}", loads all matching
// effects, validates required items, and applies effects atomically.
func (p *AdventureGameLocationChoiceProcessor) processObjectChoice(
	l logger.Logger,
	gameInstanceRec *game_record.GameInstance,
	characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance,
	objectChoice string,
) error {
	l = l.WithFunctionContext("processObjectChoice")

	parts := strings.SplitN(objectChoice, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid object_choice format: %q", objectChoice)
	}
	instanceID := parts[0]
	actionType := parts[1]

	// Load object instance.
	objectInstance, err := p.Domain.GetAdventureGameLocationObjectInstanceRec(instanceID, nil)
	if err != nil {
		return fmt.Errorf("failed to get object instance >%s<: %w", instanceID, err)
	}

	// Load all effects for this object and action type.
	allEffects, err := p.Domain.GetManyAdventureGameLocationObjectEffectRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameLocationObjectEffectAdventureGameLocationObjectID, Val: objectInstance.AdventureGameLocationObjectID},
			{Col: adventure_game_record.FieldAdventureGameLocationObjectEffectActionType, Val: actionType},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get effects for object: %w", err)
	}

	// Filter effects matching current state.
	var matchingEffects []*adventure_game_record.AdventureGameLocationObjectEffect
	for _, effect := range allEffects {
		if effect.RequiredAdventureGameLocationObjectStateID.Valid && effect.RequiredAdventureGameLocationObjectStateID.String != objectInstance.CurrentAdventureGameLocationObjectStateID {
			continue
		}
		matchingEffects = append(matchingEffects, effect)
	}

	if len(matchingEffects) == 0 {
		l.Info("no matching effects for object >%s< action >%s< state >%s<", instanceID, actionType, objectInstance.CurrentAdventureGameLocationObjectStateID)
		return nil
	}

	// Validate required items — any effect's required item blocks all effects.
	for _, effect := range matchingEffects {
		if !effect.RequiredAdventureGameItemID.Valid || effect.RequiredAdventureGameItemID.String == "" {
			continue
		}
		hasItem, _, err := p.characterHasItemForEffect(characterInstanceRec.ID, effect.RequiredAdventureGameItemID.String)
		if err != nil {
			return fmt.Errorf("failed to check required item: %w", err)
		}
		if !hasItem {
			l.Info("character does not have required item for object interaction")
			return fmt.Errorf("required item not in inventory")
		}
	}

	// Apply all matching effects atomically.
	var resultDescriptions []string
	for _, effect := range matchingEffects {
		desc, err := p.applyObjectEffect(l, gameInstanceRec, characterInstanceRec, objectInstance, effect)
		if err != nil {
			l.Warn("failed to apply effect >%s< >%v<", effect.ID, err)
			return fmt.Errorf("failed to apply object effect: %w", err)
		}
		if desc != "" {
			resultDescriptions = append(resultDescriptions, desc)
		}
	}

	// Append a combined world event with all result descriptions.
	if len(resultDescriptions) > 0 {
		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryWorld,
			Icon:     turnsheet.TurnEventIconWorld,
			Message:  strings.Join(resultDescriptions, " "),
		})
		if _, saveErr := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec); saveErr != nil {
			l.Warn("failed to save object interaction events >%v<", saveErr)
		}
	}

	return nil
}
