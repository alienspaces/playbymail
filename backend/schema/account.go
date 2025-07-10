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
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountQueryParams struct {
	QueryParamsPagination
	AccountResponseData
}
