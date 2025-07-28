package api

import "time"

// AdventureGameLocationResponseData -
type AdventureGameLocationResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationResponse struct {
	Data       *AdventureGameLocationResponseData `json:"data"`
	Error      *ResponseError                     `json:"error,omitempty"`
	Pagination *ResponsePagination                `json:"pagination,omitempty"`
}

type AdventureGameLocationCollectionResponse struct {
	Data       []*AdventureGameLocationResponseData `json:"data"`
	Error      *ResponseError                       `json:"error,omitempty"`
	Pagination *ResponsePagination                  `json:"pagination,omitempty"`
}

type AdventureGameLocationRequest struct {
	Request
	GameID      string `json:"game_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AdventureGameLocationQueryParams struct {
	QueryParamsPagination
	AdventureGameLocationResponseData
}
