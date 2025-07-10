package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func GameLocationRecordToResponseData(l logger.Logger, rec *record.GameLocation) (schema.GameLocationResponseData, error) {
	l.Debug("mapping game_location record to response data")
	data := schema.GameLocationResponseData{
		ID:          rec.ID,
		GameID:      rec.GameID,
		Name:        rec.Name,
		Description: rec.Description,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:   nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func GameLocationRecordToResponse(l logger.Logger, rec *record.GameLocation) (schema.GameLocationResponse, error) {
	data, err := GameLocationRecordToResponseData(l, rec)
	if err != nil {
		return schema.GameLocationResponse{}, err
	}
	return schema.GameLocationResponse{
		Data: &data,
	}, nil
}

func GameLocationRecordsToCollectionResponse(l logger.Logger, recs []*record.GameLocation) (schema.GameLocationCollectionResponse, error) {
	var data []*schema.GameLocationResponseData
	for _, rec := range recs {
		d, err := GameLocationRecordToResponseData(l, rec)
		if err != nil {
			return schema.GameLocationCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.GameLocationCollectionResponse{
		Data: data,
	}, nil
}

// GameLocationRequestToRecord maps a GameLocationRequest to a record.GameLocation
func GameLocationRequestToRecord(l logger.Logger, req *schema.GameLocationRequest, rec *record.GameLocation) (*record.GameLocation, error) {
	if rec == nil {
		rec = &record.GameLocation{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping game_location request to record")
	rec.GameID = req.GameID
	rec.Name = req.Name
	rec.Description = req.Description
	return rec, nil
}
