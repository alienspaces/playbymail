package turn_sheet_processor

import (
	"fmt"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

const defaultUnarmedAttackDamage = 5

// ResolveEquipmentStats returns (weaponDamage, armorDefense) for a character
// based on their currently equipped items. Used by both the monster encounter
// processor and the location choice processor (flee penalty).
func ResolveEquipmentStats(l logger.Logger, d *domain.Domain, characterInstanceID string) (weaponDamage, armorDefense int, err error) {
	weaponDamage = defaultUnarmedAttackDamage
	armorDefense = 0

	inventoryItems, err := d.GetAdventureGameItemInstanceRecsByCharacterInstance(characterInstanceID)
	if err != nil {
		return weaponDamage, armorDefense, fmt.Errorf("failed to get inventory: %w", err)
	}

	for _, itemInstance := range inventoryItems {
		if !itemInstance.IsEquipped || !itemInstance.EquipmentSlot.Valid {
			continue
		}

		slot := itemInstance.EquipmentSlot.String

		itemDef, err := d.GetAdventureGameItemRec(itemInstance.AdventureGameItemID, nil)
		if err != nil {
			l.Warn("failed to get item definition >%s< >%v<", itemInstance.AdventureGameItemID, err)
			continue
		}

		if slot == adventure_game_record.AdventureGameItemEquipmentSlotWeapon {
			weaponDamage = itemDef.Damage
		} else if slot != "" {
			armorDefense += itemDef.Defense
		}
	}

	return weaponDamage, armorDefense, nil
}

// HasAggressiveCreaturesAtLocation returns true if any alive aggressive creature
// instances exist at the given location.
func HasAggressiveCreaturesAtLocation(d *domain.Domain, gameInstanceID, locationInstanceID string) (bool, error) {
	creatureInstances, err := d.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, Val: locationInstanceID},
		},
	})
	if err != nil {
		return false, fmt.Errorf("failed to get creature instances at location: %w", err)
	}

	for _, ci := range creatureInstances {
		if ci.Health <= 0 {
			continue
		}
		creatureDef, err := d.GetAdventureGameCreatureRec(ci.AdventureGameCreatureID, nil)
		if err != nil {
			continue
		}
		if creatureDef.Disposition == adventure_game_record.AdventureGameCreatureDispositionAggressive {
			return true, nil
		}
	}
	return false, nil
}

// GetAliveCreaturesAtLocation returns turnsheet creature entries for all alive
// creature instances at the given location.
func GetAliveCreaturesAtLocation(l logger.Logger, d *domain.Domain, gameInstanceID, locationInstanceID string) ([]turnsheet.LocationCreature, error) {
	creatureInstances, err := d.GetManyAdventureGameCreatureInstanceRecs(&coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceGameInstanceID, Val: gameInstanceID},
			{Col: adventure_game_record.FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID, Val: locationInstanceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query creature instances: %w", err)
	}

	var creatures []turnsheet.LocationCreature
	for _, inst := range creatureInstances {
		if inst.Health <= 0 {
			continue
		}
		creatureRec, err := d.GetAdventureGameCreatureRec(inst.AdventureGameCreatureID, nil)
		if err != nil {
			l.Warn("failed to get creature definition >%s< >%v<", inst.AdventureGameCreatureID, err)
			continue
		}
		creatures = append(creatures, turnsheet.LocationCreature{
			Name:        creatureRec.Name,
			Description: creatureRec.Description,
			Disposition: creatureRec.Disposition,
		})
	}
	return creatures, nil
}

// ReadTurnEventsForCategories reads turn events from the character instance without
// clearing them, filters to only the given categories, and excludes flee_context events.
// Call ClearTurnEvents once all processors have built their sheets.
func ReadTurnEventsForCategories(l logger.Logger, d *domain.Domain, characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance, categories ...string) []turnsheet.TurnEvent {
	events, err := turnsheet.ReadTurnEvents(characterInstanceRec)
	if err != nil {
		l.Warn("failed to read turn events >%v<", err)
		return nil
	}
	return turnsheet.FilterTurnEventsByCategory(events, categories...)
}

// ClearTurnEvents reads and clears all turn events from the character instance and
// persists the cleared state. Call this once after all processors have built their sheets.
func ClearTurnEvents(l logger.Logger, d *domain.Domain, characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance) {
	_, err := turnsheet.ReadAndClearTurnEvents(characterInstanceRec)
	if err != nil {
		l.Warn("failed to read turn events for clearing >%v<", err)
		return
	}
	if _, saveErr := d.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec); saveErr != nil {
		l.Warn("failed to clear turn events on character instance >%v<", saveErr)
	}
}
