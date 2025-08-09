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

func GameConfigurationRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameConfiguration) (*game_record.GameConfiguration, error) {
	l.Debug("mapping game configuration request to record")

	var req game_schema.GameConfigurationRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameType = req.GameType
		rec.ConfigKey = req.ConfigKey
		rec.Description = req.Description
		rec.ValueType = req.ValueType
		rec.DefaultValue = nullstring.FromStringPtr(req.DefaultValue)
		rec.IsRequired = req.IsRequired
		rec.IsGlobal = req.IsGlobal
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.GameType = req.GameType
		rec.ConfigKey = req.ConfigKey
		rec.Description = req.Description
		rec.ValueType = req.ValueType
		rec.DefaultValue = nullstring.FromStringPtr(req.DefaultValue)
		rec.IsRequired = req.IsRequired
		rec.IsGlobal = req.IsGlobal
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameConfigurationRecordToResponseData(l logger.Logger, rec *game_record.GameConfiguration) (*game_schema.GameConfiguration, error) {
	l.Debug("mapping game configuration record to response data")
	return &game_schema.GameConfiguration{
		ID:           rec.ID,
		GameType:     rec.GameType,
		ConfigKey:    rec.ConfigKey,
		Description:  rec.Description,
		ValueType:    rec.ValueType,
		DefaultValue: nullstring.ToStringPtr(rec.DefaultValue),
		IsRequired:   rec.IsRequired,
		IsGlobal:     rec.IsGlobal,
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func GameConfigurationRecordToResponse(l logger.Logger, rec *game_record.GameConfiguration) (*game_schema.GameConfigurationResponse, error) {
	data, err := GameConfigurationRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameConfigurationResponse{
		Data: data,
	}, nil
}

func GameConfigurationRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameConfiguration) (game_schema.GameConfigurationCollectionResponse, error) {
	data := []*game_schema.GameConfiguration{}
	for _, rec := range recs {
		d, err := GameConfigurationRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameConfigurationCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameConfigurationCollectionResponse{
		Data: data,
	}, nil
}
