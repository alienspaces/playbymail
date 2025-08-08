package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func GameInstanceRecordToResponseData(l logger.Logger, rec *game_record.GameInstance) (game_schema.GameInstance, error) {
	data := game_schema.GameInstance{
		ID:                  rec.ID,
		GameID:              rec.GameID,
		Status:              rec.Status,
		CurrentTurn:         rec.CurrentTurn,
		LastTurnProcessedAt: nulltime.ToTimePtr(rec.LastTurnProcessedAt),
		NextTurnDueAt:       nulltime.ToTimePtr(rec.NextTurnDueAt),
		StartedAt:           nulltime.ToTimePtr(rec.StartedAt),
		CompletedAt:         nulltime.ToTimePtr(rec.CompletedAt),
		CreatedAt:           rec.CreatedAt,
		UpdatedAt:           nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:           nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func GameInstanceRecordToResponse(l logger.Logger, rec *game_record.GameInstance) (game_schema.GameInstanceResponse, error) {
	data, err := GameInstanceRecordToResponseData(l, rec)
	if err != nil {
		return game_schema.GameInstanceResponse{}, err
	}
	return game_schema.GameInstanceResponse{
		Data: &data,
	}, nil
}

func GameInstanceRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameInstance) (game_schema.GameInstanceCollectionResponse, error) {
	data := []*game_schema.GameInstance{}
	for _, rec := range recs {
		d, err := GameInstanceRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return game_schema.GameInstanceCollectionResponse{
		Data: data,
	}, nil
}

func GameInstanceRequestToRecord(l logger.Logger, req *game_schema.GameInstanceRequest, rec *game_record.GameInstance) (*game_record.GameInstance, error) {
	if req.GameID != "" {
		rec.GameID = req.GameID
	}
	if req.Status != "" {
		rec.Status = req.Status
	}
	if req.CurrentTurn > 0 {
		rec.CurrentTurn = req.CurrentTurn
	}
	if req.LastTurnProcessedAt != nil {
		rec.LastTurnProcessedAt = nulltime.FromTimePtr(req.LastTurnProcessedAt)
	}
	if req.NextTurnDueAt != nil {
		rec.NextTurnDueAt = nulltime.FromTimePtr(req.NextTurnDueAt)
	}
	if req.StartedAt != nil {
		rec.StartedAt = nulltime.FromTimePtr(req.StartedAt)
	}
	if req.CompletedAt != nil {
		rec.CompletedAt = nulltime.FromTimePtr(req.CompletedAt)
	}
	return rec, nil
}
