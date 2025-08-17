package game_record

// GameParameter represents configuration parameters that are available for
// different game types. These define the structure and constraints for
// parameters that can be configured when creating games and game instances.
//
// Game parameters work in a two-tier system:
// 1. Game parameters - parameters that are available per game type
// 2. Instance parameters - parameter values that are set per game instance
//
// Game parameters are used to define the structure and constraints for
// parameters that can be configured when creating games.
//
// Instance parameters are used to set the parameter values for a specific
// game instance.
type GameParameter struct {
	// GameType specifies which game type this parameter applies to
	// (e.g., "adventure", "strategy", etc.)
	GameType string

	// ConfigKey is the unique identifier for this parameter within the game type
	// (e.g., "character_lives", "max_players", "turn_duration")
	ConfigKey string

	// Description provides a human-readable explanation of what this parameter
	// controls in the game
	Description string

	// ValueType defines the data type for this parameter's value
	// Valid types: "string", "integer", "boolean"
	ValueType string

	// DefaultValue is the fallback value used by the game engine when no
	// value is explicitly set. This is optional - the game engine may have
	// its own internal defaults
	DefaultValue string
}
