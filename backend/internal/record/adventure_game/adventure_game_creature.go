package adventure_game_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameCreature = "adventure_game_creature"

const (
	FieldAdventureGameCreatureID          = "id"
	FieldAdventureGameCreatureGameID      = "game_id"
	FieldAdventureGameCreatureName        = "name"
	FieldAdventureGameCreatureDescription = "description"
)

type AdventureGameCreature struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *AdventureGameCreature) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameCreatureGameID] = r.GameID
	args[FieldAdventureGameCreatureName] = r.Name
	args[FieldAdventureGameCreatureDescription] = r.Description
	return args
}
