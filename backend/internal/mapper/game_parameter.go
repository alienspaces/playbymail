package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func GameParameterRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameParameter) (*game_record.GameParameter, error) {
	l.Debug("mapping game parameter request to record")

	var req game_schema.GameParameterRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameType = req.GameType
		rec.ConfigKey = req.ConfigKey
		rec.ValueType = req.ValueType
		rec.DefaultValue = nullstring.FromStringPtr(req.DefaultValue)
		rec.IsRequired = req.IsRequired
		rec.Description = nullstring.FromStringPtr(req.Description)
		rec.IsGlobal = req.IsGlobal
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.GameType = req.GameType
		rec.ConfigKey = req.ConfigKey
		rec.ValueType = req.ValueType
		rec.DefaultValue = nullstring.FromStringPtr(req.DefaultValue)
		rec.IsRequired = req.IsRequired
		rec.Description = nullstring.FromStringPtr(req.Description)
		rec.IsGlobal = req.IsGlobal
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameParameterRecordToResponseData(l logger.Logger, rec *game_record.GameParameter) (*game_schema.GameParameter, error) {
	l.Debug("mapping game parameter record to response data")
	return &game_schema.GameParameter{
		ID:           rec.ID,
		GameType:     rec.GameType,
		ConfigKey:    rec.ConfigKey,
		Description:  nullstring.ToStringPtr(rec.Description),
		ValueType:    rec.ValueType,
		DefaultValue: nullstring.ToStringPtr(rec.DefaultValue),
		IsRequired:   rec.IsRequired,
		IsGlobal:     rec.IsGlobal,
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func GameParameterRecordToResponse(l logger.Logger, rec *game_record.GameParameter) (*game_schema.GameParameterResponse, error) {
	data, err := GameParameterRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameParameterResponse{
		Data: data,
	}, nil
}

func GameParameterRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameParameter) (game_schema.GameParameterCollectionResponse, error) {
	data := []*game_schema.GameParameter{}
	for _, rec := range recs {
		d, err := GameParameterRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameParameterCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameParameterCollectionResponse{
		Data: data,
	}, nil
}
