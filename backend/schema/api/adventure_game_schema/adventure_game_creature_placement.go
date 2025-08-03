package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameCreaturePlacementResponseData -
type AdventureGameCreaturePlacementResponseData struct {
	ID                      string     `json:"id"`
	GameID                  string     `json:"game_id"`
	AdventureGameCreatureID string     `json:"adventure_game_creature_id"`
	AdventureGameLocationID string     `json:"adventure_game_location_id"`
	InitialCount            int        `json:"initial_count"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`
	DeletedAt               *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameCreaturePlacementResponse struct {
	Data       *AdventureGameCreaturePlacementResponseData `json:"data"`
	Error      *common_schema.ResponseError                `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination           `json:"pagination,omitempty"`
}

type AdventureGameCreaturePlacementCollectionResponse struct {
	Data       []*AdventureGameCreaturePlacementResponseData `json:"data"`
	Error      *common_schema.ResponseError                  `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination             `json:"pagination,omitempty"`
}

type AdventureGameCreaturePlacementRequest struct {
	common_schema.Request
	AdventureGameCreatureID string `json:"adventure_game_creature_id"`
	AdventureGameLocationID string `json:"adventure_game_location_id"`
	InitialCount            int    `json:"initial_count,omitempty"`
}

type AdventureGameCreaturePlacementQueryParams struct {
	common_schema.QueryParamsPagination
	AdventureGameCreaturePlacementResponseData
}
