package turnsheet

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/scanner"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

// OrdersScannedDataSchemaName is the filename of the JSON schema for orders scanned_data (under schema/turnsheet/mecha/).
const OrdersScannedDataSchemaName = "orders.schema.json"

const defaultOrdersInstructions = "For each mech in your squad, choose a sector to move to and/or a target to attack, then return this form."
const ordersTemplatePath = "turnsheet/mecha_game_orders.template"

// DefaultOrdersInstructions returns the default instruction text for orders turn sheets.
func DefaultOrdersInstructions() string {
	return defaultOrdersInstructions
}

// MechWeaponEntry is a single weapon fitted to a mech.
type MechWeaponEntry struct {
	WeaponID     string `json:"weapon_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Damage       int    `json:"damage"`
	HeatCost     int    `json:"heat_cost"`
	RangeBand    string `json:"range_band,omitempty"`
	SlotLocation string `json:"slot_location,omitempty"`
	// AmmoCapacity is the weapon's designer-configured ammo cost per
	// trigger-pull pool allocation. Zero means the weapon does not
	// consume ammo and the mech's shared ammo pool is not decremented
	// when this weapon fires.
	AmmoCapacity int `json:"ammo_capacity"`
}

// MechEquipmentEntry is a single equipment item fitted to a mech, shown on
// the orders / management sheets so the player can see which mounts are
// accounted for by equipment vs weapons.
type MechEquipmentEntry struct {
	EquipmentID  string `json:"equipment_id,omitempty"`
	Name         string `json:"name,omitempty"`
	EffectKind   string `json:"effect_kind,omitempty"`
	Magnitude    int    `json:"magnitude"`
	HeatCost     int    `json:"heat_cost"`
	MountSize    string `json:"mount_size,omitempty"`
	SlotLocation string `json:"slot_location,omitempty"`
}

// MechOrderEntry is a single mech's orders for a turn.
type MechOrderEntry struct {
	MechInstanceID             string            `json:"mech_instance_id"`
	MechCallsign               string            `json:"mech_callsign,omitempty"`
	MechStatus                 string            `json:"mech_status,omitempty"`
	CurrentSectorName          string            `json:"current_sector_name,omitempty"`
	MoveToSectorInstanceID     string            `json:"move_to_sector_instance_id,omitempty"`
	AttackTargetMechInstanceID string            `json:"attack_target_mech_instance_id,omitempty"`
	ChassisName                string            `json:"chassis_name,omitempty"`
	ChassisClass               string            `json:"chassis_class,omitempty"`
	CurrentArmor               int               `json:"current_armor"`
	MaxArmor                   int               `json:"max_armor"`
	CurrentStructure           int               `json:"current_structure"`
	MaxStructure               int               `json:"max_structure"`
	CurrentHeat                int               `json:"current_heat"`
	HeatCapacity               int               `json:"heat_capacity"`
	Speed                      int               `json:"speed"`
	PilotSkill                 int               `json:"pilot_skill"`
	IsRefitting                bool              `json:"is_refitting"`
	Weapons                    []MechWeaponEntry `json:"weapons,omitempty"`
	// Equipment lists non-weapon items mounted on this mech's chassis,
	// rendered alongside Weapons so players can see combined slot usage.
	Equipment []MechEquipmentEntry `json:"equipment,omitempty"`
	// AmmoRemaining is the mech's current shared ammo pool, shown only
	// for mechs carrying at least one weapon with ammo_capacity > 0.
	AmmoRemaining int `json:"ammo_remaining,omitempty"`
	// AmmoCapacity is the mech's maximum ammo pool derived from weapon
	// ammo_capacity totals + ammo_bin magnitudes, used alongside
	// AmmoRemaining to render "x/y" ammo indicators on the turn sheet.
	AmmoCapacity int `json:"ammo_capacity,omitempty"`
	// ReachableSectors lists sectors this mech can reach within its speed budget.
	ReachableSectors []SectorOption `json:"reachable_sectors,omitempty"`
}

// SectorOption represents a sector available for movement.
type SectorOption struct {
	SectorInstanceID string `json:"sector_instance_id"`
	SectorName       string `json:"sector_name"`
}

// EnemyMechOption represents an enemy mech that can be attacked.
type EnemyMechOption struct {
	MechInstanceID string `json:"mech_instance_id"`
	Callsign       string `json:"callsign"`
	SectorName     string `json:"sector_name"`
}

// OrdersData is the data model for a mecha orders turn sheet.
type OrdersData struct {
	TurnSheetTemplateData

	// Squad information
	SquadName string `json:"squad_name,omitempty"`

	// Mechs in this squad with their current state and available orders
	SquadMechs []MechOrderEntry `json:"squad_mechs,omitempty"`

	// Available adjacent sectors for movement
	AvailableSectors []SectorOption `json:"available_sectors,omitempty"`

	// Visible enemy mechs that can be targeted
	EnemyMechs []EnemyMechOption `json:"enemy_mechs,omitempty"`
}

// OrdersScanData represents scanned orders data submitted by the player.
type OrdersScanData struct {
	MechOrders []ScannedMechOrder `json:"mech_orders,omitempty"`
}

// ScannedMechOrder is a single mech's orders extracted from the scanned turn sheet.
type ScannedMechOrder struct {
	MechInstanceID             string `json:"mech_instance_id"`
	MoveToSectorInstanceID     string `json:"move_to_sector_instance_id,omitempty"`
	AttackTargetMechInstanceID string `json:"attack_target_mech_instance_id,omitempty"`
}

// MechaGameOrdersProcessor implements the DocumentProcessor interface for mecha orders sheets.
type MechaGameOrdersProcessor struct {
	*BaseProcessor
}

// NewMechaGameOrdersProcessor creates a new mecha orders processor.
func NewMechaGameOrdersProcessor(l logger.Logger, cfg config.Config) (*MechaGameOrdersProcessor, error) {
	baseProcessor, err := NewBaseProcessor(l, cfg)
	if err != nil {
		return nil, err
	}
	return &MechaGameOrdersProcessor{
		BaseProcessor: baseProcessor,
	}, nil
}

// GeneratePreviewData generates dummy data for a mecha orders turn sheet preview.
func (p *MechaGameOrdersProcessor) GeneratePreviewData(ctx context.Context, l logger.Logger, gameRec *game_record.Game, backgroundImage *string) ([]byte, error) {
	l = l.WithFunctionContext("MechaGameOrdersProcessor/GeneratePreviewData")

	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode("preview-turn-sheet-id")
	if err != nil {
		l.Warn("failed to generate turn sheet code >%v<", err)
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	turnNumber := 1
	title := "Mech Orders"
	instructions := defaultOrdersInstructions
	data := OrdersData{
		TurnSheetTemplateData: TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr(gameRec.GameType),
			TurnNumber:            &turnNumber,
			TurnSheetTitle:        &title,
			TurnSheetDescription:  convert.Ptr(gameRec.Description),
			TurnSheetInstructions: &instructions,
			TurnSheetCode:         convert.Ptr(turnSheetCode),
			TurnEvents: []TurnEvent{
				{
					Category: TurnEventCategoryMovement,
					Icon:     TurnEventIconMovement,
					Message:  "Hammer advanced to Central Wastes.",
				},
				{
					Category: TurnEventCategoryCombat,
					Icon:     TurnEventIconCombat,
					Message:  "Anvil fired Light Pulse Cannon at Stalker — hit for 3 damage.",
				},
			},
		},
		SquadName: "Preview Squad",
		SquadMechs: []MechOrderEntry{
			{
				MechInstanceID:    "preview-mech-1",
				MechCallsign:      "Hammer",
				MechStatus:        "operational",
				CurrentSectorName: "Central Wastes",
				ChassisName:       "Scout",
				ChassisClass:      "light",
				CurrentArmor:      55,
				MaxArmor:          72,
				CurrentStructure:  32,
				MaxStructure:      32,
				CurrentHeat:       4,
				HeatCapacity:      18,
				Speed:             7,
				PilotSkill:        4,
				Weapons: []MechWeaponEntry{
					{Name: "Light Pulse Cannon", Damage: 3, HeatCost: 1, RangeBand: "short", SlotLocation: "left-arm"},
					{Name: "Chaingun", Damage: 2, HeatCost: 0, RangeBand: "short", SlotLocation: "right-arm"},
				},
			},
			{
				MechInstanceID:    "preview-mech-2",
				MechCallsign:      "Anvil",
				MechStatus:        "operational",
				CurrentSectorName: "Central Wastes",
				ChassisName:       "Sentinel",
				ChassisClass:      "medium",
				CurrentArmor:      130,
				MaxArmor:          130,
				CurrentStructure:  65,
				MaxStructure:      65,
				CurrentHeat:       0,
				HeatCapacity:      28,
				Speed:             4,
				PilotSkill:        4,
				Weapons: []MechWeaponEntry{
					{Name: "Pulse Cannon", Damage: 5, HeatCost: 3, RangeBand: "medium", SlotLocation: "left-torso"},
					{Name: "Rocket Pack", Damage: 8, HeatCost: 3, RangeBand: "short", SlotLocation: "right-arm"},
				},
			},
		},
		AvailableSectors: []SectorOption{
			{SectorInstanceID: "preview-sector-1", SectorName: "Northern Ridge"},
			{SectorInstanceID: "preview-sector-2", SectorName: "Southern Flats"},
		},
		EnemyMechs: []EnemyMechOption{
			{MechInstanceID: "enemy-mech-1", Callsign: "Stalker", SectorName: "Northern Ridge"},
		},
	}

	if backgroundImage != nil {
		data.BackgroundImage = backgroundImage
	}

	return json.Marshal(data)
}

// GenerateTurnSheet generates a mecha orders turn sheet document.
func (p *MechaGameOrdersProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, sheetData []byte) ([]byte, error) {
	l = l.WithFunctionContext("MechaGameOrdersProcessor/GenerateTurnSheet")

	var data OrdersData
	if err := json.Unmarshal(sheetData, &data); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to parse sheet data: %w", err)
	}

	if err := p.ValidateBaseTemplateData(&data.TurnSheetTemplateData); err != nil {
		l.Warn("failed to validate base template data >%v<", err)
		return nil, fmt.Errorf("template data validation failed: %w", err)
	}

	if data.TurnSheetInstructions == nil || strings.TrimSpace(*data.TurnSheetInstructions) == "" {
		instructions := defaultOrdersInstructions
		data.TurnSheetInstructions = &instructions
	}

	if data.TurnSheetTitle == nil || strings.TrimSpace(*data.TurnSheetTitle) == "" {
		title := "Mech Orders"
		data.TurnSheetTitle = &title
	}

	return p.GenerateDocument(ctx, format, ordersTemplatePath, &data)
}

// ScanTurnSheet extracts mech orders from the uploaded document using the hosted scanner.
func (p *MechaGameOrdersProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	l = l.WithFunctionContext("MechaGameOrdersProcessor/ScanTurnSheet")

	if len(imageData) == 0 {
		return nil, fmt.Errorf("empty image data provided")
	}

	templateData := p.resolveOrdersTemplateData(sheetData)

	templateImage, err := p.renderTemplatePreview(ctx, ordersTemplatePath, templateData)
	if err != nil {
		l.Warn("failed to generate template preview >%v<", err)
		return nil, fmt.Errorf("failed to generate template preview: %w", err)
	}

	if len(templateImage) == 0 {
		return nil, fmt.Errorf("template preview generation returned empty image")
	}

	req := scanner.StructuredScanRequest{
		Instructions:       buildOrdersInstructions(),
		AdditionalContext:  buildOrdersContext(templateData),
		TemplateImage:      templateImage,
		TemplateImageMIME:  "image/png",
		FilledImage:        imageData,
		ExpectedJSONSchema: ordersExpectedSchema(),
	}

	raw, err := p.Scanner.ExtractStructuredData(ctx, req)
	if err != nil {
		l.Warn("structured extraction failed >%v<", err)
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	var scanData OrdersScanData
	if err := json.Unmarshal(raw, &scanData); err != nil {
		l.Warn("failed to decode structured response >%v<", err)
		return nil, fmt.Errorf("failed to decode structured response: %w", err)
	}

	normalizeOrdersScanData(&scanData)

	return json.Marshal(scanData)
}

func (p *MechaGameOrdersProcessor) resolveOrdersTemplateData(sheetData []byte) *OrdersData {
	var data OrdersData
	if len(sheetData) > 0 {
		if err := json.Unmarshal(sheetData, &data); err != nil {
			return defaultOrdersTemplateData()
		}
		return &data
	}
	return defaultOrdersTemplateData()
}

func defaultOrdersTemplateData() *OrdersData {
	title := "Mech Orders"
	instructions := defaultOrdersInstructions
	return &OrdersData{
		TurnSheetTemplateData: TurnSheetTemplateData{
			TurnSheetTitle:        &title,
			TurnSheetInstructions: &instructions,
		},
	}
}

func ordersExpectedSchema() map[string]any {
	return map[string]any{
		"mech_orders": []map[string]any{
			{
				"mech_instance_id":               "",
				"move_to_sector_instance_id":     "",
				"attack_target_mech_instance_id": "",
			},
		},
	}
}

func buildOrdersInstructions() string {
	return `You are comparing two images of a PlayByMail "Mech Orders" form.
- Image 1 is the blank reference form showing all available mechs, sectors, and enemy targets.
- Image 2 is the completed form containing the player's handwritten orders.
Extract the player's orders for each mech and return them as JSON with the key "mech_orders".

Each entry in "mech_orders" must contain:
- "mech_instance_id": the mech ID exactly as printed on the form
- "move_to_sector_instance_id": the sector instance ID the player chose to move to, or empty string if staying in place
- "attack_target_mech_instance_id": the enemy mech instance ID the player chose to attack, or empty string if no attack

Return empty strings for fields the player left blank. Do not invent or guess values.`
}

func buildOrdersContext(data *OrdersData) []string {
	var ctx []string
	if data != nil {
		if data.GameName != nil {
			ctx = append(ctx, fmt.Sprintf("Game name: %s", strings.TrimSpace(*data.GameName)))
		}
		if data.SquadName != "" {
			ctx = append(ctx, fmt.Sprintf("Squad name: %s", data.SquadName))
		}
		if data.TurnNumber != nil {
			ctx = append(ctx, fmt.Sprintf("Turn number: %d", *data.TurnNumber))
		}
	}
	ctx = append(ctx,
		"The JSON must only contain the requested keys.",
		"Return empty strings when the player left a field blank.",
		"Copy mech_instance_id, sector instance IDs, and enemy mech IDs exactly as printed — do not modify them.",
	)
	return ctx
}

func normalizeOrdersScanData(data *OrdersScanData) {
	for i := range data.MechOrders {
		data.MechOrders[i].MechInstanceID = strings.TrimSpace(data.MechOrders[i].MechInstanceID)
		data.MechOrders[i].MoveToSectorInstanceID = strings.TrimSpace(data.MechOrders[i].MoveToSectorInstanceID)
		data.MechOrders[i].AttackTargetMechInstanceID = strings.TrimSpace(data.MechOrders[i].AttackTargetMechInstanceID)
	}
}
