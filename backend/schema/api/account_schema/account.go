package account_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AccountResponseData -
type AccountResponseData struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type AccountResponse struct {
	Data       *AccountResponseData              `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type AccountCollectionResponse struct {
	Data       []*AccountResponseData            `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type AccountRequest struct {
	common_schema.Request
	Email *string `json:"email,omitempty"`
}

type AccountQueryParams struct {
	common_schema.QueryParamsPagination
	AccountResponseData
}

// RequestAuthRequest maps to the /request-auth request schema
type RequestAuthRequest struct {
	Email string `json:"email"`
}

type RequestAuthResponse struct {
	Status string `json:"status"`
}

// VerifyAuthRequest maps to the /verify-auth request schema
type VerifyAuthRequest struct {
	Email             string `json:"email"`
	VerificationToken string `json:"verification_token"`
}

type VerifyAuthResponse struct {
	SessionToken string `json:"session_token"`
}

// AccountContactResponseData -
type AccountContactResponseData struct {
	ID                 string     `json:"id"`
	AccountID          string     `json:"account_id"`
	Name               string     `json:"name"`
	PostalAddressLine1 string     `json:"postal_address_line1"`
	PostalAddressLine2 string     `json:"postal_address_line2,omitempty"`
	StateProvince      string     `json:"state_province"`
	Country            string     `json:"country"`
	PostalCode         string     `json:"postal_code"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

type AccountContactResponse struct {
	Data       *AccountContactResponseData       `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type AccountContactCollectionResponse struct {
	Data       []*AccountContactResponseData     `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type AccountContactRequest struct {
	common_schema.Request
	Name               string `json:"name"`
	PostalAddressLine1 string `json:"postal_address_line1"`
	PostalAddressLine2 string `json:"postal_address_line2,omitempty"`
	StateProvince      string `json:"state_province"`
	Country            string `json:"country"`
	PostalCode         string `json:"postal_code"`
}

type AccountContactQueryParams struct {
	common_schema.QueryParamsPagination
	AccountContactResponseData
}
