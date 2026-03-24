package turnsheet

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/record"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/scanner"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
	"gitlab.com/alienspaces/playbymail/internal/utils/turnsheetutil"
)

const defaultMechaJoinGameInstructions = "Fill out your account information and commander name, then return this form to join the game."
const mechaJoinGameTemplatePath = "turnsheet/mecha_join_game.template"

// DefaultMechaJoinGameInstructions returns the default instruction text for mecha join game turn sheets.
func DefaultMechaJoinGameInstructions() string {
	return defaultMechaJoinGameInstructions
}

// MechaJoinGameScanData captures the fields extracted from a scanned mecha join turn sheet.
// It embeds the generic JoinGameScanData and adds mecha specific fields.
type MechaJoinGameScanData struct {
	JoinGameScanData
	CommanderName string `json:"commander_name"`
}

// Validate ensures required fields are present in the scanned data.
func (d *MechaJoinGameScanData) Validate() error {
	if err := d.JoinGameScanData.Validate(); err != nil {
		return err
	}
	if d.CommanderName == "" {
		return fmt.Errorf("commander name is required")
	}
	return nil
}

// MechaJoinGameProcessor implements the DocumentProcessor interface for mecha join sheets.
type MechaJoinGameProcessor struct {
	*BaseProcessor
}

// NewMechaJoinGameProcessor creates a new mecha join game processor.
func NewMechaJoinGameProcessor(l logger.Logger, cfg config.Config) (*MechaJoinGameProcessor, error) {
	baseProcessor, err := NewBaseProcessor(l, cfg)
	if err != nil {
		return nil, err
	}
	return &MechaJoinGameProcessor{
		BaseProcessor: baseProcessor,
	}, nil
}

// createMechaJoinGameData creates join game data from a mecha record and turn sheet code.
func createMechaJoinGameData(gameRec *game_record.Game, turnSheetCode string) JoinGameData {
	return JoinGameData{
		TurnSheetTemplateData: TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr(gameRec.GameType),
			TurnNumber:            convert.Ptr(0),
			AccountName:           nil,
			TurnSheetTitle:        convert.Ptr("Join Game"),
			TurnSheetDescription:  convert.Ptr(gameRec.Description),
			TurnSheetInstructions: convert.Ptr(DefaultMechaJoinGameInstructions()),
			TurnSheetDeadline:     nil,
			TurnSheetCode:         convert.Ptr(turnSheetCode),
			HideNarrative:         true,
		},
		GameDescription: gameRec.Description,
	}
}

// GeneratePreviewData generates dummy data for a mecha join game turn sheet preview.
func (p *MechaJoinGameProcessor) GeneratePreviewData(ctx context.Context, l logger.Logger, gameRec *game_record.Game, backgroundImage *string) ([]byte, error) {
	l = l.WithFunctionContext("MechaJoinGameProcessor/GeneratePreviewData")

	turnSheetCode, err := turnsheetutil.GenerateJoinGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		l.Warn("failed to generate join turn sheet code >%v<", err)
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	turnSheetData := createMechaJoinGameData(gameRec, turnSheetCode)

	if backgroundImage != nil && *backgroundImage != "" {
		turnSheetData.BackgroundImage = backgroundImage
	}

	return json.Marshal(turnSheetData)
}

// GenerateTurnSheet generates a mecha join game turn sheet document.
func (p *MechaJoinGameProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, sheetData []byte) ([]byte, error) {
	l = l.WithFunctionContext("MechaJoinGameProcessor/GenerateTurnSheet")

	var data JoinGameData
	if err := json.Unmarshal(sheetData, &data); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to parse sheet data: %w", err)
	}

	if err := p.ValidateBaseTemplateData(&data.TurnSheetTemplateData); err != nil {
		l.Warn("failed to validate base template data >%v<", err)
		return nil, fmt.Errorf("template data validation failed: %w", err)
	}

	if data.TurnSheetInstructions == nil || strings.TrimSpace(*data.TurnSheetInstructions) == "" {
		instruction := defaultMechaJoinGameInstructions
		data.TurnSheetInstructions = &instruction
	}

	if data.TurnSheetTitle == nil || strings.TrimSpace(*data.TurnSheetTitle) == "" {
		title := "Join Game"
		data.TurnSheetTitle = &title
	}

	if data.TurnSheetDescription == nil || strings.TrimSpace(*data.TurnSheetDescription) == "" {
		if data.GameDescription != "" {
			desc := data.GameDescription
			data.TurnSheetDescription = &desc
		}
	}

	return p.GenerateDocument(ctx, format, mechaJoinGameTemplatePath, &data)
}

// ScanTurnSheet extracts join game information from the uploaded document using the hosted scanner.
func (p *MechaJoinGameProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	l = l.WithFunctionContext("MechaJoinGameProcessor/ScanTurnSheet")

	if len(imageData) == 0 {
		l.Warn("empty image data provided")
		return nil, fmt.Errorf("empty image data provided")
	}

	templateData := p.resolveMechaJoinGameTemplateData(sheetData)

	templateImage, err := p.renderTemplatePreview(ctx, mechaJoinGameTemplatePath, templateData)
	if err != nil {
		l.Warn("failed to generate template preview >%v<", err)
		return nil, fmt.Errorf("failed to generate template preview: %w", err)
	}

	if len(templateImage) == 0 {
		return nil, fmt.Errorf("template preview generation returned empty image")
	}

	req := scanner.StructuredScanRequest{
		Instructions:       buildMechaJoinGameInstructions(),
		AdditionalContext:  buildMechaJoinGameContext(templateData),
		TemplateImage:      templateImage,
		TemplateImageMIME:  "image/png",
		FilledImage:        imageData,
		ExpectedJSONSchema: mechaJoinGameExpectedSchema(),
	}

	raw, err := p.Scanner.ExtractStructuredData(ctx, req)
	if err != nil {
		l.Warn("structured extraction failed >%v<", err)
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	var scanData MechaJoinGameScanData
	if err := json.Unmarshal(raw, &scanData); err != nil {
		l.Warn("failed to decode structured response >%v<", err)
		return nil, fmt.Errorf("failed to decode structured response: %w", err)
	}

	normalizeMechaJoinGameScanData(&scanData)

	if err := scanData.Validate(); err != nil {
		l.Warn("failed to validate scan data >%v<", err)
		return nil, err
	}

	return json.Marshal(scanData)
}

func (p *MechaJoinGameProcessor) resolveMechaJoinGameTemplateData(sheetData []byte) *JoinGameData {
	var data JoinGameData
	if len(sheetData) > 0 {
		if err := json.Unmarshal(sheetData, &data); err != nil {
			return defaultMechaJoinGameTemplateData()
		}
		return &data
	}
	return defaultMechaJoinGameTemplateData()
}

func defaultMechaJoinGameTemplateData() *JoinGameData {
	title := "Join Game"
	instructions := DefaultMechaJoinGameInstructions()
	return &JoinGameData{
		TurnSheetTemplateData: TurnSheetTemplateData{
			TurnSheetTitle:        &title,
			TurnSheetInstructions: &instructions,
		},
	}
}

func mechaJoinGameExpectedSchema() map[string]any {
	return map[string]any{
		"email":                "",
		"name":                 "",
		"postal_address_line1": "",
		"postal_address_line2": "",
		"state_province":       "",
		"country":              "",
		"postal_code":          "",
		"delivery_method":      "",
		"commander_name":       "",
	}
}

func buildMechaJoinGameInstructions() string {
	return `You are comparing two images of a PlayByMail "Join Game" form for a mecha.
- Image 1 is the blank reference form.
- Image 2 is the completed form containing handwriting.
Extract the player's answers and return them as JSON with the keys:
email, name, postal_address_line1, postal_address_line2, state_province, country, postal_code, delivery_method, commander_name.

For delivery_method: extract the chosen option as one of "email", "local", or "post". Leave empty if no delivery selection is present.

NOTE: Some forms may not include postal address fields when the game only supports email or local pickup delivery.

IMPORTANT: For email addresses, pay special attention to the domain portion.
Common email domains include: gmail.com, yahoo.com, hotmail.com, outlook.com, etc.
Copy the email address exactly as written, including the @ symbol and full domain name.

For all other fields, copy the player's spelling exactly and leave values blank when fields are empty.`
}

func buildMechaJoinGameContext(data *JoinGameData) []string {
	var ctx []string
	if data != nil {
		if data.GameName != nil {
			ctx = append(ctx, fmt.Sprintf("Game name: %s", strings.TrimSpace(*data.GameName)))
		}
		if data.GameDescription != "" {
			ctx = append(ctx, fmt.Sprintf("Game description: %s", data.GameDescription))
		}
	}
	ctx = append(ctx,
		"The JSON must only contain the requested keys.",
		"Return an empty string when the player left a field blank.",
		"Email addresses must be extracted with complete accuracy - pay close attention to the domain name.",
	)
	return ctx
}

func normalizeMechaJoinGameScanData(data *MechaJoinGameScanData) {
	data.Email = strings.TrimSpace(data.Email)
	data.Email = removeIncorrectEmailPeriods(data.Email)
	data.Email = correctEmailDomainOCR(data.Email)
	data.Name = strings.TrimSpace(data.Name)
	data.PostalAddressLine1 = strings.TrimSpace(data.PostalAddressLine1)
	data.PostalAddressLine2 = strings.TrimSpace(data.PostalAddressLine2)
	data.StateProvince = strings.TrimSpace(data.StateProvince)
	data.Country = strings.TrimSpace(data.Country)
	data.PostalCode = strings.TrimSpace(data.PostalCode)
	data.DeliveryMethod = strings.ToLower(strings.TrimSpace(data.DeliveryMethod))
	data.CommanderName = strings.TrimSpace(data.CommanderName)
}
