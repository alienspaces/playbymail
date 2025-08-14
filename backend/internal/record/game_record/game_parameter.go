package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameParameter
const (
	TableGameParameter string = "game_parameter"
)

const (
	FieldGameParameterID           string = "id"
	FieldGameParameterGameType     string = "game_type"
	FieldGameParameterConfigKey    string = "config_key"
	FieldGameParameterDescription  string = "description"
	FieldGameParameterValueType    string = "value_type"
	FieldGameParameterDefaultValue string = "default_value"
	FieldGameParameterIsRequired   string = "is_required"
	FieldGameParameterIsGlobal     string = "is_global"
	FieldGameParameterCreatedAt    string = "created_at"
	FieldGameParameterUpdatedAt    string = "updated_at"
	FieldGameParameterDeletedAt    string = "deleted_at"
)

// GameParameter represents configuration parameters that are available for
// different game types. These define the structure and constraints for
// parameters that can be configured when creating games and game instances.
//
// Game parameter configurations work in a three-tier system:
// 1. Configuration (this record) - defines what parameters are available
// 2. Game parameters - set actual values when creating a game
// 3. Instance parameters - can override game values when creating instances
type GameParameter struct {
	record.Record

	// GameType specifies which game type this parameter applies to
	// (e.g., "adventure", "strategy", etc.)
	GameType string `db:"game_type"`

	// ConfigKey is the unique identifier for this parameter within the game type
	// (e.g., "character_lives", "max_players", "turn_duration")
	ConfigKey string `db:"config_key"`

	// Description provides a human-readable explanation of what this parameter
	// controls in the game
	Description sql.NullString `db:"description"`

	// ValueType defines the data type for this parameter's value
	// Valid types: "string", "integer", "boolean", "json"
	ValueType string `db:"value_type"`

	// DefaultValue is the fallback value used by the game engine when no
	// value is explicitly set. This is optional - the game engine may have
	// its own internal defaults
	DefaultValue sql.NullString `db:"default_value"`

	// IsRequired determines whether this parameter MUST be set in the game
	// studio when creating a game of this type. If true, the game cannot be
	// created without providing a value for this parameter
	IsRequired bool `db:"is_required"`

	// IsGlobal controls whether this parameter can be overridden at the
	// instance level:
	// - true: The parameter value is fixed at the game level and cannot be
	//   changed when creating instances (global/immutable setting)
	// - false: The parameter can be overridden with different values when
	//   creating individual game instances (instance-specific setting)
	IsGlobal bool `db:"is_global"`
}

// ToNamedArgs -
func (r *GameParameter) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameParameterGameType] = r.GameType
	args[FieldGameParameterConfigKey] = r.ConfigKey
	args[FieldGameParameterDescription] = r.Description
	args[FieldGameParameterValueType] = r.ValueType
	args[FieldGameParameterDefaultValue] = r.DefaultValue
	args[FieldGameParameterIsRequired] = r.IsRequired
	args[FieldGameParameterIsGlobal] = r.IsGlobal
	return args
}
