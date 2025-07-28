package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func AdventureGameItemPlacementRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameItemPlacement) (api.AdventureGameItemPlacementResponseData, error) {
	l.Debug("mapping adventure_game_item_placement record to response data")
	data := api.AdventureGameItemPlacementResponseData{
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

func AdventureGameItemPlacementRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameItemPlacement) (api.AdventureGameItemPlacementResponse, error) {
	data, err := AdventureGameItemPlacementRecordToResponseData(l, rec)
	if err != nil {
		return api.AdventureGameItemPlacementResponse{}, err
	}
	return api.AdventureGameItemPlacementResponse{
		Data: &data,
	}, nil
}

func AdventureGameItemPlacementRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameItemPlacement) (api.AdventureGameItemPlacementCollectionResponse, error) {
	data := []*api.AdventureGameItemPlacementResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameItemPlacementRecordToResponseData(l, rec)
		if err != nil {
			return api.AdventureGameItemPlacementCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AdventureGameItemPlacementCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameItemPlacementRequestToRecord(l logger.Logger, req *api.AdventureGameItemPlacementRequest, rec *adventure_game_record.AdventureGameItemPlacement) (*adventure_game_record.AdventureGameItemPlacement, error) {
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
