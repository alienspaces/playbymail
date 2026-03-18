package adventure_game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

const TableAdventureGameLocationLinkRequirement = "adventure_game_location_link_requirement"

// Purpose values
const (
	AdventureGameLocationLinkRequirementPurposeTraverse = "traverse"
	AdventureGameLocationLinkRequirementPurposeVisible  = "visible"
)

// Condition values for item-based requirements
const (
	AdventureGameLocationLinkRequirementConditionInInventory = "in_inventory"
	AdventureGameLocationLinkRequirementConditionEquipped    = "equipped"
)

// Condition values for creature-based requirements
const (
	AdventureGameLocationLinkRequirementConditionDeadAtLocation      = "dead_at_location"
	AdventureGameLocationLinkRequirementConditionNoneAliveAtLocation  = "none_alive_at_location"
	AdventureGameLocationLinkRequirementConditionNoneAliveInGame      = "none_alive_in_game"
)

const (
	FieldAdventureGameLocationLinkRequirementID                          = "id"
	FieldAdventureGameLocationLinkRequirementGameID                      = "game_id"
	FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID = "adventure_game_location_link_id"
	FieldAdventureGameLocationLinkRequirementAdventureGameItemID         = "adventure_game_item_id"
	FieldAdventureGameLocationLinkRequirementAdventureGameCreatureID     = "adventure_game_creature_id"
	FieldAdventureGameLocationLinkRequirementPurpose                     = "purpose"
	FieldAdventureGameLocationLinkRequirementCondition                   = "condition"
	FieldAdventureGameLocationLinkRequirementQuantity                    = "quantity"
)

// AdventureGameLocationLinkRequirement specifies conditions required to traverse or see a location link.
// Exactly one of AdventureGameItemID or AdventureGameCreatureID must be set.
// Multiple rows for the same link + purpose = AND (all conditions must be satisfied).
type AdventureGameLocationLinkRequirement struct {
	record.Record
	GameID                      string         `db:"game_id"`
	AdventureGameLocationLinkID string         `db:"adventure_game_location_link_id"`
	AdventureGameItemID         sql.NullString `db:"adventure_game_item_id"`
	AdventureGameCreatureID     sql.NullString `db:"adventure_game_creature_id"`
	Purpose                     string         `db:"purpose"`
	Condition                   string         `db:"condition"`
	Quantity                    int            `db:"quantity"`
}

func (r *AdventureGameLocationLinkRequirement) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameLocationLinkRequirementGameID] = r.GameID
	args[FieldAdventureGameLocationLinkRequirementAdventureGameLocationLinkID] = r.AdventureGameLocationLinkID
	args[FieldAdventureGameLocationLinkRequirementAdventureGameItemID] = r.AdventureGameItemID
	args[FieldAdventureGameLocationLinkRequirementAdventureGameCreatureID] = r.AdventureGameCreatureID
	args[FieldAdventureGameLocationLinkRequirementPurpose] = r.Purpose
	args[FieldAdventureGameLocationLinkRequirementCondition] = r.Condition
	args[FieldAdventureGameLocationLinkRequirementQuantity] = r.Quantity
	return args
}
