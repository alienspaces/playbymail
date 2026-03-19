package adventure_game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameLocationObject = "adventure_game_location_object"

const (
	FieldAdventureGameLocationObjectID                                       = "id"
	FieldAdventureGameLocationObjectGameID                                   = "game_id"
	FieldAdventureGameLocationObjectAdventureGameLocationID                  = "adventure_game_location_id"
	FieldAdventureGameLocationObjectName                                     = "name"
	FieldAdventureGameLocationObjectDescription                              = "description"
	FieldAdventureGameLocationObjectInitialAdventureGameLocationObjectStateID = "initial_adventure_game_location_object_state_id"
	FieldAdventureGameLocationObjectIsHidden                                 = "is_hidden"
)

type AdventureGameLocationObject struct {
	record.Record
	GameID                                    string         `db:"game_id"`
	AdventureGameLocationID                   string         `db:"adventure_game_location_id"`
	Name                                      string         `db:"name"`
	Description                               string         `db:"description"`
	InitialAdventureGameLocationObjectStateID sql.NullString `db:"initial_adventure_game_location_object_state_id"`
	IsHidden                                  bool           `db:"is_hidden"`
}

func (r *AdventureGameLocationObject) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationObjectGameID] = r.GameID
	args[FieldAdventureGameLocationObjectAdventureGameLocationID] = r.AdventureGameLocationID
	args[FieldAdventureGameLocationObjectName] = r.Name
	args[FieldAdventureGameLocationObjectDescription] = r.Description
	args[FieldAdventureGameLocationObjectInitialAdventureGameLocationObjectStateID] = r.InitialAdventureGameLocationObjectStateID
	args[FieldAdventureGameLocationObjectIsHidden] = r.IsHidden
	return args
}
