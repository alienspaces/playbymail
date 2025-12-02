package adventure_game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameItemInstance = "adventure_game_item_instance"

const (
	FieldAdventureGameItemInstanceID                               = "id"
	FieldAdventureGameItemInstanceGameID                           = "game_id"
	FieldAdventureGameItemInstanceGameInstanceID                   = "game_instance_id"
	FieldAdventureGameItemInstanceAdventureGameItemID              = "adventure_game_item_id"
	FieldAdventureGameItemInstanceAdventureGameLocationInstanceID  = "adventure_game_location_instance_id"
	FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID = "adventure_game_character_instance_id"
	FieldAdventureGameItemInstanceAdventureGameCreatureInstanceID  = "adventure_game_creature_instance_id"
	FieldAdventureGameItemInstanceIsEquipped                       = "is_equipped"
	FieldAdventureGameItemInstanceEquipmentSlot                    = "equipment_slot"
	FieldAdventureGameItemInstanceIsUsed                           = "is_used"
	FieldAdventureGameItemInstanceUsesRemaining                    = "uses_remaining"
)

// GameItemInstance represents a specific instance of a game item, which may be at a location, in a character inventory, or in a creature inventory.
type AdventureGameItemInstance struct {
	record.Record
	GameID                           string         `db:"game_id"`
	GameInstanceID                   string         `db:"game_instance_id"`
	AdventureGameItemID              string         `db:"adventure_game_item_id"`
	AdventureGameLocationInstanceID  sql.NullString `db:"adventure_game_location_instance_id"`
	AdventureGameCharacterInstanceID sql.NullString `db:"adventure_game_character_instance_id"`
	AdventureGameCreatureInstanceID  sql.NullString `db:"adventure_game_creature_instance_id"`
	IsEquipped                       bool           `db:"is_equipped"`
	EquipmentSlot                    sql.NullString `db:"equipment_slot"`
	IsUsed                           bool           `db:"is_used"`
	UsesRemaining                    int            `db:"uses_remaining"`
}

func (r *AdventureGameItemInstance) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameItemInstanceGameID] = r.GameID
	args[FieldAdventureGameItemInstanceGameInstanceID] = r.GameInstanceID
	args[FieldAdventureGameItemInstanceAdventureGameItemID] = r.AdventureGameItemID
	args[FieldAdventureGameItemInstanceAdventureGameLocationInstanceID] = r.AdventureGameLocationInstanceID
	args[FieldAdventureGameItemInstanceAdventureGameCharacterInstanceID] = r.AdventureGameCharacterInstanceID
	args[FieldAdventureGameItemInstanceAdventureGameCreatureInstanceID] = r.AdventureGameCreatureInstanceID
	args[FieldAdventureGameItemInstanceIsEquipped] = r.IsEquipped
	args[FieldAdventureGameItemInstanceEquipmentSlot] = r.EquipmentSlot
	args[FieldAdventureGameItemInstanceIsUsed] = r.IsUsed
	args[FieldAdventureGameItemInstanceUsesRemaining] = r.UsesRemaining
	return args
}
