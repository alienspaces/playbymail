package turn_sheet_processor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/turnsheet"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
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
	// This ensures items are properly moved before equipping.
	// charNeedsUpdate is set whenever the character record is mutated (health change or turn event
	// appended) so we only issue one UpdateAdventureGameCharacterInstanceRec at the end.
	charNeedsUpdate := false

	// Unequip items first
	for _, itemInstanceID := range scanData.Unequip {
		l.Info("unequipping item >%s<", itemInstanceID)
		name := p.resolveItemName(l, itemInstanceID)
		_, err := p.Domain.UnequipAdventureGameItemInstanceRec(characterInstanceRec.ID, itemInstanceID)
		if err != nil {
			l.Warn("failed to unequip item >%s< >%v<", itemInstanceID, err)
			return fmt.Errorf("failed to unequip item %s: %w", itemInstanceID, err)
		}
		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryInventory,
			Icon:     turnsheet.TurnEventIconInventory,
			Message:  fmt.Sprintf("You moved %s to your backpack.", name),
		})
		charNeedsUpdate = true
	}

	// Drop items
	for _, itemInstanceID := range scanData.Drop {
		l.Info("dropping item >%s<", itemInstanceID)
		name := p.resolveItemName(l, itemInstanceID)
		_, err := p.Domain.DropAdventureGameItemInstanceRec(characterInstanceRec.ID, itemInstanceID)
		if err != nil {
			l.Warn("failed to drop item >%s< >%v<", itemInstanceID, err)
			return fmt.Errorf("failed to drop item %s: %w", itemInstanceID, err)
		}
		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryInventory,
			Icon:     turnsheet.TurnEventIconInventory,
			Message:  fmt.Sprintf("You dropped %s.", name),
		})
		charNeedsUpdate = true
	}

	// Pick up items
	for _, itemInstanceID := range scanData.PickUp {
		l.Info("picking up item >%s<", itemInstanceID)
		name := p.resolveItemName(l, itemInstanceID)
		_, err := p.Domain.PickUpAdventureGameItemInstanceRec(characterInstanceRec.ID, itemInstanceID)
		if err != nil {
			l.Warn("failed to pick up item >%s< >%v<", itemInstanceID, err)
			return fmt.Errorf("failed to pick up item %s: %w", itemInstanceID, err)
		}
		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryInventory,
			Icon:     turnsheet.TurnEventIconInventory,
			Message:  fmt.Sprintf("You picked up %s.", name),
		})
		charNeedsUpdate = true
	}

	// Equip items last (after pick up, so items are in inventory).
	// If an equip action targets a location item (not yet in inventory), auto-pick it up first
	// so the player can equip ground items in a single turn action.
	for _, action := range scanData.Equip {
		itemRec, err := p.Domain.GetAdventureGameItemInstanceRec(action.ItemInstanceID, nil)
		if err != nil {
			l.Warn("failed to get item instance >%s< >%v<", action.ItemInstanceID, err)
			return fmt.Errorf("failed to get item instance %s: %w", action.ItemInstanceID, err)
		}
		if !itemRec.AdventureGameCharacterInstanceID.Valid &&
			itemRec.AdventureGameLocationInstanceID.Valid &&
			itemRec.AdventureGameLocationInstanceID.String == characterInstanceRec.AdventureGameLocationInstanceID.String {
			l.Info("auto picking up location item >%s< before equip", action.ItemInstanceID)
			_, err := p.Domain.PickUpAdventureGameItemInstanceRec(characterInstanceRec.ID, action.ItemInstanceID)
			if err != nil {
				l.Warn("failed to pick up item >%s< before equip >%v<", action.ItemInstanceID, err)
				return fmt.Errorf("failed to pick up item %s before equip: %w", action.ItemInstanceID, err)
			}
		}

		// Resolve the correct equipment slot from the item definition.
		// The HTML form sends every equip action with DefaultEquipSlot ("weapon"),
		// but each item has its own defined slot (armor, jewelry, etc.).
		slot := action.Slot
		itemDef, err := p.Domain.GetAdventureGameItemRec(itemRec.AdventureGameItemID, nil)
		if err != nil {
			l.Warn("failed to get item definition >%s< >%v<", itemRec.AdventureGameItemID, err)
			return fmt.Errorf("failed to get item definition for item %s: %w", action.ItemInstanceID, err)
		}
		if itemDef.EquipmentSlot != nil {
			slot = *itemDef.EquipmentSlot
		}

		l.Info("equipping item >%s< to slot >%s<", action.ItemInstanceID, slot)
		_, err = p.Domain.EquipAdventureGameItemInstanceRec(characterInstanceRec.ID, action.ItemInstanceID, slot)
		if err != nil {
			l.Warn("failed to equip item >%s< >%v<", action.ItemInstanceID, err)
			return fmt.Errorf("failed to equip item %s: %w", action.ItemInstanceID, err)
		}
		_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
			Category: turnsheet.TurnEventCategoryInventory,
			Icon:     turnsheet.TurnEventIconInventory,
			Message:  fmt.Sprintf("You equipped %s.", itemDef.Name),
		})
		charNeedsUpdate = true
	}

	// Process Use item actions.
	for _, itemInstanceID := range scanData.Use {
		l.Info("using item >%s<", itemInstanceID)

		itemInstance, err := p.Domain.GetAdventureGameItemInstanceRec(itemInstanceID, nil)
		if err != nil {
			l.Warn("failed to get item instance >%s< >%v<", itemInstanceID, err)
			return fmt.Errorf("failed to get item instance %s: %w", itemInstanceID, err)
		}

		itemDef, err := p.Domain.GetAdventureGameItemRec(itemInstance.AdventureGameItemID, nil)
		if err != nil {
			l.Warn("failed to get item definition >%s< >%v<", itemInstance.AdventureGameItemID, err)
			return fmt.Errorf("failed to get item definition %s: %w", itemInstanceID, err)
		}

		if !itemDef.CanBeUsed {
			l.Warn("item >%s< is not usable — skipping", itemInstanceID)
			continue
		}
		if itemInstance.UsesRemaining <= 0 {
			l.Warn("item >%s< has no uses remaining — skipping", itemInstanceID)
			continue
		}

		// Apply heal effect to character.
		if itemDef.HealAmount > 0 {
			characterInstanceRec.Health += itemDef.HealAmount
			if characterInstanceRec.Health > 100 {
				characterInstanceRec.Health = 100
			}
			l.Info("item >%s< healed character for >%d< (health now %d)", itemDef.Name, itemDef.HealAmount, characterInstanceRec.Health)
			_ = turnsheet.AppendTurnEvent(characterInstanceRec, turnsheet.TurnEvent{
				Category: turnsheet.TurnEventCategoryInventory,
				Icon:     turnsheet.TurnEventIconHeal,
				Message:  fmt.Sprintf("You used a %s and recovered %d health.", itemDef.Name, itemDef.HealAmount),
			})
		}

		// Decrement uses and mark as used when exhausted.
		itemInstance.UsesRemaining--
		if itemInstance.UsesRemaining <= 0 {
			itemInstance.IsUsed = true
		}
		if _, err := p.Domain.UpdateAdventureGameItemInstanceRec(itemInstance); err != nil {
			l.Warn("failed to update item instance after use >%v<", err)
		}
		charNeedsUpdate = true
	}

	// Persist character record if health or turn events changed.
	if charNeedsUpdate {
		_, err := p.Domain.UpdateAdventureGameCharacterInstanceRec(characterInstanceRec)
		if err != nil {
			return fmt.Errorf("failed to save character instance after inventory actions: %w", err)
		}
	}

	l.Info("successfully processed inventory management actions for character >%s<", characterInstanceRec.ID)

	return nil
}

// resolveItemName looks up the item definition name for a given item instance ID.
// Returns "an item" as a fallback if the lookup fails, so event generation is non-fatal.
func (p *AdventureGameInventoryManagementProcessor) resolveItemName(l logger.Logger, itemInstanceID string) string {
	itemRec, err := p.Domain.GetAdventureGameItemInstanceRec(itemInstanceID, nil)
	if err != nil {
		l.Warn("resolveItemName: failed to get item instance >%s< >%v<", itemInstanceID, err)
		return "an item"
	}
	itemDef, err := p.Domain.GetAdventureGameItemRec(itemRec.AdventureGameItemID, nil)
	if err != nil {
		l.Warn("resolveItemName: failed to get item definition >%s< >%v<", itemRec.AdventureGameItemID, err)
		return "an item"
	}
	return itemDef.Name
}

// CreateNextTurnSheet creates a new inventory management turn sheet for a character (implements TurnSheetProcessor interface)
func (p *AdventureGameInventoryManagementProcessor) CreateNextTurnSheet(ctx context.Context, gameInstanceRec *game_record.GameInstance, characterInstanceRec *adventure_game_record.AdventureGameCharacterInstance) (*game_record.GameTurnSheet, error) {
	l := p.Logger.WithFunctionContext("AdventureGameInventoryManagementProcessor/CreateNextTurnSheet")

	l.Info("creating inventory management turn sheet for character >%s<", characterInstanceRec.ID)

	// Step 1: Get character's current location instance
	locationInstanceRec, err := p.Domain.GetAdventureGameLocationInstanceRec(characterInstanceRec.AdventureGameLocationInstanceID.String, nil)
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
	accountUserRec, err := p.Domain.GetAccountUserRec(characterRec.AccountUserID, nil)
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

	// Step 7: Check for aggressive creatures at the current location.
	hasAggressiveCreatures, err := HasAggressiveCreaturesAtLocation(p.Domain, gameInstanceRec.ID, locationInstanceRec.ID)
	if err != nil {
		l.Warn("failed to check for creatures at location >%v<", err)
	}

	// Step 7a: Get items at current location (empty if aggressive creatures are present).
	var locationItems []*adventure_game_record.AdventureGameItemInstance
	if !hasAggressiveCreatures {
		locationItems, err = p.Domain.GetAdventureGameItemInstanceRecsByLocationInstance(locationInstanceRec.ID)
		if err != nil {
			l.Warn("failed to get location items >%v<", err)
			return nil, fmt.Errorf("failed to get location items: %w", err)
		}
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
			case slot == adventure_game_record.AdventureGameItemEquipmentSlotWeapon:
				displaySlot = adventure_game_record.AdventureGameItemEquipmentSlotWeapon
			case strings.HasPrefix(slot, "armor_"), slot == adventure_game_record.AdventureGameItemEquipmentSlotArmor:
				displaySlot = adventure_game_record.AdventureGameItemEquipmentSlotArmor
			case strings.HasPrefix(slot, "clothing_"), slot == adventure_game_record.AdventureGameItemEquipmentSlotClothing:
				displaySlot = adventure_game_record.AdventureGameItemEquipmentSlotClothing
			case strings.HasPrefix(slot, "jewelry_"), slot == adventure_game_record.AdventureGameItemEquipmentSlotJewelry:
				displaySlot = adventure_game_record.AdventureGameItemEquipmentSlotJewelry
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
			CanUse:          itemDef.CanBeUsed,
			UsesRemaining:   itemInstance.UsesRemaining,
		}
		inventoryItemList = append(inventoryItemList, inventoryItem)
	}

	// Sort inventory items: equipped items first (weapon, armor, jewelry order), then unequipped
	equippedItems := make([]turnsheet.InventoryItem, 0)
	unequippedItems := make([]turnsheet.InventoryItem, 0)

	// Define slot priority for sorting equipped items
	slotPriority := map[string]int{
		adventure_game_record.AdventureGameItemEquipmentSlotWeapon:   1,
		adventure_game_record.AdventureGameItemEquipmentSlotArmor:    2,
		adventure_game_record.AdventureGameItemEquipmentSlotClothing: 3,
		adventure_game_record.AdventureGameItemEquipmentSlotJewelry:  4,
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
	inventoryItemList = make([]turnsheet.InventoryItem, 0, len(equippedItems)+len(unequippedItems))
	inventoryItemList = append(inventoryItemList, equippedItems...)
	inventoryItemList = append(inventoryItemList, unequippedItems...)

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
			CanEquip:        itemDef.CanBeEquipped,
		}
		locationItemList = append(locationItemList, locationItem)
	}

	// Step 10: Generate turn sheet code for template rendering
	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		l.Warn("failed to generate turn sheet code >%v<", err)
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	// Step 10a: Load background image for inventory management (game-level image)
	var backgroundImage *string
	bgImageURL, err := p.Domain.GetAdventureGameInventoryTurnSheetImageDataURL(gameRec.ID)
	if err != nil {
		l.Warn("failed to get turn sheet background image >%v<", err)
	} else if bgImageURL != "" {
		backgroundImage = &bgImageURL
		l.Info("loaded background image for inventory management turn sheet, length >%d<", len(bgImageURL))
	} else {
		l.Info("no background image found for inventory management turn sheet")
	}

	// Step 11: Read inventory events for this sheet. Events are cleared after all processors run.
	displayEvents := ReadTurnEventsForCategories(l, p.Domain, characterInstanceRec, turnsheet.TurnEventCategoryInventory)

	// Step 12: Create sheet data
	sheetData := turnsheet.InventoryManagementData{
		TurnSheetTemplateData: turnsheet.TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr("adventure"),
			TurnNumber:            convert.Ptr(gameInstanceRec.CurrentTurn),
			AccountName:           convert.Ptr(accountUserRec.Email),
			TurnSheetTitle:        convert.Ptr("Inventory Management"),
			TurnSheetDescription:  convert.Ptr(fmt.Sprintf("Manage your inventory and equipment. Carrying %d/%d items.", len(inventoryItemList), characterInstanceRec.InventoryCapacity)),
			TurnSheetInstructions: convert.Ptr(turnsheet.DefaultInventoryManagementInstructions()),
			TurnSheetCode:         convert.Ptr(turnSheetCode),
			BackgroundImage:       backgroundImage,
			TurnEvents:            displayEvents,
		},
		CharacterName:          characterRec.Name,
		CurrentLocationName:    locationRec.Name,
		InventoryCapacity:      characterInstanceRec.InventoryCapacity,
		InventoryCount:         len(inventoryItemList),
		CurrentInventory:       inventoryItemList,
		EquipmentSlots:         equipmentSlots,
		LocationItems:          locationItemList,
		HasAggressiveCreatures: hasAggressiveCreatures,
	}

	sheetDataBytes, err := json.Marshal(sheetData)
	if err != nil {
		l.Warn("failed to marshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to marshal sheet data: %w", err)
	}

	// Step 13: Create turn sheet record
	turnSheet := &game_record.GameTurnSheet{
		GameID:           gameInstanceRec.GameID,
		AccountID:        accountUserRec.AccountID,
		AccountUserID:    characterRec.AccountUserID,
		TurnNumber:       gameInstanceRec.CurrentTurn,
		SheetType:        adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement,
		SheetOrder:       adventure_game_record.AdventureGameSheetOrderForType(adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement),
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
