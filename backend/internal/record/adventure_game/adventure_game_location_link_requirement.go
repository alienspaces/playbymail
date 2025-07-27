package adventure_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameLocationLinkRequirement = "adventure_game_location_link_requirement"

const (
	FieldAdventureGameLocationLinkRequirementID                          = "id"
	FieldAdventureGameLocationLinkRequirementGameID                      = "game_id"
	FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID = "adventure_game_location_link_id"
	FieldAdventureGameLocationLinkRequirementAdventureGameItemID         = "adventure_game_item_id"
	FieldAdventureGameLocationLinkRequirementQuantity                    = "quantity"
)

// AdventureGameLocationLinkRequirement specifies which items (and how many) are required to traverse a location link.
type AdventureGameLocationLinkRequirement struct {
	record.Record
	GameID                      string `db:"game_id"`
	AdventureGameLocationLinkID string `db:"adventure_game_location_link_id"`
	AdventureGameItemID         string `db:"adventure_game_item_id"`
	Quantity                    int    `db:"quantity"`
}

func (r *AdventureGameLocationLinkRequirement) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationLinkRequirementGameID] = r.GameID
	args[FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID] = r.AdventureGameLocationLinkID
	args[FieldAdventureGameLocationLinkRequirementAdventureGameItemID] = r.AdventureGameItemID
	args[FieldAdventureGameLocationLinkRequirementQuantity] = r.Quantity
	return args
}
