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

func GameRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.Game) (*game_record.Game, error) {

	l.Debug("mapping game request to record")

	var req game_schema.GameRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Name = req.Name
		rec.GameType = req.GameType
		rec.TurnDurationHours = req.TurnDurationHours
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.GameType = req.GameType
		rec.TurnDurationHours = req.TurnDurationHours
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameRecordToResponseData(l logger.Logger, rec *game_record.Game) (*game_schema.GameResponseData, error) {
	l.Debug("mapping game record to response data")
	return &game_schema.GameResponseData{
		ID:                rec.ID,
		Name:              rec.Name,
		GameType:          rec.GameType,
		TurnDurationHours: rec.TurnDurationHours,
		CreatedAt:         rec.CreatedAt,
		UpdatedAt:         nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func GameRecordToResponse(l logger.Logger, rec *game_record.Game) (*game_schema.GameResponse, error) {
	l.Debug("mapping game record to response")
	data, err := GameRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameResponse{
		Data: data,
	}, nil
}

func GameRecordsToCollectionResponse(l logger.Logger, recs []*game_record.Game) (game_schema.GameCollectionResponse, error) {
	l.Debug("mapping game records to collection response")
	data := []*game_schema.GameResponseData{}
	for _, rec := range recs {
		d, err := GameRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameCollectionResponse{
		Data: data,
	}, nil
}
