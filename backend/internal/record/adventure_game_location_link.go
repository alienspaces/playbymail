package record

import (
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
	FieldAdventureGameLocationLinkName                        = "name"
	FieldAdventureGameLocationLinkCreatedAt                   = "created_at"
	FieldAdventureGameLocationLinkUpdatedAt                   = "updated_at"
)

type AdventureGameLocationLink struct {
	record.Record
	GameID                      string `db:"game_id"`
	FromAdventureGameLocationID string `db:"from_adventure_game_location_id"`
	ToAdventureGameLocationID   string `db:"to_adventure_game_location_id"`
	Description                 string `db:"description"`
	Name                        string `db:"name"`
}

func (r *AdventureGameLocationLink) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationLinkGameID] = r.GameID
	args[FieldAdventureGameLocationLinkFromAdventureGameLocationID] = r.FromAdventureGameLocationID
	args[FieldAdventureGameLocationLinkToAdventureGameLocationID] = r.ToAdventureGameLocationID
	args[FieldAdventureGameLocationLinkDescription] = r.Description
	args[FieldAdventureGameLocationLinkName] = r.Name
	return args
}
