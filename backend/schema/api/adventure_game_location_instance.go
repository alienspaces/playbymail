package api

import "time"

// AdventureGameLocationInstanceResponseData -
type AdventureGameLocationInstanceResponseData struct {
	ID             string     `json:"id"`
	GameID         string     `json:"game_id"`
	GameInstanceID string     `json:"game_instance_id"`
	GameLocationID string     `json:"game_location_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationInstanceResponse struct {
	Data       *AdventureGameLocationInstanceResponseData `json:"data"`
	Error      *ResponseError                             `json:"error,omitempty"`
	Pagination *ResponsePagination                        `json:"pagination,omitempty"`
}

type AdventureGameLocationInstanceCollectionResponse struct {
	Data       []*AdventureGameLocationInstanceResponseData `json:"data"`
	Error      *ResponseError                               `json:"error,omitempty"`
	Pagination *ResponsePagination                          `json:"pagination,omitempty"`
}

type AdventureGameLocationInstanceRequest struct {
	Request
	GameID         string `json:"game_id"`
	GameInstanceID string `json:"game_instance_id"`
	GameLocationID string `json:"game_location_id"`
}
