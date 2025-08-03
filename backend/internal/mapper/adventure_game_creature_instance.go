package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/schema/api/adventure_game_schema"

	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func AdventureGameCreatureInstanceRecordToResponseData(rec *adventure_game_record.AdventureGameCreatureInstance) *adventure_game_schema.AdventureGameCreatureInstanceResponseData {
	if rec == nil {
		return nil
	}
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
	}
}

func AdventureGameCreatureInstanceRecordToResponse(rec *adventure_game_record.AdventureGameCreatureInstance) *adventure_game_schema.AdventureGameCreatureInstanceResponse {
	return &adventure_game_schema.AdventureGameCreatureInstanceResponse{
		Data: AdventureGameCreatureInstanceRecordToResponseData(rec),
	}
}

func AdventureGameCreatureInstanceRecordsToCollectionResponse(recs []*adventure_game_record.AdventureGameCreatureInstance) *adventure_game_schema.AdventureGameCreatureInstanceCollectionResponse {
	data := []*adventure_game_schema.AdventureGameCreatureInstanceResponseData{}
	for _, rec := range recs {
		data = append(data, AdventureGameCreatureInstanceRecordToResponseData(rec))
	}
	return &adventure_game_schema.AdventureGameCreatureInstanceCollectionResponse{
		Data: data,
	}
}
