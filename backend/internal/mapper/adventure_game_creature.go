package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameCreatureRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCreature) (adventure_game_schema.AdventureGameCreatureResponseData, error) {
	l.Debug("mapping adventure_game_creature record to response data")
	data := adventure_game_schema.AdventureGameCreatureResponseData{
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

func AdventureGameCreatureRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCreature) (adventure_game_schema.AdventureGameCreatureResponse, error) {
	data, err := AdventureGameCreatureRecordToResponseData(l, rec)
	if err != nil {
		return adventure_game_schema.AdventureGameCreatureResponse{}, err
	}
	return adventure_game_schema.AdventureGameCreatureResponse{
		Data: &data,
	}, nil
}

func AdventureGameCreatureRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCreature) (adventure_game_schema.AdventureGameCreatureCollectionResponse, error) {
	data := []*adventure_game_schema.AdventureGameCreatureResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameCreatureRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameCreatureCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return adventure_game_schema.AdventureGameCreatureCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameCreatureRequestToRecord(l logger.Logger, req *adventure_game_schema.AdventureGameCreatureRequest, rec *adventure_game_record.AdventureGameCreature) (*adventure_game_record.AdventureGameCreature, error) {
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
