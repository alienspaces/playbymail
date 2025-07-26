package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/convert"
	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

// Account

func AccountRequestToRecord(l logger.Logger, r *http.Request, rec *record.Account) (*record.Account, error) {
	l.Debug("mapping account request to record")

	var req schema.AccountRequest
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

func AccountRecordToResponseData(l logger.Logger, rec *record.Account) (schema.AccountResponseData, error) {
	l.Debug("mapping account record to response data")

	data := schema.AccountResponseData{
		ID:        rec.ID,
		Email:     rec.Email,
		Name:      rec.Name,
		CreatedAt: rec.CreatedAt,
		UpdatedAt: nulltime.ToTimePtr(rec.UpdatedAt),
	}

	return data, nil
}

func AccountRecordToResponse(l logger.Logger, rec *record.Account) (schema.AccountResponse, error) {
	data, err := AccountRecordToResponseData(l, rec)
	if err != nil {
		return schema.AccountResponse{}, err
	}
	return schema.AccountResponse{
		Data: &data,
	}, nil
}

func AccountRecordsToCollectionResponse(l logger.Logger, recs []*record.Account) (schema.AccountCollectionResponse, error) {
	data := []*schema.AccountResponseData{}
	for _, rec := range recs {
		d, err := AccountRecordToResponseData(l, rec)
		if err != nil {
			return schema.AccountCollectionResponse{}, err
		}
		data = append(data, &d)
	}
	return schema.AccountCollectionResponse{
		Data: data,
	}, nil
}

// Authentication

func MapRequestAuthRequestToDomain(req *schema.RequestAuthRequest) string {
	return req.Email
}

func MapRequestAuthResponse(status string) *schema.RequestAuthResponse {
	return &schema.RequestAuthResponse{Status: status}
}

func MapVerifyAuthRequestToDomain(req *schema.VerifyAuthRequest) (string, string) {
	return req.Email, req.VerificationToken
}

func MapVerifyAuthResponse(token string) *schema.VerifyAuthResponse {
	return &schema.VerifyAuthResponse{SessionToken: token}
}
