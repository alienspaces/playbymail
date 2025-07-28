package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameInstanceConfigurationSchema
const (
	TableGameInstanceConfigurationSchema string = "game_instance_configuration_schema"
)

const (
	FieldGameInstanceConfigurationSchemaID              string = "id"
	FieldGameInstanceConfigurationSchemaGameType        string = "game_type"
	FieldGameInstanceConfigurationSchemaConfigKey       string = "config_key"
	FieldGameInstanceConfigurationSchemaValueType       string = "value_type"
	FieldGameInstanceConfigurationSchemaDefaultValue    string = "default_value"
	FieldGameInstanceConfigurationSchemaIsRequired      string = "is_required"
	FieldGameInstanceConfigurationSchemaDescription     string = "description"
	FieldGameInstanceConfigurationSchemaUIHint          string = "ui_hint"
	FieldGameInstanceConfigurationSchemaValidationRules string = "validation_rules"
	FieldGameInstanceConfigurationSchemaCreatedAt       string = "created_at"
	FieldGameInstanceConfigurationSchemaUpdatedAt       string = "updated_at"
	FieldGameInstanceConfigurationSchemaDeletedAt       string = "deleted_at"
)

type GameInstanceConfigurationSchema struct {
	record.Record
	GameType        string         `db:"game_type"`
	ConfigKey       string         `db:"config_key"`
	ValueType       string         `db:"value_type"`
	DefaultValue    sql.NullString `db:"default_value"`
	IsRequired      bool           `db:"is_required"`
	Description     sql.NullString `db:"description"`
	UIHint          sql.NullString `db:"ui_hint"`
	ValidationRules sql.NullString `db:"validation_rules"`
}

func (r *GameInstanceConfigurationSchema) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameInstanceConfigurationSchemaGameType] = r.GameType
	args[FieldGameInstanceConfigurationSchemaConfigKey] = r.ConfigKey
	args[FieldGameInstanceConfigurationSchemaValueType] = r.ValueType
	args[FieldGameInstanceConfigurationSchemaDefaultValue] = r.DefaultValue
	args[FieldGameInstanceConfigurationSchemaIsRequired] = r.IsRequired
	args[FieldGameInstanceConfigurationSchemaDescription] = r.Description
	args[FieldGameInstanceConfigurationSchemaUIHint] = r.UIHint
	args[FieldGameInstanceConfigurationSchemaValidationRules] = r.ValidationRules
	return args
} 