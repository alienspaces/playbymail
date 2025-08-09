package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nullbool"
	"gitlab.com/alienspaces/playbymail/core/nullint32"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

// GameInstanceConfigurationRequestToRecord maps a request to a record for consistency
func GameInstanceConfigurationRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameInstanceConfiguration) (*game_record.GameInstanceConfiguration, error) {
	l.Debug("mapping game_instance_configuration request to record")

	var req game_schema.GameInstanceConfigurationRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameInstanceID = req.GameInstanceID
		rec.ConfigKey = req.ConfigKey
		rec.ValueType = req.ValueType
		rec.StringValue = nullstring.FromStringPtr(req.StringValue)
		rec.IntegerValue = nullint32.FromInt32Ptr(req.IntegerValue)
		rec.BooleanValue = nullbool.FromBoolPtr(req.BooleanValue)
		rec.JSONValue = nullstring.FromStringPtr(req.JSONValue)
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.GameInstanceID = req.GameInstanceID
		rec.ConfigKey = req.ConfigKey
		rec.ValueType = req.ValueType
		rec.StringValue = nullstring.FromStringPtr(req.StringValue)
		rec.IntegerValue = nullint32.FromInt32Ptr(req.IntegerValue)
		rec.BooleanValue = nullbool.FromBoolPtr(req.BooleanValue)
		rec.JSONValue = nullstring.FromStringPtr(req.JSONValue)
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameInstanceConfigurationRecordToResponseData(l logger.Logger, rec *game_record.GameInstanceConfiguration) (*game_schema.GameInstanceConfiguration, error) {
	l.Debug("mapping game_instance_configuration record to response data")
	data := &game_schema.GameInstanceConfiguration{
		ID:             rec.ID,
		GameInstanceID: rec.GameInstanceID,
		ConfigKey:      rec.ConfigKey,
		ValueType:      rec.ValueType,
		StringValue:    nil,
		IntegerValue:   nil,
		BooleanValue:   nil,
		JSONValue:      nil,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:      nulltime.ToTimePtr(rec.DeletedAt),
	}

	data.StringValue = nullstring.ToStringPtr(rec.StringValue)

	intValue, err := nullint32.ToInt32Ptr(rec.IntegerValue)
	if err != nil {
		return nil, err
	}

	data.IntegerValue = intValue
	data.BooleanValue = nullbool.ToBoolPtr(rec.BooleanValue)
	data.JSONValue = nullstring.ToStringPtr(rec.JSONValue)

	return data, nil
}

func GameInstanceConfigurationRecordToResponse(l logger.Logger, rec *game_record.GameInstanceConfiguration) (*game_schema.GameInstanceConfigurationResponse, error) {
	l.Debug("mapping game_instance_configuration record to response")
	data, err := GameInstanceConfigurationRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameInstanceConfigurationResponse{
		Data: data,
	}, nil
}

func GameInstanceConfigurationRecsToCollectionResponse(l logger.Logger, recs []*game_record.GameInstanceConfiguration) (game_schema.GameInstanceConfigurationCollectionResponse, error) {
	l.Debug("mapping game_instance_configuration records to collection response")
	data := []*game_schema.GameInstanceConfiguration{}
	for _, rec := range recs {
		d, err := GameInstanceConfigurationRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameInstanceConfigurationCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameInstanceConfigurationCollectionResponse{
		Data: data,
	}, nil
}
