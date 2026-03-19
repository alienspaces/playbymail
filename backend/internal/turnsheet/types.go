package turnsheet

import "time"

type DocumentFormat string

const (
	DocumentFormatPDF  DocumentFormat = "pdf"
	DocumentFormatHTML DocumentFormat = "html"
)

// TurnEvent categories
const (
	TurnEventCategoryCombat    = "combat"
	TurnEventCategoryInventory = "inventory"
	TurnEventCategoryMovement  = "movement"
	TurnEventCategoryWorld     = "world"
	TurnEventCategoryFlee      = "flee"
	// flee_context is an internal category used to pass flee state between processors
	TurnEventCategoryFleeContext = "flee_context"
)

// TurnEvent icons (unicode)
const (
	TurnEventIconCombat    = "⚔️"
	TurnEventIconInventory = "🎒"
	TurnEventIconMovement  = "👣"
	TurnEventIconWorld     = "🌍"
	TurnEventIconFlee      = "💨"
	TurnEventIconDeath     = "💀"
	TurnEventIconHeal      = "💚"
)

// TurnEvent represents a narrative event that occurred during turn processing.
// Events are stored in character_instance.last_turn_events and displayed on the next turn's sheet.
type TurnEvent struct {
	Category string `json:"category"` // "combat", "inventory", "movement", "world", "flee", "flee_context"
	Icon     string `json:"icon"`     // unicode emoji
	Message  string `json:"message"`  // human-readable narrative
}

// FleeContext is stored as a TurnEvent with category "flee_context" to pass flee state
// from the location choice processor to the encounter processor for the next turn.
type FleeContext struct {
	PreviousLocationID   string                `json:"previous_location_id"`
	PreviousLocationName string                `json:"previous_location_name"`
	Creatures            []FleeContextCreature `json:"creatures"`
}

// FleeContextCreature stores minimal creature info for flee narrative display.
type FleeContextCreature struct {
	InstanceID        string `json:"instance_id"`
	Name              string `json:"name"`
	AttackMethod      string `json:"attack_method"`
	AttackDescription string `json:"attack_description"`
	DamageDealt       int    `json:"damage_dealt"`
}

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

	// Narrative events from the previous turn, displayed in the "What Happened" section
	TurnEvents []TurnEvent `json:"turn_events,omitempty"`
}
