package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nullint32"
	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/game_record"
	"gitlab.com/alienspaces/playbymail/schema/api/game_schema"
)

func GameSubscriptionRequestToRecord(l logger.Logger, r *http.Request, rec *game_record.GameSubscription) (*game_record.GameSubscription, error) {
	l.Debug("mapping game_subscription request to record")

	var req game_schema.GameSubscriptionRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.GameID = req.GameID
		rec.AccountID = req.AccountID
		rec.AccountUserID = nullstring.FromStringPtr(req.AccountUserID)
		rec.AccountUserContactID = nullstring.FromStringPtr(req.AccountUserContactID)
		rec.SubscriptionType = req.SubscriptionType
		rec.InstanceLimit = nullint32.FromInt32Ptr(req.InstanceLimit)
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.GameID = req.GameID
		rec.AccountID = req.AccountID
		rec.AccountUserID = nullstring.FromStringPtr(req.AccountUserID)
		rec.AccountUserContactID = nullstring.FromStringPtr(req.AccountUserContactID)
		rec.SubscriptionType = req.SubscriptionType
		rec.InstanceLimit = nullint32.FromInt32Ptr(req.InstanceLimit)
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func GameSubscriptionRecordToResponseData(l logger.Logger, rec *game_record.GameSubscription, instanceIDs []string) (*game_schema.GameSubscriptionResponseData, error) {
	l.Debug("mapping game_subscription record to response data")
	instanceLimitPtr, _ := nullint32.ToInt32Ptr(rec.InstanceLimit)
	return &game_schema.GameSubscriptionResponseData{
		ID:               rec.ID,
		GameID:           rec.GameID,
		AccountID:        rec.AccountID,
		AccountUserID:    nullstring.ToStringPtr(rec.AccountUserID),
		GameInstanceIDs:  instanceIDs,
		SubscriptionType: rec.SubscriptionType,
		InstanceLimit:    instanceLimitPtr,
		Status:           rec.Status,
		CreatedAt:        rec.CreatedAt,
		UpdatedAt:        nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:        nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func GameSubscriptionRecordToResponse(l logger.Logger, rec *game_record.GameSubscription, instanceIDs []string) (*game_schema.GameSubscriptionResponse, error) {
	l.Debug("mapping game_subscription record to response")
	data, err := GameSubscriptionRecordToResponseData(l, rec, instanceIDs)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameSubscriptionResponse{
		Data: data,
	}, nil
}

func GameSubscriptionViewRecordToResponseData(l logger.Logger, rec *game_record.GameSubscriptionView) (*game_schema.GameSubscriptionResponseData, error) {
	l.Debug("mapping game_subscription_view record to response data")

	instanceLimitPtr, _ := nullint32.ToInt32Ptr(rec.InstanceLimit)

	return &game_schema.GameSubscriptionResponseData{
		ID:               rec.ID,
		GameID:           rec.GameID,
		AccountID:        rec.AccountID,
		GameInstanceIDs:  rec.GameInstanceIDs,
		SubscriptionType: rec.SubscriptionType,
		InstanceLimit:    instanceLimitPtr,
		Status:           rec.Status,
		CreatedAt:        rec.CreatedAt,
		UpdatedAt:        nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:        nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

// GameSubscriptionViewRecordToResponse maps a view record to a response
func GameSubscriptionViewRecordToResponse(l logger.Logger, rec *game_record.GameSubscriptionView) (*game_schema.GameSubscriptionResponse, error) {
	l.Debug("mapping game_subscription_view record to response")
	data, err := GameSubscriptionViewRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &game_schema.GameSubscriptionResponse{
		Data: data,
	}, nil
}

// GameSubscriptionViewRecordsToCollectionResponse maps view records to collection response
// This replaces the anti-pattern of passing a function callback
func GameSubscriptionViewRecordsToCollectionResponse(l logger.Logger, recs []*game_record.GameSubscriptionView) (game_schema.GameSubscriptionCollectionResponse, error) {
	l.Debug("mapping game_subscription_view records to collection response")
	data := []*game_schema.GameSubscriptionResponseData{}
	for _, rec := range recs {
		d, err := GameSubscriptionViewRecordToResponseData(l, rec)
		if err != nil {
			return game_schema.GameSubscriptionCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return game_schema.GameSubscriptionCollectionResponse{
		Data: data,
	}, nil
}
