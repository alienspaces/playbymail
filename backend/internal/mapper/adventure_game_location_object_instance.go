package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

// NOTE: adventure_game_location_object_instance records are created by the game instance
// creation process and are not created or updated directly by the user.

func AdventureGameLocationObjectInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameLocationObjectInstance) (*adventure_game_schema.AdventureGameLocationObjectInstanceResponseData, error) {
	l.Debug("mapping adventure_game_location_object_instance record to response data")
	return &adventure_game_schema.AdventureGameLocationObjectInstanceResponseData{
		ID:                              rec.ID,
		GameID:                          rec.GameID,
		GameInstanceID:                  rec.GameInstanceID,
		AdventureGameLocationObjectID:   rec.AdventureGameLocationObjectID,
		AdventureGameLocationInstanceID: rec.AdventureGameLocationInstanceID,
		CurrentState:                    rec.CurrentState,
		IsVisible:                       rec.IsVisible,
		CreatedAt:                       rec.CreatedAt,
		UpdatedAt:                       nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:                       nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameLocationObjectInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameLocationObjectInstance) (*adventure_game_schema.AdventureGameLocationObjectInstanceResponse, error) {
	l.Debug("mapping adventure_game_location_object_instance record to response")
	data, err := AdventureGameLocationObjectInstanceRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameLocationObjectInstanceResponse{
		Data: data,
	}, nil
}

func AdventureGameLocationObjectInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameLocationObjectInstance) (adventure_game_schema.AdventureGameLocationObjectInstanceCollectionResponse, error) {
	l.Debug("mapping adventure_game_location_object_instance records to collection response")
	data := []*adventure_game_schema.AdventureGameLocationObjectInstanceResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameLocationObjectInstanceRecordToResponseData(l, rec)
		if err != nil {
			return adventure_game_schema.AdventureGameLocationObjectInstanceCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return adventure_game_schema.AdventureGameLocationObjectInstanceCollectionResponse{
		Data: data,
	}, nil
}
