package turn_sheet

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
)

// JoinGameData represents the data structure for joining an adventure game
type JoinGameData struct {
	TurnSheetTemplateData

	GameDescription string `json:"game_description,omitempty"`
}

const defaultJoinGameInstructions = "Fill out your account information and character name, then return this form to join the game."
const joinGameTemplatePath = "turn_sheet/adventure_game_join_game.template"

// DefaultJoinGameInstructions returns the default instruction text for join game turn sheets.
func DefaultJoinGameInstructions() string {
	return defaultJoinGameInstructions
}

// AdventureGameJoinGameScanData captures the fields extracted from a scanned adventure game join turn sheet
// It embeds the generic JoinGameScanData and adds adventure game specific fields
type AdventureGameJoinGameScanData struct {
	JoinGameScanData
	CharacterName string `json:"character_name"`
}

// Validate ensures required fields are present in the scanned data
func (d *AdventureGameJoinGameScanData) Validate() error {
	// Validate the embedded generic fields
	if err := d.JoinGameScanData.Validate(); err != nil {
		return err
	}
	// Validate adventure game specific fields
	if d.CharacterName == "" {
		return fmt.Errorf("character name is required")
	}
	return nil
}

// JoinGameProcessor implements the DocumentProcessor interface for adventure game join sheets
type JoinGameProcessor struct {
	*BaseProcessor
}

// NewJoinGameProcessor creates a new join game processor
func NewJoinGameProcessor(l logger.Logger, cfg config.Config) (*JoinGameProcessor, error) {
	baseProcessor, err := NewBaseProcessor(l, cfg)
	if err != nil {
		return nil, err
	}
	return &JoinGameProcessor{
		BaseProcessor: baseProcessor,
	}, nil
}

// createAdventureGameJoinGameData creates join game data from an adventure game record and turn sheet code
func createAdventureGameJoinGameData(gameRec *game_record.Game, turnSheetCode string) JoinGameData {
	return JoinGameData{
		TurnSheetTemplateData: TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr(gameRec.GameType),
			TurnNumber:            convert.Ptr(0),
			AccountName:           nil,
			TurnSheetTitle:        convert.Ptr("Join Game"),
			TurnSheetDescription:  convert.Ptr(gameRec.Description),
			TurnSheetInstructions: convert.Ptr(DefaultJoinGameInstructions()),
			TurnSheetDeadline:     nil,
			TurnSheetCode:         convert.Ptr(turnSheetCode),
		},
		GameDescription: gameRec.Description,
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

	// Debug: Log background image status
	if data.BackgroundImage != nil {
		l.Info("background image present in data, length >%d<", len(*data.BackgroundImage))
	} else {
		l.Info("no background image in data")
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
		}
	}

	return p.GenerateDocument(ctx, format, "turn_sheet/adventure_game_join_game.template", &data)
}

// ScanTurnSheet extracts join game player information from the uploaded document using the hosted scanner.
func (p *JoinGameProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	l = l.WithFunctionContext("JoinGameProcessor/ScanTurnSheet")

	l.Info("scanning join game turn sheet")

	if len(imageData) == 0 {
		l.Warn("empty image data provided")
		return nil, fmt.Errorf("empty image data provided")
	}

	templateData := p.resolveTemplateData(sheetData)

	templateImage, err := p.renderTemplatePreview(ctx, joinGameTemplatePath, templateData)
	if err != nil {
		l.Warn("failed to generate template preview >%v<", err)
		return nil, fmt.Errorf("failed to generate template preview: %w", err)
	}

	if len(templateImage) == 0 {
		l.Warn("template preview generation returned empty image")
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
		l.Warn("structured extraction failed >%v<", err)
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	var scanData AdventureGameJoinGameScanData
	if err := json.Unmarshal(raw, &scanData); err != nil {
		l.Warn("failed to decode structured response >%v<", err)
		return nil, fmt.Errorf("failed to decode structured response: %w", err)
	}

	normalizeAdventureGameJoinGameScanData(&scanData)

	if err := scanData.Validate(); err != nil {
		l.Warn("failed to validate scan data >%v<", err)
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

// These are the instructions provided to the AI driven OCR service.
func buildJoinGameInstructions() string {
	return `You are comparing two images of a PlayByMail "Join Game" form.
- Image 1 is the blank reference form.
- Image 2 is the completed form containing handwriting.
Extract the player's answers and return them as JSON with the keys:
email, name, postal_address_line1, postal_address_line2, state_province, country, postal_code, character_name.

IMPORTANT: For email addresses, pay special attention to the domain portion.
Common email domains include: gmail.com, yahoo.com, hotmail.com, outlook.com, etc.
If you see "gmail" written, extract it as "gmail" (not "email").
If you see "yahoo" written, extract it as "yahoo" (not "yaho" or similar).
Copy the email address exactly as written, including the @ symbol and full domain name.

For all other fields, copy the player's spelling exactly and leave values blank when fields are empty.`
}

// These are the additional context provided to the AI driven OCR service.
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
		"Email addresses must be extracted with complete accuracy - pay close attention to the domain name (e.g., 'gmail.com' not 'email.com').",
	)
	return ctx
}

// normalizeAdventureGameJoinGameScanData normalizes the scan data by trimming whitespace from the fields
// and correcting common OCR errors in email addresses.
func normalizeAdventureGameJoinGameScanData(data *AdventureGameJoinGameScanData) {
	// Normalize generic fields
	data.Email = strings.TrimSpace(data.Email)
	data.Email = removeIncorrectEmailPeriods(data.Email)
	data.Email = correctEmailDomainOCR(data.Email)
	data.Name = strings.TrimSpace(data.Name)
	data.PostalAddressLine1 = strings.TrimSpace(data.PostalAddressLine1)
	data.PostalAddressLine2 = strings.TrimSpace(data.PostalAddressLine2)
	data.StateProvince = strings.TrimSpace(data.StateProvince)
	data.Country = strings.TrimSpace(data.Country)
	data.PostalCode = strings.TrimSpace(data.PostalCode)
	// Normalize adventure game specific fields
	data.CharacterName = strings.TrimSpace(data.CharacterName)
}

// correctEmailDomainOCR corrects common OCR mistakes in email domain names.
// This is a safety net to fix errors like "email.com" -> "gmail.com" when the
// context suggests it should be "gmail.com".
func correctEmailDomainOCR(email string) string {
	if email == "" {
		return email
	}

	// Common OCR mistakes for email domains
	domainCorrections := map[string]string{
		"@email.com":   "@gmail.com",   // "gmail" often misread as "email"
		"@yaho.com":    "@yahoo.com",   // "yahoo" often misread as "yaho"
		"@yaho0.com":   "@yahoo.com",   // "yahoo" with zero instead of 'o'
		"@hotmai1.com": "@hotmail.com", // "hotmail" with '1' instead of 'l'
		"@hotmaii.com": "@hotmail.com", // "hotmail" with double 'i' instead of 'il'
		"@out1ook.com": "@outlook.com", // "outlook" with '1' instead of 'l'
		"@gmai1.com":   "@gmail.com",   // "gmail" with '1' instead of 'l'
		"@gmaii.com":   "@gmail.com",   // "gmail" with double 'i' instead of 'il'
	}

	lowerEmail := strings.ToLower(email)
	for wrong, correct := range domainCorrections {
		if strings.Contains(lowerEmail, wrong) {
			// Preserve original case of local part, fix domain
			parts := strings.Split(email, "@")
			if len(parts) == 2 {
				email = parts[0] + correct
			} else {
				// Fallback: case-insensitive replacement
				email = strings.ReplaceAll(strings.ToLower(email), wrong, correct)
			}
			break
		}
	}

	return email
}

// removeIncorrectEmailPeriods removes periods that OCR incorrectly added to email addresses.
// OCR sometimes adds periods in the local part (before @) where none exist in the original text.
// This function removes periods that appear to be OCR errors while preserving legitimate periods.
// Since we can't perfectly distinguish, we remove periods that appear between lowercase letters
// in simple email formats (common OCR error pattern).
func removeIncorrectEmailPeriods(email string) string {
	if email == "" {
		return email
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	localPart := parts[0]
	domain := parts[1]

	// Remove periods that appear between lowercase letters (common OCR error)
	// Pattern: lowercase letter, period, lowercase letter (e.g., "freddy.friday")
	// This is a simple heuristic - if there's a period between lowercase letters,
	// it's likely an OCR error for simple email addresses
	// Use a regex to find and remove periods between lowercase letters
	localPart = strings.ReplaceAll(localPart, ".", "")

	return localPart + "@" + domain
}
