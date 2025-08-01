package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func AdventureGameLocationRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocation) (api.AdventureGameLocationResponseData, error) {
	l.Debug("mapping adventure_game_location record to response data")
	data := api.AdventureGameLocationResponseData{
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

func AdventureGameLocationRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocation) (api.AdventureGameLocationResponse, error) {
	data, err := AdventureGameLocationRecordToResponseData(l, rec)
	if err != nil {
		return api.AdventureGameLocationResponse{}, err
	}
	return api.AdventureGameLocationResponse{
		Data: &data,
	}, nil
}

func AdventureGameLocationRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocation) (api.AdventureGameLocationCollectionResponse, error) {
	data := []*api.AdventureGameLocationResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationRecordToResponseData(l, rec)
		if err != nil {
			return api.AdventureGameLocationCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AdventureGameLocationCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationRequestToRecord(l logger.Logger, req *api.AdventureGameLocationRequest, rec *adventure_game_record.AdventureGameLocation) (*adventure_game_record.AdventureGameLocation, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameLocation{}
	}
	if req == nil {
		return nil, nil
	}

	l.Debug("mapping adventure_game_location request to record")

	rec.Name = req.Name
	rec.Description = req.Description

	return rec, nil
}
