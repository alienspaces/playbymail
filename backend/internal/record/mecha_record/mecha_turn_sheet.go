package mecha_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaTurnSheet string = "mecha_turn_sheet"
)

const (
	FieldMechaTurnSheetID                         string = "id"
	FieldMechaTurnSheetGameID                     string = "game_id"
	FieldMechaTurnSheetMechaLanceInstanceID string = "mecha_lance_instance_id"
	FieldMechaTurnSheetGameTurnSheetID            string = "game_turn_sheet_id"
	FieldMechaTurnSheetCreatedAt                  string = "created_at"
	FieldMechaTurnSheetUpdatedAt                  string = "updated_at"
	FieldMechaTurnSheetDeletedAt                  string = "deleted_at"
)

const (
	MechaTurnSheetTypeJoinGame        string = "mecha_join_game"
	MechaTurnSheetTypeOrders          string = "mecha_orders"
	MechaTurnSheetTypeLanceManagement string = "mecha_lance_management"
)

// MechaTurnSheetProcessingOrder defines the order in which mecha turn sheets
// are processed during turn resolution. Management resolves before orders so
// refitting mechs are flagged before movement is applied.
// The join game sheet is excluded; it is handled through the subscription workflow.
var MechaTurnSheetProcessingOrder = []string{
	MechaTurnSheetTypeLanceManagement, // 1 - apply repair/refit/swap orders
	MechaTurnSheetTypeOrders,          // 2 - process lance movement orders
}

// MechaSheetOrderForType returns the 1-indexed processing order
// for a mecha turn sheet type. Returns 0 if the type is not
// in the processing order (e.g. join_game).
func MechaSheetOrderForType(sheetType string) int {
	for i, t := range MechaTurnSheetProcessingOrder {
		if t == sheetType {
			return i + 1
		}
	}
	return 0
}

// MechaTurnSheetPresentationOrder defines the order in which mecha turn sheets
// are presented to the player in the UI. Orders are shown first (primary action),
// management is secondary.
var MechaTurnSheetPresentationOrder = []string{
	MechaTurnSheetTypeOrders,          // 1 - submit orders for all mechs
	MechaTurnSheetTypeLanceManagement, // 2 - manage repairs and refits
}

// MechaSheetPresentationOrderForType returns the 1-indexed presentation
// order for a mecha turn sheet type. Returns 0 if not in the order.
func MechaSheetPresentationOrderForType(sheetType string) int {
	for i, t := range MechaTurnSheetPresentationOrder {
		if t == sheetType {
			return i + 1
		}
	}
	return 0
}

// MechaSheetTypes is the set of all mecha sheet type strings.
var MechaSheetTypes = set.New(
	MechaTurnSheetTypeJoinGame,
	MechaTurnSheetTypeOrders,
	MechaTurnSheetTypeLanceManagement,
)

type MechaTurnSheet struct {
	record.Record
	GameID                     string `db:"game_id"`
	MechaLanceInstanceID string `db:"mecha_lance_instance_id"`
	GameTurnSheetID            string `db:"game_turn_sheet_id"`
}

func (r *MechaTurnSheet) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaTurnSheetGameID] = r.GameID
	args[FieldMechaTurnSheetMechaLanceInstanceID] = r.MechaLanceInstanceID
	args[FieldMechaTurnSheetGameTurnSheetID] = r.GameTurnSheetID
	return args
}
