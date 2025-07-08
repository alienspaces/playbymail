package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableLocationLink = "location_link"

const (
	FieldLocationLinkID             = "id"
	FieldLocationLinkFromLocationID = "from_location_id"
	FieldLocationLinkToLocationID   = "to_location_id"
	FieldLocationLinkDescription    = "description"
	FieldLocationLinkName           = "name"
	FieldLocationLinkCreatedAt      = "created_at"
	FieldLocationLinkUpdatedAt      = "updated_at"
)

type LocationLink struct {
	record.Record
	FromLocationID string `db:"from_location_id"`
	ToLocationID   string `db:"to_location_id"`
	Description    string `db:"description"`
	Name           string `db:"name"`
}

func (r *LocationLink) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldLocationLinkFromLocationID] = r.FromLocationID
	args[FieldLocationLinkToLocationID] = r.ToLocationID
	args[FieldLocationLinkDescription] = r.Description
	args[FieldLocationLinkName] = r.Name
	return args
}
