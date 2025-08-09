package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"
)

// NOTE: Adventure game creature instance records are created by the game instance creation process
// and are not created or updated by the user.

func AdventureGameCreatureInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameCreatureInstance) (*adventure_game_schema.AdventureGameCreatureInstanceResponseData, error) {
	l.Debug("mapping adventure_game_creature_instance record to response data")
	return &adventure_game_schema.AdventureGameCreatureInstanceResponseData{
		ID:                     rec.ID,
		GameID:                 rec.GameID,
		GameCreatureID:         rec.AdventureGameCreatureID,
		GameInstanceID:         rec.GameInstanceID,
		GameLocationInstanceID: rec.AdventureGameLocationInstanceID,
		IsAlive:                rec.Health > 0, // Convert health to boolean for now
		CreatedAt:              rec.CreatedAt,
		UpdatedAt:              nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:              nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AdventureGameCreatureInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameCreatureInstance) (*adventure_game_schema.AdventureGameCreatureInstanceResponse, error) {
	l.Debug("mapping adventure_game_creature_instance record to response")
	data, err := AdventureGameCreatureInstanceRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &adventure_game_schema.AdventureGameCreatureInstanceResponse{
		Data: data,
	}, nil
}

func AdventureGameCreatureInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameCreatureInstance) (*adventure_game_schema.AdventureGameCreatureInstanceCollectionResponse, error) {
	l.Debug("mapping adventure_game_creature_instance records to collection response")
	data := []*adventure_game_schema.AdventureGameCreatureInstanceResponseData{}
	for _, rec := range recs {
		d, err := AdventureGameCreatureInstanceRecordToResponseData(l, rec)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return &adventure_game_schema.AdventureGameCreatureInstanceCollectionResponse{
		Data: data,
	}, nil
}
