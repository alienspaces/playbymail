package game_record

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

const GameTypeAdventure = "adventure"

type Game struct {
	record.Record
	Name     string `db:"name"`
	GameType string `db:"game_type"`
}

func (r *Game) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameName] = r.Name
	args["game_type"] = r.GameType
	return args
}
