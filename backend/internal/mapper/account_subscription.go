package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nullstring"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/schema/api/account_subscription_schema"
)

func AccountSubscriptionRequestToRecord(l logger.Logger, r *http.Request, rec *account_record.AccountSubscription) (*account_record.AccountSubscription, error) {
	l.Debug("mapping account_subscription request to record")

	var req account_subscription_schema.AccountSubscriptionRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.AccountID = nullstring.FromString(req.AccountID)
		rec.SubscriptionType = req.SubscriptionType
		if req.SubscriptionPeriod != nil {
			rec.SubscriptionPeriod = *req.SubscriptionPeriod
		}
		if req.Status != nil {
			rec.Status = *req.Status
		}
		if req.AutoRenew != nil {
			rec.AutoRenew = *req.AutoRenew
		}
		if req.ExpiresAt != nil {
			rec.ExpiresAt = nulltime.FromTimePtr(req.ExpiresAt)
		}
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.AccountID = nullstring.FromString(req.AccountID)
		rec.SubscriptionType = req.SubscriptionType
		if req.SubscriptionPeriod != nil {
			rec.SubscriptionPeriod = *req.SubscriptionPeriod
		}
		if req.Status != nil {
			rec.Status = *req.Status
		}
		if req.AutoRenew != nil {
			rec.AutoRenew = *req.AutoRenew
		}
		if req.ExpiresAt != nil {
			rec.ExpiresAt = nulltime.FromTimePtr(req.ExpiresAt)
		}
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AccountSubscriptionRecordToResponseData(l logger.Logger, rec *account_record.AccountSubscription) (*account_subscription_schema.AccountSubscriptionResponseData, error) {
	l.Debug("mapping account_subscription record to response data")
	return &account_subscription_schema.AccountSubscriptionResponseData{
		ID:                 rec.ID,
		AccountID:          rec.AccountID.String, // Use String value
		SubscriptionType:   rec.SubscriptionType,
		SubscriptionPeriod: rec.SubscriptionPeriod,
		Status:             rec.Status,
		AutoRenew:          rec.AutoRenew,
		ExpiresAt:          nulltime.ToTimePtr(rec.ExpiresAt),
		CreatedAt:          rec.CreatedAt,
		UpdatedAt:          nulltime.ToTimePtr(rec.UpdatedAt),
		DeletedAt:          nulltime.ToTimePtr(rec.DeletedAt),
	}, nil
}

func AccountSubscriptionRecordToResponse(l logger.Logger, rec *account_record.AccountSubscription) (*account_subscription_schema.AccountSubscriptionResponse, error) {
	l.Debug("mapping account_subscription record to response")
	data, err := AccountSubscriptionRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &account_subscription_schema.AccountSubscriptionResponse{
		Data: data,
	}, nil
}

func AccountSubscriptionRecordsToCollectionResponse(l logger.Logger, recs []*account_record.AccountSubscription) (account_subscription_schema.AccountSubscriptionCollectionResponse, error) {
	l.Debug("mapping account_subscription records to collection response")
	data := []*account_subscription_schema.AccountSubscriptionResponseData{}
	for _, rec := range recs {
		d, err := AccountSubscriptionRecordToResponseData(l, rec)
		if err != nil {
			return account_subscription_schema.AccountSubscriptionCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return account_subscription_schema.AccountSubscriptionCollectionResponse{
		Data: data,
	}, nil
}
