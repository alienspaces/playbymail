package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func LocationRecordToResponseData(l logger.Logger, rec *record.Location) (schema.LocationResponseData, error) {
	l.Debug("mapping location record to response data")
	data := schema.LocationResponseData{
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

func LocationRecordToResponse(rec *record.Location) *schema.LocationResponse {
	data, _ := LocationRecordToResponseData(nil, rec)
	return &schema.LocationResponse{
		LocationResponseData: &data,
	}
}

// LocationRequestToRecord maps a LocationRequest to a record.Location
func LocationRequestToRecord(l logger.Logger, req *schema.LocationRequest, rec *record.Location) (*record.Location, error) {
	if rec == nil {
		rec = &record.Location{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping location request to record")
	rec.GameID = req.GameID
	rec.Name = req.Name
	rec.Description = req.Description
	return rec, nil
}
