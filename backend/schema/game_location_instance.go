package schema

import "time"

// GameLocationInstanceResponseData -
type GameLocationInstanceResponseData struct {
	ID             string     `json:"id"`
	GameID         string     `json:"game_id"`
	GameInstanceID string     `json:"game_instance_id"`
	GameLocationID string     `json:"game_location_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type GameLocationInstanceResponse struct {
	Data       *GameLocationInstanceResponseData `json:"data"`
	Error      *ResponseError                    `json:"error,omitempty"`
	Pagination *ResponsePagination               `json:"pagination,omitempty"`
}

type GameLocationInstanceCollectionResponse struct {
	Data       []*GameLocationInstanceResponseData `json:"data"`
	Error      *ResponseError                      `json:"error,omitempty"`
	Pagination *ResponsePagination                 `json:"pagination,omitempty"`
}

type GameLocationInstanceRequest struct {
	Request
	GameID         string `json:"game_id"`
	GameInstanceID string `json:"game_instance_id"`
	GameLocationID string `json:"game_location_id"`
}
