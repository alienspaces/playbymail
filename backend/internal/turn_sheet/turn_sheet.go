package turn_sheet

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/generator"
	"gitlab.com/alienspaces/playbymail/internal/scanner"
	"gitlab.com/alienspaces/playbymail/internal/utils/config"
)

// Trigger CI backend tests (One)

// BaseProcessor provides common functionality for all turn sheet processors
type BaseProcessor struct {
	Scanner      *scanner.ImageScanner
	Generator    *generator.PDFGenerator
	Log          logger.Logger
	Config       config.Config
	TemplatePath string
}

// NewBaseProcessor creates a new base processor
func NewBaseProcessor(l logger.Logger, cfg config.Config) *BaseProcessor {

	scannerInstance := scanner.NewImageScanner(l)
	generatorInstance := generator.NewPDFGenerator(l)

	templatePath := "./backend/templates"
	templatePath = cfg.TemplatesPath

	return &BaseProcessor{
		Scanner:      scannerInstance,
		Generator:    generatorInstance,
		Log:          l,
		Config:       cfg,
		TemplatePath: templatePath,
	}
}

// ExtractTurnSheetCode extracts the turn sheet code from OCR text
// This is common across all turn sheet types
func (bp *BaseProcessor) ExtractTurnSheetCode(text string) (string, error) {
	log := bp.Log.WithFunctionContext("BaseProcessor/ExtractTurnSheetCode")

	log.Info("extracting turn sheet code from text")
	log.Debug("searching for turn sheet code in text: >%s<", text)

	var allCodes []string

	// Look for turn sheet code patterns - try with case-insensitive matching
	patterns := []string{
		`[Tt]urn [Ss]heet [Cc]ode:\s*([A-Z0-9\-]+)`,
		`[Cc]ode:\s*([A-Z0-9\-]+)`,
		`[Tt]urn [Cc]ode:\s*([A-Z0-9\-]+)`,
		`[Ss]heet [Cc]ode:\s*([A-Z0-9\-]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) > 1 {
				code := strings.TrimSpace(match[1])
				if len(code) > 0 {
					allCodes = append(allCodes, code)
				}
			}
		}
	}

	// Try a more flexible pattern - look for codes that appear AFTER "Turn Sheet Code" label
	turnSheetCodeLabelPattern := regexp.MustCompile(`[Tt]urn [Ss]heet [Cc]ode:\s*.*?([A-Z0-9]{6,12})`)
	matches := turnSheetCodeLabelPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 1 && len(match[1]) >= 6 {
			allCodes = append(allCodes, match[1])
		}
	}

	// Fallback: look for any long alphanumeric string (6+ characters)
	flexiblePattern := regexp.MustCompile(`([A-Z0-9]{6,12})`)
	longMatches := flexiblePattern.FindAllString(text, -1)
	for _, match := range longMatches {
		if len(match) >= 6 {
			allCodes = append(allCodes, match)
		}
	}

	// Return the best matching code
	// Prefer codes that look like game codes (alphanumeric, not just numeric)
	if len(allCodes) > 0 {
		var bestCode string
		var bestScore int

		for _, code := range allCodes {
			score := 0
			// Prefer alphanumeric over pure numeric
			hasLetters := false
			hasNumbers := false
			for _, char := range code {
				if char >= 'A' && char <= 'Z' {
					hasLetters = true
				} else if char >= '0' && char <= '9' {
					hasNumbers = true
				}
			}
			if hasLetters && hasNumbers {
				score += 10 // Alphanumeric gets bonus
			} else if hasLetters {
				score += 5 // Letters only
			}

			// Prefer reasonable length (6-12 characters)
			if len(code) >= 6 && len(code) <= 12 {
				score += 5
			}

			// Longer codes get slight bonus
			score += len(code) / 2

			if score > bestScore {
				bestScore = score
				bestCode = code
			}
		}

		log.Info("extracted turn sheet code >%s< from candidates: %v", bestCode, allCodes)
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
