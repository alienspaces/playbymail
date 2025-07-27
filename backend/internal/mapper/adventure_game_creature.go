package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AdventureGameCreatureRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCreature) (schema.AdventureGameCreatureResponseData, error) {
	l.Debug("mapping adventure_game_creature record to response data")
	data := schema.AdventureGameCreatureResponseData{
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

func AdventureGameCreatureRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCreature) (schema.AdventureGameCreatureResponse, error) {
	data, err := AdventureGameCreatureRecordToResponseData(l, rec)
	if err != nil {
		return schema.AdventureGameCreatureResponse{}, err
	}
	return schema.AdventureGameCreatureResponse{
		Data: &data,
	}, nil
}

func AdventureGameCreatureRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCreature) (schema.AdventureGameCreatureCollectionResponse, error) {
	data := []*schema.AdventureGameCreatureResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameCreatureRecordToResponseData(l, rec)
		if err != nil {
			return schema.AdventureGameCreatureCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.AdventureGameCreatureCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameCreatureRequestToRecord(l logger.Logger, req *schema.AdventureGameCreatureRequest, rec *adventure_game_record.AdventureGameCreature) (*adventure_game_record.AdventureGameCreature, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameCreature{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping adventure_game_creature request to record")
	rec.GameID = req.GameID
	rec.Name = req.Name
	rec.Description = req.Description
	return rec, nil
}
