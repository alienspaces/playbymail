package turn_sheet

import (
	"time"
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

	// Game instance data
	TurnNumber *int `json:"turn_number"`

	// Account data
	AccountName *string `json:"account_name"`

	// Background images
	BackgroundTop    *string `json:"background_top"`
	BackgroundMiddle *string `json:"background_middle"`
	BackgroundBottom *string `json:"background_bottom"`

	// Turn sheet
	TurnSheetDeadline *time.Time `json:"turn_sheet_deadline"`
	TurnSheetCode     *string    `json:"turn_sheet_code"`
}
