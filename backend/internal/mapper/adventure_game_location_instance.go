package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func AdventureGameLocationInstanceRequestToRecord(l logger.Logger, req *api.AdventureGameLocationInstanceRequest, rec *adventure_game_record.AdventureGameLocationInstance) (*adventure_game_record.AdventureGameLocationInstance, error) {
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

func AdventureGameLocationInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationInstance) (api.AdventureGameLocationInstanceResponseData, error) {
	l.Debug("mapping adventure_game_location_instance record to response data")
	data := api.AdventureGameLocationInstanceResponseData{
		ID:             rec.ID,
		GameID:         rec.GameID,
		GameInstanceID: rec.GameInstanceID,
		GameLocationID: rec.AdventureGameLocationID,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:      nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameLocationInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationInstance) (api.AdventureGameLocationInstanceResponse, error) {
	data, err := AdventureGameLocationInstanceRecordToResponseData(l, rec)
	if err != nil {
		return api.AdventureGameLocationInstanceResponse{}, err
	}
	return api.AdventureGameLocationInstanceResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationInstance) (api.AdventureGameLocationInstanceCollectionResponse, error) {
	l.Debug("mapping adventure_game_location_instance records to collection response")
	data := []*api.AdventureGameLocationInstanceResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationInstanceRecordToResponseData(l, rec)
		if err != nil {
			return api.AdventureGameLocationInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AdventureGameLocationInstanceCollectionResponse{
		Data: data,
	}, nil
}
