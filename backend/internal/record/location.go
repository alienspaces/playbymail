package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

// Location
const (
	TableLocation string = "location"
)

const (
	FieldLocationID          string = "id"
	FieldLocationGameID      string = "game_id"
	FieldLocationName        string = "name"
	FieldLocationDescription string = "description"
)

type Location struct {
	record.Record
	GameID      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (r *Location) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldLocationGameID] = r.GameID
	args[FieldLocationName] = r.Name
	args[FieldLocationDescription] = r.Description
	return args
}
