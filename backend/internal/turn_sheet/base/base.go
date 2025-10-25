package base

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/generator"
	"gitlab.com/alienspaces/playbymail/internal/scanner"
)

// BaseProcessor provides common functionality for all turn sheet processors
type BaseProcessor struct {
	Scanner   *scanner.ImageScanner
	Generator *generator.PDFGenerator
	Log       logger.Logger
}

// NewBaseProcessor creates a new base processor
func NewBaseProcessor(l logger.Logger) *BaseProcessor {

	scannerInstance := scanner.NewImageScanner(l)
	generatorInstance := generator.NewPDFGenerator(l)

	return &BaseProcessor{
		Scanner:   scannerInstance,
		Generator: generatorInstance,
		Log:       l,
	}
}

// ExtractTurnSheetCode extracts the turn sheet code from OCR text
// This is common across all turn sheet types
func (bp *BaseProcessor) ExtractTurnSheetCode(text string) (string, error) {
	log := bp.Log.WithFunctionContext("BaseProcessor/ExtractTurnSheetCode")

	log.Info("extracting turn sheet code from text")

	// Look for turn sheet code patterns
	patterns := []string{
		`Turn Sheet Code:\s*([A-Z0-9\-]+)`,
		`Code:\s*([A-Z0-9\-]+)`,
		`Turn Code:\s*([A-Z0-9\-]+)`,
		`Sheet Code:\s*([A-Z0-9\-]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			code := strings.TrimSpace(matches[1])
			if len(code) > 0 {
				log.Info("extracted turn sheet code >%s<", code)
				return code, nil
			}
		}
	}

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

	// Parse the turn sheet code from extracted text
	turnSheetCode, err := bp.ExtractTurnSheetCode(text)
	if err != nil {
		log.Warn("failed to extract turn sheet code from text >%v<", err)
		return "", fmt.Errorf("turn sheet code parsing failed: %w", err)
	}

	return turnSheetCode, nil
}

// ExtractTextFromImage delegates to the scanner for OCR
func (bp *BaseProcessor) ExtractTextFromImage(ctx context.Context, imageData []byte) (string, error) {
	return bp.Scanner.ExtractTextFromImage(ctx, imageData)
}
