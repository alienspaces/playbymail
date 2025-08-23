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
		rec.ParameterKey = req.ParameterKey
		// Convert any type to string for database storage
		if req.ParameterValue != nil {
			rec.ParameterValue = nullstring.FromString(fmt.Sprintf("%v", req.ParameterValue))
		}
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.ParameterKey = req.ParameterKey
		// Convert any type to string for database storage
		if req.ParameterValue != nil {
			rec.ParameterValue = nullstring.FromString(fmt.Sprintf("%v", req.ParameterValue))
		}
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
		ParameterKey:   rec.ParameterKey,
		ParameterValue: nullstring.ToString(rec.ParameterValue),
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      nulltime.ToTimePtr(rec.UpdatedAt),
	}

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
