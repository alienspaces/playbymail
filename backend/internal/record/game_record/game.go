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
	FieldGameID                string = "id"
	FieldGameAccountID         string = "account_id"
	FieldGameName              string = "name"
	FieldGameDescription       string = "description"
	FieldGameType              string = "game_type"
	FieldGameTurnDurationHours string = "turn_duration_hours"
)

const GameTypeAdventure = "adventure"

type Game struct {
	record.Record
	AccountID         string `db:"account_id"`
	Name              string `db:"name"`
	Description       string `db:"description"`
	GameType          string `db:"game_type"`
	TurnDurationHours int    `db:"turn_duration_hours"`
}

func (r *Game) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameAccountID] = r.AccountID
	args[FieldGameName] = r.Name
	args[FieldGameDescription] = r.Description
	args[FieldGameType] = r.GameType
	args[FieldGameTurnDurationHours] = r.TurnDurationHours
	return args
}
