package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/adventure_game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func AdventureGameInstanceRecordToResponseData(l logger.Logger, rec *adventure_game_record.AdventureGameInstance) (api.AdventureGameInstance, error) {
	data := api.AdventureGameInstance{
		ID:                  rec.ID,
		GameID:              rec.GameID,
		Status:              rec.Status,
		CurrentTurn:         rec.CurrentTurn,
		MaxTurns:            rec.MaxTurns,
		TurnDeadlineHours:   rec.TurnDeadlineHours,
		LastTurnProcessedAt: rec.LastTurnProcessedAt,
		NextTurnDeadline:    rec.NextTurnDeadline,
		StartedAt:           rec.StartedAt,
		CompletedAt:         rec.CompletedAt,
		GameConfig:          rec.GameConfig,
		CreatedAt:           rec.CreatedAt,
		UpdatedAt:           nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:           nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func AdventureGameInstanceRecordToResponse(l logger.Logger, rec *adventure_game_record.AdventureGameInstance) (api.AdventureGameInstanceResponse, error) {
	data, err := AdventureGameInstanceRecordToResponseData(l, rec)
	if err != nil {
		return api.AdventureGameInstanceResponse{}, err
	}
	return api.AdventureGameInstanceResponse{
		Data: &data,
	}, nil
}

func AdventureGameInstanceRecordsToCollectionResponse(l logger.Logger, recs []*adventure_game_record.AdventureGameInstance) (api.AdventureGameInstanceCollectionResponse, error) {
	data := []*api.AdventureGameInstance{}
	for _, rec := range recs {
		d, err := AdventureGameInstanceRecordToResponseData(l, rec)
		if err != nil {
			return api.AdventureGameInstanceCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AdventureGameInstanceCollectionResponse{
		Data: data,
	}, nil
}

func AdventureGameInstanceRequestToRecord(l logger.Logger, req *api.AdventureGameInstanceRequest, rec *adventure_game_record.AdventureGameInstance) (*adventure_game_record.AdventureGameInstance, error) {
	if req.GameID != "" {
		rec.GameID = req.GameID
	}
	if req.Status != "" {
		rec.Status = req.Status
	}
	if req.CurrentTurn > 0 {
		rec.CurrentTurn = req.CurrentTurn
	}
	if req.MaxTurns != nil {
		rec.MaxTurns = req.MaxTurns
	}
	if req.TurnDeadlineHours > 0 {
		rec.TurnDeadlineHours = req.TurnDeadlineHours
	}
	if req.LastTurnProcessedAt != nil {
		rec.LastTurnProcessedAt = req.LastTurnProcessedAt
	}
	if req.NextTurnDeadline != nil {
		rec.NextTurnDeadline = req.NextTurnDeadline
	}
	if req.StartedAt != nil {
		rec.StartedAt = req.StartedAt
	}
	if req.CompletedAt != nil {
		rec.CompletedAt = req.CompletedAt
	}
	if req.GameConfig != nil {
		rec.GameConfig = req.GameConfig
	}
	return rec, nil
}
