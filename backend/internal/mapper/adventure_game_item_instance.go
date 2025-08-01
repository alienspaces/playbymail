package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func AdventureGameItemInstanceRequestToRecord(l logger.Logger, req *api.AdventureGameItemInstanceRequest, rec *adventure_game_record.AdventureGameItemInstance) (*adventure_game_record.AdventureGameItemInstance, error) {
	if rec == nil {
		rec = &adventure_game_record.AdventureGameItemInstance{}
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

func AdventureGameItemInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameItemInstance) (api.AdventureGameItemInstanceResponseData, error) {
	l.Debug("mapping adventure_game_item_instance record to response data")
	data := api.AdventureGameItemInstanceResponseData{
		ID:                      rec.ID,
		GameID:                  rec.GameID,
		GameItemID:              rec.AdventureGameItemID,
		GameInstanceID:          rec.GameInstanceID,
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

func AdventureGameItemInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameItemInstance) (api.AdventureGameItemInstanceResponse, error) {
	data, err := AdventureGameItemInstanceRecordToResponseData(l, rec)
	if err != nil {
		return api.AdventureGameItemInstanceResponse{}, err
	}
	return api.AdventureGameItemInstanceResponse{
		Data: &data,
	}, nil
}

func AdventureGameItemInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameItemInstance) (api.AdventureGameItemInstanceCollectionResponse, error) {
	data := []*api.AdventureGameItemInstanceResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameItemInstanceRecordToResponseData(l, rec)
		if err != nil {
			return api.AdventureGameItemInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AdventureGameItemInstanceCollectionResponse{
		Data: data,
	}, nil
}
