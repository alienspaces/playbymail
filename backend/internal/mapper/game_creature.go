package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func GameCreatureRecordToResponseData(l logger.Logger, rec *record.GameCreature) (schema.GameCreatureResponseData, error) {
	l.Debug("mapping game_creature record to response data")
	data := schema.GameCreatureResponseData{
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

func GameCreatureRecordToResponse(l logger.Logger, rec *record.GameCreature) (schema.GameCreatureResponse, error) {
	data, err := GameCreatureRecordToResponseData(l, rec)
	if err != nil {
		return schema.GameCreatureResponse{}, err
	}
	return schema.GameCreatureResponse{
		Data: &data,
	}, nil
}

func GameCreatureRecordsToCollectionResponse(l logger.Logger, recs []*record.GameCreature) (schema.GameCreatureCollectionResponse, error) {
	var data []*schema.GameCreatureResponseData
	for _, rec := range recs {
		d, err := GameCreatureRecordToResponseData(l, rec)
		if err != nil {
			return schema.GameCreatureCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.GameCreatureCollectionResponse{
		Data: data,
	}, nil
}

func GameCreatureRequestToRecord(l logger.Logger, req *schema.GameCreatureRequest, rec *record.GameCreature) (*record.GameCreature, error) {
	if rec == nil {
		rec = &record.GameCreature{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping game_creature request to record")
	rec.GameID = req.GameID
	rec.Name = req.Name
	rec.Description = req.Description
	return rec, nil
}
