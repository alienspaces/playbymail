package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func GameInstanceRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameInstance) (*game_record.GameInstance, error) {
	l.Debug("mapping game_instance request to record")

	var req game_schema.GameInstanceRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameID = req.GameID
		rec.Status = req.Status
		rec.CurrentTurn = req.CurrentTurn
		rec.LastTurnProcessedAt = nulltime.FromTimePtr(req.LastTurnProcessedAt)
		rec.NextTurnDueAt = nulltime.FromTimePtr(req.NextTurnDueAt)
		rec.StartedAt = nulltime.FromTimePtr(req.StartedAt)
		rec.CompletedAt = nulltime.FromTimePtr(req.CompletedAt)
	case server.HttpMethodPut, server.HttpMethodPatch:
		if req.Status != "" {
			rec.Status = req.Status
		}
		if req.CurrentTurn != 0 {
			rec.CurrentTurn = req.CurrentTurn
		}
		rec.LastTurnProcessedAt = nulltime.FromTimePtr(req.LastTurnProcessedAt)
		rec.NextTurnDueAt = nulltime.FromTimePtr(req.NextTurnDueAt)
		rec.StartedAt = nulltime.FromTimePtr(req.StartedAt)
		rec.CompletedAt = nulltime.FromTimePtr(req.CompletedAt)
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameInstanceRecordToResponseData(l logger.Logger, rec *game_record.GameInstance) (*game_schema.GameInstanceResponseData, error) {
	l.Debug("mapping game_instance record to response data")
	return &game_schema.GameInstanceResponseData{
		ID:                  rec.ID,
		GameID:              rec.GameID,
		GameSubscriptionID:  rec.GameSubscriptionID,
		Status:              rec.Status,
		CurrentTurn:         rec.CurrentTurn,
		LastTurnProcessedAt: nulltime.ToTimePtr(rec.LastTurnProcessedAt),
		NextTurnDueAt:       nulltime.ToTimePtr(rec.NextTurnDueAt),
		StartedAt:           nulltime.ToTimePtr(rec.StartedAt),
		CompletedAt:         nulltime.ToTimePtr(rec.CompletedAt),
		CreatedAt:           rec.CreatedAt,
		UpdatedAt:           nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func GameInstanceRecordToResponse(l logger.Logger, rec *game_record.GameInstance) (*game_schema.GameInstanceResponse, error) {
	l.Debug("mapping game_instance record to response")
	data, err := GameInstanceRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameInstanceResponse{
		Data: data,
	}, nil
}

func GameInstanceRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameInstance) (game_schema.GameInstanceCollectionResponse, error) {
	l.Debug("mapping game_instance records to collection response")
	data := []*game_schema.GameInstanceResponseData{}
	for _, rec := range recs {
		d, err := GameInstanceRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameInstanceCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameInstanceCollectionResponse{
		Data: data,
	}, nil
}
