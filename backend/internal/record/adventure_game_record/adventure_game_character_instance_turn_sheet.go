package adventure_game_record

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"gitlab.com/alienspaces/playbymail/core/record"
)

// AdventureGameCharacterInstanceTurnSheet
const (
	TableAdventureGameCharacterInstanceTurnSheet string = "adventure_game_character_instance_turn_sheet"
)

const (
	FieldAdventureGameCharacterInstanceTurnSheetID                               string = "id"
	FieldAdventureGameCharacterInstanceTurnSheetAdventureGameCharacterInstanceID string = "adventure_game_character_instance_id"
	FieldAdventureGameCharacterInstanceTurnSheetGameTurnSheetID                  string = "game_turn_sheet_id"
	FieldAdventureGameCharacterInstanceTurnSheetCreatedAt                        string = "created_at"
)

type AdventureGameCharacterInstanceTurnSheet struct {
	record.Record
	AdventureGameCharacterInstanceID string       `db:"adventure_game_character_instance_id"`
	GameTurnSheetID                  string       `db:"game_turn_sheet_id"`
	CreatedAt                        sql.NullTime `db:"created_at"`
}

func (r *AdventureGameCharacterInstanceTurnSheet) ToNamedArgs() pgx.NamedArgs {
	args := r.Record.ToNamedArgs()
	args[FieldAdventureGameCharacterInstanceTurnSheetAdventureGameCharacterInstanceID] = r.AdventureGameCharacterInstanceID
	args[FieldAdventureGameCharacterInstanceTurnSheetGameTurnSheetID] = r.GameTurnSheetID
	args[FieldAdventureGameCharacterInstanceTurnSheetCreatedAt] = r.CreatedAt
	return args
}
