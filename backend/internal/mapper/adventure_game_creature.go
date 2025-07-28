package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func AdventureGameCreatureRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCreature) (api.AdventureGameCreatureResponseData, error) {
	l.Debug("mapping adventure_game_creature record to response data")
	data := api.AdventureGameCreatureResponseData{
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

func AdventureGameCreatureRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCreature) (api.AdventureGameCreatureResponse, error) {
	data, err := AdventureGameCreatureRecordToResponseData(l, rec)
	if err != nil {
		return api.AdventureGameCreatureResponse{}, err
	}
	return api.AdventureGameCreatureResponse{
		Data: &data,
	}, nil
}

func AdventureGameCreatureRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCreature) (api.AdventureGameCreatureCollectionResponse, error) {
	data := []*api.AdventureGameCreatureResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameCreatureRecordToResponseData(l, rec)
		if err != nil {
			return api.AdventureGameCreatureCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AdventureGameCreatureCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameCreatureRequestToRecord(l logger.Logger, req *api.AdventureGameCreatureRequest, rec *adventure_game_record.AdventureGameCreature) (*adventure_game_record.AdventureGameCreature, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameCreature{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping adventure_game_creature request to record")
	rec.Name = req.Name
	rec.Description = req.Description
	return rec, nil
}
