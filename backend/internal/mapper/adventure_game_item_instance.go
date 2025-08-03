package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameItemInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameItemInstance) (adventure_game_schema.AdventureGameItemInstanceResponseData, error) {
	l.Debug("mapping adventure_game_item_instance record to response data")
	data := adventure_game_schema.AdventureGameItemInstanceResponseData{
		ID:                               rec.ID,
		GameID:                           rec.GameID,
		GameInstanceID:                   rec.GameInstanceID,
		AdventureGameItemID:              rec.AdventureGameItemID,
		AdventureGameCharacterInstanceID: nullstring.ToString(rec.AdventureGameCharacterInstanceID),
		AdventureGameCreatureInstanceID:  nullstring.ToString(rec.AdventureGameCreatureInstanceID),
		AdventureGameLocationInstanceID:  nullstring.ToString(rec.AdventureGameLocationInstanceID),
		IsEquipped:                       rec.IsEquipped,
		IsUsed:                           rec.IsUsed,
		UsesRemaining:                    rec.UsesRemaining,
		CreatedAt:                        rec.CreatedAt,
		UpdatedAt:                        nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:                        nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameItemInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameItemInstance) (adventure_game_schema.AdventureGameItemInstanceResponse, error) {
	data, err := AdventureGameItemInstanceRecordToResponseData(l, rec)
	if err != nil {
		return adventure_game_schema.AdventureGameItemInstanceResponse{}, err
	}
	return adventure_game_schema.AdventureGameItemInstanceResponse{
		Data: &data,
	}, nil
}

func AdventureGameItemInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameItemInstance) (adventure_game_schema.AdventureGameItemInstanceCollectionResponse, error) {
	data := []*adventure_game_schema.AdventureGameItemInstanceResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameItemInstanceRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameItemInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return adventure_game_schema.AdventureGameItemInstanceCollectionResponse{
		Data: data,
	}, nil
}
