package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func GameRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.Game) (*game_record.Game, error) {

	l.Debug("mapping game request to record")

	var req api.GameRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Name = req.Name
		rec.GameType = req.GameType
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.Name = req.Name
		rec.GameType = req.GameType
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameRecordToResponseData(l logger.Logger, rec *game_record.Game) (api.GameResponseData, error) {

	l.Debug("mapping game record to response data")

	data := api.GameResponseData{
		ID:        rec.ID,
		Name:      rec.Name,
		GameType:  rec.GameType,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
	}

	return data, nil
}

func GameRecordToResponse(l logger.Logger, rec *game_record.Game) (api.GameResponse, error) {
	data, err := GameRecordToResponseData(l, rec)
	if err != nil {
		return api.GameResponse{}, err
	}
	return api.GameResponse{
		Data: &data,
	}, nil
}

func GameRecordsToCollectionResponse(l logger.Logger, recs []*game_record.Game) (api.GameCollectionResponse, error) {
	data := []*api.GameResponseData{}
	for _, rec := range recs {
		d, err := GameRecordToResponseData(l, rec)
		if err != nil {
			return api.GameCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.GameCollectionResponse{
		Data: data,
	}, nil
}
