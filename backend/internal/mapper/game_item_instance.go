package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func GameItemInstanceRequestToRecord(l logger.Logger, req *schema.GameItemInstanceRequest, rec *record.GameItemInstance) (*record.GameItemInstance, error) {
	if rec == nil {
		rec = &record.GameItemInstance{}
	}
	if req == nil {
		return nil, nil
	}
	l.Debug("mapping game_item_instance request to record")
	rec.GameID = req.GameID
	rec.GameItemID = req.GameItemID
	rec.GameInstanceID = req.GameInstanceID
	rec.GameLocationInstanceID = nullstring.FromString(req.GameLocationInstanceID)
	rec.GameCharacterInstanceID = nullstring.FromString(req.GameCharacterInstanceID)
	rec.GameCreatureInstanceID = nullstring.FromString(req.GameCreatureInstanceID)
	rec.IsEquipped = req.IsEquipped
	rec.IsUsed = req.IsUsed
	rec.UsesRemaining = req.UsesRemaining
	return rec, nil
}

func GameItemInstanceRecordToResponseData(l logger.Logger, rec *record.GameItemInstance) (schema.GameItemInstanceResponseData, error) {
	l.Debug("mapping game_item_instance record to response data")
	data := schema.GameItemInstanceResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		GameItemID:              rec.GameItemID,
		GameInstanceID:          rec.GameInstanceID,
		GameLocationInstanceID:  nullstring.ToString(rec.GameLocationInstanceID),
		GameCharacterInstanceID: nullstring.ToString(rec.GameCharacterInstanceID),
		GameCreatureInstanceID:  nullstring.ToString(rec.GameCreatureInstanceID),
		IsEquipped:              rec.IsEquipped,
		IsUsed:                  rec.IsUsed,
		UsesRemaining:           rec.UsesRemaining,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func GameItemInstanceRecordToResponse(l logger.Logger, rec *record.GameItemInstance) (schema.GameItemInstanceResponse, error) {
	data, err := GameItemInstanceRecordToResponseData(l, rec)
	if err != nil {
		return schema.GameItemInstanceResponse{}, err
	}
	return schema.GameItemInstanceResponse{
		Data: &data,
	}, nil
}

func GameItemInstanceRecordsToCollectionResponse(l logger.Logger, recs []*record.GameItemInstance) (schema.GameItemInstanceCollectionResponse, error) {
	var data []*schema.GameItemInstanceResponseData
	for _, rec := range recs {
		d, err := GameItemInstanceRecordToResponseData(l, rec)
		if err != nil {
			return schema.GameItemInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.GameItemInstanceCollectionResponse{
		Data: data,
	}, nil
}
