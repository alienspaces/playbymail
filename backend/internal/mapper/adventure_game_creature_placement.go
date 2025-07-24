package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AdventureGameCreaturePlacementRecordToResponseData(l logger.Logger, rec *record.AdventureGameCreaturePlacement) (schema.AdventureGameCreaturePlacementResponseData, error) {
	l.Debug("mapping adventure_game_creature_placement record to response data")
	data := schema.AdventureGameCreaturePlacementResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		AdventureGameCreatureID: rec.AdventureGameCreatureID,
		AdventureGameLocationID: rec.AdventureGameLocationID,
		InitialCount:            rec.InitialCount,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameCreaturePlacementRecordsToCollectionResponse(l logger.Logger, recs []*record.AdventureGameCreaturePlacement) (*schema.AdventureGameCreaturePlacementCollectionResponse, error) {
	l.Debug("mapping adventure_game_creature_placement records to collection response")
	data := []*schema.AdventureGameCreaturePlacementResponseData{}
	for _, rec := range recs {
		item, err := AdventureGameCreaturePlacementRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, &item)
	}
	return &schema.AdventureGameCreaturePlacementCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameCreaturePlacementRecordToResponse(l logger.Logger, rec *record.AdventureGameCreaturePlacement) (*schema.AdventureGameCreaturePlacementResponse, error) {
	l.Debug("mapping adventure_game_creature_placement record to response")
	data, err := AdventureGameCreaturePlacementRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &schema.AdventureGameCreaturePlacementResponse{
		Data: &data,
	}, nil
}

func AdventureGameCreaturePlacementRequestToRecord(l logger.Logger, req *schema.AdventureGameCreaturePlacementRequest, rec *record.AdventureGameCreaturePlacement) (*record.AdventureGameCreaturePlacement, error) {
	if rec == nil {
		rec = &record.AdventureGameCreaturePlacement{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping adventure_game_creature_placement request to record")
	rec.AdventureGameCreatureID = req.AdventureGameCreatureID
	rec.AdventureGameLocationID = req.AdventureGameLocationID
	rec.InitialCount = req.InitialCount
	return rec, nil
}
