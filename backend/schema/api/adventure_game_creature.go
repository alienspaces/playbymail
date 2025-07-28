package api

import "time"

// AdventureGameCreatureResponseData -
type AdventureGameCreatureResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameCreatureResponse struct {
	Data       *AdventureGameCreatureResponseData `json:"data"`
	Error      *ResponseError                     `json:"error,omitempty"`
	Pagination *ResponsePagination                `json:"pagination,omitempty"`
}

type AdventureGameCreatureCollectionResponse struct {
	Data       []*AdventureGameCreatureResponseData `json:"data"`
	Error      *ResponseError                       `json:"error,omitempty"`
	Pagination *ResponsePagination                  `json:"pagination,omitempty"`
}

type AdventureGameCreatureRequest struct {
	Request
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AdventureGameCreatureQueryParams struct {
	QueryParamsPagination
	AdventureGameCreatureResponseData
}
