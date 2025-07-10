package schema

import "time"

// GameCharacterResponseData -
type GameCharacterResponseData struct {
	ID        string     `json:"id"`
	GameID    string     `json:"game_id"`
	AccountID string     `json:"account_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type GameCharacterResponse struct {
	Data       *GameCharacterResponseData `json:"data"`
	Error      *ResponseError             `json:"error,omitempty"`
	Pagination *ResponsePagination        `json:"pagination,omitempty"`
}

type GameCharacterCollectionResponse struct {
	Data       []*GameCharacterResponseData `json:"data"`
	Error      *ResponseError               `json:"error,omitempty"`
	Pagination *ResponsePagination          `json:"pagination,omitempty"`
}

type GameCharacterRequest struct {
	Request
	GameID    string `json:"game_id"`
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
}

type GameCharacterQueryParams struct {
	QueryParamsPagination
	GameCharacterResponseData
}
