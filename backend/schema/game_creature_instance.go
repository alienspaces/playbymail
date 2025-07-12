package schema

import "time"

// GameCreatureInstanceResponseData -
type GameCreatureInstanceResponseData struct {
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

type GameCreatureInstanceResponse struct {
	Data       *GameCreatureInstanceResponseData `json:"data"`
	Error      *ResponseError                    `json:"error,omitempty"`
	Pagination *ResponsePagination               `json:"pagination,omitempty"`
}

type GameCreatureInstanceCollectionResponse struct {
	Data       []*GameCreatureInstanceResponseData `json:"data"`
	Error      *ResponseError                      `json:"error,omitempty"`
	Pagination *ResponsePagination                 `json:"pagination,omitempty"`
}

type GameCreatureInstanceRequest struct {
	Request
	GameID                 string `json:"game_id"`
	GameCreatureID         string `json:"game_creature_id"`
	GameInstanceID         string `json:"game_instance_id"`
	GameLocationInstanceID string `json:"game_location_instance_id"`
	IsAlive                bool   `json:"is_alive"`
}
