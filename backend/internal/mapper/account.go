package mapper

import (
	"fmt"
	"net/http"

	"gitlab.com/alienspaces/playbymail/core/nulltime"
	"gitlab.com/alienspaces/playbymail/core/server"
	"gitlab.com/alienspaces/playbymail/core/type/logger"
	"gitlab.com/alienspaces/playbymail/internal/record"
	"gitlab.com/alienspaces/playbymail/schema"
)

func AccountRequestToRecord(l logger.Logger, r *http.Request, rec *record.Account) (*record.Account, error) {
	l.Debug("mapping account request to record")

	var req schema.AccountRequest
	_, err := server.ReadRequest(l, r, &req)
	if err != nil {
		return nil, err
	}

	switch server.HttpMethod(r.Method) {
	case server.HttpMethodPost:
		rec.Email = req.Email
		rec.Name = req.Name
	case server.HttpMethodPut, server.HttpMethodPatch:
		rec.Email = req.Email
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
	var data []*schema.AccountResponseData
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

// RequestAuthRequest maps to the /request-auth request schema
//
//go:generate jsonschema -o ../../schema/account.request-auth.request.schema.json -t RequestAuthRequest
//go:generate jsonschema -o ../../schema/account.request-auth.response.schema.json -t RequestAuthResponse
type RequestAuthRequest struct {
	Email string `json:"email"`
}
type RequestAuthResponse struct {
	Status string `json:"status"`
}

// VerifyAuthRequest maps to the /verify-auth request schema
//
//go:generate jsonschema -o ../../schema/account.verify-auth.request.schema.json -t VerifyAuthRequest
//go:generate jsonschema -o ../../schema/account.verify-auth.response.schema.json -t VerifyAuthResponse
type VerifyAuthRequest struct {
	Email             string `json:"email"`
	VerificationToken string `json:"verification_token"`
}
type VerifyAuthResponse struct {
	SessionToken string `json:"session_token"`
}

// Mapper functions (expand as needed)
func MapRequestAuthRequestToDomain(req *RequestAuthRequest) string {
	return req.Email
}
func MapRequestAuthResponse(status string) *RequestAuthResponse {
	return &RequestAuthResponse{Status: status}
}
func MapVerifyAuthRequestToDomain(req *VerifyAuthRequest) (string, string) {
	return req.Email, req.VerificationToken
}
func MapVerifyAuthResponse(token string) *VerifyAuthResponse {
	return &VerifyAuthResponse{SessionToken: token}
}
