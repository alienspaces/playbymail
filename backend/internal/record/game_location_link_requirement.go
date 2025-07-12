package record

import (
	"github.com/jackc/pgx/v5"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableGameLocationLinkRequirement = "game_location_link_requirement"

const (
	FieldGameLocationLinkRequirementID                 = "id"
	FieldGameLocationLinkRequirementGameID             = "game_id"
	FieldGameLocationLinkRequirementGameLocationLinkID = "game_location_link_id"
	FieldGameLocationLinkRequirementGameItemID         = "game_item_id"
	FieldGameLocationLinkRequirementQuantity           = "quantity"
)

// GameLocationLinkRequirement specifies which items (and how many) are required to traverse a location link.
type GameLocationLinkRequirement struct {
	record.Record
	GameID             string `db:"game_id" json:"gameId"`
	GameLocationLinkID string `db:"game_location_link_id" json:"gameLocationLinkId"`
	GameItemID         string `db:"game_item_id" json:"gameItemId"`
	Quantity           int    `db:"quantity" json:"quantity"`
}

func (r *GameLocationLinkRequirement) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldGameLocationLinkRequirementGameID] = r.GameID
	args[FieldGameLocationLinkRequirementGameLocationLinkID] = r.GameLocationLinkID
	args[FieldGameLocationLinkRequirementGameItemID] = r.GameItemID
	args[FieldGameLocationLinkRequirementQuantity] = r.Quantity
	return args
}
