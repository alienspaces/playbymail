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
	GameName *string
	GameType *string

	// Game instance data
	TurnNumber *int

	// Account data
	AccountName *string

	// Background images
	BackgroundTop    *string
	BackgroundMiddle *string
	BackgroundBottom *string

	// Turn sheet
	TurnSheetDeadline *time.Time
	TurnSheetCode     *string
}
