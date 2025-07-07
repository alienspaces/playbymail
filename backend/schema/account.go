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
	Response
	*AccountResponseData
}

type AccountCollectionResponse = []*AccountResponseData

type AccountRequest struct {
	Request
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountQueryParams struct {
	QueryParamsPagination
	AccountResponseData
}
