package adventure_game_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// GameLocation
const (
	TableAdventureGameLocation string = "adventure_game_location"
)

const (
	FieldAdventureGameLocationID                 string = "id"
	FieldAdventureGameLocationGameID             string = "game_id"
	FieldAdventureGameLocationName               string = "name"
	FieldAdventureGameLocationDescription        string = "description"
	FieldAdventureGameLocationIsStartingLocation string = "is_starting_location"
)

type AdventureGameLocation struct {
	record.Record
	GameID             string `db:"game_id"`
	Name               string `db:"name"`
	Description        string `db:"description"`
	IsStartingLocation bool   `db:"is_starting_location"`
}

func (r *AdventureGameLocation) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationGameID] = r.GameID
	args[FieldAdventureGameLocationName] = r.Name
	args[FieldAdventureGameLocationDescription] = r.Description
	args[FieldAdventureGameLocationIsStartingLocation] = r.IsStartingLocation
	return args
}
