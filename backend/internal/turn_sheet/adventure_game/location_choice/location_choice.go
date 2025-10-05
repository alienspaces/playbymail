package location_choice

import (
	"time"

	"gitlab.com/alienspaces/playbymail/internal/turn_sheet"
)

// LocationOption represents a location choice option for the player
type LocationOption struct {
	LocationID              string
	LocationLinkName        string
	LocationLinkDescription string
}

type LocationChoiceData struct {
	turn_sheet.TurnSheetData

	// Current location information
	LocationName        string
	LocationDescription string

	// Available location options
	LocationOptions []LocationOption

	// Turn deadline information
	TurnDeadline *time.Time
}
