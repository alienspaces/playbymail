package turn_sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
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

// ScanTurnSheet extracts join game player information from the uploaded document
func (p *JoinGameProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	l = l.WithFunctionContext("JoinGameProcessor/ScanTurnSheet")

	l.Info("scanning join game turn sheet")

	// Extract text from image data
	text, err := p.ExtractTextFromImage(ctx, imageData)
	if err != nil {
		l.Warn("failed to extract text from image >%v<", err)
		return nil, fmt.Errorf("text extraction failed: %w", err)
	}

	// Parse the join game form fields from the extracted text
	scanData, err := p.parseJoinGameText(l, text)
	if err != nil {
		l.Warn("failed to parse join game turn sheet >%v<", err)
		return nil, err
	}

	if err := scanData.Validate(); err != nil {
		l.Warn("validation failed for join game scan data >%v<", err)
		return nil, err
	}

	return json.Marshal(scanData)
}

func (p *JoinGameProcessor) parseJoinGameText(_ logger.Logger, text string) (*JoinGameScanData, error) {

	lines := strings.Split(text, "\n")
	scanData := &JoinGameScanData{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		switch {
		case strings.Contains(lower, "email") && !strings.Contains(lower, "character") && scanData.Email == "":
			if value := extractValue(trimmed); value != "" {
				scanData.Email = value
			}
		case strings.Contains(lower, "name") && !strings.Contains(lower, "character") && scanData.Name == "":
			if value := extractValue(trimmed); value != "" {
				scanData.Name = value
			}
		case strings.Contains(lower, "address line 1") && scanData.PostalAddressLine1 == "":
			if value := extractValue(trimmed); value != "" {
				scanData.PostalAddressLine1 = value
			}
		case strings.Contains(lower, "address line 2") && scanData.PostalAddressLine2 == "":
			if value := extractValue(trimmed); value != "" {
				scanData.PostalAddressLine2 = value
			}
		case strings.Contains(lower, "state") || strings.Contains(lower, "province"):
			if value := extractValue(trimmed); value != "" {
				scanData.StateProvince = value
			}
		case strings.Contains(lower, "country") && scanData.Country == "":
			if value := extractValue(trimmed); value != "" {
				scanData.Country = value
			}
		case (strings.Contains(lower, "post code") || strings.Contains(lower, "postcode")) && scanData.PostalCode == "":
			if value := extractValue(trimmed); value != "" {
				scanData.PostalCode = value
			}
		case strings.Contains(lower, "character name") && scanData.CharacterName == "":
			if value := extractValue(trimmed); value != "" {
				scanData.CharacterName = value
			}
		}
	}

	return scanData, nil
}

func extractValue(line string) string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
