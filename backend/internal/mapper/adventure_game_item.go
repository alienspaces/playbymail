package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameItemRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameItem) (adventure_game_schema.AdventureGameItemResponseData, error) {
	l.Debug("mapping adventure_game_item record to response data")
	data := adventure_game_schema.AdventureGameItemResponseData{
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

func AdventureGameItemRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameItem) (adventure_game_schema.AdventureGameItemResponse, error) {
	data, err := AdventureGameItemRecordToResponseData(l, rec)
	if err != nil {
		return adventure_game_schema.AdventureGameItemResponse{}, err
	}
	return adventure_game_schema.AdventureGameItemResponse{
		Data: &data,
	}, nil
}

func AdventureGameItemRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameItem) (adventure_game_schema.AdventureGameItemCollectionResponse, error) {
	data := []*adventure_game_schema.AdventureGameItemResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameItemRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameItemCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return adventure_game_schema.AdventureGameItemCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameItemRequestToRecord(l logger.Logger, req *adventure_game_schema.AdventureGameItemRequest, rec *adventure_game_record.AdventureGameItem) (*adventure_game_record.AdventureGameItem, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameItem{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping adventure_game_item request to record")
	rec.Name = req.Name
	rec.Description = req.Description
	return rec, nil
}
