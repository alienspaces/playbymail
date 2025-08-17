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
	FieldGameInstanceParameterParameterKey   string = "parameter_key"
	FieldGameInstanceParameterParameterValue string = "parameter_value"
	FieldGameInstanceParameterCreatedAt      string = "created_at"
	FieldGameInstanceParameterUpdatedAt      string = "updated_at"
	FieldGameInstanceParameterDeletedAt      string = "deleted_at"
)

// GameInstanceParameter represents runtime parameter values for specific game instances.
// These parameters override or extend the default values defined in game parameters.
type GameInstanceParameter struct {
	record.Record
	GameInstanceID string         `db:"game_instance_id"`
	ParameterKey   string         `db:"parameter_key"`
	ParameterValue sql.NullString `db:"parameter_value"`
}

// ToNamedArgs -
func (r *GameInstanceParameter) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameInstanceParameterGameInstanceID] = r.GameInstanceID
	args[FieldGameInstanceParameterParameterKey] = r.ParameterKey
	args[FieldGameInstanceParameterParameterValue] = r.ParameterValue
	return args
}
