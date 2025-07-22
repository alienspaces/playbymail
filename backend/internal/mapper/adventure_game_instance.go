package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AdventureGameInstanceRecordToResponseData(l logger.Logger, rec *record.AdventureGameInstance) (schema.AdventureGameInstance, error) {
	data := schema.AdventureGameInstance{
		ID:        rec.ID,
		GameID:    rec.GameID,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt: nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameInstanceRecordToResponse(l logger.Logger, rec *record.AdventureGameInstance) (schema.AdventureGameInstanceResponse, error) {
	data, err := AdventureGameInstanceRecordToResponseData(l, rec)
	if err != nil {
		return schema.AdventureGameInstanceResponse{}, err
	}
	return schema.AdventureGameInstanceResponse{
		Data: &data,
	}, nil
}

func AdventureGameInstanceRecordsToCollectionResponse(l logger.Logger, recs []*record.AdventureGameInstance) (schema.AdventureGameInstanceCollectionResponse, error) {
	var data []*schema.AdventureGameInstance
	for _, rec := range recs {
		d, err := AdventureGameInstanceRecordToResponseData(l, rec)
		if err != nil {
			return schema.AdventureGameInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.AdventureGameInstanceCollectionResponse{
		Data: data,
	}, nil
}
