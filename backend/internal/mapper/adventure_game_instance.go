package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

func AdventureGameInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameInstance) (adventure_game_schema.AdventureGameInstanceResponseData, error) {
	data := adventure_game_schema.AdventureGameInstanceResponseData{
		ID:        rec.ID,
		GameID:    rec.GameID,
		AccountID: "", // TODO: This field doesn't exist in the record, needs to be added
		Status:    rec.Status,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt: nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameInstance) (adventure_game_schema.AdventureGameInstanceResponse, error) {
	data, err := AdventureGameInstanceRecordToResponseData(l, rec)
	if err != nil {
		return adventure_game_schema.AdventureGameInstanceResponse{}, err
	}
	return adventure_game_schema.AdventureGameInstanceResponse{
		Data: &data,
	}, nil
}

func AdventureGameInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameInstance) (adventure_game_schema.AdventureGameInstanceCollectionResponse, error) {
	data := []*adventure_game_schema.AdventureGameInstanceResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameInstanceRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return adventure_game_schema.AdventureGameInstanceCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameInstanceRequestToRecord(l logger.Logger, req *adventure_game_schema.AdventureGameInstanceRequest, rec *adventure_game_record.AdventureGameInstance) (*adventure_game_record.AdventureGameInstance, error) {
	// TODO: AccountID field doesn't exist in record, needs to be added
	if req.Status != "" {
		rec.Status = req.Status
	}
	return rec, nil
}
