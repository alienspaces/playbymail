package game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameConfiguration
const (
	TableGameConfiguration string = "game_configuration"
)

const (
	FieldGameConfigurationID              string = "id"
	FieldGameConfigurationGameType        string = "game_type"
	FieldGameConfigurationConfigKey       string = "config_key"
	FieldGameConfigurationValueType       string = "value_type"
	FieldGameConfigurationDefaultValue    string = "default_value"
	FieldGameConfigurationIsRequired      string = "is_required"
	FieldGameConfigurationDescription     string = "description"
	FieldGameConfigurationUIHint          string = "ui_hint"
	FieldGameConfigurationValidationRules string = "validation_rules"
	FieldGameConfigurationCreatedAt       string = "created_at"
	FieldGameConfigurationUpdatedAt       string = "updated_at"
	FieldGameConfigurationDeletedAt       string = "deleted_at"
)

// GameConfiguration -
type GameConfiguration struct {
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

// ToNamedArgs -
func (r *GameConfiguration) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameConfigurationGameType] = r.GameType
	args[FieldGameConfigurationConfigKey] = r.ConfigKey
	args[FieldGameConfigurationValueType] = r.ValueType
	args[FieldGameConfigurationDefaultValue] = r.DefaultValue
	args[FieldGameConfigurationIsRequired] = r.IsRequired
	args[FieldGameConfigurationDescription] = r.Description
	args[FieldGameConfigurationUIHint] = r.UIHint
	args[FieldGameConfigurationValidationRules] = r.ValidationRules
	return args
}
