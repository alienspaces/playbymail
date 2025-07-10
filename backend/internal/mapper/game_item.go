package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func GameItemRecordToResponseData(l logger.Logger, rec *record.GameItem) (schema.GameItemResponseData, error) {
	l.Debug("mapping game_item record to response data")
	data := schema.GameItemResponseData{
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

func GameItemRecordToResponse(l logger.Logger, rec *record.GameItem) (schema.GameItemResponse, error) {
	data, err := GameItemRecordToResponseData(l, rec)
	if err != nil {
		return schema.GameItemResponse{}, err
	}
	return schema.GameItemResponse{
		Data: &data,
	}, nil
}

func GameItemRecordsToCollectionResponse(l logger.Logger, recs []*record.GameItem) (schema.GameItemCollectionResponse, error) {
	var data []*schema.GameItemResponseData
	for _, rec := range recs {
		d, err := GameItemRecordToResponseData(l, rec)
		if err != nil {
			return schema.GameItemCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.GameItemCollectionResponse{
		Data: data,
	}, nil
}

func GameItemRequestToRecord(l logger.Logger, req *schema.GameItemRequest, rec *record.GameItem) (*record.GameItem, error) {
	if rec == nil {
		rec = &record.GameItem{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping game_item request to record")
	rec.GameID = req.GameID
	rec.Name = req.Name
	rec.Description = req.Description
	return rec, nil
}
