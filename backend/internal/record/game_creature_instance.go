package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	FieldGameCreatureInstanceID                 = "id"
	FieldGameCreatureInstanceGameID             = "game_id"
	FieldGameCreatureInstanceGameCreatureID     = "game_creature_id"
	FieldGameCreatureInstanceGameInstanceID     = "game_instance_id"
	FieldGameCreatureInstanceGameLocationInstID = "game_location_instance_id"
	FieldGameCreatureInstanceIsAlive            = "is_alive"
	FieldGameCreatureInstanceCreatedAt          = "created_at"
	FieldGameCreatureInstanceUpdatedAt          = "updated_at"
	FieldGameCreatureInstanceDeletedAt          = "deleted_at"
)

const TableGameCreatureInstance = "game_creature_instance"

type GameCreatureInstance struct {
	record.Record
	GameID                 string `db:"game_id"`
	GameCreatureID         string `db:"game_creature_id"`
	GameInstanceID         string `db:"game_instance_id"`
	GameLocationInstanceID string `db:"game_location_instance_id"`
	IsAlive                bool   `db:"is_alive"`
}

func (r *GameCreatureInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameCreatureInstanceGameID] = r.GameID
	args[FieldGameCreatureInstanceGameCreatureID] = r.GameCreatureID
	args[FieldGameCreatureInstanceGameInstanceID] = r.GameInstanceID
	args[FieldGameCreatureInstanceGameLocationInstID] = r.GameLocationInstanceID
	args[FieldGameCreatureInstanceIsAlive] = r.IsAlive
	return args
}
