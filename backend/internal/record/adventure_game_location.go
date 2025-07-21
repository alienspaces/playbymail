package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameLocation
const (
	TableAdventureGameLocation string = "adventure_game_location"
)

const (
	FieldAdventureGameLocationID          string = "id"
	FieldAdventureGameLocationGameID      string = "game_id"
	FieldAdventureGameLocationName        string = "name"
	FieldAdventureGameLocationDescription string = "description"
)

type AdventureGameLocation struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *AdventureGameLocation) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationGameID] = r.GameID
	args[FieldAdventureGameLocationName] = r.Name
	args[FieldAdventureGameLocationDescription] = r.Description
	return args
}
