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

// JoinGameData represents the data structure for joining an adventure game
type JoinGameData struct {
	TurnSheetTemplateData

	GameDescription string `json:"game_description,omitempty"`
}

const defaultJoinGameInstructions = "Fill out your account information and character name, then return this form to join the game."

// DefaultJoinGameInstructions returns the default instruction text for join game turn sheets.
func DefaultJoinGameInstructions() string {
	return defaultJoinGameInstructions
}

// JoinGameScanData captures the fields extracted from a scanned join game turn sheet
type JoinGameScanData struct {
	Email              string `json:"email"`
	Name               string `json:"name"`
	PostalAddressLine1 string `json:"postal_address_line1"`
	PostalAddressLine2 string `json:"postal_address_line2,omitempty"`
	StateProvince      string `json:"state_province"`
	Country            string `json:"country"`
	PostalCode         string `json:"postal_code"`
	CharacterName      string `json:"character_name"`
}

// Validate ensures required fields are present in the scanned data
func (d *JoinGameScanData) Validate() error {
	switch {
	case d.Email == "":
		return fmt.Errorf("email is required")
	case d.Name == "":
		return fmt.Errorf("name is required")
	case d.PostalAddressLine1 == "":
		return fmt.Errorf("postal address line 1 is required")
	case d.StateProvince == "":
		return fmt.Errorf("state or province is required")
	case d.Country == "":
		return fmt.Errorf("country is required")
	case d.PostalCode == "":
		return fmt.Errorf("post code is required")
	case d.CharacterName == "":
		return fmt.Errorf("character name is required")
	default:
		return nil
	}
}

// JoinGameProcessor implements the DocumentProcessor interface for adventure game join sheets
type JoinGameProcessor struct {
	*BaseProcessor
}

// NewJoinGameProcessor creates a new join game processor
func NewJoinGameProcessor(l logger.Logger, cfg config.Config) *JoinGameProcessor {
	return &JoinGameProcessor{
		BaseProcessor: NewBaseProcessor(l, cfg),
	}
}

// GenerateTurnSheet generates a join game turn sheet document
func (p *JoinGameProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, sheetData []byte) ([]byte, error) {
	l = l.WithFunctionContext("JoinGameProcessor/GenerateTurnSheet")

	l.Info("generating join game turn sheet")

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
		instruction := defaultJoinGameInstructions
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
		} else if data.GameName != nil && strings.TrimSpace(*data.GameName) != "" {
			desc := fmt.Sprintf("Welcome to %s! Welcome to the PlayByMail Adventure!", *data.GameName)
			data.TurnSheetDescription = &desc
		}
	}

	return p.GenerateDocument(ctx, format, "turn_sheet/adventure_game_join_game.template", &data)
}

const joinGameTemplatePath = "turn_sheet/adventure_game_join_game.template"

// ScanTurnSheet extracts join game player information from the uploaded document using the hosted scanner.
func (p *JoinGameProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	l = l.WithFunctionContext("JoinGameProcessor/ScanTurnSheet")

	l.Info("scanning join game turn sheet")

	if len(imageData) == 0 {
		return nil, fmt.Errorf("empty image data provided")
	}

	templateData := p.resolveTemplateData(sheetData)

	templateImage, err := p.renderTemplatePreview(ctx, joinGameTemplatePath, templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to generate template preview: %w", err)
	}
	if len(templateImage) == 0 {
		return nil, fmt.Errorf("template preview generation returned empty image")
	}

	req := scanner.StructuredScanRequest{
		Instructions:       buildJoinGameInstructions(),
		AdditionalContext:  buildJoinGameContext(templateData),
		TemplateImage:      templateImage,
		TemplateImageMIME:  "image/png",
		FilledImage:        imageData,
		ExpectedJSONSchema: joinGameExpectedSchema(),
	}

	raw, err := p.Scanner.ExtractStructuredData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	var scanData JoinGameScanData
	if err := json.Unmarshal(raw, &scanData); err != nil {
		return nil, fmt.Errorf("failed to decode structured response: %w", err)
	}

	normalizeJoinGameScanData(&scanData)

	if err := scanData.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(scanData)
}

func (p *JoinGameProcessor) resolveTemplateData(sheetData []byte) *JoinGameData {
	var data JoinGameData
	if len(sheetData) > 0 {
		if err := json.Unmarshal(sheetData, &data); err != nil {
			return defaultJoinGameTemplateData()
		}
		return &data
	}
	return defaultJoinGameTemplateData()
}

func defaultJoinGameTemplateData() *JoinGameData {
	title := "Join Game"
	instructions := DefaultJoinGameInstructions()
	return &JoinGameData{
		TurnSheetTemplateData: TurnSheetTemplateData{
			TurnSheetTitle:        &title,
			TurnSheetInstructions: &instructions,
		},
	}
}

func joinGameExpectedSchema() map[string]any {
	return map[string]any{
		"email":                "",
		"name":                 "",
		"postal_address_line1": "",
		"postal_address_line2": "",
		"state_province":       "",
		"country":              "",
		"postal_code":          "",
		"character_name":       "",
	}
}

func buildJoinGameInstructions() string {
	return `You are comparing two images of a PlayByMail "Join Game" form.
- Image 1 is the blank reference form.
- Image 2 is the completed form containing handwriting.
Extract the player's answers and return them as JSON with the keys:
email, name, postal_address_line1, postal_address_line2, state_province, country, postal_code, character_name.
Copy the player's spelling exactly and leave values blank when fields are empty.`
}

func buildJoinGameContext(data *JoinGameData) []string {
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
	)
	return ctx
}

func normalizeJoinGameScanData(data *JoinGameScanData) {
	data.Email = strings.TrimSpace(data.Email)
	data.Name = strings.TrimSpace(data.Name)
	data.PostalAddressLine1 = strings.TrimSpace(data.PostalAddressLine1)
	data.PostalAddressLine2 = strings.TrimSpace(data.PostalAddressLine2)
	data.StateProvince = strings.TrimSpace(data.StateProvince)
	data.Country = strings.TrimSpace(data.Country)
	data.PostalCode = strings.TrimSpace(data.PostalCode)
	data.CharacterName = strings.TrimSpace(data.CharacterName)
}
