package adventure_game_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameLocationObjectState = "adventure_game_location_object_state"

const (
	FieldAdventureGameLocationObjectStateID                            = "id"
	FieldAdventureGameLocationObjectStateGameID                        = "game_id"
	FieldAdventureGameLocationObjectStateAdventureGameLocationObjectID = "adventure_game_location_object_id"
	FieldAdventureGameLocationObjectStateName                          = "name"
	FieldAdventureGameLocationObjectStateDescription                   = "description"
	FieldAdventureGameLocationObjectStateSortOrder                     = "sort_order"
)

// AdventureGameLocationObjectState defines a discrete named state an object can be in.
// States are scoped per object and identified by name (e.g. "intact", "activated", "broken").
type AdventureGameLocationObjectState struct {
	record.Record
	GameID                        string `db:"game_id"`
	AdventureGameLocationObjectID string `db:"adventure_game_location_object_id"`
	Name                          string `db:"name"`
	Description                   string `db:"description"`
	SortOrder                     int    `db:"sort_order"`
}

func (r *AdventureGameLocationObjectState) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationObjectStateGameID] = r.GameID
	args[FieldAdventureGameLocationObjectStateAdventureGameLocationObjectID] = r.AdventureGameLocationObjectID
	args[FieldAdventureGameLocationObjectStateName] = r.Name
	args[FieldAdventureGameLocationObjectStateDescription] = r.Description
	args[FieldAdventureGameLocationObjectStateSortOrder] = r.SortOrder
	return args
}
