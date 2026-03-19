package turnsheet

import (
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

// ScannedDataSchemaLocation is the subdirectory (relative to SchemaPath) where scanned_data
// JSON schemas for turn sheet types live, e.g. turnsheet/adventure_game.
const ScannedDataSchemaLocation = "turnsheet/adventure_game"

// scannedDataSchemaNameBySheetType maps turn sheet types to their scanned_data schema filename.
// When a schema exists, callers (e.g. the save handler) can validate request scanned_data against it.
var scannedDataSchemaNameBySheetType = map[string]string{
	adventure_game_record.AdventureGameTurnSheetTypeLocationChoice:      LocationChoiceScannedDataSchemaName,
	adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement: InventoryManagementScannedDataSchemaName,
	adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter:   MonsterEncounterScannedDataSchemaName,
}

// ScannedDataSchemaName returns the JSON schema filename for the given sheet type's scanned_data,
// or empty string if the sheet type has no schema. The schema lives under ScannedDataSchemaLocation.
func ScannedDataSchemaName(sheetType string) string {
	return scannedDataSchemaNameBySheetType[sheetType]
}
