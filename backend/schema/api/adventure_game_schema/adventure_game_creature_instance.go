package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

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
	Error      *common_schema.ResponseError               `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination          `json:"pagination,omitempty"`
}

type AdventureGameCreatureInstanceCollectionResponse struct {
	Data       []*AdventureGameCreatureInstanceResponseData `json:"data"`
	Error      *common_schema.ResponseError                 `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination            `json:"pagination,omitempty"`
}
