package adventuregamegenerator

import (
	"context"

	"gitlab.com/alienspaces/playbymail/core/log"
	"gitlab.com/alienspaces/playbymail/internal/domain"
	"gitlab.com/alienspaces/playbymail/internal/generator"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// AdventureGameLocationGenerator generates location choice turn sheets for adventure games
type AdventureGameLocationGenerator struct {
	Logger       *log.Log
	Domain       *domain.Domain
	PDFGenerator *generator.PDFGenerator
}

// AdventureGameLocationTurnSheetData represents the data structure for adventure game location choice turn sheets
type AdventureGameLocationTurnSheetData struct {
	// Current location information
	CurrentLocationRec        *adventure_game_record.AdventureGameLocation       `json:"current_location_rec"`
	AvailableLocationLinkRecs []*adventure_game_record.AdventureGameLocationLink `json:"available_location_links"`

	// Turn information
	NextTurnDeadline string `json:"next_turn_deadline,omitempty"`
}

// NewAdventureGameLocationGenerator creates a new adventure game location generator
func NewAdventureGameLocationGenerator(l *log.Log, d *domain.Domain, templateDir, outputDir string) *AdventureGameLocationGenerator {
	pdfGen := generator.NewPDFGenerator(l, templateDir, outputDir)
	return &AdventureGameLocationGenerator{
		Logger:       l,
		Domain:       d,
		PDFGenerator: pdfGen,
	}
}

// GenerateLocationChoiceTurnSheet generates a location choice turn sheet PDF
func (g *AdventureGameLocationGenerator) GenerateLocationChoiceTurnSheet(ctx context.Context, data generator.TemplateData, turnSheetCode string) ([]byte, error) {

	_, ok := data.TurnSheetData.(*AdventureGameLocationTurnSheetData)
	if !ok {
		return nil, domain.Internal("invalid turn sheet data")
	}

	// Add the code to the template data
	data.TurnSheetCode = turnSheetCode

	templatePath := "templates/location_choice.template"
	return g.PDFGenerator.GeneratePDF(ctx, templatePath, data)
}

// GenerateLocationChoiceTurnSheetToFile generates and saves a location choice turn sheet PDF
func (g *AdventureGameLocationGenerator) GenerateLocationChoiceTurnSheetToFile(ctx context.Context, data generator.TemplateData, turnSheetCode, filename string) error {

	_, ok := data.TurnSheetData.(*AdventureGameLocationTurnSheetData)
	if !ok {
		return domain.Internal("invalid turn sheet data")
	}

	// Add the code to the template data
	data.TurnSheetCode = turnSheetCode

	templatePath := "templates/location_choice.template"
	return g.PDFGenerator.GeneratePDFToFile(ctx, templatePath, data, filename)
}

// GenerateLocationChoiceTurnSheetHTML generates HTML for a location choice turn sheet
func (g *AdventureGameLocationGenerator) GenerateLocationChoiceTurnSheetHTML(ctx context.Context, data generator.TemplateData, turnSheetCode string) (string, error) {

	_, ok := data.TurnSheetData.(*AdventureGameLocationTurnSheetData)
	if !ok {
		return "", domain.Internal("invalid turn sheet data")
	}

	// Add the code to the template data
	data.TurnSheetCode = turnSheetCode

	templatePath := "templates/location_choice.template"
	return g.PDFGenerator.GenerateHTML(ctx, templatePath, data)
}
