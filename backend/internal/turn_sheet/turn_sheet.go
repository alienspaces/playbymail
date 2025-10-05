package turn_sheet

type TurnSheetData struct {
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

	// Dynamic content sections
	Header map[string]string
	Footer map[string]string

	// Turn sheet identification
	TurnSheetCode *string
}
