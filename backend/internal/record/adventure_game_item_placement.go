package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameItemPlacement = "adventure_game_item_placement"

const (
	FieldAdventureGameItemPlacementID                      = "id"
	FieldAdventureGameItemPlacementGameID                  = "game_id"
	FieldAdventureGameItemPlacementAdventureGameItemID     = "adventure_game_item_id"
	FieldAdventureGameItemPlacementAdventureGameLocationID = "adventure_game_location_id"
	FieldAdventureGameItemPlacementInitialCount            = "initial_count"
)

type AdventureGameItemPlacement struct {
	record.Record
	GameID                  string `db:"game_id"`
	AdventureGameItemID     string `db:"adventure_game_item_id"`
	AdventureGameLocationID string `db:"adventure_game_location_id"`
	InitialCount            int    `db:"initial_count"`
}

func (r *AdventureGameItemPlacement) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameItemPlacementGameID] = r.GameID
	args[FieldAdventureGameItemPlacementAdventureGameItemID] = r.AdventureGameItemID
	args[FieldAdventureGameItemPlacementAdventureGameLocationID] = r.AdventureGameLocationID
	args[FieldAdventureGameItemPlacementInitialCount] = r.InitialCount
	return args
}
