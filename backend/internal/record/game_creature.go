package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableGameCreature = "game_creature"

const (
	FieldGameCreatureID          = "id"
	FieldGameCreatureGameID      = "game_id"
	FieldGameCreatureName        = "name"
	FieldGameCreatureDescription = "description"
)

type GameCreature struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *GameCreature) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameCreatureGameID] = r.GameID
	args[FieldGameCreatureName] = r.Name
	args[FieldGameCreatureDescription] = r.Description
	return args
}
