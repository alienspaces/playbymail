package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableGameItem = "game_item"

const (
	FieldGameItemID          = "id"
	FieldGameItemGameID      = "game_id"
	FieldGameItemName        = "name"
	FieldGameItemDescription = "description"
)

type GameItem struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *GameItem) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameItemGameID] = r.GameID
	args[FieldGameItemName] = r.Name
	args[FieldGameItemDescription] = r.Description
	return args
}
