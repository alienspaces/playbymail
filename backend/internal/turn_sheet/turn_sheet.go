package turn_sheet

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/generator"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/scanner"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// GetTurnSheetJoinGameData creates join game data for a game based on its game type
func GetTurnSheetJoinGameData(gameRec *game_record.Game, turnSheetCode string) (JoinGameData, error) {
	switch gameRec.GameType {
	case game_record.GameTypeAdventure:
		return createAdventureGameJoinGameData(gameRec, turnSheetCode), nil
	default:
		return JoinGameData{}, fmt.Errorf("join game turn sheets not supported for game type: %s", gameRec.GameType)
	}
}

// BaseProcessor provides common functionality for all turn sheet processors
type BaseProcessor struct {
	Scanner      *scanner.ImageScanner
	Generator    *generator.PDFGenerator
	Log          logger.Logger
	Config       config.Config
	TemplatePath string
}

// NewBaseProcessor creates a new base processor
func NewBaseProcessor(l logger.Logger, cfg config.Config) (*BaseProcessor, error) {
	l = l.WithFunctionContext("NewBaseProcessor")

	l.Info("creating base processor")

	scannerInstance, err := scanner.NewImageScanner(l, cfg)
	if err != nil {
		l.Warn("failed to create image scanner >%v<", err)
		return nil, fmt.Errorf("failed to create image scanner: %w", err)
	}

	generatorInstance, err := generator.NewPDFGenerator(l)
	if err != nil {
		l.Warn("failed to create PDF generator >%v<", err)
		return nil, fmt.Errorf("failed to create PDF generator: %w", err)
	}

	templatePath := "./backend/templates"
	templatePath = cfg.TemplatesPath

	return &BaseProcessor{
		Scanner:      scannerInstance,
		Generator:    generatorInstance,
		Log:          l,
		Config:       cfg,
		TemplatePath: templatePath,
	}, nil
}

// GenerateDocument renders a template in the requested document format.
func (bp *BaseProcessor) GenerateDocument(ctx context.Context, format DocumentFormat, templatePath string, data any) ([]byte, error) {
	l := bp.Log.WithFunctionContext("BaseProcessor/GenerateDocument")

	bp.Generator.SetTemplatePath(bp.TemplatePath)

	switch format {
	case DocumentFormatPDF, "":
		l.Info("generating PDF document template=%s", templatePath)
		return bp.Generator.GeneratePDF(ctx, templatePath, data)
	case DocumentFormatHTML:
		l.Info("generating HTML document template=%s", templatePath)
		html, err := bp.Generator.GenerateHTML(ctx, templatePath, data)
		if err != nil {
			return nil, err
		}
		return []byte(html), nil
	default:
		return nil, fmt.Errorf("unsupported document format: %s", format)
	}
}

// ExtractTurnSheetCode extracts the turn sheet code from OCR text
// This is common across all turn sheet types
// Turn sheet codes are base64 URL-encoded JSON strings (100+ characters)
func (bp *BaseProcessor) ExtractTurnSheetCode(text string) (string, error) {
	log := bp.Log.WithFunctionContext("BaseProcessor/ExtractTurnSheetCode")

	log.Info("extracting turn sheet code from text")
	log.Debug("searching for turn sheet code in text: >%s<", text)

	var allCodes []string

	// Primary pattern: Look for codes after "Turn Sheet Code:" label
	// This captures both short test codes and long base64 production codes
	// Base64 URL encoding uses: A-Za-z0-9_- and may have = padding
	// The code may be on the same line or the next line after the label
	labelPatterns := []string{
		// Match "Turn Sheet Code:" followed by code on same line or next line
		`[Tt]urn\s+[Ss]heet\s+[Cc]ode:\s*\**\s*\n?\s*([A-Za-z0-9_\-=]+)`,
		// Match just "Code:" followed by code
		`[Cc]ode:\s*\**\s*\n?\s*([A-Za-z0-9_\-=]+)`,
	}

	for _, pattern := range labelPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) > 1 {
				// Clean up the code - remove any whitespace/newlines that OCR might have added
				code := strings.TrimSpace(match[1])
				code = strings.ReplaceAll(code, " ", "")
				code = strings.ReplaceAll(code, "\n", "")
				code = strings.ReplaceAll(code, "\r", "")
				if len(code) >= 6 { // Minimum reasonable code length
					allCodes = append(allCodes, code)
				}
			}
		}
	}

	// Secondary: Look for long base64-like strings (50+ characters) anywhere in text
	// This catches cases where the label wasn't recognized but the code is present
	longBase64Pattern := regexp.MustCompile(`([A-Za-z0-9_\-]{50,}[A-Za-z0-9_\-=]*)`)
	longMatches := longBase64Pattern.FindAllString(text, -1)
	for _, match := range longMatches {
		// Clean up whitespace
		cleanMatch := strings.ReplaceAll(match, " ", "")
		cleanMatch = strings.ReplaceAll(cleanMatch, "\n", "")
		if len(cleanMatch) >= 50 {
			allCodes = append(allCodes, cleanMatch)
		}
	}

	// Return the best matching code
	if len(allCodes) > 0 {
		var bestCode string
		var bestScore int

		for _, code := range allCodes {
			score := 0

			// Strongly prefer longer codes - base64 encoded JSON is typically 100+ chars
			if len(code) >= 100 {
				score += 100
			} else if len(code) >= 50 {
				score += 50
			}

			// Add length bonus
			score += len(code)

			// Check if it looks like valid base64 JSON (starts with "ey" which is "{" in base64)
			if strings.HasPrefix(code, "ey") {
				score += 50 // JSON objects start with "{" which encodes to "ey"
			}

			if score > bestScore {
				bestScore = score
				bestCode = code
			}
		}

		log.Info("extracted turn sheet code >%s< (length: %d) from candidates: %v", bestCode, len(bestCode), allCodes)
		return bestCode, nil
	}

	log.Warn("no turn sheet code found in text: >%s<", text)
	return "", fmt.Errorf("no turn sheet code found in text")
}

// ParseTurnSheetCodeFromImage extracts and parses a turn sheet code from
// scanned image data
func (bp *BaseProcessor) ParseTurnSheetCodeFromImage(ctx context.Context, imageData []byte) (string, error) {
	log := bp.Log.WithFunctionContext("BaseProcessor/ParseTurnSheetCodeFromImage")

	log.Info("parsing turn sheet code from image data")

	// Extract all text from image
	text, err := bp.Scanner.ExtractTextFromImage(ctx, imageData)
	if err != nil {
		log.Warn("failed to extract text from image >%v<", err)
		return "", err
	}

	// Log the full OCR text for debugging
	log.Debug("full OCR text extracted: >%s<", text)

	// Parse the turn sheet code from extracted text
	turnSheetCode, err := bp.ExtractTurnSheetCode(text)
	if err != nil {
		log.Warn("failed to extract turn sheet code from text >%v<", err)
		log.Debug("attempting to extract from full text: >%s<", text)
		return "", fmt.Errorf("turn sheet code parsing failed: %w", err)
	}

	return turnSheetCode, nil
}

// ExtractTextFromImage delegates to the scanner for OCR
func (bp *BaseProcessor) ExtractTextFromImage(ctx context.Context, imageData []byte) (string, error) {
	return bp.Scanner.ExtractTextFromImage(ctx, imageData)
}

// renderTemplatePreview renders the specified template with the provided data
// and returns a PNG representation suitable for sending to hosted OCR services.
func (bp *BaseProcessor) renderTemplatePreview(ctx context.Context, templatePath string, data any) ([]byte, error) {
	l := bp.Log.WithFunctionContext("BaseProcessor/renderTemplatePreview")

	bp.Generator.SetTemplatePath(bp.TemplatePath)

	png, err := bp.Generator.GeneratePNG(ctx, templatePath, data)
	if err != nil {
		l.Warn("failed to render template preview >%v<", err)
		return nil, err
	}

	return png, nil
}

// ValidateBaseTemplateData validates the required base fields for turn sheet generation
func (bp *BaseProcessor) ValidateBaseTemplateData(data *TurnSheetTemplateData) error {
	l := bp.Log.WithFunctionContext("BaseProcessor/ValidateBaseTemplateData")

	if data == nil {
		return fmt.Errorf("template data is nil")
	}

	// Game name is required
	if data.GameName == nil || *data.GameName == "" {
		l.Warn("game name is missing or empty")
		return fmt.Errorf("game name is required")
	}

	// Turn number is required
	if data.TurnNumber == nil {
		l.Warn("turn number is missing")
		return fmt.Errorf("turn number is required")
	}

	// Turn sheet code is required
	if data.TurnSheetCode == nil || *data.TurnSheetCode == "" {
		l.Warn("turn sheet code is missing or empty")
		return fmt.Errorf("turn sheet code is required")
	}

	return nil
}

// JoinGameScanData captures the generic fields extracted from a scanned join game turn sheet
// This includes email and contact information that is common across all game types
type JoinGameScanData struct {
	Email              string `json:"email"`
	Name               string `json:"name"`
	PostalAddressLine1 string `json:"postal_address_line1"`
	PostalAddressLine2 string `json:"postal_address_line2,omitempty"`
	StateProvince      string `json:"state_province"`
	Country            string `json:"country"`
	PostalCode         string `json:"postal_code"`
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
	default:
		return nil
	}
}
