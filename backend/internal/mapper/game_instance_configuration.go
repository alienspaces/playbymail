package mapper

import (
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func GameInstanceConfigurationRecordToResponseData(l logger.Logger, rec *game_record.GameInstanceConfiguration) (api.GameInstanceConfiguration, error) {
	data := api.GameInstanceConfiguration{
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

	if rec.StringValue.Valid {
		data.StringValue = &rec.StringValue.String
	}
	if rec.IntegerValue.Valid {
		value := int(rec.IntegerValue.Int32)
		data.IntegerValue = &value
	}
	if rec.BooleanValue.Valid {
		data.BooleanValue = &rec.BooleanValue.Bool
	}
	if rec.JSONValue.Valid {
		data.JSONValue = &rec.JSONValue.String
	}

	return data, nil
}

func GameInstanceConfigurationRecordToResponse(l logger.Logger, rec *game_record.GameInstanceConfiguration) (api.GameInstanceConfigurationResponse, error) {
	data, err := GameInstanceConfigurationRecordToResponseData(l, rec)
	if err != nil {
		return api.GameInstanceConfigurationResponse{}, err
	}
	return api.GameInstanceConfigurationResponse{
		Data: &data,
	}, nil
}

func GameInstanceConfigurationRecsToCollectionResponse(l logger.Logger, recs []*game_record.GameInstanceConfiguration) (api.GameInstanceConfigurationCollectionResponse, error) {
	data := []*api.GameInstanceConfiguration{}
	for _, rec := range recs {
		d, err := GameInstanceConfigurationRecordToResponseData(l, rec)
		if err != nil {
			return api.GameInstanceConfigurationCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.GameInstanceConfigurationCollectionResponse{
		Data: data,
	}, nil
}

func GameInstanceConfigurationRequestToRecord(l logger.Logger, req *api.GameInstanceConfigurationRequest, rec *game_record.GameInstanceConfiguration) (*game_record.GameInstanceConfiguration, error) {
	if req.GameInstanceID != "" {
		rec.GameInstanceID = req.GameInstanceID
	}
	if req.ConfigKey != "" {
		rec.ConfigKey = req.ConfigKey
	}
	if req.ValueType != "" {
		rec.ValueType = req.ValueType
	}
	if req.StringValue != nil {
		rec.StringValue.String = *req.StringValue
		rec.StringValue.Valid = true
	}
	if req.IntegerValue != nil {
		rec.IntegerValue.Int32 = int32(*req.IntegerValue)
		rec.IntegerValue.Valid = true
	}
	if req.BooleanValue != nil {
		rec.BooleanValue.Bool = *req.BooleanValue
		rec.BooleanValue.Valid = true
	}
	if req.JSONValue != nil {
		rec.JSONValue.String = *req.JSONValue
		rec.JSONValue.Valid = true
	}
	return rec, nil
}
