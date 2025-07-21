package record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameLocationInstance = "adventure_game_location_instance"

const (
	FieldAdventureGameLocationInstanceID                      = "id"
	FieldAdventureGameLocationInstanceGameID                  = "game_id"
	FieldAdventureGameLocationInstanceAdventureGameInstanceID = "adventure_game_instance_id"
	FieldAdventureGameLocationInstanceAdventureGameLocationID = "adventure_game_location_id"
)

type AdventureGameLocationInstance struct {
	record.Record
	GameID                  string `db:"game_id"`
	AdventureGameInstanceID string `db:"adventure_game_instance_id"`
	AdventureGameLocationID string `db:"adventure_game_location_id"`
}

func (r *AdventureGameLocationInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationInstanceGameID] = r.GameID
	args[FieldAdventureGameLocationInstanceAdventureGameInstanceID] = r.AdventureGameInstanceID
	args[FieldAdventureGameLocationInstanceAdventureGameLocationID] = r.AdventureGameLocationID
	return args
}
