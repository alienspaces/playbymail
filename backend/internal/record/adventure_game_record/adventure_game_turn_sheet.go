package adventure_game_record

import (
	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/collection/set"
	"gitlab.com/alienspaces/playbymail/core/record"
)

const (
	TableAdventureGameTurnSheet = "adventure_game_turn_sheet"

	FieldAdventureGameTurnSheetID                               = "id"
	FieldAdventureGameTurnSheetGameID                           = "game_id"
	FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID = "adventure_game_character_instance_id"
	FieldAdventureGameTurnSheetGameTurnSheetID                  = "game_turn_sheet_id"
	FieldAdventureGameTurnSheetCreatedAt                        = "created_at"
	FieldAdventureGameTurnSheetUpdatedAt                        = "updated_at"
	FieldAdventureGameTurnSheetDeletedAt                        = "deleted_at"
)

// Turn sheet type constants for different game types
const (
	// Adventure game sheet types
	AdventureGameTurnSheetTypeLocationChoice      = "adventure_game_location_choice"
	AdventureGameTurnSheetTypeJoinGame            = "adventure_game_join_game"
	AdventureGameTurnSheetTypeInventoryManagement = "adventure_game_inventory_management"
	AdventureGameTurnSheetTypeCombat              = "adventure_game_combat"
	AdventureGameTurnSheetTypePuzzle              = "adventure_game_puzzle"
	AdventureGameTurnSheetTypeMonster             = "adventure_game_monster"
)

// AdventureGameTurnSheetProcessingOrder defines the order in which
// adventure game turn sheets are processed during turn resolution.
// Position determines SheetOrder (1-indexed).
// The join game sheet is excluded; it is handled through the
// subscription workflow, not turn processing.
var AdventureGameTurnSheetProcessingOrder = []string{
	AdventureGameTurnSheetTypeInventoryManagement, // 1 - manage items at current location before moving
	AdventureGameTurnSheetTypeLocationChoice,      // 2 - move to a new location
}

// AdventureGameSheetOrderForType returns the 1-indexed processing order
// for an adventure game turn sheet type. Returns 0 if the type is not
// in the processing order (e.g. join_game).
func AdventureGameSheetOrderForType(sheetType string) int {
	for i, t := range AdventureGameTurnSheetProcessingOrder {
		if t == sheetType {
			return i + 1
		}
	}
	return 0
}

// AdventureGameTurnSheetPresentationOrder defines the order in which
// adventure game turn sheets are presented to the player in the UI.
// This is intentionally different from the processing order: players
// see location choice first (where am I going?) then inventory (what do I take?).
var AdventureGameTurnSheetPresentationOrder = []string{
	AdventureGameTurnSheetTypeLocationChoice,      // 1 - shown first: choose where to go
	AdventureGameTurnSheetTypeInventoryManagement, // 2 - shown second: manage items for the turn
}

// AdventureGameSheetPresentationOrderForType returns the 1-indexed presentation
// order for an adventure game turn sheet type. Returns 0 if the type is not
// in the presentation order.
func AdventureGameSheetPresentationOrderForType(sheetType string) int {
	for i, t := range AdventureGameTurnSheetPresentationOrder {
		if t == sheetType {
			return i + 1
		}
	}
	return 0
}

// AdventureGameSheetTypes is the set of all adventure game sheet types
var AdventureGameSheetTypes = set.New(
	AdventureGameTurnSheetTypeLocationChoice,
	AdventureGameTurnSheetTypeJoinGame,
	AdventureGameTurnSheetTypeInventoryManagement,
	AdventureGameTurnSheetTypeCombat,
	AdventureGameTurnSheetTypePuzzle,
	AdventureGameTurnSheetTypeMonster,
)

type AdventureGameTurnSheet struct {
	record.Record
	GameID                           string `db:"game_id"`
	AdventureGameCharacterInstanceID string `db:"adventure_game_character_instance_id"`
	GameTurnSheetID                  string `db:"game_turn_sheet_id"`
}

func (r *AdventureGameTurnSheet) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameTurnSheetGameID] = r.GameID
	args[FieldAdventureGameTurnSheetAdventureGameCharacterInstanceID] = r.AdventureGameCharacterInstanceID
	args[FieldAdventureGameTurnSheetGameTurnSheetID] = r.GameTurnSheetID
	return args
}
