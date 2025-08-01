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

func GameSubscriptionRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameSubscription) (*game_record.GameSubscription, error) {
	l.Debug("mapping game_subscription request to record")

	var req api.GameSubscriptionRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameID = req.GameID
		rec.AccountID = req.AccountID
		rec.SubscriptionType = req.SubscriptionType
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.GameID = req.GameID
		rec.AccountID = req.AccountID
		rec.SubscriptionType = req.SubscriptionType
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameSubscriptionRecordToResponseData(l logger.Logger, rec *game_record.GameSubscription) (api.GameSubscriptionResponseData, error) {
	l.Debug("mapping game_subscription record to response data")
	data := api.GameSubscriptionResponseData{
		ID:               rec.ID,
		GameID:           rec.GameID,
		AccountID:        rec.AccountID,
		SubscriptionType: rec.SubscriptionType,
		CreatedAt:        rec.CreatedAt,
		UpdatedAt:        nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:        nulltime.ToTimePtr(rec.DeletedAt),
	}
	return data, nil
}

func GameSubscriptionRecordToResponse(l logger.Logger, rec *game_record.GameSubscription) (api.GameSubscriptionResponse, error) {
	data, err := GameSubscriptionRecordToResponseData(l, rec)
	if err != nil {
		return api.GameSubscriptionResponse{}, err
	}
	return api.GameSubscriptionResponse{
		Data: &data,
	}, nil
}

func GameSubscriptionRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameSubscription) (api.GameSubscriptionCollectionResponse, error) {
	data := []*api.GameSubscriptionResponseData{}
	for _, rec := range recs {
		d, err := GameSubscriptionRecordToResponseData(l, rec)
		if err != nil {
			return api.GameSubscriptionCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.GameSubscriptionCollectionResponse{
		Data: data,
	}, nil
}
