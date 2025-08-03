package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameLocationInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationInstance) (adventure_game_schema.AdventureGameLocationInstanceResponseData, error) {
	l.Debug("mapping adventure_game_location_instance record to response data")
	data := adventure_game_schema.AdventureGameLocationInstanceResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		GameInstanceID:          rec.GameInstanceID,
		AdventureGameLocationID: rec.AdventureGameLocationID,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameLocationInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationInstance) (adventure_game_schema.AdventureGameLocationInstanceResponse, error) {
	data, err := AdventureGameLocationInstanceRecordToResponseData(l, rec)
	if err != nil {
		return adventure_game_schema.AdventureGameLocationInstanceResponse{}, err
	}
	return adventure_game_schema.AdventureGameLocationInstanceResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationInstance) (adventure_game_schema.AdventureGameLocationInstanceCollectionResponse, error) {
	l.Debug("mapping adventure_game_location_instance records to collection response")
	data := []*adventure_game_schema.AdventureGameLocationInstanceResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationInstanceRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameLocationInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return adventure_game_schema.AdventureGameLocationInstanceCollectionResponse{
		Data: data,
	}, nil
}
