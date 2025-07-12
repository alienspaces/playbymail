package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableGameLocationInstance = "game_location_instance"

const (
	FieldGameLocationInstanceID             = "id"
	FieldGameLocationInstanceGameID         = "game_id"
	FieldGameLocationInstanceGameInstanceID = "game_instance_id"
	FieldGameLocationInstanceGameLocationID = "game_location_id"
)

type GameLocationInstance struct {
	record.Record
	GameID         string `db:"game_id"`
	GameInstanceID string `db:"game_instance_id"`
	GameLocationID string `db:"game_location_id"`
}

func (r *GameLocationInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameLocationInstanceGameID] = r.GameID
	args[FieldGameLocationInstanceGameInstanceID] = r.GameInstanceID
	args[FieldGameLocationInstanceGameLocationID] = r.GameLocationID
	return args
}
