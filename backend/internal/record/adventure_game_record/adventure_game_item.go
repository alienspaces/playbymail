package adventure_game_record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameItem = "adventure_game_item"

const (
	FieldAdventureGameItemID          = "id"
	FieldAdventureGameItemGameID      = "game_id"
	FieldAdventureGameItemName        = "name"
	FieldAdventureGameItemDescription = "description"
	FieldAdventureGameItemCanBeEquipped = "can_be_equipped"
	FieldAdventureGameItemCategory    = "item_category"
	FieldAdventureGameItemEquipmentSlot = "equipment_slot"
)

type AdventureGameItem struct {
	record.Record
	GameID        string  `db:"game_id"`
	Name          string  `db:"name"`
	Description   string  `db:"description"`
	CanBeEquipped bool    `db:"can_be_equipped"`
	ItemCategory  *string `db:"item_category"`
	EquipmentSlot *string `db:"equipment_slot"`
}

func (r *AdventureGameItem) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameItemGameID] = r.GameID
	args[FieldAdventureGameItemName] = r.Name
	args[FieldAdventureGameItemDescription] = r.Description
	args[FieldAdventureGameItemCanBeEquipped] = r.CanBeEquipped
	args[FieldAdventureGameItemCategory] = r.ItemCategory
	args[FieldAdventureGameItemEquipmentSlot] = r.EquipmentSlot
	return args
}
