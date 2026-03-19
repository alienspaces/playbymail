package adventure_game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameLocationLink = "adventure_game_location_link"

const (
	FieldAdventureGameLocationLinkID                          = "id"
	FieldAdventureGameLocationLinkGameID                      = "game_id"
	FieldAdventureGameLocationLinkFromAdventureGameLocationID = "from_adventure_game_location_id"
	FieldAdventureGameLocationLinkToAdventureGameLocationID   = "to_adventure_game_location_id"
	FieldAdventureGameLocationLinkDescription                 = "description"
	FieldAdventureGameLocationLinkLockedDescription           = "locked_description"
	FieldAdventureGameLocationLinkTraversalDescription        = "traversal_description"
	FieldAdventureGameLocationLinkName                        = "name"
	FieldAdventureGameLocationLinkCreatedAt                   = "created_at"
	FieldAdventureGameLocationLinkUpdatedAt                   = "updated_at"
)

type AdventureGameLocationLink struct {
	record.Record
	GameID                      string         `db:"game_id"`
	FromAdventureGameLocationID string         `db:"from_adventure_game_location_id"`
	ToAdventureGameLocationID   string         `db:"to_adventure_game_location_id"`
	Name                        string         `db:"name"`
	Description                 string         `db:"description"`
	LockedDescription           sql.NullString `db:"locked_description"`
	TraversalDescription        sql.NullString `db:"traversal_description"`
}

func (r *AdventureGameLocationLink) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationLinkGameID] = r.GameID
	args[FieldAdventureGameLocationLinkFromAdventureGameLocationID] = r.FromAdventureGameLocationID
	args[FieldAdventureGameLocationLinkToAdventureGameLocationID] = r.ToAdventureGameLocationID
	args[FieldAdventureGameLocationLinkName] = r.Name
	args[FieldAdventureGameLocationLinkDescription] = r.Description
	args[FieldAdventureGameLocationLinkLockedDescription] = r.LockedDescription
	args[FieldAdventureGameLocationLinkTraversalDescription] = r.TraversalDescription
	return args
}
