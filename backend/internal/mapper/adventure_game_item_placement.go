package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameItemPlacementRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameItemPlacement) (adventure_game_schema.AdventureGameItemPlacementResponseData, error) {
	data := adventure_game_schema.AdventureGameItemPlacementResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		AdventureGameItemID:     rec.AdventureGameItemID,     // Map old field name to new
		AdventureGameLocationID: rec.AdventureGameLocationID, // Map old field name to new
		InitialCount:            rec.InitialCount,            // Map old field name to new
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameItemPlacementRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameItemPlacement) (*adventure_game_schema.AdventureGameItemPlacementCollectionResponse, error) {
	data := []*adventure_game_schema.AdventureGameItemPlacementResponseData{}
	for _, rec := range recs {
		item, err := AdventureGameItemPlacementRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, &item)
	}
	return &adventure_game_schema.AdventureGameItemPlacementCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameItemPlacementRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameItemPlacement) (*adventure_game_schema.AdventureGameItemPlacementResponse, error) {
	data, err := AdventureGameItemPlacementRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameItemPlacementResponse{
		Data: &data,
	}, nil
}

func AdventureGameItemPlacementRequestToRecord(l logger.Logger, req *adventure_game_schema.AdventureGameItemPlacementRequest, rec *adventure_game_record.AdventureGameItemPlacement) (*adventure_game_record.AdventureGameItemPlacement, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameItemPlacement{}
	}
	if req == nil {
		return nil, nil
	}
	rec.AdventureGameItemID = req.AdventureGameItemID
	rec.AdventureGameLocationID = req.AdventureGameLocationID
	rec.InitialCount = req.InitialCount
	return rec, nil
}
