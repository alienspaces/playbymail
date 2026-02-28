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

type LocationChoiceData struct {
	TurnSheetTemplateData

	// Current location information
	LocationName        string `json:"location_name"`
	LocationDescription string `json:"location_description"`

	// Available location options
	LocationOptions []LocationOption `json:"location_options"`
}

// LocationOption represents a location choice option for the player
type LocationOption struct {
	LocationID              string `json:"location_id"`
	LocationLinkName        string `json:"location_link_name"`
	LocationLinkDescription string `json:"location_link_description"`
}

// LocationChoiceScanData represents the scanned data from a location choice turn sheet
type LocationChoiceScanData struct {
	Choices []string `json:"choices"`
}

const defaultLocationChoiceInstructions = "Select your next location and return this form by the deadline to continue your adventure."

// DefaultLocationChoiceInstructions returns the default instruction text for location choice turn sheets.
func DefaultLocationChoiceInstructions() string {
	return defaultLocationChoiceInstructions
}

// LocationChoiceProcessor implements the DocumentProcessor interface for location choice turn sheets
type LocationChoiceProcessor struct {
	*BaseProcessor
}

// NewLocationChoiceProcessor creates a new location choice processor
func NewLocationChoiceProcessor(l logger.Logger, cfg config.Config) (*LocationChoiceProcessor, error) {
	baseProcessor, err := NewBaseProcessor(l, cfg)
	if err != nil {
		return nil, err
	}
	return &LocationChoiceProcessor{
		BaseProcessor: baseProcessor,
	}, nil
}

// GeneratePreviewData generates dummy data for a location choice turn sheet preview
func (p *LocationChoiceProcessor) GeneratePreviewData(ctx context.Context, l logger.Logger, gameRec *game_record.Game, backgroundImage *string) ([]byte, error) {
	l = l.WithFunctionContext("LocationChoiceProcessor/GeneratePreviewData")

	l.Info("generating location choice preview data")

	// Generate join turn sheet code for this game
	turnSheetCode, err := turnsheetutil.GeneratePlayGameTurnSheetCode(record.NewRecordID())
	if err != nil {
		l.Warn("failed to generate play game turn sheet code >%v<", err)
		return nil, fmt.Errorf("failed to generate turn sheet code: %w", err)
	}

	// Create dummy location choice data
	locationChoiceData := LocationChoiceData{
		TurnSheetTemplateData: TurnSheetTemplateData{
			GameName:              convert.Ptr(gameRec.Name),
			GameType:              convert.Ptr(gameRec.GameType),
			TurnNumber:            convert.Ptr(5),
			AccountName:           nil, // Not needed for preview
			TurnSheetTitle:        convert.Ptr("Location Choice"),
			TurnSheetDescription:  convert.Ptr("Select your next destination."),
			TurnSheetInstructions: convert.Ptr(DefaultLocationChoiceInstructions()),
			TurnSheetDeadline:     nil,
			TurnSheetCode:         convert.Ptr(turnSheetCode),
		},
		LocationName:        "Dark Forest",
		LocationDescription: "You are standing at the edge of a dark, ominous forest. The trees are tall and gnarled, blocking out most of the sunlight. You hear strange noises coming from within.",
		LocationOptions: []LocationOption{
			{
				LocationID:              "loc-1",
				LocationLinkName:        "Enter the Forest",
				LocationLinkDescription: "Brave the dangers of the forest.",
			},
			{
				LocationID:              "loc-2",
				LocationLinkName:        "Follow the River",
				LocationLinkDescription: "Walk along the river bank towards the mountains.",
			},
			{
				LocationID:              "loc-3",
				LocationLinkName:        "Return to Village",
				LocationLinkDescription: "Head back to the safety of the village.",
			},
		},
	}

	// Set background image if provided
	if backgroundImage != nil && *backgroundImage != "" {
		locationChoiceData.BackgroundImage = backgroundImage
	}

	return json.Marshal(locationChoiceData)
}

// GenerateTurnSheet generates a location choice turn sheet document
func (p *LocationChoiceProcessor) GenerateTurnSheet(ctx context.Context, l logger.Logger, format DocumentFormat, sheetData []byte) ([]byte, error) {
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

	if locationChoiceData.TurnSheetInstructions == nil || strings.TrimSpace(*locationChoiceData.TurnSheetInstructions) == "" {
		instruction := defaultLocationChoiceInstructions
		locationChoiceData.TurnSheetInstructions = &instruction
	}

	if locationChoiceData.TurnSheetTitle == nil || strings.TrimSpace(*locationChoiceData.TurnSheetTitle) == "" {
		if locationChoiceData.LocationName != "" {
			title := locationChoiceData.LocationName
			locationChoiceData.TurnSheetTitle = &title
		}
	}

	if locationChoiceData.TurnSheetDescription == nil || strings.TrimSpace(*locationChoiceData.TurnSheetDescription) == "" {
		if locationChoiceData.LocationDescription != "" {
			desc := locationChoiceData.LocationDescription
			locationChoiceData.TurnSheetDescription = &desc
		}
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

	return p.GenerateDocument(ctx, format, locationChoiceTemplatePath, &locationChoiceData)
}

const locationChoiceTemplatePath = "turnsheet/adventure_game_location_choice.template"

// ScanTurnSheet scans a location choice turn sheet and extracts player choices using hosted OCR
func (p *LocationChoiceProcessor) ScanTurnSheet(ctx context.Context, l logger.Logger, sheetData []byte, imageData []byte) ([]byte, error) {
	l = l.WithFunctionContext("LocationChoiceProcessor/ScanTurnSheet")

	l.Info("scanning location choice turn sheet")

	if len(imageData) == 0 {
		l.Warn("empty image data provided")
		return nil, fmt.Errorf("empty image data provided")
	}

	var locationChoiceData LocationChoiceData
	if err := json.Unmarshal(sheetData, &locationChoiceData); err != nil {
		l.Warn("failed to unmarshal sheet data >%v<", err)
		return nil, fmt.Errorf("failed to parse sheet data: %w", err)
	}

	if len(locationChoiceData.LocationOptions) == 0 {
		l.Warn("no location options supplied in sheet data")
		return nil, fmt.Errorf("no location options supplied in sheet data")
	}

	templateImage, err := p.renderTemplatePreview(ctx, locationChoiceTemplatePath, &locationChoiceData)
	if err != nil {
		l.Warn("failed to generate template preview >%v<", err)
		return nil, fmt.Errorf("failed to generate template preview: %w", err)
	}
	if len(templateImage) == 0 {
		l.Warn("template preview generation returned empty image")
		return nil, fmt.Errorf("template preview generation returned empty image")
	}

	expected := map[string]any{
		"choices": []string{},
	}

	req := scanner.StructuredScanRequest{
		Instructions:       buildLocationChoiceInstructions(),
		AdditionalContext:  buildLocationChoiceContext(&locationChoiceData),
		TemplateImage:      templateImage,
		TemplateImageMIME:  "image/png",
		FilledImage:        imageData,
		ExpectedJSONSchema: expected,
	}

	raw, err := p.Scanner.ExtractStructuredData(ctx, req)
	if err != nil {
		l.Warn("structured extraction failed >%v<", err)
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	var scanData LocationChoiceScanData
	if err := json.Unmarshal(raw, &scanData); err != nil {
		return nil, fmt.Errorf("failed to decode structured location choices: %w", err)
	}

	if err := validateLocationChoices(&locationChoiceData, &scanData); err != nil {
		return nil, err
	}

	return json.Marshal(scanData)
}

// These are the instructions provided to the AI driven OCR service.
func buildLocationChoiceInstructions() string {
	return `Compare the blank template image with the completed turn sheet.
Determine which checkbox/circle is marked by the player.
Respond with JSON containing a "choices" array of location_id values (strings).
Use the provided reference list to map the printed location names to their ids.
If no boxes are marked, return an empty array.`
}

// These are the additional context provided to the AI driven OCR service.
func buildLocationChoiceContext(data *LocationChoiceData) []string {
	var ctx []string
	if data != nil {
		for _, option := range data.LocationOptions {
			ctx = append(ctx, fmt.Sprintf("location_id=%s label=%s description=%s",
				option.LocationID,
				strings.TrimSpace(option.LocationLinkName),
				strings.TrimSpace(option.LocationLinkDescription),
			))
		}
	}
	return ctx
}

func validateLocationChoices(sheetData *LocationChoiceData, scanData *LocationChoiceScanData) error {
	if scanData == nil {
		return fmt.Errorf("no scan data provided")
	}

	validIDs := make(map[string]bool)
	for _, opt := range sheetData.LocationOptions {
		if opt.LocationID != "" {
			validIDs[opt.LocationID] = true
		}
	}

	for _, choice := range scanData.Choices {
		if !validIDs[choice] {
			return fmt.Errorf("invalid location_id returned: %s", choice)
		}
	}

	return nil
}
