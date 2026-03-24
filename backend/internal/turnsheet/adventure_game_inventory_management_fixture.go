package turnsheet

import (
	"time"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// AdventureGameInventoryManagementFixture returns the sample rendering fixture for the
// adventure game inventory management turn sheet.
func AdventureGameInventoryManagementFixture() DevFixture {
	return DevFixture{
		TemplatePath:   "turnsheet/adventure_game_inventory_management.template",
		OutputBaseName: "adventure_game_inventory_management_turnsheet",
		BackgroundFile: "background-dungeon.png",
		MakeData: func(bg, code string) any {
			deadline := time.Now().Add(7 * 24 * time.Hour)
			return &InventoryManagementData{
				TurnSheetTemplateData: TurnSheetTemplateData{
					GameName:              strPtr("The Enchanted Forest Adventure"),
					GameType:              strPtr("adventure"),
					TurnNumber:            intPtr(1),
					AccountName:           strPtr("Test Player"),
					TurnSheetTitle:        strPtr("Inventory Management"),
					TurnSheetInstructions: strPtr(DefaultInventoryManagementInstructions()),
					TurnSheetCode:         strPtr(code),
					TurnSheetDeadline:     &deadline,
					BackgroundImage:       &bg,
					TurnEvents: []TurnEvent{
						{Category: TurnEventCategorySystem, Icon: TurnEventIconSystem, Message: "Aria looted the chest in the Crystal Caverns."},
						{Category: TurnEventCategorySystem, Icon: TurnEventIconSystem, Message: "Found a Healing Potion and a Scroll of Light."},
						{Category: TurnEventCategorySystem, Icon: TurnEventIconSystem, Message: "Dropped the Rusty Dagger to make space."},
					},
				},
				CharacterName:       "Aria the Mage",
				CurrentLocationName: "Mystic Grove",
				InventoryCapacity:   10,
				InventoryCount:      5,
				CurrentInventory: []InventoryItem{
					{ItemInstanceID: "item-1", ItemName: "Crystal Key", ItemDescription: "A glowing crystal key that hums with ancient magic", IsEquipped: false, CanEquip: false},
					{ItemInstanceID: "item-2", ItemName: "Iron Sword", ItemDescription: "A sturdy iron sword with a leather-wrapped hilt", IsEquipped: true, EquipmentSlot: "weapon", CanEquip: true},
					{ItemInstanceID: "item-3", ItemName: "Leather Armor", ItemDescription: "Basic leather protection, worn but serviceable", IsEquipped: true, EquipmentSlot: "armor", CanEquip: true},
					{ItemInstanceID: "item-4", ItemName: "Healing Potion", ItemDescription: "A red potion that restores health when consumed", IsEquipped: false, CanEquip: false},
					{ItemInstanceID: "item-5", ItemName: "Magic Ring", ItemDescription: "A silver ring imbued with protective magic", IsEquipped: false, CanEquip: true},
				},
				EquipmentSlots: EquipmentSlots{
					Weapon: &EquippedItem{ItemInstanceID: "item-2", ItemName: "Iron Sword", SlotName: "weapon"},
					Armor:  &EquippedItem{ItemInstanceID: "item-3", ItemName: "Leather Armor", SlotName: "armor"},
				},
				LocationItems: []LocationItem{
					{ItemInstanceID: "item-6", ItemName: "Shadow Cloak", ItemDescription: "A dark cloak that seems to blend with shadows", CanEquip: true},
					{ItemInstanceID: "item-7", ItemName: "Wind Charm", ItemDescription: "A small charm that whispers with the wind", CanEquip: false},
				},
			}
		},
		NewProcessor: func(l logger.Logger, cfg config.Config) (TurnSheetProcessor, error) {
			return NewInventoryManagementProcessor(l, cfg)
		},
	}
}
