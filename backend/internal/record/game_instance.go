package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableGameInstance = "game_instance"

const (
	FieldGameInstanceID     = "id"
	FieldGameInstanceGameID = "game_id"
)

type GameInstance struct {
	record.Record
	GameID string `db:"game_id"`
}

func (r *GameInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameInstanceGameID] = r.GameID
	return args
}
