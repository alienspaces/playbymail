package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/schema/api/account_schema"
)

func AccountUserRequestToRecord(l logger.Logger, r *http.Request, rec *account_record.AccountUser) (*account_record.AccountUser, error) {
	l.Debug("mapping account user request to record")

	var req account_schema.AccountUserRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Email = convert.String(req.Email)
		if req.Status != nil {
			rec.Status = *req.Status
		} else {
			rec.Status = account_record.AccountUserStatusPendingApproval
		}
	case server.HttpMethodPut, server.HttpMethodPatch:
		if req.Status != nil {
			rec.Status = *req.Status
		}
		// Email typically not changeable via simple update, but if needed:
		// if req.Email != nil { rec.Email = *req.Email }
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AccountUserRecordToResponseData(l logger.Logger, rec *account_record.AccountUser) (*account_schema.AccountUserResponseData, error) {
	l.Debug("mapping account user record to response data")
	return &account_schema.AccountUserResponseData{
		ID:        rec.ID,
		AccountID: rec.AccountID,
		Email:     rec.Email,
		Status:    rec.Status,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func AccountUserRecordToResponse(l logger.Logger, rec *account_record.AccountUser) (*account_schema.AccountUserResponse, error) {
	l.Debug("mapping account user record to response")
	data, err := AccountUserRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &account_schema.AccountUserResponse{
		Data: data,
	}, nil
}

func AccountUserRecordsToCollectionResponse(l logger.Logger, recs []*account_record.AccountUser) (account_schema.AccountUserCollectionResponse, error) {
	l.Debug("mapping account user records to collection response")
	data := []*account_schema.AccountUserResponseData{}
	for _, rec := range recs {
		d, err := AccountUserRecordToResponseData(l, rec)
		if err != nil {
			return account_schema.AccountUserCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return account_schema.AccountUserCollectionResponse{
		Data: data,
	}, nil
}
