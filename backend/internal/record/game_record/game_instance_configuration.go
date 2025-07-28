package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameInstanceConfiguration
const (
	TableGameInstanceConfiguration string = "game_instance_configuration"
)

const (
	FieldGameInstanceConfigurationID             string = "id"
	FieldGameInstanceConfigurationGameInstanceID string = "game_instance_id"
	FieldGameInstanceConfigurationConfigKey      string = "config_key"
	FieldGameInstanceConfigurationValueType      string = "value_type"
	FieldGameInstanceConfigurationStringValue    string = "string_value"
	FieldGameInstanceConfigurationIntegerValue   string = "integer_value"
	FieldGameInstanceConfigurationBooleanValue   string = "boolean_value"
	FieldGameInstanceConfigurationJSONValue      string = "json_value"
	FieldGameInstanceConfigurationCreatedAt      string = "created_at"
	FieldGameInstanceConfigurationUpdatedAt      string = "updated_at"
	FieldGameInstanceConfigurationDeletedAt      string = "deleted_at"
)

// GameInstanceConfiguration -
type GameInstanceConfiguration struct {
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
func (r *GameInstanceConfiguration) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameInstanceConfigurationGameInstanceID] = r.GameInstanceID
	args[FieldGameInstanceConfigurationConfigKey] = r.ConfigKey
	args[FieldGameInstanceConfigurationValueType] = r.ValueType
	args[FieldGameInstanceConfigurationStringValue] = r.StringValue
	args[FieldGameInstanceConfigurationIntegerValue] = r.IntegerValue
	args[FieldGameInstanceConfigurationBooleanValue] = r.BooleanValue
	args[FieldGameInstanceConfigurationJSONValue] = r.JSONValue
	return args
}
