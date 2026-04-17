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
	FieldMechaTurnSheetMechaSquadInstanceID string = "mecha_squad_instance_id"
	FieldMechaTurnSheetGameTurnSheetID            string = "game_turn_sheet_id"
	FieldMechaTurnSheetCreatedAt                  string = "created_at"
	FieldMechaTurnSheetUpdatedAt                  string = "updated_at"
	FieldMechaTurnSheetDeletedAt                  string = "deleted_at"
)

const (
	MechaTurnSheetTypeJoinGame        string = "mecha_join_game"
	MechaTurnSheetTypeOrders          string = "mecha_orders"
	MechaTurnSheetTypeSquadManagement string = "mecha_squad_management"
)

// MechaTurnSheetProcessingOrder defines the order in which mecha turn sheets
// are processed during turn resolution. Management resolves before orders so
// refitting mechs are flagged before movement is applied.
// The join game sheet is excluded; it is handled through the subscription workflow.
var MechaTurnSheetProcessingOrder = []string{
	MechaTurnSheetTypeSquadManagement, // 1 - apply repair/refit/swap orders
	MechaTurnSheetTypeOrders,          // 2 - process squad movement orders
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
	MechaTurnSheetTypeSquadManagement, // 2 - manage repairs and refits
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
	MechaTurnSheetTypeSquadManagement,
)

type MechaTurnSheet struct {
	record.Record
	GameID               string `db:"game_id"`
	MechaSquadInstanceID string `db:"mecha_squad_instance_id"`
	GameTurnSheetID      string `db:"game_turn_sheet_id"`
}

func (r *MechaTurnSheet) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaTurnSheetGameID] = r.GameID
	args[FieldMechaTurnSheetMechaSquadInstanceID] = r.MechaSquadInstanceID
	args[FieldMechaTurnSheetGameTurnSheetID] = r.GameTurnSheetID
	return args
}
