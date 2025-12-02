package domain

import (
	"database/sql"

	"gitlab.com/alienspaces/playbymail/core/domain"
	coresql "gitlab.com/alienspaces/playbymail/core/sql"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// GetAdventureGameItemInstanceRecsByCharacterInstance gets all item instances in a character's inventory
func (m *Domain) GetAdventureGameItemInstanceRecsByCharacterInstance(characterInstanceID string) ([]*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("GetAdventureGameItemInstanceRecsByCharacterInstance")

	l.Debug("getting item instances for character instance ID >%s<", characterInstanceID)

	if err := domain.ValidateUUIDField("character_instance_id", characterInstanceID); err != nil {
		return nil, err
	}

	opts := &coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceID},
		},
	}

	return m.GetManyAdventureGameItemInstanceRecs(opts)
}

// GetAdventureGameItemInstanceRecsByLocationInstance gets all item instances at a location
func (m *Domain) GetAdventureGameItemInstanceRecsByLocationInstance(locationInstanceID string) ([]*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("GetAdventureGameItemInstanceRecsByLocationInstance")

	l.Debug("getting item instances for location instance ID >%s<", locationInstanceID)

	if err := domain.ValidateUUIDField("location_instance_id", locationInstanceID); err != nil {
		return nil, err
	}

	opts := &coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameLocationInstanceID, Val: locationInstanceID},
		},
	}

	return m.GetManyAdventureGameItemInstanceRecs(opts)
}

// PickUpAdventureGameItemInstanceRec moves an item instance from a location to a character's inventory
func (m *Domain) PickUpAdventureGameItemInstanceRec(characterInstanceID, itemInstanceID string) (*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("PickUpAdventureGameItemInstanceRec")

	l.Debug("picking up item instance >%s< for character instance >%s<", itemInstanceID, characterInstanceID)

	// Validate inputs
	if err := domain.ValidateUUIDField("character_instance_id", characterInstanceID); err != nil {
		return nil, err
	}
	if err := domain.ValidateUUIDField("item_instance_id", itemInstanceID); err != nil {
		return nil, err
	}

	// Get character instance to check inventory capacity
	characterRec, err := m.GetAdventureGameCharacterInstanceRec(characterInstanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	// Get item instance with lock
	itemRec, err := m.GetAdventureGameItemInstanceRec(itemInstanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	// Validate item is at a location (not already in inventory)
	if !itemRec.AdventureGameLocationInstanceID.Valid {
		return nil, InvalidField("item_instance", itemInstanceID, "item is not at a location")
	}

	// Check inventory capacity
	inventoryItems, err := m.GetAdventureGameItemInstanceRecsByCharacterInstance(characterInstanceID)
	if err != nil {
		return nil, err
	}

	if len(inventoryItems) >= characterRec.InventoryCapacity {
		return nil, InvalidField("inventory", characterInstanceID, "inventory is at capacity")
	}

	// Move item to character inventory
	itemRec.AdventureGameCharacterInstanceID = nullstring.FromString(characterInstanceID)
	itemRec.AdventureGameLocationInstanceID = sql.NullString{Valid: false}

	updatedRec, err := m.UpdateAdventureGameItemInstanceRec(itemRec)
	if err != nil {
		return nil, err
	}

	return updatedRec, nil
}

// DropAdventureGameItemInstanceRec moves an item instance from a character's inventory to their current location
func (m *Domain) DropAdventureGameItemInstanceRec(characterInstanceID, itemInstanceID string) (*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("DropAdventureGameItemInstanceRec")

	l.Debug("dropping item instance >%s< from character instance >%s<", itemInstanceID, characterInstanceID)

	// Validate inputs
	if err := domain.ValidateUUIDField("character_instance_id", characterInstanceID); err != nil {
		return nil, err
	}
	if err := domain.ValidateUUIDField("item_instance_id", itemInstanceID); err != nil {
		return nil, err
	}

	// Get character instance to get current location
	characterRec, err := m.GetAdventureGameCharacterInstanceRec(characterInstanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	if characterRec.AdventureGameLocationInstanceID == "" {
		return nil, InvalidField("character_instance", characterInstanceID, "character is not at a location")
	}

	// Get item instance with lock
	itemRec, err := m.GetAdventureGameItemInstanceRec(itemInstanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	// Validate item is in character's inventory
	if !itemRec.AdventureGameCharacterInstanceID.Valid || itemRec.AdventureGameCharacterInstanceID.String != characterInstanceID {
		return nil, InvalidField("item_instance", itemInstanceID, "item is not in character's inventory")
	}

	// If equipped, unequip first
	if itemRec.IsEquipped {
		itemRec.IsEquipped = false
		itemRec.EquipmentSlot = sql.NullString{Valid: false}
	}

	// Move item to location
	itemRec.AdventureGameLocationInstanceID = nullstring.FromString(characterRec.AdventureGameLocationInstanceID)
	itemRec.AdventureGameCharacterInstanceID = sql.NullString{Valid: false}

	updatedRec, err := m.UpdateAdventureGameItemInstanceRec(itemRec)
	if err != nil {
		return nil, err
	}

	return updatedRec, nil
}

// EquipAdventureGameItemInstanceRec equips an item instance to a character
func (m *Domain) EquipAdventureGameItemInstanceRec(characterInstanceID, itemInstanceID, equipmentSlot string) (*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("EquipAdventureGameItemInstanceRec")

	l.Debug("equipping item instance >%s< to character instance >%s< in slot >%s<", itemInstanceID, characterInstanceID, equipmentSlot)

	// Validate inputs
	if err := domain.ValidateUUIDField("character_instance_id", characterInstanceID); err != nil {
		return nil, err
	}
	if err := domain.ValidateUUIDField("item_instance_id", itemInstanceID); err != nil {
		return nil, err
	}
	if equipmentSlot == "" {
		return nil, RequiredField("equipment_slot")
	}

	// Get item instance with lock
	itemRec, err := m.GetAdventureGameItemInstanceRec(itemInstanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	// Validate item is in character's inventory
	if !itemRec.AdventureGameCharacterInstanceID.Valid || itemRec.AdventureGameCharacterInstanceID.String != characterInstanceID {
		return nil, InvalidField("item_instance", itemInstanceID, "item is not in character's inventory")
	}

	// Get item definition to check if it can be equipped
	itemDef, err := m.GetAdventureGameItemRec(itemRec.AdventureGameItemID, nil)
	if err != nil {
		return nil, err
	}

	if !itemDef.CanBeEquipped {
		return nil, InvalidField("item", itemRec.AdventureGameItemID, "item cannot be equipped")
	}

	// Validate equipment slot matches item's equipment slot
	if itemDef.EquipmentSlot != nil && *itemDef.EquipmentSlot != equipmentSlot {
		return nil, InvalidField("equipment_slot", equipmentSlot, "equipment slot does not match item's slot")
	}

	// Check if slot is already occupied
	// Note: We need to query for items with matching equipment slot
	// Since equipment_slot is nullable, we'll get all equipped items and filter
	opts := &coresql.Options{
		Params: []coresql.Param{
			{Col: adventure_game_record.FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID, Val: characterInstanceID},
			{Col: adventure_game_record.FieldAdventureGameItemInstanceIsEquipped, Val: true},
		},
		Limit: 100, // Get all equipped items to check for slot conflicts
	}

	equippedItems, err := m.GetManyAdventureGameItemInstanceRecs(opts)
	if err != nil {
		return nil, err
	}

	// Filter for items in the same slot
	var slotOccupiedItems []*adventure_game_record.AdventureGameItemInstance
	for _, item := range equippedItems {
		if item.ID != itemInstanceID && item.EquipmentSlot.Valid && item.EquipmentSlot.String == equipmentSlot {
			slotOccupiedItems = append(slotOccupiedItems, item)
		}
	}

	// If slot is occupied, unequip existing item
	if len(slotOccupiedItems) > 0 {
		for _, equippedItem := range slotOccupiedItems {
			equippedItem.IsEquipped = false
			equippedItem.EquipmentSlot = sql.NullString{Valid: false}
			_, err := m.UpdateAdventureGameItemInstanceRec(equippedItem)
			if err != nil {
				return nil, err
			}
		}
	}

	// Equip the item
	itemRec.IsEquipped = true
	itemRec.EquipmentSlot = nullstring.FromString(equipmentSlot)

	updatedRec, err := m.UpdateAdventureGameItemInstanceRec(itemRec)
	if err != nil {
		return nil, err
	}

	return updatedRec, nil
}

// UnequipAdventureGameItemInstanceRec unequips an item instance from a character
func (m *Domain) UnequipAdventureGameItemInstanceRec(characterInstanceID, itemInstanceID string) (*adventure_game_record.AdventureGameItemInstance, error) {
	l := m.Logger("UnequipAdventureGameItemInstanceRec")

	l.Debug("unequipping item instance >%s< from character instance >%s<", itemInstanceID, characterInstanceID)

	// Validate inputs
	if err := domain.ValidateUUIDField("character_instance_id", characterInstanceID); err != nil {
		return nil, err
	}
	if err := domain.ValidateUUIDField("item_instance_id", itemInstanceID); err != nil {
		return nil, err
	}

	// Get item instance with lock
	itemRec, err := m.GetAdventureGameItemInstanceRec(itemInstanceID, coresql.ForUpdateNoWait)
	if err != nil {
		return nil, err
	}

	// Validate item is in character's inventory
	if !itemRec.AdventureGameCharacterInstanceID.Valid || itemRec.AdventureGameCharacterInstanceID.String != characterInstanceID {
		return nil, InvalidField("item_instance", itemInstanceID, "item is not in character's inventory")
	}

	// Validate item is equipped
	if !itemRec.IsEquipped {
		return nil, InvalidField("item_instance", itemInstanceID, "item is not equipped")
	}

	// Unequip the item
	itemRec.IsEquipped = false
	itemRec.EquipmentSlot = sql.NullString{Valid: false}

	updatedRec, err := m.UpdateAdventureGameItemInstanceRec(itemRec)
	if err != nil {
		return nil, err
	}

	return updatedRec, nil
}

