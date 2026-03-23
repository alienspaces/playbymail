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

// OrdersScannedDataSchemaName is the filename of the JSON schema for orders scanned_data (under schema/turnsheet/mech_wargame/).
const OrdersScannedDataSchemaName = "orders.schema.json"

const defaultOrdersInstructions = "For each mech in your lance, choose a sector to move to and/or a target to attack, then return this form."
const ordersTemplatePath = "turnsheet/mech_wargame_orders.template"

// DefaultOrdersInstructions returns the default instruction text for orders turn sheets.
func DefaultOrdersInstructions() string {
	return defaultOrdersInstructions
}

// MechOrderEntry is a single mech's orders for a turn.
type MechOrderEntry struct {
	MechInstanceID              string `json:"mech_instance_id"`
	MechCallsign                string `json:"mech_callsign,omitempty"`
	MechStatus                  string `json:"mech_status,omitempty"`
	CurrentSectorName           string `json:"current_sector_name,omitempty"`
	MoveToSectorInstanceID      string `json:"move_to_sector_instance_id,omitempty"`
	AttackTargetMechInstanceID  string `json:"attack_target_mech_instance_id,omitempty"`
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

// OrdersData is the data model for a mech wargame orders turn sheet.
type OrdersData struct {
	TurnSheetTemplateData

	// Lance information
	LanceName string `json:"lance_name,omitempty"`

	// Mechs in this lance with their current state and available orders
	LanceMechs []MechOrderEntry `json:"lance_mechs,omitempty"`

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

// MechWargameOrdersProcessor implements the DocumentProcessor interface for mech wargame orders sheets.
type MechWargameOrdersProcessor struct {
	*BaseProcessor
}

// NewMechWargameOrdersProcessor creates a new mech wargame orders processor.
func NewMechWargameOrdersProcessor(l logger.Logger, cfg config.Config) (*MechWargameOrdersProcessor, error) {
	baseProcessor, err := NewBaseProcessor(l, cfg)
	if err != nil {
		return nil, err
	}
	return &MechWargameOrdersProcessor{
		BaseProcessor: baseProcessor,
	}, nil
}

// GeneratePreviewData generates dummy data for a mech wargame orders turn sheet preview.
func (p *MechWargameOrdersProcessor) GeneratePreviewData(ctx context.Context, l logger.Logger, gameRec *game_record.Game, backgroundImage *string) ([]byte, error) {
	l = l.WithFunctionContext("MechWargameOrdersProcessor/GeneratePreviewData")

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
		},
		LanceName: "Preview Lance",
		LanceMechs: []MechOrderEntry{
			{
				MechInstanceID:    "preview-mech-1",
				MechCallsign:      "Hammer",
				MechStatus:        "operational",
				CurrentSectorName: "Central Wastes",
			},
			{
				MechInstanceID:    "preview-mech-2",
				MechCallsign:      "Anvil",
				MechStatus:        "operational",
				CurrentSectorName: "Central Wastes",
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

// GenerateTurnSheet generates a mech wargame orders turn sheet document.
func (p *MechWargameOrdersProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, sheetData []byte) ([]byte, error) {
	l = l.WithFunctionContext("MechWargameOrdersProcessor/GenerateTurnSheet")

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
func (p *MechWargameOrdersProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	l = l.WithFunctionContext("MechWargameOrdersProcessor/ScanTurnSheet")

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

func (p *MechWargameOrdersProcessor) resolveOrdersTemplateData(sheetData []byte) *OrdersData {
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
				"mech_instance_id":              "",
				"move_to_sector_instance_id":    "",
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
		if data.LanceName != "" {
			ctx = append(ctx, fmt.Sprintf("Lance name: %s", data.LanceName))
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
