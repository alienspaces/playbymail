package turn_sheet_processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
)

// AdventureGameInventoryManagementProcessor processes inventory management turn sheet business logic for adventure games
type AdventureGameInventoryManagementProcessor struct {
	Logger logger.Logger
	Domain *domain.Domain
}

// NewAdventureGameInventoryManagementProcessor creates a new adventure game inventory management processor
func NewAdventureGameInventoryManagementProcessor(l logger.Logger, d *domain.Domain) (*AdventureGameInventoryManagementProcessor, error) {
	l = l.WithFunctionContext("NewAdventureGameInventoryManagementProcessor")

	p := &AdventureGameInventoryManagementProcessor{
		Logger: l,
		Domain: d,
	}
	return p, nil
}

// GetSheetType returns the sheet type this processor handles (implements TurnSheetProcessor interface)
func (p *AdventureGameInventoryManagementProcessor) GetSheetType() string {
	return adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement
}

// ProcessTurnSheetResponse processes a single turn sheet response (implements TurnSheetProcessor interface)
func (p *AdventureGameInventoryManagementProcessor) ProcessTurnSheetResponse(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance, turnSheet *game_record.GameTurnSheet) error {
	l := p.Logger.WithFunctionContext("AdventureGameInventoryManagementProcessor/ProcessTurnSheetResponse")

	l.Info("processing inventory management for turn sheet >%s< for character >%s<", turnSheet.ID, characterInstanceRec.ID)

	// Verify this is an inventory management sheet
	if turnSheet.SheetType != adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement {
		l.Warn("expected inventory management sheet type, got >%s<", turnSheet.SheetType)
		return fmt.Errorf("invalid sheet type: expected %s, got %s", adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement, turnSheet.SheetType)
	}

	// Step 1: Parse the player's inventory actions from ScannedData
	var scanData turnsheet.InventoryManagementScanData
	if err := json.Unmarshal(turnSheet.ScannedData, &scanData); err != nil {
		l.Warn("failed to unmarshal scanned data >%v<", err)
		return fmt.Errorf("failed to parse scanned data: %w", err)
	}

	// Step 2: Process actions in order: unequip → drop → pick up → equip
	// This ensures items are properly moved before equipping

	// Unequip items first
	for _, itemInstanceID := range scanData.Unequip {
		l.Info("unequipping item >%s<", itemInstanceID)
		_, err := p.Domain.UnequipAdventureGameItemInstanceRec(characterInstanceRec.ID, itemInstanceID)
		if err != nil {
			l.Warn("failed to unequip item >%s< >%v<", itemInstanceID, err)
			return fmt.Errorf("failed to unequip item %s: %w", itemInstanceID, err)
		}
	}

	// Drop items
	for _, itemInstanceID := range scanData.Drop {
		l.Info("dropping item >%s<", itemInstanceID)
		_, err := p.Domain.DropAdventureGameItemInstanceRec(characterInstanceRec.ID, itemInstanceID)
		if err != nil {
			l.Warn("failed to drop item >%s< >%v<", itemInstanceID, err)
			return fmt.Errorf("failed to drop item %s: %w", itemInstanceID, err)
		}
	}

	// Pick up items
	for _, itemInstanceID := range scanData.PickUp {
		l.Info("picking up item >%s<", itemInstanceID)
		_, err := p.Domain.PickUpAdventureGameItemInstanceRec(characterInstanceRec.ID, itemInstanceID)
		if err != nil {
			l.Warn("failed to pick up item >%s< >%v<", itemInstanceID, err)
			return fmt.Errorf("failed to pick up item %s: %w", itemInstanceID, err)
		}
	}

	// Equip items last (after pick up, so items are in inventory)
	for _, action := range scanData.Equip {
		l.Info("equipping item >%s< to slot >%s<", action.ItemInstanceID, action.Slot)
		_, err := p.Domain.EquipAdventureGameItemInstanceRec(characterInstanceRec.ID, action.ItemInstanceID, action.Slot)
		if err != nil {
			l.Warn("failed to equip item >%s< >%v<", action.ItemInstanceID, err)
			return fmt.Errorf("failed to equip item %s: %w", action.ItemInstanceID, err)
		}
	}

	l.Info("successfully processed inventory management actions for character >%s<", characterInstanceRec.ID)

	return nil
}

// CreateNextTurnSheet creates a new inventory management turn sheet for a character (implements TurnSheetProcessor interface)
func (p *AdventureGameInventoryManagementProcessor) CreateNextTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("AdventureGameInventoryManagementProcessor/CreateNextTurnSheet")

	l.Info("creating inventory management turn sheet for character >%s<", characterInstanceRec.ID)

	// Step 1: Get character's current location instance
	locationInstanceRec, err := p.Domain.GetAdventureGameLocationInstanceRec(characterInstanceRec.AdventureGameLocationInstanceID, nil)
	if err != nil {
		l.Warn("failed to get character's current location >%v<", err)
		return nil, fmt.Errorf("failed to get character's current location: %w", err)
	}

	// Step 2: Get the location definition
	locationRec, err := p.Domain.GetAdventureGameLocationRec(locationInstanceRec.AdventureGameLocationID, nil)
	if err != nil {
		l.Warn("failed to get location definition >%v<", err)
		return nil, fmt.Errorf("failed to get location definition: %w", err)
	}

	// Step 3: Get character definition
	characterRec, err := p.Domain.GetAdventureGameCharacterRec(characterInstanceRec.AdventureGameCharacterID, nil)
	if err != nil {
		l.Warn("failed to get character >%v<", err)
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	// Step 4: Get account for name
	accountRec, err := p.Domain.GetAccountRec(characterRec.AccountUserID, nil)
	if err != nil {
		l.Warn("failed to get account >%v<", err)
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Step 5: Get game for game name
	gameRec, err := p.Domain.GetGameRec(gameInstanceRec.GameID, nil)
	if err != nil {
		l.Warn("failed to get game >%v<", err)
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Step 6: Get character's inventory
	inventoryItems, err := p.Domain.GetAdventureGameItemInstanceRecsByCharacterInstance(characterInstanceRec.ID)
	if err != nil {
		l.Warn("failed to get character inventory >%v<", err)
		return nil, fmt.Errorf("failed to get character inventory: %w", err)
	}

	// Step 7: Get items at current location
	locationItems, err := p.Domain.GetAdventureGameItemInstanceRecsByLocationInstance(locationInstanceRec.ID)
	if err != nil {
		l.Warn("failed to get location items >%v<", err)
		return nil, fmt.Errorf("failed to get location items: %w", err)
	}

	// Step 8: Build inventory items list with item definitions
	inventoryItemList := make([]turnsheet.InventoryItem, 0, len(inventoryItems))
	equipmentSlots := turnsheet.EquipmentSlots{}

	for _, itemInstance := range inventoryItems {
		// Get item definition
		itemDef, err := p.Domain.GetAdventureGameItemRec(itemInstance.AdventureGameItemID, nil)
		if err != nil {
			l.Warn("failed to get item definition >%s< >%v<", itemInstance.AdventureGameItemID, err)
			continue
		}

		// Determine equipment slot for display (simplified grouping)
		var displaySlot string
		if itemInstance.IsEquipped && itemInstance.EquipmentSlot.Valid {
			slot := itemInstance.EquipmentSlot.String
			// Map specific slots to display slots
			switch {
			case slot == "weapon":
				displaySlot = "weapon"
			case strings.HasPrefix(slot, "armor_"):
				displaySlot = "armor"
			case strings.HasPrefix(slot, "clothing_"):
				displaySlot = "clothing"
			case strings.HasPrefix(slot, "jewelry_"):
				displaySlot = "jewelry"
			default:
				displaySlot = slot
			}
		}

		inventoryItem := turnsheet.InventoryItem{
			ItemInstanceID:  itemInstance.ID,
			ItemName:        itemDef.Name,
			ItemDescription: itemDef.Description,
			IsEquipped:      itemInstance.IsEquipped,
			EquipmentSlot:   displaySlot,
			CanEquip:        itemDef.CanBeEquipped,
		}
		inventoryItemList = append(inventoryItemList, inventoryItem)
	}

	// Sort inventory items: equipped items first (weapon, armor, jewelry order), then unequipped
	equippedItems := make([]turnsheet.InventoryItem, 0)
	unequippedItems := make([]turnsheet.InventoryItem, 0)

	// Define slot priority for sorting equipped items
	slotPriority := map[string]int{
		"weapon":   1,
		"armor":    2,
		"clothing": 3,
		"jewelry":  4,
	}

	for _, item := range inventoryItemList {
		if item.IsEquipped {
			equippedItems = append(equippedItems, item)
		} else {
			unequippedItems = append(unequippedItems, item)
		}
	}

	// Sort equipped items by slot priority
	for i := 0; i < len(equippedItems)-1; i++ {
		for j := i + 1; j < len(equippedItems); j++ {
			priorityI := slotPriority[equippedItems[i].EquipmentSlot]
			priorityJ := slotPriority[equippedItems[j].EquipmentSlot]
			if priorityI == 0 {
				priorityI = 99 // Unknown slots go to end
			}
			if priorityJ == 0 {
				priorityJ = 99
			}
			if priorityI > priorityJ {
				equippedItems[i], equippedItems[j] = equippedItems[j], equippedItems[i]
			}
		}
	}

	// Combine: equipped first, then unequipped
	inventoryItemList = append(equippedItems, unequippedItems...)

	// Step 9: Build location items list
	locationItemList := make([]turnsheet.LocationItem, 0, len(locationItems))
	for _, itemInstance := range locationItems {
		// Get item definition
		itemDef, err := p.Domain.GetAdventureGameItemRec(itemInstance.AdventureGameItemID, nil)
		if err != nil {
			l.Warn("failed to get item definition >%s< >%v<", itemInstance.AdventureGameItemID, err)
			continue
		}

		locationItem := turnsheet.LocationItem{
			ItemInstanceID:  itemInstance.ID,
			ItemName:        itemDef.Name,
			ItemDescription: itemDef.Description,
		}
		locationItemList = append(locationItemList, locationItem)
	}

	// Step 10: Create sheet data
	sheetData := turnsheet.InventoryManagementData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr("adventure"),
			TurnNumber:            convert.Ptr(gameInstanceRec.CurrentTurn + 1),
			AccountName:           convert.Ptr(accountRec.Email),
			TurnSheetTitle:        convert.Ptr("Inventory Management"),
			TurnSheetDescription:  convert.Ptr(fmt.Sprintf("Manage your inventory and equipment. Carrying %d/%d items.", len(inventoryItemList), characterInstanceRec.InventoryCapacity)),
			TurnSheetInstructions: convert.Ptr(turnsheet.DefaultInventoryManagementInstructions()),
		},
		CharacterName:       characterRec.Name,
		CurrentLocationName: locationRec.Name,
		InventoryCapacity:   characterInstanceRec.InventoryCapacity,
		InventoryCount:      len(inventoryItemList),
		CurrentInventory:    inventoryItemList,
		EquipmentSlots:      equipmentSlots,
		LocationItems:       locationItemList,
	}

	sheetDataBytes, err := json.Marshal(sheetData)
	if err != nil {
		l.Warn("failed to marshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to marshal sheet data: %w", err)
	}

	// Step 11: Create turn sheet record
	turnSheet := &game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		AccountID:        accountRec.AccountID,
		AccountUserID:    characterRec.AccountUserID,
		TurnNumber:       gameInstanceRec.CurrentTurn + 1,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement,
		SheetOrder:       1,
		SheetData:        json.RawMessage(sheetDataBytes),
		IsCompleted:      false,
		ProcessingStatus: game_record.TurnSheetProcessingStatusPending,
	}
	turnSheet.GameInstanceID = sql.NullString{String: gameInstanceRec.ID, Valid: true}

	// Create the turn sheet record
	createdTurnSheetRec, err := p.Domain.CreateGameTurnSheetRec(turnSheet)
	if err != nil {
		l.Warn("failed to create turn sheet record >%v<", err)
		return nil, fmt.Errorf("failed to create turn sheet record: %w", err)
	}

	// Link it to the character via AdventureGameTurnSheet
	adventureTurnSheet := &adventure_game_record.AdventureGameTurnSheet{
		GameID:                           gameInstanceRec.GameID,
		AdventureGameCharacterInstanceID: characterInstanceRec.ID,
		GameTurnSheetID:                  createdTurnSheetRec.ID,
	}

	_, err = p.Domain.CreateAdventureGameTurnSheetRec(adventureTurnSheet)
	if err != nil {
		l.Warn("failed to create adventure game turn sheet record >%v<", err)
		return nil, fmt.Errorf("failed to create adventure game turn sheet record: %w", err)
	}

	l.Info("created inventory management turn sheet >%s< for character >%s< with %d inventory items and %d location items",
		createdTurnSheetRec.ID, characterInstanceRec.ID, len(inventoryItemList), len(locationItemList))

	return createdTurnSheetRec, nil
}
