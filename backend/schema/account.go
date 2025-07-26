package schema

import (
	"time"
)

// AccountResponseData -
type AccountResponseData struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type AccountResponse struct {
	Data       *AccountResponseData `json:"data"`
	Error      *ResponseError       `json:"error,omitempty"`
	Pagination *ResponsePagination  `json:"pagination,omitempty"`
}

type AccountCollectionResponse struct {
	Data       []*AccountResponseData `json:"data"`
	Error      *ResponseError         `json:"error,omitempty"`
	Pagination *ResponsePagination    `json:"pagination,omitempty"`
}

type AccountRequest struct {
	Request
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
}

type AccountQueryParams struct {
	QueryParamsPagination
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
