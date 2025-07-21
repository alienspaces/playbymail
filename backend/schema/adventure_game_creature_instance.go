package schema

import "time"

// AdventureGameCreatureInstanceResponseData -
type AdventureGameCreatureInstanceResponseData struct {
	ID                     string     `json:"id"`
	GameID                 string     `json:"game_id"`
	GameCreatureID         string     `json:"game_creature_id"`
	GameInstanceID         string     `json:"game_instance_id"`
	GameLocationInstanceID string     `json:"game_location_instance_id"`
	IsAlive                bool       `json:"is_alive"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at,omitempty"`
	DeletedAt              *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameCreatureInstanceResponse struct {
	Data       *AdventureGameCreatureInstanceResponseData `json:"data"`
	Error      *ResponseError                             `json:"error,omitempty"`
	Pagination *ResponsePagination                        `json:"pagination,omitempty"`
}

type AdventureGameCreatureInstanceCollectionResponse struct {
	Data       []*AdventureGameCreatureInstanceResponseData `json:"data"`
	Error      *ResponseError                               `json:"error,omitempty"`
	Pagination *ResponsePagination                          `json:"pagination,omitempty"`
}

type AdventureGameCreatureInstanceRequest struct {
	Request
	GameID                 string `json:"game_id"`
	GameCreatureID         string `json:"game_creature_id"`
	GameInstanceID         string `json:"game_instance_id"`
	GameLocationInstanceID string `json:"game_location_instance_id"`
	IsAlive                bool   `json:"is_alive"`
}
