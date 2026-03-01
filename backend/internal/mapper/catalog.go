package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/catalog_schema"
)

// CatalogGameInstanceRecordToData maps a game instance record to catalog instance data.
func CatalogGameInstanceRecordToData(rec *game_record.GameInstance, playerCount int) *catalog_schema.CatalogGameInstanceData {
	return &catalog_schema.CatalogGameInstanceData{
		ID:                    rec.ID,
		RequiredPlayerCount:   rec.RequiredPlayerCount,
		PlayerCount:           playerCount,
		DeliveryPhysicalPost:  rec.DeliveryPhysicalPost,
		DeliveryPhysicalLocal: rec.DeliveryPhysicalLocal,
		DeliveryEmail:         rec.DeliveryEmail,
	}
}

// CatalogGameRecordToResponseData maps a game record and its available instances to catalog response data.
func CatalogGameRecordToResponseData(l logger.Logger, rec *game_record.Game, instances []*catalog_schema.CatalogGameInstanceData) *catalog_schema.CatalogGameResponseData {
	l.Debug("mapping game record to catalog response data")
	if instances == nil {
		instances = []*catalog_schema.CatalogGameInstanceData{}
	}
	return &catalog_schema.CatalogGameResponseData{
		ID:                 rec.ID,
		Name:               rec.Name,
		Description:        rec.Description,
		GameType:           rec.GameType,
		TurnDurationHours:  rec.TurnDurationHours,
		AvailableInstances: instances,
	}
}
