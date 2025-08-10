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

// GameInstanceParameterRequestToRecord maps a request to a record for consistency
func GameInstanceParameterRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameInstanceParameter) (*game_record.GameInstanceParameter, error) {
	l.Debug("mapping game_instance_parameter request to record")

	var req game_schema.GameInstanceParameterRequest
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

func GameInstanceParameterRecordToResponseData(l logger.Logger, rec *game_record.GameInstanceParameter) (*game_schema.GameInstanceParameter, error) {
	l.Debug("mapping game_instance_parameter record to response data")
	data := &game_schema.GameInstanceParameter{
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

func GameInstanceParameterRecordToResponse(l logger.Logger, rec *game_record.GameInstanceParameter) (*game_schema.GameInstanceParameterResponse, error) {
	l.Debug("mapping game_instance_parameter record to response")
	data, err := GameInstanceParameterRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameInstanceParameterResponse{
		Data: data,
	}, nil
}

func GameInstanceParameterRecsToCollectionResponse(l logger.Logger, recs []*game_record.GameInstanceParameter) (game_schema.GameInstanceParameterCollectionResponse, error) {
	l.Debug("mapping game_instance_parameter records to collection response")
	data := []*game_schema.GameInstanceParameter{}
	for _, rec := range recs {
		d, err := GameInstanceParameterRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameInstanceParameterCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameInstanceParameterCollectionResponse{
		Data: data,
	}, nil
}
