package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableGameCharacterInstance string = "game_character_instance"
)

const (
	FieldGameCharacterInstanceID                     string = "id"
	FieldGameCharacterInstanceGameID                 string = "game_id"
	FieldGameCharacterInstanceGameInstanceID         string = "game_instance_id"
	FieldGameCharacterInstanceGameCharacterID        string = "game_character_id"
	FieldGameCharacterInstanceGameLocationInstanceID string = "game_location_instance_id"
	FieldGameCharacterInstanceHealth                 string = "health"
	FieldGameCharacterInstanceCreatedAt              string = "created_at"
	FieldGameCharacterInstanceUpdatedAt              string = "updated_at"
	FieldGameCharacterInstanceDeletedAt              string = "deleted_at"
)

type GameCharacterInstance struct {
	record.Record
	GameID                 string `db:"game_id"`
	GameInstanceID         string `db:"game_instance_id"`
	GameCharacterID        string `db:"game_character_id"`
	GameLocationInstanceID string `db:"game_location_instance_id"`
	Health                 int    `db:"health"`
}

func (r *GameCharacterInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameCharacterInstanceGameID] = r.GameID
	args[FieldGameCharacterInstanceGameInstanceID] = r.GameInstanceID
	args[FieldGameCharacterInstanceGameCharacterID] = r.GameCharacterID
	args[FieldGameCharacterInstanceGameLocationInstanceID] = r.GameLocationInstanceID
	args[FieldGameCharacterInstanceHealth] = r.Health
	return args
}
