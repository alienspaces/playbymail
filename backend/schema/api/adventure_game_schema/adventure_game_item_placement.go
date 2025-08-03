package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameItemPlacementResponseData -
type AdventureGameItemPlacementResponseData struct {
	ID                      string     `json:"id"`
	GameID                  string     `json:"game_id"`
	AdventureGameItemID     string     `json:"adventure_game_item_id"`
	AdventureGameLocationID string     `json:"adventure_game_location_id"`
	InitialCount            int        `json:"initial_count"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`
	DeletedAt               *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameItemPlacementResponse struct {
	Data       *AdventureGameItemPlacementResponseData `json:"data"`
	Error      *common_schema.ResponseError            `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination       `json:"pagination,omitempty"`
}

type AdventureGameItemPlacementCollectionResponse struct {
	Data       []*AdventureGameItemPlacementResponseData `json:"data"`
	Error      *common_schema.ResponseError              `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination         `json:"pagination,omitempty"`
}

type AdventureGameItemPlacementRequest struct {
	common_schema.Request
	AdventureGameItemID     string `json:"adventure_game_item_id"`
	AdventureGameLocationID string `json:"adventure_game_location_id"`
	InitialCount            int    `json:"initial_count"`
}

type AdventureGameItemPlacementQueryParams struct {
	common_schema.QueryParamsPagination
	AdventureGameItemPlacementResponseData
}
