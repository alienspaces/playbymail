package mapper

import (
	"database/sql"

	"gitlab.com/alienspaces/playbymail/core/nullbool"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func ManagerGameInstanceViewRecordToResponseData(l logger.Logger, rec *game_record.ManagerGameInstanceView) (*game_schema.ManagerGameInstanceResponseData, error) {
	l.Debug("mapping manager_game_instance_view record to response data")

	return &game_schema.ManagerGameInstanceResponseData{
		GameID:                rec.GameID,
		GameName:              rec.GameName,
		GameType:              rec.GameType,
		GameDescription:       rec.GameDescription,
		GameSubscriptionID:    rec.GameSubscriptionID,
		GameInstanceID:        nullstring.ToStringPtr(rec.GameInstanceID),
		InstanceStatus:        nullstring.ToStringPtr(rec.InstanceStatus),
		CurrentTurn:           nullInt32ToIntPtr(rec.CurrentTurn),
		RequiredPlayerCount:   nullInt32ToIntPtr(rec.RequiredPlayerCount),
		DeliveryEmail:         nullbool.ToBoolPtr(rec.DeliveryEmail),
		DeliveryPhysicalPost:  nullbool.ToBoolPtr(rec.DeliveryPhysicalPost),
		DeliveryPhysicalLocal: nullbool.ToBoolPtr(rec.DeliveryPhysicalLocal),
		IsClosedTesting:       nullbool.ToBoolPtr(rec.IsClosedTesting),
		StartedAt:             nulltime.ToTimePtr(rec.StartedAt),
		NextTurnDueAt:         nulltime.ToTimePtr(rec.NextTurnDueAt),
		InstanceCreatedAt:     nulltime.ToTimePtr(rec.InstanceCreatedAt),
		CreatedAt:             rec.CreatedAt,
		UpdatedAt:             nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func ManagerGameInstanceViewRecordsToCollectionResponse(l logger.Logger, recs []*game_record.ManagerGameInstanceView) (game_schema.ManagerGameInstanceCollectionResponse, error) {
	l.Debug("mapping manager_game_instance_view records to collection response")
	data := []*game_schema.ManagerGameInstanceResponseData{}
	for _, rec := range recs {
		d, err := ManagerGameInstanceViewRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.ManagerGameInstanceCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.ManagerGameInstanceCollectionResponse{
		Data: data,
	}, nil
}

func nullInt32ToIntPtr(n sql.NullInt32) *int {
	if !n.Valid {
		return nil
	}
	v := int(n.Int32)
	return &v
}
