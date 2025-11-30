package turn_sheet

import "time"

type DocumentFormat string

const (
	DocumentFormatPDF  DocumentFormat = "pdf"
	DocumentFormatHTML DocumentFormat = "html"
)

// Turn sheet types
const (
	TurnSheetTypeLocationChoice string = "location_choice"
)

// TurnSheetData represents the data for a turn sheet
//
// All turn sheet types use this same data structure
type TurnSheetTemplateData struct {
	// Game data
	GameName *string `json:"game_name"`
	GameType *string `json:"game_type"`
	// Display data
	TurnSheetTitle       *string `json:"turn_sheet_title"`
	TurnSheetDescription *string `json:"turn_sheet_description"`

	// Game instance data
	TurnNumber *int `json:"turn_number"`

	// Account data
	AccountName *string `json:"account_name"`

	// Background image (single image covering the page)
	BackgroundImage *string `json:"background_image"`

	// Turn sheet
	TurnSheetInstructions *string    `json:"turn_sheet_instructions"`
	TurnSheetDeadline     *time.Time `json:"turn_sheet_deadline"`
	TurnSheetCode         *string    `json:"turn_sheet_code"`
}
