package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/schema/api"

	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
)

func AdventureGameCreatureInstanceRequestToRecord(req *api.AdventureGameCreatureInstanceRequest, rec *adventure_game_record.AdventureGameCreatureInstance) (*adventure_game_record.AdventureGameCreatureInstance, error) {
	if req == nil {
		return nil, nil
	}

	health := 0
	if req.IsAlive {
		health = 100 // Default health value when alive
	}

	rec.AdventureGameCreatureID = req.GameCreatureID
	rec.AdventureGameLocationInstanceID = req.GameLocationInstanceID
	rec.Health = health

	return rec, nil
}

func AdventureGameCreatureInstanceRecordToResponseData(rec *adventure_game_record.AdventureGameCreatureInstance) *api.AdventureGameCreatureInstanceResponseData {
	if rec == nil {
		return nil
	}
	return &api.AdventureGameCreatureInstanceResponseData{
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

func AdventureGameCreatureInstanceRecordToResponse(rec *adventure_game_record.AdventureGameCreatureInstance) *api.AdventureGameCreatureInstanceResponse {
	return &api.AdventureGameCreatureInstanceResponse{
		Data: AdventureGameCreatureInstanceRecordToResponseData(rec),
	}
}

func AdventureGameCreatureInstanceRecordsToCollectionResponse(recs []*adventure_game_record.AdventureGameCreatureInstance) *api.AdventureGameCreatureInstanceCollectionResponse {
	data := []*api.AdventureGameCreatureInstanceResponseData{}
	for _, rec := range recs {
		data = append(data, AdventureGameCreatureInstanceRecordToResponseData(rec))
	}
	return &api.AdventureGameCreatureInstanceCollectionResponse{
		Data: data,
	}
}
