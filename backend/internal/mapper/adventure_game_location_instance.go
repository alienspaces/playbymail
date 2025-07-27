package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AdventureGameLocationInstanceRequestToRecord(l logger.Logger, req *schema.AdventureGameLocationInstanceRequest, rec *adventure_game_record.AdventureGameLocationInstance) (*adventure_game_record.AdventureGameLocationInstance, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocationInstance{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping adventure_game_location_instance request to record")

	rec.AdventureGameLocationID = req.GameLocationID

	return rec, nil
}

func AdventureGameLocationInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationInstance) (schema.AdventureGameLocationInstanceResponseData, error) {
	l.Debug("mapping adventure_game_location_instance record to response data")
	data := schema.AdventureGameLocationInstanceResponseData{
		ID:             rec.ID,
		GameID:         rec.GameID,
		GameInstanceID: rec.AdventureGameInstanceID,
		GameLocationID: rec.AdventureGameLocationID,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:      nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameLocationInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationInstance) (schema.AdventureGameLocationInstanceResponse, error) {
	data, err := AdventureGameLocationInstanceRecordToResponseData(l, rec)
	if err != nil {
		return schema.AdventureGameLocationInstanceResponse{}, err
	}
	return schema.AdventureGameLocationInstanceResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationInstance) (schema.AdventureGameLocationInstanceCollectionResponse, error) {
	l.Debug("mapping adventure_game_location_instance records to collection response")
	data := []*schema.AdventureGameLocationInstanceResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationInstanceRecordToResponseData(l, rec)
		if err != nil {
			return schema.AdventureGameLocationInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.AdventureGameLocationInstanceCollectionResponse{
		Data: data,
	}, nil
}
