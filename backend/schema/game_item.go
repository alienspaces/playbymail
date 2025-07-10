package schema

import "time"

// GameItemResponseData -
type GameItemResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type GameItemResponse struct {
	Data       *GameItemResponseData `json:"data"`
	Error      *ResponseError        `json:"error,omitempty"`
	Pagination *ResponsePagination   `json:"pagination,omitempty"`
}

type GameItemCollectionResponse struct {
	Data       []*GameItemResponseData `json:"data"`
	Error      *ResponseError          `json:"error,omitempty"`
	Pagination *ResponsePagination     `json:"pagination,omitempty"`
}

type GameItemRequest struct {
	Request
	GameID      string `json:"game_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GameItemQueryParams struct {
	QueryParamsPagination
	GameItemResponseData
}
