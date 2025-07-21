package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AdventureGameItemInstanceRequestToRecord(l logger.Logger, req *schema.AdventureGameItemInstanceRequest, rec *record.AdventureGameItemInstance) (*record.AdventureGameItemInstance, error) {
	if rec == nil {
		rec = &record.AdventureGameItemInstance{}
	}
	if req == nil {
		return nil, nil
	}

	l.Debug("mapping adventure_game_item_instance request to record")

	rec.AdventureGameItemID = req.GameItemID
	rec.AdventureGameLocationInstanceID = nullstring.FromString(req.GameLocationInstanceID)
	rec.AdventureGameCharacterInstanceID = nullstring.FromString(req.GameCharacterInstanceID)
	rec.AdventureGameCreatureInstanceID = nullstring.FromString(req.GameCreatureInstanceID)
	rec.IsEquipped = req.IsEquipped
	rec.IsUsed = req.IsUsed
	rec.UsesRemaining = req.UsesRemaining

	return rec, nil
}

func AdventureGameItemInstanceRecordToResponseData(l logger.Logger, rec *record.AdventureGameItemInstance) (schema.AdventureGameItemInstanceResponseData, error) {
	l.Debug("mapping adventure_game_item_instance record to response data")
	data := schema.AdventureGameItemInstanceResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		GameItemID:              rec.AdventureGameItemID,
		GameInstanceID:          rec.AdventureGameInstanceID,
		GameLocationInstanceID:  nullstring.ToString(rec.AdventureGameLocationInstanceID),
		GameCharacterInstanceID: nullstring.ToString(rec.AdventureGameCharacterInstanceID),
		GameCreatureInstanceID:  nullstring.ToString(rec.AdventureGameCreatureInstanceID),
		IsEquipped:              rec.IsEquipped,
		IsUsed:                  rec.IsUsed,
		UsesRemaining:           rec.UsesRemaining,
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameItemInstanceRecordToResponse(l logger.Logger, rec *record.AdventureGameItemInstance) (schema.AdventureGameItemInstanceResponse, error) {
	data, err := AdventureGameItemInstanceRecordToResponseData(l, rec)
	if err != nil {
		return schema.AdventureGameItemInstanceResponse{}, err
	}
	return schema.AdventureGameItemInstanceResponse{
		Data: &data,
	}, nil
}

func AdventureGameItemInstanceRecordsToCollectionResponse(l logger.Logger, recs []*record.AdventureGameItemInstance) (schema.AdventureGameItemInstanceCollectionResponse, error) {
	var data []*schema.AdventureGameItemInstanceResponseData
	for _, rec := range recs {
		d, err := AdventureGameItemInstanceRecordToResponseData(l, rec)
		if err != nil {
			return schema.AdventureGameItemInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.AdventureGameItemInstanceCollectionResponse{
		Data: data,
	}, nil
}
