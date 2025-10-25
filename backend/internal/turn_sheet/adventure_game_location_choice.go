package turn_sheet

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"gitlab.com/alienspaces/playbymail/core/config"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
)

type LocationChoiceData struct {
	TurnSheetTemplateData

	// Current location information
	LocationName        string
	LocationDescription string

	// Available location options
	LocationOptions []LocationOption
}

// LocationOption represents a location choice option for the player
type LocationOption struct {
	LocationID              string
	LocationLinkName        string
	LocationLinkDescription string
}

// LocationChoiceScanData represents the scanned data from a location choice turn sheet
type LocationChoiceScanData struct {
	Choices []string
}

// LocationChoiceProcessor implements the DocumentProcessor interface for location choice turn sheets
type LocationChoiceProcessor struct {
	*BaseProcessor
}

// NewLocationChoiceProcessor creates a new location choice processor
func NewLocationChoiceProcessor(l logger.Logger, cfg *config.Config) *LocationChoiceProcessor {
	return &LocationChoiceProcessor{
		BaseProcessor: NewBaseProcessor(l, cfg),
	}
}

// GenerateTurnSheet generates a location choice turn sheet PDF
func (p *LocationChoiceProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, data any) ([]byte, error) {
	l = l.WithFunctionContext("LocationChoiceProcessor/GenerateTurnSheet")

	l.Info("generating location choice turn sheet")

	// Validate and assert data type
	locationChoiceData, ok := data.(*LocationChoiceData)
	if !ok {
		l.Warn("invalid data type for location choice turn sheet, expected *LocationChoiceData")
		return nil, fmt.Errorf("invalid data type: expected *LocationChoiceData, got %T", data)
	}

	// Validate base template data
	if err := p.ValidateBaseTemplateData(&locationChoiceData.TurnSheetTemplateData); err != nil {
		l.Warn("failed to validate base template data >%v<", err)
		return nil, fmt.Errorf("template data validation failed: %w", err)
	}

	// Validate location-specific data
	if locationChoiceData.LocationName == "" {
		l.Warn("location name is missing")
		return nil, fmt.Errorf("location name is required")
	}

	if len(locationChoiceData.LocationOptions) == 0 {
		l.Warn("no location options provided")
		return nil, fmt.Errorf("at least one location option is required")
	}

	// Set the template path on the generator
	p.Generator.SetTemplatePath(p.TemplatePath)

	// Generate PDF using the location choice template
	templatePath := "turn_sheet/adventure_game.location_choice.template"

	return p.Generator.GeneratePDF(ctx, templatePath, data)
}

// ScanTurnSheet scans a location choice turn sheet and extracts player choices
func (p *LocationChoiceProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, imageData []byte, sheetData any) (any, error) {
	log := l.WithFunctionContext("LocationChoiceProcessor/ScanTurnSheet")

	log.Info("scanning location choice turn sheet")

	// Parse sheet data to get location information
	sheetDataMap, ok := sheetData.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid sheet data format for location choice turn sheet")
	}

	// Extract text from image using base processor
	text, err := p.ExtractTextFromImage(ctx, imageData)
	if err != nil {
		log.Warn("failed to extract text from image >%v<", err)
		return nil, fmt.Errorf("text extraction failed: %w", err)
	}

	// Parse location choices from extracted text using the sheet data
	log.Debug("full OCR text extracted: >%s<", text)
	choices, err := p.parseLocationChoicesWithSheetData(text, sheetDataMap)
	if err != nil {
		log.Warn("failed to parse location choices >%v<", err)
		return nil, fmt.Errorf("location choice parsing failed: %w", err)
	}

	log.Info("extracted %d location choices", len(choices))

	return &LocationChoiceScanData{
		Choices: choices,
	}, nil
}

// ParseLocationChoicesWithSheetData parses location choices using the actual sheet data
func (p *LocationChoiceProcessor) parseLocationChoicesWithSheetData(text string, sheetData map[string]any) ([]string, error) {
	// Extract location data from sheet data
	locations, ok := sheetData["locations"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no locations found in sheet data")
	}

	// Check if locations array is empty
	if len(locations) == 0 {
		return nil, fmt.Errorf("no locations found in sheet data")
	}

	// Create a map of location names for validation
	locationNames := make(map[string]string)
	for _, loc := range locations {
		if locMap, ok := loc.(map[string]interface{}); ok {
			if name, exists := locMap["name"]; exists {
				if nameStr, ok := name.(string); ok {
					locationNames[nameStr] = nameStr
					// Also add lowercase version for matching
					locationNames[strings.ToLower(nameStr)] = nameStr
				}
			}
		}
	}

	var choices []string
	seen := make(map[string]bool)

	// Look for checked/marked location options using OCR patterns
	patterns := []string{
		// OCR variations - checked boxes from real scanned images
		`O/\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`, // OCR reads filled checkbox as O/
		`Sf\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`, // OCR reads ☑ as Sf
		`S\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,  // OCR reads ☑ as S

		// Original patterns - checked boxes
		`☑\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,
		`\[X\]\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,
		`✓\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(text, -1)

		for _, match := range matches {
			if len(match) > 1 {
				choice := strings.TrimSpace(match[1])

				// Only accept choices that look like location names (2+ words or single capitalized word)
				// Filter out common false matches
				if len(choice) > 0 &&
					(strings.Contains(choice, " ") || (len(choice) > 3 && choice[0] >= 'A' && choice[0] <= 'Z')) &&
					!strings.Contains(strings.ToLower(choice), "submit") &&
					!strings.Contains(strings.ToLower(choice), "deadline") &&
					!strings.Contains(strings.ToLower(choice), "turn") &&
					!strings.Contains(strings.ToLower(choice), "sheet") &&
					!strings.Contains(strings.ToLower(choice), "code") &&
					!strings.HasPrefix(strings.ToLower(choice), "f ") &&
					!strings.HasPrefix(strings.ToLower(choice), "f_") &&
					!strings.HasPrefix(strings.ToLower(choice), "h") &&
					!strings.HasPrefix(strings.ToLower(choice), "u") {

					// Check if this choice matches any of the actual location names (case-insensitive)
					var matchedLocation string
					var exists bool

					// First try exact match
					if matchedLocation, exists = locationNames[choice]; !exists {
						// Then try case-insensitive match
						if matchedLocation, exists = locationNames[strings.ToLower(choice)]; !exists {
							// Finally try case-insensitive match with original case
							for originalName := range locationNames {
								if strings.EqualFold(originalName, choice) {
									matchedLocation = originalName
									exists = true
									break
								}
							}
						}
					}

					if exists {
						// Convert to location ID format (lowercase, spaces to underscores)
						locationID := strings.ToLower(strings.ReplaceAll(matchedLocation, " ", "_"))

						// Only add if we haven't seen this location before
						if !seen[locationID] {
							choices = append(choices, locationID)
							seen[locationID] = true
						}
					}
				}
			}
		}
	}

	if len(choices) == 0 {
		return nil, fmt.Errorf("no valid location choices found in text")
	}

	return choices, nil
}
