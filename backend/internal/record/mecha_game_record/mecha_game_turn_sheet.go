package mecha_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableMechaGameTurnSheet string = "mecha_game_turn_sheet"
)

const (
	FieldMechaGameTurnSheetID                         string = "id"
	FieldMechaGameTurnSheetGameID                     string = "game_id"
	FieldMechaGameTurnSheetMechaGameSquadInstanceID string = "mecha_game_squad_instance_id"
	FieldMechaGameTurnSheetGameTurnSheetID            string = "game_turn_sheet_id"
	FieldMechaGameTurnSheetCreatedAt                  string = "created_at"
	FieldMechaGameTurnSheetUpdatedAt                  string = "updated_at"
	FieldMechaGameTurnSheetDeletedAt                  string = "deleted_at"
)

const (
	MechaGameTurnSheetTypeJoinGame        string = "mecha_game_join_game"
	MechaGameTurnSheetTypeOrders          string = "mecha_game_orders"
	MechaGameTurnSheetTypeSquadManagement string = "mecha_game_squad_management"
)

// MechaGameTurnSheetProcessingOrder defines the order in which mecha turn sheets
// are processed during turn resolution. Management resolves before orders so
// refitting mechs are flagged before movement is applied.
// The join game sheet is excluded; it is handled through the subscription workflow.
var MechaGameTurnSheetProcessingOrder = []string{
	MechaGameTurnSheetTypeSquadManagement, // 1 - apply repair/refit/swap orders
	MechaGameTurnSheetTypeOrders,          // 2 - process squad movement orders
}

// MechaGameSheetOrderForType returns the 1-indexed processing order
// for a mecha turn sheet type. Returns 0 if the type is not
// in the processing order (e.g. join_game).
func MechaGameSheetOrderForType(sheetType string) int {
	for i, t := range MechaGameTurnSheetProcessingOrder {
		if t == sheetType {
			return i + 1
		}
	}
	return 0
}

// MechaGameTurnSheetPresentationOrder defines the order in which mecha turn sheets
// are presented to the player in the UI. Orders are shown first (primary action),
// management is secondary.
var MechaGameTurnSheetPresentationOrder = []string{
	MechaGameTurnSheetTypeOrders,          // 1 - submit orders for all mechs
	MechaGameTurnSheetTypeSquadManagement, // 2 - manage repairs and refits
}

// MechaGameSheetPresentationOrderForType returns the 1-indexed presentation
// order for a mecha turn sheet type. Returns 0 if not in the order.
func MechaGameSheetPresentationOrderForType(sheetType string) int {
	for i, t := range MechaGameTurnSheetPresentationOrder {
		if t == sheetType {
			return i + 1
		}
	}
	return 0
}

// MechaGameSheetTypes is the set of all mecha sheet type strings.
var MechaGameSheetTypes = set.New(
	MechaGameTurnSheetTypeJoinGame,
	MechaGameTurnSheetTypeOrders,
	MechaGameTurnSheetTypeSquadManagement,
)

type MechaGameTurnSheet struct {
	record.Record
	GameID               string `db:"game_id"`
	MechaGameSquadInstanceID string `db:"mecha_game_squad_instance_id"`
	GameTurnSheetID      string `db:"game_turn_sheet_id"`
}

func (r *MechaGameTurnSheet) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldMechaGameTurnSheetGameID] = r.GameID
	args[FieldMechaGameTurnSheetMechaGameSquadInstanceID] = r.MechaGameSquadInstanceID
	args[FieldMechaGameTurnSheetGameTurnSheetID] = r.GameTurnSheetID
	return args
}
