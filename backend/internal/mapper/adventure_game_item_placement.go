package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AdventureGameItemPlacementRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameItemPlacement) (schema.AdventureGameItemPlacementResponseData, error) {
	l.Debug("mapping adventure_game_item_placement record to response data")
	data := schema.AdventureGameItemPlacementResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		AdventureGameItemID:     rec.AdventureGameItemID,
		AdventureGameLocationID: rec.AdventureGameLocationID,
		InitialCount:            rec.InitialCount,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameItemPlacementRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameItemPlacement) (schema.AdventureGameItemPlacementResponse, error) {
	data, err := AdventureGameItemPlacementRecordToResponseData(l, rec)
	if err != nil {
		return schema.AdventureGameItemPlacementResponse{}, err
	}
	return schema.AdventureGameItemPlacementResponse{
		Data: &data,
	}, nil
}

func AdventureGameItemPlacementRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameItemPlacement) (schema.AdventureGameItemPlacementCollectionResponse, error) {
	data := []*schema.AdventureGameItemPlacementResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameItemPlacementRecordToResponseData(l, rec)
		if err != nil {
			return schema.AdventureGameItemPlacementCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.AdventureGameItemPlacementCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameItemPlacementRequestToRecord(l logger.Logger, req *schema.AdventureGameItemPlacementRequest, rec *adventure_game_record.AdventureGameItemPlacement) (*adventure_game_record.AdventureGameItemPlacement, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameItemPlacement{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping adventure_game_item_placement request to record")
	rec.AdventureGameItemID = req.AdventureGameItemID
	rec.AdventureGameLocationID = req.AdventureGameLocationID
	rec.InitialCount = req.InitialCount
	return rec, nil
}
