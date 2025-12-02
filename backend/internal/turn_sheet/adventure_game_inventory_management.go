package turn_sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/scanner"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// InventoryManagementData represents the data structure for inventory management turn sheets
type InventoryManagementData struct {
	TurnSheetTemplateData

	// Character information
	CharacterName        string `json:"character_name"`
	CurrentLocationName  string `json:"current_location_name"`
	Health               int    `json:"health"`
	InventoryCapacity    int    `json:"inventory_capacity"`
	InventoryCount       int    `json:"inventory_count"`

	// Current inventory items
	CurrentInventory []InventoryItem `json:"current_inventory"`

	// Equipment slots (simplified to 4 slots for turn sheet display)
	EquipmentSlots EquipmentSlots `json:"equipment_slots"`

	// Items available at current location
	LocationItems []LocationItem `json:"location_items"`
}

// InventoryItem represents an item in the character's inventory
type InventoryItem struct {
	ItemInstanceID string `json:"item_instance_id"`
	ItemName       string `json:"item_name"`
	ItemDescription string `json:"item_description,omitempty"`
	IsEquipped     bool   `json:"is_equipped"`
	EquipmentSlot  string `json:"equipment_slot,omitempty"`
	CanEquip       bool   `json:"can_equip"`
}

// EquipmentSlots represents the equipped items in simplified slots
type EquipmentSlots struct {
	Weapon   *EquippedItem `json:"weapon,omitempty"`
	Armor    *EquippedItem `json:"armor,omitempty"`
	Clothing *EquippedItem `json:"clothing,omitempty"`
	Jewelry  *EquippedItem `json:"jewelry,omitempty"`
}

// EquippedItem represents an item equipped in a slot
type EquippedItem struct {
	ItemInstanceID string `json:"item_instance_id"`
	ItemName       string `json:"item_name"`
	SlotName       string `json:"slot_name"`
}

// LocationItem represents an item available at the current location
type LocationItem struct {
	ItemInstanceID string `json:"item_instance_id"`
	ItemName       string `json:"item_name"`
	ItemDescription string `json:"item_description,omitempty"`
}

// InventoryManagementScanData represents the scanned data from an inventory management turn sheet
type InventoryManagementScanData struct {
	PickUp  []string              `json:"pick_up,omitempty"`
	Drop    []string              `json:"drop,omitempty"`
	Equip   []EquipAction         `json:"equip,omitempty"`
	Unequip []string              `json:"unequip,omitempty"`
}

// EquipAction represents an equip action with slot
type EquipAction struct {
	ItemInstanceID string `json:"item_instance_id"`
	Slot           string `json:"slot"`
}

const defaultInventoryManagementInstructions = "Manage your inventory by checking boxes to pick up items, drop items, equip items, or unequip items. Return this form by the deadline."

const inventoryManagementTemplatePath = "turn_sheet/adventure_game_inventory_management.template"

// DefaultInventoryManagementInstructions returns the default instruction text for inventory management turn sheets.
func DefaultInventoryManagementInstructions() string {
	return defaultInventoryManagementInstructions
}

// InventoryManagementProcessor implements the DocumentProcessor interface for inventory management turn sheets
type InventoryManagementProcessor struct {
	*BaseProcessor
}

// NewInventoryManagementProcessor creates a new inventory management processor
func NewInventoryManagementProcessor(l logger.Logger, cfg config.Config) (*InventoryManagementProcessor, error) {
	baseProcessor, err := NewBaseProcessor(l, cfg)
	if err != nil {
		return nil, err
	}
	return &InventoryManagementProcessor{
		BaseProcessor: baseProcessor,
	}, nil
}

// GenerateTurnSheet generates an inventory management turn sheet document
func (p *InventoryManagementProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, sheetData []byte) ([]byte, error) {
	l = l.WithFunctionContext("InventoryManagementProcessor/GenerateTurnSheet")

	l.Info("generating inventory management turn sheet")

	// Unmarshal sheet data
	var inventoryData InventoryManagementData
	if err := json.Unmarshal(sheetData, &inventoryData); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to parse sheet data: %w", err)
	}

	// Validate base template data
	if err := p.ValidateBaseTemplateData(&inventoryData.TurnSheetTemplateData); err != nil {
		l.Warn("failed to validate base template data >%v<", err)
		return nil, fmt.Errorf("template data validation failed: %w", err)
	}

	// Set default instructions if not provided
	if inventoryData.TurnSheetInstructions == nil || strings.TrimSpace(*inventoryData.TurnSheetInstructions) == "" {
		instruction := defaultInventoryManagementInstructions
		inventoryData.TurnSheetInstructions = &instruction
	}

	// Set default title if not provided
	if inventoryData.TurnSheetTitle == nil || strings.TrimSpace(*inventoryData.TurnSheetTitle) == "" {
		title := "Inventory Management"
		inventoryData.TurnSheetTitle = &title
	}

	// Set default description if not provided
	if inventoryData.TurnSheetDescription == nil || strings.TrimSpace(*inventoryData.TurnSheetDescription) == "" {
		desc := fmt.Sprintf("Manage your inventory and equipment. Carrying %d/%d items.", inventoryData.InventoryCount, inventoryData.InventoryCapacity)
		inventoryData.TurnSheetDescription = &desc
	}

	// Validate inventory-specific data
	if inventoryData.CharacterName == "" {
		l.Warn("character name is missing")
		return nil, fmt.Errorf("character name is required")
	}

	// Generate document using the inventory management template
	return p.GenerateDocument(ctx, format, inventoryManagementTemplatePath, &inventoryData)
}

// ScanTurnSheet scans an inventory management turn sheet and extracts player actions using hosted OCR
func (p *InventoryManagementProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	l = l.WithFunctionContext("InventoryManagementProcessor/ScanTurnSheet")

	l.Info("scanning inventory management turn sheet")

	if len(imageData) == 0 {
		l.Warn("empty image data provided")
		return nil, fmt.Errorf("empty image data provided")
	}

	var inventoryData InventoryManagementData
	if err := json.Unmarshal(sheetData, &inventoryData); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to parse sheet data: %w", err)
	}

	templateImage, err := p.renderTemplatePreview(ctx, inventoryManagementTemplatePath, &inventoryData)
	if err != nil {
		l.Warn("failed to generate template preview >%v<", err)
		return nil, fmt.Errorf("failed to generate template preview: %w", err)
	}
	if len(templateImage) == 0 {
		l.Warn("template preview generation returned empty image")
		return nil, fmt.Errorf("template preview generation returned empty image")
	}

	expected := map[string]any{
		"pick_up":  []string{},
		"drop":     []string{},
		"equip":    []map[string]string{},
		"unequip":  []string{},
	}

	req := scanner.StructuredScanRequest{
		Instructions:       buildInventoryManagementInstructions(),
		AdditionalContext:  buildInventoryManagementContext(&inventoryData),
		TemplateImage:      templateImage,
		TemplateImageMIME:  "image/png",
		FilledImage:        imageData,
		ExpectedJSONSchema: expected,
	}

	raw, err := p.Scanner.ExtractStructuredData(ctx, req)
	if err != nil {
		l.Warn("structured extraction failed >%v<", err)
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	var scanData InventoryManagementScanData
	if err := json.Unmarshal(raw, &scanData); err != nil {
		return nil, fmt.Errorf("failed to decode structured inventory actions: %w", err)
	}

	if err := validateInventoryActions(&inventoryData, &scanData); err != nil {
		return nil, err
	}

	return json.Marshal(scanData)
}

// buildInventoryManagementInstructions returns instructions for the AI-driven OCR service
func buildInventoryManagementInstructions() string {
	return `Compare the blank template image with the completed turn sheet.
Determine which checkboxes are marked by the player for:
- Pick Up: Items at location to pick up
- Drop: Items in inventory to drop
- Equip: Items in inventory to equip (with slot selection)
- Unequip: Equipped items to unequip
Respond with JSON containing arrays of item_instance_id values for each action.
For equip actions, include both item_instance_id and slot.`
}

// buildInventoryManagementContext returns additional context for the AI-driven OCR service
func buildInventoryManagementContext(data *InventoryManagementData) []string {
	var ctx []string
	if data != nil {
		ctx = append(ctx, fmt.Sprintf("Character: %s", data.CharacterName))
		ctx = append(ctx, fmt.Sprintf("Location: %s", data.CurrentLocationName))
		ctx = append(ctx, fmt.Sprintf("Inventory: %d/%d items", data.InventoryCount, data.InventoryCapacity))
		
		if len(data.CurrentInventory) > 0 {
			ctx = append(ctx, "Current Inventory:")
			for _, item := range data.CurrentInventory {
				ctx = append(ctx, fmt.Sprintf("  - %s (ID: %s, Equipped: %v)", item.ItemName, item.ItemInstanceID, item.IsEquipped))
			}
		}
		
		if len(data.LocationItems) > 0 {
			ctx = append(ctx, "Items at Location:")
			for _, item := range data.LocationItems {
				ctx = append(ctx, fmt.Sprintf("  - %s (ID: %s)", item.ItemName, item.ItemInstanceID))
			}
		}
	}
	return ctx
}

// validateInventoryActions validates the scanned inventory actions
func validateInventoryActions(sheetData *InventoryManagementData, scanData *InventoryManagementScanData) error {
	if scanData == nil {
		return fmt.Errorf("no scan data provided")
	}

	// Build maps of valid item instance IDs
	inventoryItemIDs := make(map[string]bool)
	for _, item := range sheetData.CurrentInventory {
		inventoryItemIDs[item.ItemInstanceID] = true
	}

	locationItemIDs := make(map[string]bool)
	for _, item := range sheetData.LocationItems {
		locationItemIDs[item.ItemInstanceID] = true
	}

	// Validate pick up actions
	for _, itemID := range scanData.PickUp {
		if !locationItemIDs[itemID] {
			return fmt.Errorf("invalid item_instance_id for pick up: %s", itemID)
		}
	}

	// Validate drop actions
	for _, itemID := range scanData.Drop {
		if !inventoryItemIDs[itemID] {
			return fmt.Errorf("invalid item_instance_id for drop: %s", itemID)
		}
	}

	// Validate equip actions
	for _, action := range scanData.Equip {
		if !inventoryItemIDs[action.ItemInstanceID] {
			return fmt.Errorf("invalid item_instance_id for equip: %s", action.ItemInstanceID)
		}
		// Validate slot
		validSlots := []string{"weapon", "armor", "clothing", "jewelry"}
		valid := false
		for _, slot := range validSlots {
			if action.Slot == slot {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid equipment slot: %s", action.Slot)
		}
	}

	// Validate unequip actions
	for _, itemID := range scanData.Unequip {
		if !inventoryItemIDs[itemID] {
			return fmt.Errorf("invalid item_instance_id for unequip: %s", itemID)
		}
	}

	return nil
}

