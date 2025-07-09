package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableLocationLink = "location_link"

const (
	FieldLocationLinkID                 = "id"
	FieldLocationLinkFromGameLocationID = "from_game_location_id"
	FieldLocationLinkToGameLocationID   = "to_game_location_id"
	FieldLocationLinkDescription        = "description"
	FieldLocationLinkName               = "name"
	FieldLocationLinkCreatedAt          = "created_at"
	FieldLocationLinkUpdatedAt          = "updated_at"
)

type LocationLink struct {
	record.Record
	FromGameLocationID string `db:"from_game_location_id"`
	ToGameLocationID   string `db:"to_game_location_id"`
	Description        string `db:"description"`
	Name               string `db:"name"`
}

func (r *LocationLink) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldLocationLinkFromGameLocationID] = r.FromGameLocationID
	args[FieldLocationLinkToGameLocationID] = r.ToGameLocationID
	args[FieldLocationLinkDescription] = r.Description
	args[FieldLocationLinkName] = r.Name
	return args
}
