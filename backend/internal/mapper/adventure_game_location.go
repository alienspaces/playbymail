package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AdventureGameLocationRecordToResponseData(l logger.Logger, rec *record.AdventureGameLocation) (schema.AdventureGameLocationResponseData, error) {
	l.Debug("mapping adventure_game_location record to response data")
	data := schema.AdventureGameLocationResponseData{
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

func AdventureGameLocationRecordToResponse(l logger.Logger, rec *record.AdventureGameLocation) (schema.AdventureGameLocationResponse, error) {
	data, err := AdventureGameLocationRecordToResponseData(l, rec)
	if err != nil {
		return schema.AdventureGameLocationResponse{}, err
	}
	return schema.AdventureGameLocationResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationRecordsToCollectionResponse(l logger.Logger, recs []*record.AdventureGameLocation) (schema.AdventureGameLocationCollectionResponse, error) {
	data := []*schema.AdventureGameLocationResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationRecordToResponseData(l, rec)
		if err != nil {
			return schema.AdventureGameLocationCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.AdventureGameLocationCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationRequestToRecord(l logger.Logger, req *schema.AdventureGameLocationRequest, rec *record.AdventureGameLocation) (*record.AdventureGameLocation, error) {
	if rec == nil {
		rec = &record.AdventureGameLocation{}
	}
	if req == nil {
		return nil, nil
	}

	l.Debug("mapping adventure_game_location request to record")

	rec.Name = req.Name
	rec.Description = req.Description

	return rec, nil
}
