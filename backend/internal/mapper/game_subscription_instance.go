package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/internal/utils/stringutil"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func GameSubscriptionInstanceRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameSubscriptionInstance) (*game_record.GameSubscriptionInstance, error) {
	l.Debug("mapping game_subscription_instance request to record")

	var req game_schema.GameSubscriptionInstanceRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameSubscriptionID = req.GameSubscriptionID
		rec.GameInstanceID = req.GameInstanceID
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.GameSubscriptionID = req.GameSubscriptionID
		rec.GameInstanceID = req.GameInstanceID
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameSubscriptionInstanceRecordToResponseData(l logger.Logger, rec *game_record.GameSubscriptionInstance) (*game_schema.GameSubscriptionInstanceResponseData, error) {
	l.Debug("mapping game_subscription_instance record to response data")

	// Mask the token for display purposes
	var maskedToken *string
	if nullstring.IsValid(rec.TurnSheetToken) {
		masked := stringutil.MaskSensitiveValue(nullstring.ToString(rec.TurnSheetToken))
		maskedToken = &masked
	}

	return &game_schema.GameSubscriptionInstanceResponseData{
		ID:                      rec.ID,
		AccountID:               rec.AccountID,
		GameSubscriptionID:      rec.GameSubscriptionID,
		GameInstanceID:          rec.GameInstanceID,
		TurnSheetTokenMasked:    maskedToken,
		TurnSheetTokenExpiresAt: nulltime.ToTimePtr(rec.TurnSheetTokenExpiresAt),
		CreatedAt:               rec.CreatedAt,
		UpdatedAt:               nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:               nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func GameSubscriptionInstanceRecordToResponse(l logger.Logger, rec *game_record.GameSubscriptionInstance) (*game_schema.GameSubscriptionInstanceResponse, error) {
	l.Debug("mapping game_subscription_instance record to response")
	data, err := GameSubscriptionInstanceRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameSubscriptionInstanceResponse{
		Data: data,
	}, nil
}

func GameSubscriptionInstanceRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameSubscriptionInstance) (game_schema.GameSubscriptionInstanceCollectionResponse, error) {
	l.Debug("mapping game_subscription_instance records to collection response")
	data := []*game_schema.GameSubscriptionInstanceResponseData{}
	for _, rec := range recs {
		d, err := GameSubscriptionInstanceRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameSubscriptionInstanceCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameSubscriptionInstanceCollectionResponse{
		Data: data,
	}, nil
}
