package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record/account_record"
	"gitlab.com/alienspaces/playbymail/schema/api"
)

// Account

func AccountRequestToRecord(l logger.Logger, r *http.Request, rec *account_record.Account) (*account_record.Account, error) {
	l.Debug("mapping account request to record")

	var req api.AccountRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Email = convert.String(req.Email)
		rec.Name = req.Name
	case server.HttpMethodPut, server.HttpMethodPatch:
		// Email cannot be updated - only name can be changed
		rec.Name = req.Name
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AccountRecordToResponseData(l logger.Logger, rec *account_record.Account) (api.AccountResponseData, error) {
	l.Debug("mapping account record to response data")

	data := api.AccountResponseData{
		ID:        rec.ID,
		Email:     rec.Email,
		Name:      rec.Name,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
	}

	return data, nil
}

func AccountRecordToResponse(l logger.Logger, rec *account_record.Account) (api.AccountResponse, error) {
	data, err := AccountRecordToResponseData(l, rec)
	if err != nil {
		return api.AccountResponse{}, err
	}
	return api.AccountResponse{
		Data: &data,
	}, nil
}

func AccountRecordsToCollectionResponse(l logger.Logger, recs []*account_record.Account) (api.AccountCollectionResponse, error) {
	data := []*api.AccountResponseData{}
	for _, rec := range recs {
		d, err := AccountRecordToResponseData(l, rec)
		if err != nil {
			return api.AccountCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return api.AccountCollectionResponse{
		Data: data,
	}, nil
}

// Authentication

func MapRequestAuthRequestToDomain(req *api.RequestAuthRequest) string {
	return req.Email
}

func MapRequestAuthResponse(status string) *api.RequestAuthResponse {
	return &api.RequestAuthResponse{Status: status}
}

func MapVerifyAuthRequestToDomain(req *api.VerifyAuthRequest) (string, string) {
	return req.Email, req.VerificationToken
}

func MapVerifyAuthResponse(token string) *api.VerifyAuthResponse {
	return &api.VerifyAuthResponse{SessionToken: token}
}
