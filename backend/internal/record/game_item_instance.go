package record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableGameItemInstance = "game_item_instance"

const (
	FieldGameItemInstanceID                      = "id"
	FieldGameItemInstanceGameID                  = "game_id"
	FieldGameItemInstanceGameItemID              = "game_item_id"
	FieldGameItemInstanceGameInstanceID          = "game_instance_id"
	FieldGameItemInstanceGameLocationInstanceID  = "game_location_instance_id"
	FieldGameItemInstanceGameCharacterInstanceID = "game_character_instance_id"
	FieldGameItemInstanceGameCreatureInstanceID  = "game_creature_instance_id"
	FieldGameItemInstanceIsEquipped              = "is_equipped"
	FieldGameItemInstanceIsUsed                  = "is_used"
	FieldGameItemInstanceUsesRemaining           = "uses_remaining"
)

// GameItemInstance represents a specific instance of a game item, which may be at a location, in a character inventory, or in a creature inventory.
type GameItemInstance struct {
	record.Record
	GameID                  string         `db:"game_id"`
	GameItemID              string         `db:"game_item_id"`
	GameInstanceID          string         `db:"game_instance_id"`
	GameLocationInstanceID  sql.NullString `db:"game_location_instance_id"`
	GameCharacterInstanceID sql.NullString `db:"game_character_instance_id"`
	GameCreatureInstanceID  sql.NullString `db:"game_creature_instance_id"`
	IsEquipped              bool           `db:"is_equipped"`
	IsUsed                  bool           `db:"is_used"`
	UsesRemaining           *int           `db:"uses_remaining"`
}

func (r *GameItemInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameItemInstanceGameID] = r.GameID
	args[FieldGameItemInstanceGameItemID] = r.GameItemID
	args[FieldGameItemInstanceGameInstanceID] = r.GameInstanceID
	args[FieldGameItemInstanceGameLocationInstanceID] = r.GameLocationInstanceID
	args[FieldGameItemInstanceGameCharacterInstanceID] = r.GameCharacterInstanceID
	args[FieldGameItemInstanceGameCreatureInstanceID] = r.GameCreatureInstanceID
	args[FieldGameItemInstanceIsEquipped] = r.IsEquipped
	args[FieldGameItemInstanceIsUsed] = r.IsUsed
	args[FieldGameItemInstanceUsesRemaining] = r.UsesRemaining
	return args
}
