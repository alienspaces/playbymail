package api

import "time"

// AdventureGameItemResponseData -
type AdventureGameItemResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameItemResponse struct {
	Data       *AdventureGameItemResponseData `json:"data"`
	Error      *ResponseError                 `json:"error,omitempty"`
	Pagination *ResponsePagination            `json:"pagination,omitempty"`
}

type AdventureGameItemCollectionResponse struct {
	Data       []*AdventureGameItemResponseData `json:"data"`
	Error      *ResponseError                   `json:"error,omitempty"`
	Pagination *ResponsePagination              `json:"pagination,omitempty"`
}

type AdventureGameItemRequest struct {
	Request
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AdventureGameItemQueryParams struct {
	QueryParamsPagination
	AdventureGameItemResponseData
}
