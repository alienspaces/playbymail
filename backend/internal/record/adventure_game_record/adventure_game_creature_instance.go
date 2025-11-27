package adventure_game_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameCreatureInstance = "adventure_game_creature_instance"

const (
	FieldAdventureGameCreatureInstanceID                              = "id"
	FieldAdventureGameCreatureInstanceGameID                          = "game_id"
	FieldAdventureGameCreatureInstanceGameInstanceID                  = "game_instance_id"
	FieldAdventureGameCreatureInstanceAdventureGameCreatureID         = "adventure_game_creature_id"
	FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID = "adventure_game_location_instance_id"
	FieldAdventureGameCreatureInstanceHealth                          = "health"
	FieldAdventureGameCreatureInstanceCreatedAt                       = "created_at"
	FieldAdventureGameCreatureInstanceUpdatedAt                       = "updated_at"
	FieldAdventureGameCreatureInstanceDeletedAt                       = "deleted_at"
)

type AdventureGameCreatureInstance struct {
	record.Record
	GameID                          string `db:"game_id"`
	GameInstanceID                  string `db:"game_instance_id"`
	AdventureGameCreatureID         string `db:"adventure_game_creature_id"`
	AdventureGameLocationInstanceID string `db:"adventure_game_location_instance_id"`
	Health                          int    `db:"health"`
}

func (r *AdventureGameCreatureInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameCreatureInstanceGameID] = r.GameID
	args[FieldAdventureGameCreatureInstanceGameInstanceID] = r.GameInstanceID
	args[FieldAdventureGameCreatureInstanceAdventureGameCreatureID] = r.AdventureGameCreatureID
	args[FieldAdventureGameCreatureInstanceAdventureGameLocationInstanceID] = r.AdventureGameLocationInstanceID
	args[FieldAdventureGameCreatureInstanceHealth] = r.Health
	return args
}
