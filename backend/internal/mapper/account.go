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

func AccountRequestToRecord(l logger.Logger, r *http.Request, rec *account_record.Account) (*account_record.Account, error) {
	l.Debug("mapping account request to record")

	var req account_schema.AccountRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Email = convert.String(req.Email)
	case server.HttpMethodPut, server.HttpMethodPatch:
		// Email cannot be updated
	default:
		return nil, fmt.Errorf("unsupported HTTP method")
	}

	return rec, nil
}

func AccountRecordToResponseData(l logger.Logger, rec *account_record.Account) (*account_schema.AccountResponseData, error) {
	l.Debug("mapping account record to response data")
	return &account_schema.AccountResponseData{
		ID:        rec.ID,
		Email:     rec.Email,
		Status:    rec.Status,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
	}, nil
}

func AccountRecordToResponse(l logger.Logger, rec *account_record.Account) (*account_schema.AccountResponse, error) {
	l.Debug("mapping account record to response")
	data, err := AccountRecordToResponseData(l, rec)
	if err != nil {
		return nil, err
	}
	return &account_schema.AccountResponse{
		Data: data,
	}, nil
}

func AccountRecordsToCollectionResponse(l logger.Logger, recs []*account_record.Account) (account_schema.AccountCollectionResponse, error) {
	l.Debug("mapping account records to collection response")
	data := []*account_schema.AccountResponseData{}
	for _, rec := range recs {
		d, err := AccountRecordToResponseData(l, rec)
		if err != nil {
			return account_schema.AccountCollectionResponse{}, err
		}
		data = append(data, d)
	}
	return account_schema.AccountCollectionResponse{
		Data: data,
	}, nil
}

// Authentication

func MapRequestAuthRequestToDomain(req *account_schema.RequestAuthRequest) string {
	return req.Email
}

func MapRequestAuthResponse(status string) *account_schema.RequestAuthResponse {
	return &account_schema.RequestAuthResponse{Status: status}
}

func MapVerifyAuthRequestToDomain(req *account_schema.VerifyAuthRequest) (string, string) {
	return req.Email, req.VerificationToken
}

func MapVerifyAuthResponse(token string) *account_schema.VerifyAuthResponse {
	return &account_schema.VerifyAuthResponse{SessionToken: token}
}

func MapRefreshSessionResponse(status string, expiresInSeconds int) *account_schema.RefreshSessionResponse {
	return &account_schema.RefreshSessionResponse{
		Status:           status,
		ExpiresInSeconds: expiresInSeconds,
	}
}
