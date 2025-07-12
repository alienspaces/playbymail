package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/schema"

	"gitlab.com/alienspaces/playbymail/internal/record"
)

func GameCreatureInstanceRequestToRecord(req *schema.GameCreatureInstanceRequest) *record.GameCreatureInstance {
	return &record.GameCreatureInstance{
		GameID:                 req.GameID,
		GameCreatureID:         req.GameCreatureID,
		GameInstanceID:         req.GameInstanceID,
		GameLocationInstanceID: req.GameLocationInstanceID,
		IsAlive:                req.IsAlive,
	}
}

func GameCreatureInstanceRecordToResponseData(rec *record.GameCreatureInstance) *schema.GameCreatureInstanceResponseData {
	if rec == nil {
		return nil
	}
	return &schema.GameCreatureInstanceResponseData{
		ID:                     rec.ID,
		GameID:                 rec.GameID,
		GameCreatureID:         rec.GameCreatureID,
		GameInstanceID:         rec.GameInstanceID,
		GameLocationInstanceID: rec.GameLocationInstanceID,
		IsAlive:                rec.IsAlive,
		CreatedAt:              rec.CreatedAt,
		UpdatedAt:              nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:              nulltime.ToTimePtr(rec.DeletedAt),
	}
}

func GameCreatureInstanceRecordToResponse(rec *record.GameCreatureInstance) *schema.GameCreatureInstanceResponse {
	return &schema.GameCreatureInstanceResponse{
		Data: GameCreatureInstanceRecordToResponseData(rec),
	}
}

func GameCreatureInstanceRecordsToCollectionResponse(recs []*record.GameCreatureInstance) *schema.GameCreatureInstanceCollectionResponse {
	var data []*schema.GameCreatureInstanceResponseData
	for _, rec := range recs {
		data = append(data, GameCreatureInstanceRecordToResponseData(rec))
	}
	return &schema.GameCreatureInstanceCollectionResponse{
		Data: data,
	}
}
