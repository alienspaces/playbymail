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

// Ensure MechaLanceManagementProcessor implements DocumentProcessor at compile time.
var _ DocumentProcessor = (*MechaLanceManagementProcessor)(nil)

const LanceManagementScannedDataSchemaName = "lance_management.schema.json"

const defaultManagementInstructions = "Manage your lance: repair structure, swap weapons, or schedule refits. Actions consume supply points and take effect next turn."
const managementTemplatePath = "turnsheet/mecha_lance_management.template"

// DefaultManagementInstructions returns the default instruction text for management sheets.
func DefaultManagementInstructions() string {
	return defaultManagementInstructions
}

// LanceManagementData is the template data for the lance management sheet.
type LanceManagementData struct {
	TurnSheetTemplateData
	LanceName     string                 `json:"lance_name"`
	SupplyPoints  int                    `json:"supply_points"`
	Mechs         []ManagementMechEntry  `json:"mechs"`
	WeaponCatalog []CatalogWeapon        `json:"weapon_catalog"`
}

// ManagementMechEntry holds per-mech data for the management sheet.
type ManagementMechEntry struct {
	MechInstanceID   string           `json:"mech_instance_id"`
	Callsign         string           `json:"callsign"`
	ChassisName      string           `json:"chassis_name,omitempty"`
	ChassisClass     string           `json:"chassis_class,omitempty"`
	Status           string           `json:"status,omitempty"`
	IsAtDepot        bool             `json:"is_at_depot"`
	IsRefitting      bool             `json:"is_refitting"`
	CurrentArmor     int              `json:"current_armor"`
	MaxArmor         int              `json:"max_armor"`
	CurrentStructure int              `json:"current_structure"`
	MaxStructure     int              `json:"max_structure"`
	StructureDamage  int              `json:"structure_damage"`
	Weapons          []MechWeaponSlot `json:"weapons,omitempty"`
}

// MechWeaponSlot describes one weapon slot on a mech.
type MechWeaponSlot struct {
	SlotLocation      string `json:"slot_location"`
	MountSize         string `json:"mount_size,omitempty"`
	CurrentWeaponID   string `json:"current_weapon_id,omitempty"`
	CurrentWeaponName string `json:"current_weapon_name,omitempty"`
}

// CatalogWeapon is a weapon available in the game's weapon catalog.
type CatalogWeapon struct {
	WeaponID  string `json:"weapon_id"`
	Name      string `json:"name"`
	Damage    int    `json:"damage"`
	HeatCost  int    `json:"heat_cost"`
	RangeBand string `json:"range_band,omitempty"`
	MountSize string `json:"mount_size,omitempty"`
}

// LanceManagementScanData represents scanned management orders submitted by a player.
type LanceManagementScanData struct {
	MechManagementOrders []ScannedMechManagementOrder `json:"mech_management_orders,omitempty"`
}

// Validate ensures the scanned management data is minimally valid.
func (d *LanceManagementScanData) Validate() error {
	for i, order := range d.MechManagementOrders {
		if strings.TrimSpace(order.MechInstanceID) == "" {
			return fmt.Errorf("mech_management_orders[%d]: mech_instance_id is required", i)
		}
		for j, swap := range order.WeaponSwaps {
			if strings.TrimSpace(swap.SlotLocation) == "" {
				return fmt.Errorf("mech_management_orders[%d].weapon_swaps[%d]: slot_location is required", i, j)
			}
		}
	}
	return nil
}

// ScannedMechManagementOrder is one mech's management orders from a scanned sheet.
type ScannedMechManagementOrder struct {
	MechInstanceID  string               `json:"mech_instance_id"`
	RepairStructure bool                 `json:"repair_structure,omitempty"`
	WeaponSwaps     []ScannedWeaponSwap  `json:"weapon_swaps,omitempty"`
}

// ScannedWeaponSwap replaces a weapon in one slot.
type ScannedWeaponSwap struct {
	SlotLocation string `json:"slot_location"`
	NewWeaponID  string `json:"new_weapon_id"`
}

// MechaLanceManagementProcessor implements DocumentProcessor for management sheets.
type MechaLanceManagementProcessor struct {
	*BaseProcessor
}

// NewMechaLanceManagementProcessor creates a new management sheet processor.
func NewMechaLanceManagementProcessor(l logger.Logger, cfg config.Config) (*MechaLanceManagementProcessor, error) {
	baseProcessor, err := NewBaseProcessor(l, cfg)
	if err != nil {
		return nil, err
	}
	return &MechaLanceManagementProcessor{
		BaseProcessor: baseProcessor,
	}, nil
}

// GenerateTurnSheet generates a lance management turn sheet document.
func (p *MechaLanceManagementProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, sheetData []byte) ([]byte, error) {
	l = l.WithFunctionContext("MechaLanceManagementProcessor/GenerateTurnSheet")
	l.Debug("Generating lance management turn sheet")

	var data LanceManagementData
	if err := json.Unmarshal(sheetData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lance management data: %w", err)
	}

	if err := p.ValidateBaseTemplateData(&data.TurnSheetTemplateData); err != nil {
		return nil, err
	}

	return p.GenerateDocument(ctx, format, managementTemplatePath, &data)
}

// ScanTurnSheet OCR-scans a filled management turn sheet and returns
// structured JSON matching LanceManagementScanData.
func (p *MechaLanceManagementProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	l = l.WithFunctionContext("MechaLanceManagementProcessor/ScanTurnSheet")
	l.Debug("Scanning lance management turn sheet")

	if len(imageData) == 0 {
		return nil, fmt.Errorf("no image data provided for scanning")
	}

	var data LanceManagementData
	if err := json.Unmarshal(sheetData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lance management data: %w", err)
	}

	templateImage, err := p.renderTemplatePreview(ctx, managementTemplatePath, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate template preview: %w", err)
	}

	req := scanner.StructuredScanRequest{
		Instructions: `Extract the lance management orders from this turn sheet.
For each mech, identify:
- mech_instance_id: the hidden ID field value
- repair_structure: true if the repair checkbox is checked
- weapon_swaps: list of {slot_location, new_weapon_id} for any changed weapon slots`,
		TemplateImage:     templateImage,
		TemplateImageMIME: "image/png",
		FilledImage:       imageData,
		ExpectedJSONSchema: map[string]any{
			"mech_management_orders": []map[string]any{
				{
					"mech_instance_id":  "<id>",
					"repair_structure":  false,
					"weapon_swaps":      []map[string]any{{"slot_location": "<slot>", "new_weapon_id": "<id>"}},
				},
			},
		},
	}

	raw, err := p.Scanner.ExtractStructuredData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	var scanData LanceManagementScanData
	if err := json.Unmarshal(raw, &scanData); err != nil {
		l.Warn("failed to decode structured response >%v<", err)
		return nil, fmt.Errorf("failed to decode structured response: %w", err)
	}

	normalizeLanceManagementScanData(&scanData)

	if err := scanData.Validate(); err != nil {
		l.Warn("failed to validate scan data >%v<", err)
		return nil, err
	}

	return json.Marshal(scanData)
}

// GeneratePreviewData returns sample management sheet data for previewing.
func (p *MechaLanceManagementProcessor) GeneratePreviewData(ctx context.Context, l logger.Logger, gameRec *game_record.Game, backgroundImage *string) ([]byte, error) {
	l = l.WithFunctionContext("MechaLanceManagementProcessor/GeneratePreviewData")
	l.Debug("Generating lance management preview data")

	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode("preview-management-sheet-id")
	if err != nil {
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	turnNumber := 1
	title := "Lance Management"
	instructions := defaultManagementInstructions
	data := LanceManagementData{
		TurnSheetTemplateData: TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr(gameRec.GameType),
			TurnNumber:            &turnNumber,
			TurnSheetTitle:        &title,
			TurnSheetDescription:  convert.Ptr(gameRec.Description),
			TurnSheetInstructions: &instructions,
			TurnSheetCode:         convert.Ptr(turnSheetCode),
			BackgroundImage:       backgroundImage,
			TurnEvents: []TurnEvent{
				{
					Category: TurnEventCategorySystem,
					Icon:     TurnEventIconSystem,
					Message:  "Hammer field repairs restored 18 armor (72/72).",
				},
			},
		},
		LanceName:    "Preview Lance",
		SupplyPoints: 4,
		Mechs: []ManagementMechEntry{
			{
				MechInstanceID:   "preview-mech-1",
				Callsign:         "Hammer",
				ChassisName:      "Scout",
				ChassisClass:     "light",
				Status:           "operational",
				IsAtDepot:        true,
				CurrentArmor:     72,
				MaxArmor:         72,
				CurrentStructure: 28,
				MaxStructure:     32,
				StructureDamage:  4,
				Weapons: []MechWeaponSlot{
					{SlotLocation: "left-arm", CurrentWeaponID: "prev-wpn-1", CurrentWeaponName: "Light Pulse Cannon"},
					{SlotLocation: "right-arm", CurrentWeaponID: "prev-wpn-2", CurrentWeaponName: "Chaingun"},
				},
			},
			{
				MechInstanceID:   "preview-mech-2",
				Callsign:         "Anvil",
				ChassisName:      "Sentinel",
				ChassisClass:     "medium",
				Status:           "operational",
				IsAtDepot:        false,
				CurrentArmor:     100,
				MaxArmor:         130,
				CurrentStructure: 65,
				MaxStructure:     65,
				StructureDamage:  0,
				Weapons: []MechWeaponSlot{
					{SlotLocation: "left-torso", CurrentWeaponID: "prev-wpn-3", CurrentWeaponName: "Pulse Cannon"},
					{SlotLocation: "right-arm", CurrentWeaponID: "prev-wpn-4", CurrentWeaponName: "Rocket Pack"},
				},
			},
		},
		WeaponCatalog: []CatalogWeapon{
			{WeaponID: "cat-1", Name: "Light Pulse Cannon", Damage: 3, HeatCost: 1, RangeBand: "short"},
			{WeaponID: "cat-2", Name: "Chaingun", Damage: 2, HeatCost: 0, RangeBand: "short"},
			{WeaponID: "cat-3", Name: "Pulse Cannon", Damage: 5, HeatCost: 3, RangeBand: "medium"},
			{WeaponID: "cat-4", Name: "Rocket Pack", Damage: 8, HeatCost: 3, RangeBand: "short"},
		},
	}

	out, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal preview data: %w", err)
	}

	return out, nil
}

// normalizeLanceManagementScanData trims whitespace from scanned management order fields.
func normalizeLanceManagementScanData(data *LanceManagementScanData) {
	for i := range data.MechManagementOrders {
		data.MechManagementOrders[i].MechInstanceID = strings.TrimSpace(data.MechManagementOrders[i].MechInstanceID)
		for j := range data.MechManagementOrders[i].WeaponSwaps {
			data.MechManagementOrders[i].WeaponSwaps[j].SlotLocation = strings.TrimSpace(data.MechManagementOrders[i].WeaponSwaps[j].SlotLocation)
			data.MechManagementOrders[i].WeaponSwaps[j].NewWeaponID = strings.TrimSpace(data.MechManagementOrders[i].WeaponSwaps[j].NewWeaponID)
		}
	}
}
