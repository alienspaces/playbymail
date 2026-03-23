package mech_wargame_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechWargameTurnSheet string = "mech_wargame_turn_sheet"
)

const (
	FieldMechWargameTurnSheetID                         string = "id"
	FieldMechWargameTurnSheetGameID                     string = "game_id"
	FieldMechWargameTurnSheetMechWargameLanceInstanceID string = "mech_wargame_lance_instance_id"
	FieldMechWargameTurnSheetGameTurnSheetID            string = "game_turn_sheet_id"
	FieldMechWargameTurnSheetCreatedAt                  string = "created_at"
	FieldMechWargameTurnSheetUpdatedAt                  string = "updated_at"
	FieldMechWargameTurnSheetDeletedAt                  string = "deleted_at"
)

const (
	MechWargameTurnSheetTypeJoinGame string = "mech_wargame_join_game"
	MechWargameTurnSheetTypeOrders   string = "mech_wargame_orders"
)

// MechWargameTurnSheetProcessingOrder defines the order in which
// mech wargame turn sheets are processed during turn resolution.
// The join game sheet is excluded; it is handled through the
// subscription workflow, not turn processing.
var MechWargameTurnSheetProcessingOrder = []string{
	MechWargameTurnSheetTypeOrders, // 1 - process lance orders
}

// MechWargameSheetOrderForType returns the 1-indexed processing order
// for a mech wargame turn sheet type. Returns 0 if the type is not
// in the processing order (e.g. join_game).
func MechWargameSheetOrderForType(sheetType string) int {
	for i, t := range MechWargameTurnSheetProcessingOrder {
		if t == sheetType {
			return i + 1
		}
	}
	return 0
}

// MechWargameTurnSheetPresentationOrder defines the order in which
// mech wargame turn sheets are presented to the player in the UI.
var MechWargameTurnSheetPresentationOrder = []string{
	MechWargameTurnSheetTypeOrders, // 1 - submit orders for all mechs
}

// MechWargameSheetPresentationOrderForType returns the 1-indexed presentation
// order for a mech wargame turn sheet type. Returns 0 if not in the order.
func MechWargameSheetPresentationOrderForType(sheetType string) int {
	for i, t := range MechWargameTurnSheetPresentationOrder {
		if t == sheetType {
			return i + 1
		}
	}
	return 0
}

// MechWargameSheetTypes is the set of all mech wargame sheet type strings.
var MechWargameSheetTypes = set.New(
	MechWargameTurnSheetTypeJoinGame,
	MechWargameTurnSheetTypeOrders,
)

type MechWargameTurnSheet struct {
	record.Record
	GameID                     string `db:"game_id"`
	MechWargameLanceInstanceID string `db:"mech_wargame_lance_instance_id"`
	GameTurnSheetID            string `db:"game_turn_sheet_id"`
}

func (r *MechWargameTurnSheet) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechWargameTurnSheetGameID] = r.GameID
	args[FieldMechWargameTurnSheetMechWargameLanceInstanceID] = r.MechWargameLanceInstanceID
	args[FieldMechWargameTurnSheetGameTurnSheetID] = r.GameTurnSheetID
	return args
}
