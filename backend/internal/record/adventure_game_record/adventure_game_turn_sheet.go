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
	AdventureSheetTypeLocationChoice = "location_choice"
)

// GetAdventureGameSheetTypes returns the sheet types for adventure games
var AdventureGameSheetTypes = set.New(
	AdventureSheetTypeLocationChoice,
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
