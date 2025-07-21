package record

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
)

type AdventureGameItem struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *AdventureGameItem) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameItemGameID] = r.GameID
	args[FieldAdventureGameItemName] = r.Name
	args[FieldAdventureGameItemDescription] = r.Description
	return args
}
