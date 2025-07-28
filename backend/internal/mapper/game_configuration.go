package mapper

import (
	"database/sql"
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

func GameConfigurationRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameConfiguration) (*game_record.GameConfiguration, error) {
	l.Debug("mapping game configuration request to record")

	var req api.GameConfigurationRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameType = req.GameType
		rec.ConfigKey = req.ConfigKey
		rec.ValueType = req.ValueType
		rec.DefaultValue = sql.NullString{String: "", Valid: false}
		if req.DefaultValue != nil {
			rec.DefaultValue = sql.NullString{String: *req.DefaultValue, Valid: true}
		}
		rec.IsRequired = req.IsRequired
		rec.Description = sql.NullString{String: "", Valid: false}
		if req.Description != nil {
			rec.Description = sql.NullString{String: *req.Description, Valid: true}
		}
		rec.UIHint = sql.NullString{String: "", Valid: false}
		if req.UIHint != nil {
			rec.UIHint = sql.NullString{String: *req.UIHint, Valid: true}
		}
		rec.ValidationRules = sql.NullString{String: "", Valid: false}
		if req.ValidationRules != nil {
			rec.ValidationRules = sql.NullString{String: *req.ValidationRules, Valid: true}
		}
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.GameType = req.GameType
		rec.ConfigKey = req.ConfigKey
		rec.ValueType = req.ValueType
		rec.DefaultValue = sql.NullString{String: "", Valid: false}
		if req.DefaultValue != nil {
			rec.DefaultValue = sql.NullString{String: *req.DefaultValue, Valid: true}
		}
		rec.IsRequired = req.IsRequired
		rec.Description = sql.NullString{String: "", Valid: false}
		if req.Description != nil {
			rec.Description = sql.NullString{String: *req.Description, Valid: true}
		}
		rec.UIHint = sql.NullString{String: "", Valid: false}
		if req.UIHint != nil {
			rec.UIHint = sql.NullString{String: *req.UIHint, Valid: true}
		}
		rec.ValidationRules = sql.NullString{String: "", Valid: false}
		if req.ValidationRules != nil {
			rec.ValidationRules = sql.NullString{String: *req.ValidationRules, Valid: true}
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameConfigurationRecordToResponseData(l logger.Logger, rec *game_record.GameConfiguration) (api.GameConfiguration, error) {
	l.Debug("mapping game configuration record to response data")

	data := api.GameConfiguration{
		ID:              rec.ID,
		GameType:        rec.GameType,
		ConfigKey:       rec.ConfigKey,
		ValueType:       rec.ValueType,
		DefaultValue:    nil,
		IsRequired:      rec.IsRequired,
		Description:     nil,
		UIHint:          nil,
		ValidationRules: nil,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:       nulltime.ToTimePtr(rec.UpdatedAt),
	}

	if rec.DefaultValue.Valid {
		data.DefaultValue = &rec.DefaultValue.String
	}
	if rec.Description.Valid {
		data.Description = &rec.Description.String
	}
	if rec.UIHint.Valid {
		data.UIHint = &rec.UIHint.String
	}
	if rec.ValidationRules.Valid {
		data.ValidationRules = &rec.ValidationRules.String
	}

	return data, nil
}

func GameConfigurationRecordToResponse(l logger.Logger, rec *game_record.GameConfiguration) (api.GameConfigurationResponse, error) {
	data, err := GameConfigurationRecordToResponseData(l, rec)
	if err != nil {
		return api.GameConfigurationResponse{}, err
	}
	return api.GameConfigurationResponse{
		Data: &data,
	}, nil
}

func GameConfigurationRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameConfiguration) (api.GameConfigurationCollectionResponse, error) {
	data := []*api.GameConfiguration{}
	for _, rec := range recs {
		d, err := GameConfigurationRecordToResponseData(l, rec)
		if err != nil {
			return api.GameConfigurationCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.GameConfigurationCollectionResponse{
		Data: data,
	}, nil
}

// Convenience functions for the handler
func MapGameConfigurationRequestToRecord(req *api.GameConfigurationRequest) *game_record.GameConfiguration {
	rec := &game_record.GameConfiguration{}
	rec.GameType = req.GameType
	rec.ConfigKey = req.ConfigKey
	rec.ValueType = req.ValueType
	rec.DefaultValue = sql.NullString{String: "", Valid: false}
	if req.DefaultValue != nil {
		rec.DefaultValue = sql.NullString{String: *req.DefaultValue, Valid: true}
	}
	rec.IsRequired = req.IsRequired
	rec.Description = sql.NullString{String: "", Valid: false}
	if req.Description != nil {
		rec.Description = sql.NullString{String: *req.Description, Valid: true}
	}
	rec.UIHint = sql.NullString{String: "", Valid: false}
	if req.UIHint != nil {
		rec.UIHint = sql.NullString{String: *req.UIHint, Valid: true}
	}
	rec.ValidationRules = sql.NullString{String: "", Valid: false}
	if req.ValidationRules != nil {
		rec.ValidationRules = sql.NullString{String: *req.ValidationRules, Valid: true}
	}
	return rec
}

func MapGameConfigurationResponse(rec *game_record.GameConfiguration) *api.GameConfigurationResponse {
	data := api.GameConfiguration{
		ID:              rec.ID,
		GameType:        rec.GameType,
		ConfigKey:       rec.ConfigKey,
		ValueType:       rec.ValueType,
		DefaultValue:    nil,
		IsRequired:      rec.IsRequired,
		Description:     nil,
		UIHint:          nil,
		ValidationRules: nil,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:       nulltime.ToTimePtr(rec.UpdatedAt),
	}

	if rec.DefaultValue.Valid {
		data.DefaultValue = &rec.DefaultValue.String
	}
	if rec.Description.Valid {
		data.Description = &rec.Description.String
	}
	if rec.UIHint.Valid {
		data.UIHint = &rec.UIHint.String
	}
	if rec.ValidationRules.Valid {
		data.ValidationRules = &rec.ValidationRules.String
	}

	return &api.GameConfigurationResponse{
		Data: &data,
	}
}

func MapGameConfigurationCollectionResponse(recs []*game_record.GameConfiguration) *api.GameConfigurationCollectionResponse {
	data := []*api.GameConfiguration{}
	for _, rec := range recs {
		d := api.GameConfiguration{
			ID:              rec.ID,
			GameType:        rec.GameType,
			ConfigKey:       rec.ConfigKey,
			ValueType:       rec.ValueType,
			DefaultValue:    nil,
			IsRequired:      rec.IsRequired,
			Description:     nil,
			UIHint:          nil,
			ValidationRules: nil,
			CreatedAt:       rec.CreatedAt,
			UpdatedAt:       nulltime.ToTimePtr(rec.UpdatedAt),
		}

		if rec.DefaultValue.Valid {
			d.DefaultValue = &rec.DefaultValue.String
		}
		if rec.Description.Valid {
			d.Description = &rec.Description.String
		}
		if rec.UIHint.Valid {
			d.UIHint = &rec.UIHint.String
		}
		if rec.ValidationRules.Valid {
			d.ValidationRules = &rec.ValidationRules.String
		}

		data = append(data, &d)
	}
	return &api.GameConfigurationCollectionResponse{
		Data: data,
	}
}
