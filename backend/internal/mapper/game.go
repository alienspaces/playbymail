package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func GameRequestToRecord(l logger.Logger, r *http.Request, rec *record.Game) (*record.Game, error) {

	l.Debug("mapping game request to record")

	var req schema.GameRequest
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

func GameRecordToResponseData(l logger.Logger, rec *record.Game) (schema.GameResponseData, error) {

	l.Debug("mapping game record to response data")

	data := schema.GameResponseData{
		ID:        rec.ID,
		Name:      rec.Name,
		GameType:  rec.GameType,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
	}

	return data, nil
}

func GameRecordToResponse(l logger.Logger, rec *record.Game) (schema.GameResponse, error) {
	data, err := GameRecordToResponseData(l, rec)
	if err != nil {
		return schema.GameResponse{}, err
	}
	return schema.GameResponse{
		Data: &data,
	}, nil
}

func GameRecordsToCollectionResponse(l logger.Logger, recs []*record.Game) (schema.GameCollectionResponse, error) {
	data := []*schema.GameResponseData{}
	for _, rec := range recs {
		d, err := GameRecordToResponseData(l, rec)
		if err != nil {
			return schema.GameCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.GameCollectionResponse{
		Data: data,
	}, nil
}
