package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// Game
const (
	TableGame string = "game"
)

const (
	FieldGameID   string = "id"
	FieldGameName string = "name"
)

type Game struct {
	record.Record
	Name string `db:"name"`
}

func (r *Game) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args["name"] = r.Name
	return args
}
