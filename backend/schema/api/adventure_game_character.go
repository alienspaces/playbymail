package api

import "time"

// AdventureGameCharacterResponseData -
type AdventureGameCharacterResponseData struct {
	ID        string     `json:"id"`
	GameID    string     `json:"game_id"`
	AccountID string     `json:"account_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameCharacterResponse struct {
	Data       *AdventureGameCharacterResponseData `json:"data"`
	Error      *ResponseError                      `json:"error,omitempty"`
	Pagination *ResponsePagination                 `json:"pagination,omitempty"`
}

type AdventureGameCharacterCollectionResponse struct {
	Data       []*AdventureGameCharacterResponseData `json:"data"`
	Error      *ResponseError                        `json:"error,omitempty"`
	Pagination *ResponsePagination                   `json:"pagination,omitempty"`
}

type AdventureGameCharacterRequest struct {
	Request
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
}

type AdventureGameCharacterQueryParams struct {
	QueryParamsPagination
	AdventureGameCharacterResponseData
}
