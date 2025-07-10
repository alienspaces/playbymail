package schema

import "time"

// GameCreatureResponseData -
type GameCreatureResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type GameCreatureResponse struct {
	Data       *GameCreatureResponseData `json:"data"`
	Error      *ResponseError            `json:"error,omitempty"`
	Pagination *ResponsePagination       `json:"pagination,omitempty"`
}

type GameCreatureCollectionResponse struct {
	Data       []*GameCreatureResponseData `json:"data"`
	Error      *ResponseError              `json:"error,omitempty"`
	Pagination *ResponsePagination         `json:"pagination,omitempty"`
}

type GameCreatureRequest struct {
	Request
	GameID      string `json:"game_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GameCreatureQueryParams struct {
	QueryParamsPagination
	GameCreatureResponseData
}
