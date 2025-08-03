package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameLocationInstanceResponseData -
type AdventureGameLocationInstanceResponseData struct {
	ID                      string     `json:"id"`
	GameID                  string     `json:"game_id"`
	GameInstanceID          string     `json:"game_instance_id"`
	AdventureGameLocationID string     `json:"adventure_game_location_id"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`
	DeletedAt               *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationInstanceResponse struct {
	Data       *AdventureGameLocationInstanceResponseData `json:"data"`
	Error      *common_schema.ResponseError               `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination          `json:"pagination,omitempty"`
}

type AdventureGameLocationInstanceCollectionResponse struct {
	Data       []*AdventureGameLocationInstanceResponseData `json:"data"`
	Error      *common_schema.ResponseError                 `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination            `json:"pagination,omitempty"`
}
