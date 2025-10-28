package turn_sheet

import (
	"context"
	"encoding/json"
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
func (p *LocationChoiceProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte) ([]byte, error) {
	l = l.WithFunctionContext("LocationChoiceProcessor/GenerateTurnSheet")

	l.Info("generating location choice turn sheet")

	// Unmarshal sheet data
	var locationChoiceData LocationChoiceData
	if err := json.Unmarshal(sheetData, &locationChoiceData); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to parse sheet data: %w", err)
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

	return p.Generator.GeneratePDF(ctx, templatePath, &locationChoiceData)
}

// ScanTurnSheet scans a location choice turn sheet and extracts player choices
func (p *LocationChoiceProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, imageData []byte, sheetData []byte) ([]byte, error) {
	l = l.WithFunctionContext("LocationChoiceProcessor/ScanTurnSheet")

	l.Info("scanning location choice turn sheet")

	// Unmarshal sheet data to get location information
	var sheetDataMap map[string]any
	if err := json.Unmarshal(sheetData, &sheetDataMap); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to parse sheet data: %w", err)
	}

	// Extract text from image using base processor
	text, err := p.ExtractTextFromImage(ctx, imageData)
	if err != nil {
		l.Warn("failed to extract text from image >%v<", err)
		return nil, fmt.Errorf("text extraction failed: %w", err)
	}

	// Parse location choices from extracted text using the sheet data
	l.Info("full OCR text extracted: >%s<", text)
	choices, err := p.parseLocationChoicesWithSheetData(l, text, sheetDataMap)
	if err != nil {
		l.Warn("failed to parse location choices >%v<", err)
		return nil, fmt.Errorf("location choice parsing failed: %w", err)
	}

	l.Info("extracted %d location choices", len(choices))

	// Marshal scan results to JSON
	scanData := LocationChoiceScanData{
		Choices: choices,
	}
	
	scanResultBytes, err := json.Marshal(scanData)
	if err != nil {
		l.Warn("failed to marshal scan results >%v<", err)
		return nil, fmt.Errorf("failed to marshal scan results: %w", err)
	}

	return scanResultBytes, nil
}

// ParseLocationChoicesWithSheetData parses location choices using the actual sheet data
func (p *LocationChoiceProcessor) parseLocationChoicesWithSheetData(l logger.Logger, text string, sheetData map[string]any) ([]string, error) {
	l = l.WithFunctionContext("LocationChoiceProcessor/parseLocationChoicesWithSheetData")

	l.Info("parsing location choices with sheet data")

	// Extract location data from sheet data (try both "location_options" and "locations" for backwards compatibility)
	var locations []any

	if locationOptions, exists := sheetData["location_options"].([]any); exists {
		locations = locationOptions
	} else if locationsFromLegacy, exists := sheetData["locations"].([]any); exists {
		locations = locationsFromLegacy
	} else {
		return nil, fmt.Errorf("no location options found in sheet data")
	}

	if len(locations) == 0 {
		return nil, fmt.Errorf("no location options found in sheet data")
	}

	// Create a map of location names for validation
	// Key: location name (as shown on the turn sheet), Value: location ID
	locationNames := make(map[string]string)
	for _, loc := range locations {
		if locMap, ok := loc.(map[string]any); ok {
			// Try location_link_name first (from LocationOption), then name (for backward compatibility)
			var name string
			var locationID string

			if linkName, exists := locMap["location_link_name"].(string); exists {
				name = linkName
			} else if nameOnly, exists := locMap["name"].(string); exists {
				name = nameOnly
			}

			// Get location_id if available (from LocationOption)
			if id, exists := locMap["location_id"].(string); exists {
				locationID = id
			} else if name != "" {
				// Fallback: convert name to location ID format
				locationID = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
			}

			if name != "" {
				locationNames[name] = locationID
				// Also add lowercase version for matching
				locationNames[strings.ToLower(name)] = locationID
			}
		}
	}

	l.Info("location names: %v", locationNames)

	// All possible checkbox patterns (both selected and unselected)
	allPatterns := []string{
		// Selected patterns
		`Q/([A-Za-z][A-Za-z\s:]+?)(?:\n|$)`,        // OCR reads selected checkbox as Q/ (no space, includes colon)
		`☑\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,       // Standard selected checkbox
		`\[X\]\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,   // X in brackets
		`X\s+([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,       // Just an X mark
		`\\(X\\)\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`, // X in parentheses
		`vs\s+([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,      // OCR reads X as "vs"
		`✓\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,       // Checkmark
		// Unselected patterns
		`\(O\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`, // OCR reads unselected checkbox as (O
		`O/\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,  // OCR reads unselected checkbox as O/
		`Sf\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,  // OCR misreads ☑ as Sf
		`S\s*([A-Za-z][A-Za-z\s]+?)(?:\n|$)`,   // OCR misreads ☑ as S
	}

	// Collect all matches with their pattern index
	type matchResult struct {
		patternIdx int
		location   string
	}
	var allMatches []matchResult

	for i, pattern := range allPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(text, -1)
		l.Info("pattern %d (%s) found %d matches", i, pattern, len(matches))

		for _, match := range matches {
			if len(match) > 1 {
				location := strings.TrimSpace(match[1])
				// Clean up OCR artifacts like colons at the end
				location = strings.TrimSuffix(location, ":")
				location = strings.TrimSpace(location)

				if len(location) > 0 {
					allMatches = append(allMatches, matchResult{
						patternIdx: i,
						location:   location,
					})
					l.Info("found potential location: >%s< with pattern %d", location, i)
				}
			}
		}
	}

	if len(allMatches) == 0 {
		return nil, fmt.Errorf("no valid location choices found in text")
	}

	// If we have N location options and found N-1 matches with one pattern and 1 match with another,
	// the pattern with only 1 match is likely the selection
	// Count pattern frequencies
	patternCounts := make(map[int]int)
	for _, match := range allMatches {
		patternCounts[match.patternIdx]++
	}

	// Find the minority pattern (the one that appears least often)
	minorityPatternIdx := -1
	minCount := len(allMatches) // Start with max possible count
	for patternIdx, count := range patternCounts {
		if count < minCount && count > 0 {
			minCount = count
			minorityPatternIdx = patternIdx
		}
	}

	l.Info("pattern counts: %v, minority pattern: %d (count: %d)", patternCounts, minorityPatternIdx, minCount)

	// Extract choices from matches using the minority pattern
	var choices []string
	seen := make(map[string]bool)

	for _, match := range allMatches {
		if match.patternIdx == minorityPatternIdx {
			// Try to match this location name to our known locations
			var locationID string
			var exists bool

			// First try exact match
			if locationID, exists = locationNames[match.location]; !exists {
				// Then try case-insensitive match
				if locationID, exists = locationNames[strings.ToLower(match.location)]; !exists {
					// Finally try case-insensitive match with original case
					for originalName, id := range locationNames {
						if strings.EqualFold(originalName, match.location) {
							locationID = id
							exists = true
							break
						}
					}
				}
			}

			if exists {
				// Only add if we haven't seen this location before
				if !seen[locationID] {
					choices = append(choices, locationID)
					seen[locationID] = true
				}
			}
		}
	}

	if len(choices) == 0 {
		return nil, fmt.Errorf("no valid location choices found in text")
	}

	return choices, nil
}
