package adventure_game_record

import (
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableAdventureGameCharacterInstance string = "adventure_game_character_instance"
)

const (
	FieldAdventureGameCharacterInstanceID                              string = "id"
	FieldAdventureGameCharacterInstanceGameID                          string = "game_id"
	FieldAdventureGameCharacterInstanceGameInstanceID                  string = "game_instance_id"
	FieldAdventureGameCharacterInstanceAdventureGameCharacterID        string = "adventure_game_character_id"
	FieldAdventureGameCharacterInstanceAdventureGameLocationInstanceID string = "adventure_game_location_instance_id"
	FieldAdventureGameCharacterInstanceHealth                          string = "health"
	FieldAdventureGameCharacterInstanceInventoryCapacity               string = "inventory_capacity"
	FieldAdventureGameCharacterInstanceLastTurnEvents                  string = "last_turn_events"
	FieldAdventureGameCharacterInstanceCreatedAt                       string = "created_at"
	FieldAdventureGameCharacterInstanceUpdatedAt                       string = "updated_at"
	FieldAdventureGameCharacterInstanceDeletedAt                       string = "deleted_at"
)

type AdventureGameCharacterInstance struct {
	record.Record
	GameID                          string          `db:"game_id"`
	GameInstanceID                  string          `db:"game_instance_id"`
	AdventureGameCharacterID        string          `db:"adventure_game_character_id"`
	AdventureGameLocationInstanceID sql.NullString  `db:"adventure_game_location_instance_id"`
	Health                          int             `db:"health"`
	InventoryCapacity               int             `db:"inventory_capacity"`
	LastTurnEvents                  json.RawMessage `db:"last_turn_events"`
}

func (r *AdventureGameCharacterInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameCharacterInstanceGameID] = r.GameID
	args[FieldAdventureGameCharacterInstanceGameInstanceID] = r.GameInstanceID
	args[FieldAdventureGameCharacterInstanceAdventureGameCharacterID] = r.AdventureGameCharacterID
	args[FieldAdventureGameCharacterInstanceAdventureGameLocationInstanceID] = r.AdventureGameLocationInstanceID
	args[FieldAdventureGameCharacterInstanceHealth] = r.Health
	args[FieldAdventureGameCharacterInstanceInventoryCapacity] = r.InventoryCapacity
	args[FieldAdventureGameCharacterInstanceLastTurnEvents] = r.LastTurnEvents
	return args
}
