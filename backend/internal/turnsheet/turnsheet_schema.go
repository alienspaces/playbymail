package turnsheet

import (
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/internal/record/mecha_record"
)

// ScannedDataSchemaLocation is the subdirectory (relative to SchemaPath) where scanned_data
// JSON schemas for turn sheet types live, e.g. turnsheet/adventure_game.
const ScannedDataSchemaLocation = "turnsheet/adventure_game"

// MechaScannedDataSchemaLocation is the subdirectory for mecha scanned_data schemas.
const MechaScannedDataSchemaLocation = "turnsheet/mecha"

// scannedDataSchemaLocationBySheetType maps turn sheet types to their schema subdirectory.
var scannedDataSchemaLocationBySheetType = map[string]string{
	adventure_game_record.AdventureGameTurnSheetTypeLocationChoice:      ScannedDataSchemaLocation,
	adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement: ScannedDataSchemaLocation,
	adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter:   ScannedDataSchemaLocation,
	mecha_record.MechaTurnSheetTypeOrders:          MechaScannedDataSchemaLocation,
	mecha_record.MechaTurnSheetTypeLanceManagement: MechaScannedDataSchemaLocation,
}

// scannedDataSchemaNameBySheetType maps turn sheet types to their scanned_data schema filename.
// When a schema exists, callers (e.g. the save handler) can validate request scanned_data against it.
var scannedDataSchemaNameBySheetType = map[string]string{
	adventure_game_record.AdventureGameTurnSheetTypeLocationChoice:      LocationChoiceScannedDataSchemaName,
	adventure_game_record.AdventureGameTurnSheetTypeInventoryManagement: InventoryManagementScannedDataSchemaName,
	adventure_game_record.AdventureGameTurnSheetTypeCreatureEncounter:   MonsterEncounterScannedDataSchemaName,
	mecha_record.MechaTurnSheetTypeOrders:          OrdersScannedDataSchemaName,
	mecha_record.MechaTurnSheetTypeLanceManagement: LanceManagementScannedDataSchemaName,
}

// ScannedDataSchemaName returns the JSON schema filename for the given sheet type's scanned_data,
// or empty string if the sheet type has no schema. The schema lives under ScannedDataSchemaLocation.
func ScannedDataSchemaName(sheetType string) string {
	return scannedDataSchemaNameBySheetType[sheetType]
}

// ScannedDataSchemaLocationForType returns the schema subdirectory for the given sheet type's
// scanned_data, or the default ScannedDataSchemaLocation if the sheet type has no specific location.
func ScannedDataSchemaLocationForType(sheetType string) string {
	if loc, ok := scannedDataSchemaLocationBySheetType[sheetType]; ok {
		return loc
	}
	return ScannedDataSchemaLocation
}
