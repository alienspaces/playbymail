package record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableAdventureGameCharacterInstance string = "adventure_game_character_instance"
)

const (
	FieldAdventureGameCharacterInstanceID                              string = "id"
	FieldAdventureGameCharacterInstanceGameID                          string = "game_id"
	FieldAdventureGameCharacterInstanceAdventureGameInstanceID         string = "adventure_game_instance_id"
	FieldAdventureGameCharacterInstanceAdventureGameCharacterID        string = "adventure_game_character_id"
	FieldAdventureGameCharacterInstanceAdventureGameLocationInstanceID string = "adventure_game_location_instance_id"
	FieldAdventureGameCharacterInstanceHealth                          string = "health"
	FieldAdventureGameCharacterInstanceCreatedAt                       string = "created_at"
	FieldAdventureGameCharacterInstanceUpdatedAt                       string = "updated_at"
	FieldAdventureGameCharacterInstanceDeletedAt                       string = "deleted_at"
)

type AdventureGameCharacterInstance struct {
	record.Record
	GameID                          string `db:"game_id"`
	AdventureGameInstanceID         string `db:"adventure_game_instance_id"`
	AdventureGameCharacterID        string `db:"adventure_game_character_id"`
	AdventureGameLocationInstanceID string `db:"adventure_game_location_instance_id"`
	Health                          int    `db:"health"`
}

func (r *AdventureGameCharacterInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameCharacterInstanceGameID] = r.GameID
	args[FieldAdventureGameCharacterInstanceAdventureGameInstanceID] = r.AdventureGameInstanceID
	args[FieldAdventureGameCharacterInstanceAdventureGameCharacterID] = r.AdventureGameCharacterID
	args[FieldAdventureGameCharacterInstanceAdventureGameLocationInstanceID] = r.AdventureGameLocationInstanceID
	args[FieldAdventureGameCharacterInstanceHealth] = r.Health
	return args
}
