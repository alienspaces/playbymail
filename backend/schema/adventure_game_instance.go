package schema

import "time"

type AdventureGameInstance struct {
	ID        string     `json:"id"`
	GameID    string     `json:"game_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameInstanceResponse struct {
	Data       *AdventureGameInstance `json:"data"`
	Error      *ResponseError         `json:"error,omitempty"`
	Pagination *ResponsePagination    `json:"pagination,omitempty"`
}

type AdventureGameInstanceCollectionResponse struct {
	Data       []*AdventureGameInstance `json:"data"`
	Error      *ResponseError           `json:"error,omitempty"`
	Pagination *ResponsePagination      `json:"pagination,omitempty"`
}

type AdventureGameInstanceRequest struct {
	Request
	GameID string `json:"game_id"`
}
