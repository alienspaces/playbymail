package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameInstanceParameter
const (
	TableGameInstanceParameter string = "game_instance_parameter"
)

const (
	FieldGameInstanceParameterID             string = "id"
	FieldGameInstanceParameterGameInstanceID string = "game_instance_id"
	FieldGameInstanceParameterConfigKey      string = "config_key"
	FieldGameInstanceParameterValueType      string = "value_type"
	FieldGameInstanceParameterStringValue    string = "string_value"
	FieldGameInstanceParameterIntegerValue   string = "integer_value"
	FieldGameInstanceParameterBooleanValue   string = "boolean_value"
	FieldGameInstanceParameterJSONValue      string = "json_value"
	FieldGameInstanceParameterCreatedAt      string = "created_at"
	FieldGameInstanceParameterUpdatedAt      string = "updated_at"
	FieldGameInstanceParameterDeletedAt      string = "deleted_at"
)

// GameInstanceParameter -
type GameInstanceParameter struct {
	record.Record
	GameInstanceID string         `db:"game_instance_id"`
	ConfigKey      string         `db:"config_key"`
	ValueType      string         `db:"value_type"`
	StringValue    sql.NullString `db:"string_value"`
	IntegerValue   sql.NullInt32  `db:"integer_value"`
	BooleanValue   sql.NullBool   `db:"boolean_value"`
	JSONValue      sql.NullString `db:"json_value"`
}

// ToNamedArgs -
func (r *GameInstanceParameter) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameInstanceParameterGameInstanceID] = r.GameInstanceID
	args[FieldGameInstanceParameterConfigKey] = r.ConfigKey
	args[FieldGameInstanceParameterValueType] = r.ValueType
	args[FieldGameInstanceParameterStringValue] = r.StringValue
	args[FieldGameInstanceParameterIntegerValue] = r.IntegerValue
	args[FieldGameInstanceParameterBooleanValue] = r.BooleanValue
	args[FieldGameInstanceParameterJSONValue] = r.JSONValue
	return args
}
