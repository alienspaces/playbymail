package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameLocation
const (
	TableGameLocation string = "game_location"
)

const (
	FieldGameLocationID          string = "id"
	FieldGameLocationGameID      string = "game_id"
	FieldGameLocationName        string = "name"
	FieldGameLocationDescription string = "description"
)

type GameLocation struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *GameLocation) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameLocationGameID] = r.GameID
	args[FieldGameLocationName] = r.Name
	args[FieldGameLocationDescription] = r.Description
	return args
}
