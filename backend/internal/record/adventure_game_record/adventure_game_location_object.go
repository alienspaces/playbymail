package adventure_game_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameLocationObject = "adventure_game_location_object"

const (
	FieldAdventureGameLocationObjectID                     = "id"
	FieldAdventureGameLocationObjectGameID                 = "game_id"
	FieldAdventureGameLocationObjectAdventureGameLocationID = "adventure_game_location_id"
	FieldAdventureGameLocationObjectName                   = "name"
	FieldAdventureGameLocationObjectDescription            = "description"
	FieldAdventureGameLocationObjectInitialState           = "initial_state"
	FieldAdventureGameLocationObjectIsHidden               = "is_hidden"
)

type AdventureGameLocationObject struct {
	record.Record
	GameID                  string `db:"game_id"`
	AdventureGameLocationID string `db:"adventure_game_location_id"`
	Name                    string `db:"name"`
	Description             string `db:"description"`
	InitialState            string `db:"initial_state"`
	IsHidden                bool   `db:"is_hidden"`
}

func (r *AdventureGameLocationObject) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationObjectGameID] = r.GameID
	args[FieldAdventureGameLocationObjectAdventureGameLocationID] = r.AdventureGameLocationID
	args[FieldAdventureGameLocationObjectName] = r.Name
	args[FieldAdventureGameLocationObjectDescription] = r.Description
	args[FieldAdventureGameLocationObjectInitialState] = r.InitialState
	args[FieldAdventureGameLocationObjectIsHidden] = r.IsHidden
	return args
}
