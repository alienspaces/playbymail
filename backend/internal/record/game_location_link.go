package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableGameLocationLink = "game_location_link"

const (
	FieldGameLocationLinkID                 = "id"
	FieldGameLocationLinkFromGameLocationID = "from_game_location_id"
	FieldGameLocationLinkToGameLocationID   = "to_game_location_id"
	FieldGameLocationLinkDescription        = "description"
	FieldGameLocationLinkName               = "name"
	FieldGameLocationLinkCreatedAt          = "created_at"
	FieldGameLocationLinkUpdatedAt          = "updated_at"
)

type GameLocationLink struct {
	record.Record
	FromGameLocationID string `db:"from_game_location_id"`
	ToGameLocationID   string `db:"to_game_location_id"`
	Description        string `db:"description"`
	Name               string `db:"name"`
}

func (r *GameLocationLink) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameLocationLinkFromGameLocationID] = r.FromGameLocationID
	args[FieldGameLocationLinkToGameLocationID] = r.ToGameLocationID
	args[FieldGameLocationLinkDescription] = r.Description
	args[FieldGameLocationLinkName] = r.Name
	return args
}
