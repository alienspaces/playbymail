package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameCreaturePlacement = "adventure_game_creature_placement"

const (
	FieldAdventureGameCreaturePlacementID                      = "id"
	FieldAdventureGameCreaturePlacementGameID                  = "game_id"
	FieldAdventureGameCreaturePlacementAdventureGameCreatureID = "adventure_game_creature_id"
	FieldAdventureGameCreaturePlacementAdventureGameLocationID = "adventure_game_location_id"
	FieldAdventureGameCreaturePlacementInitialCount            = "initial_count"
)

type AdventureGameCreaturePlacement struct {
	record.Record
	GameID                  string `db:"game_id"`
	AdventureGameCreatureID string `db:"adventure_game_creature_id"`
	AdventureGameLocationID string `db:"adventure_game_location_id"`
	InitialCount            int    `db:"initial_count"`
}

func (r *AdventureGameCreaturePlacement) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameCreaturePlacementGameID] = r.GameID
	args[FieldAdventureGameCreaturePlacementAdventureGameCreatureID] = r.AdventureGameCreatureID
	args[FieldAdventureGameCreaturePlacementAdventureGameLocationID] = r.AdventureGameLocationID
	args[FieldAdventureGameCreaturePlacementInitialCount] = r.InitialCount
	return args
}
