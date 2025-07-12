package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func GameLocationInstanceRequestToRecord(l logger.Logger, req *schema.GameLocationInstanceRequest, rec *record.GameLocationInstance) (*record.GameLocationInstance, error) {
	if rec == nil {
		rec = &record.GameLocationInstance{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping game_location_instance request to record")
	rec.GameID = req.GameID
	rec.GameInstanceID = req.GameInstanceID
	rec.GameLocationID = req.GameLocationID
	return rec, nil
}

func GameLocationInstanceRecordToResponseData(l logger.Logger, rec *record.GameLocationInstance) (schema.GameLocationInstanceResponseData, error) {
	l.Debug("mapping game_location_instance record to response data")
	data := schema.GameLocationInstanceResponseData{
		ID:             rec.ID,
		GameID:         rec.GameID,
		GameInstanceID: rec.GameInstanceID,
		GameLocationID: rec.GameLocationID,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:      nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func GameLocationInstanceRecordToResponse(l logger.Logger, rec *record.GameLocationInstance) (schema.GameLocationInstanceResponse, error) {
	data, err := GameLocationInstanceRecordToResponseData(l, rec)
	if err != nil {
		return schema.GameLocationInstanceResponse{}, err
	}
	return schema.GameLocationInstanceResponse{
		Data: &data,
	}, nil
}

func GameLocationInstanceRecordsToCollectionResponse(l logger.Logger, recs []*record.GameLocationInstance) (schema.GameLocationInstanceCollectionResponse, error) {
	l.Debug("mapping game_location_instance records to collection response")
	var data []*schema.GameLocationInstanceResponseData
	for _, rec := range recs {
		d, err := GameLocationInstanceRecordToResponseData(l, rec)
		if err != nil {
			return schema.GameLocationInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.GameLocationInstanceCollectionResponse{
		Data: data,
	}, nil
}
