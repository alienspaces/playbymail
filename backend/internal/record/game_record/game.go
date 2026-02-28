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
	FieldGameName              string = "name"
	FieldGameDescription       string = "description"
	FieldGameType              string = "game_type"
	FieldGameTurnDurationHours string = "turn_duration_hours"
	FieldGameStatus            string = "status"
)

const (
	GameTypeAdventure string = "adventure"
)

const (
	GameStatusDraft     string = "draft"
	GameStatusPublished string = "published"
)

type Game struct {
	record.Record
	Name              string `db:"name"`
	Description       string `db:"description"`
	GameType          string `db:"game_type"`
	TurnDurationHours int    `db:"turn_duration_hours"`
	Status            string `db:"status"`
}

func (r *Game) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameName] = r.Name
	args[FieldGameDescription] = r.Description
	args[FieldGameType] = r.GameType
	args[FieldGameTurnDurationHours] = r.TurnDurationHours
	args[FieldGameStatus] = r.Status
	return args
}
