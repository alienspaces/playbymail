package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func CatalogGameInstanceViewRecordToResponseData(l logger.Logger, rec *game_record.CatalogGameInstanceView) (*game_schema.CatalogGameInstanceResponseData, error) {
	l.Debug("mapping catalog_game_instance_view record to response data")

	return &game_schema.CatalogGameInstanceResponseData{
		GameInstanceID:        rec.GameInstanceID,
		GameID:                rec.GameID,
		GameName:              rec.GameName,
		GameType:              rec.GameType,
		GameDescription:       rec.GameDescription,
		TurnDurationHours:     rec.TurnDurationHours,
		GameSubscriptionID:    rec.GameSubscriptionID,
		AccountName:           rec.AccountName,
		RequiredPlayerCount:   rec.RequiredPlayerCount,
		PlayerCount:           rec.PlayerCount,
		RemainingCapacity:     rec.RemainingCapacity,
		DeliveryEmail:         rec.DeliveryEmail,
		DeliveryPhysicalPost:  rec.DeliveryPhysicalPost,
		DeliveryPhysicalLocal: rec.DeliveryPhysicalLocal,
		IsClosedTesting:       rec.IsClosedTesting,
		CreatedAt:             rec.CreatedAt,
	}, nil
}

func CatalogGameInstanceViewRecordsToCollectionResponse(l logger.Logger, recs []*game_record.CatalogGameInstanceView) (game_schema.CatalogGameInstanceCollectionResponse, error) {
	l.Debug("mapping catalog_game_instance_view records to collection response")
	data := []*game_schema.CatalogGameInstanceResponseData{}
	for _, rec := range recs {
		d, err := CatalogGameInstanceViewRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.CatalogGameInstanceCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.CatalogGameInstanceCollectionResponse{
		Data: data,
	}, nil
}
