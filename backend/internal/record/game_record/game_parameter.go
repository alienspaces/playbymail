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

// GameParameter -
type GameParameter struct {
	record.Record
	GameType     string         `db:"game_type"`
	ConfigKey    string         `db:"config_key"`
	Description  sql.NullString `db:"description"`
	ValueType    string         `db:"value_type"`
	DefaultValue sql.NullString `db:"default_value"`
	IsRequired   bool           `db:"is_required"`
	IsGlobal     bool           `db:"is_global"`
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
