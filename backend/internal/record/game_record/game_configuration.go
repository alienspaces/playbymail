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
	FieldGameConfigurationID           string = "id"
	FieldGameConfigurationGameType     string = "game_type"
	FieldGameConfigurationConfigKey    string = "config_key"
	FieldGameConfigurationDescription  string = "description"
	FieldGameConfigurationValueType    string = "value_type"
	FieldGameConfigurationDefaultValue string = "default_value"
	FieldGameConfigurationIsRequired   string = "is_required"
	FieldGameConfigurationIsGlobal     string = "is_global"
	FieldGameConfigurationCreatedAt    string = "created_at"
	FieldGameConfigurationUpdatedAt    string = "updated_at"
	FieldGameConfigurationDeletedAt    string = "deleted_at"
)

// GameConfiguration -
type GameConfiguration struct {
	record.Record
	GameType     string         `db:"game_type"`
	ConfigKey    string         `db:"config_key"`
	Description  string         `db:"description"`
	ValueType    string         `db:"value_type"`
	DefaultValue sql.NullString `db:"default_value"`
	IsRequired   bool           `db:"is_required"`
	IsGlobal     bool           `db:"is_global"`
}

// ToNamedArgs -
func (r *GameConfiguration) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameConfigurationGameType] = r.GameType
	args[FieldGameConfigurationConfigKey] = r.ConfigKey
	args[FieldGameConfigurationDescription] = r.Description
	args[FieldGameConfigurationValueType] = r.ValueType
	args[FieldGameConfigurationDefaultValue] = r.DefaultValue
	args[FieldGameConfigurationIsRequired] = r.IsRequired
	args[FieldGameConfigurationIsGlobal] = r.IsGlobal
	return args
}
