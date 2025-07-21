package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/schema"

	"gitlab.com/alienspaces/playbymail/internal/record"
)

func AdventureGameCreatureInstanceRequestToRecord(req *schema.AdventureGameCreatureInstanceRequest, rec *record.AdventureGameCreatureInstance) (*record.AdventureGameCreatureInstance, error) {
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

func AdventureGameCreatureInstanceRecordToResponseData(rec *record.AdventureGameCreatureInstance) *schema.AdventureGameCreatureInstanceResponseData {
	if rec == nil {
		return nil
	}
	return &schema.AdventureGameCreatureInstanceResponseData{
		ID:                     rec.ID,
		GameID:                 rec.GameID,
		GameCreatureID:         rec.AdventureGameCreatureID,
		GameInstanceID:         rec.AdventureGameInstanceID,
		GameLocationInstanceID: rec.AdventureGameLocationInstanceID,
		IsAlive:                rec.Health > 0, // Convert health to boolean for now
		CreatedAt:              rec.CreatedAt,
		UpdatedAt:              nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:              nulltime.ToTimePtr(rec.DeletedAt),
	}
}

func AdventureGameCreatureInstanceRecordToResponse(rec *record.AdventureGameCreatureInstance) *schema.AdventureGameCreatureInstanceResponse {
	return &schema.AdventureGameCreatureInstanceResponse{
		Data: AdventureGameCreatureInstanceRecordToResponseData(rec),
	}
}

func AdventureGameCreatureInstanceRecordsToCollectionResponse(recs []*record.AdventureGameCreatureInstance) *schema.AdventureGameCreatureInstanceCollectionResponse {
	var data []*schema.AdventureGameCreatureInstanceResponseData
	for _, rec := range recs {
		data = append(data, AdventureGameCreatureInstanceRecordToResponseData(rec))
	}
	return &schema.AdventureGameCreatureInstanceCollectionResponse{
		Data: data,
	}
}
