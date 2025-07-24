package schema

import "time"

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
	Error      *ResponseError                              `json:"error,omitempty"`
	Pagination *ResponsePagination                         `json:"pagination,omitempty"`
}

type AdventureGameCreaturePlacementCollectionResponse struct {
	Data       []*AdventureGameCreaturePlacementResponseData `json:"data"`
	Error      *ResponseError                                `json:"error,omitempty"`
	Pagination *ResponsePagination                           `json:"pagination,omitempty"`
}

type AdventureGameCreaturePlacementRequest struct {
	Request
	AdventureGameCreatureID string `json:"adventure_game_creature_id"`
	AdventureGameLocationID string `json:"adventure_game_location_id"`
	InitialCount            int    `json:"initial_count"`
}

type AdventureGameCreaturePlacementQueryParams struct {
	QueryParamsPagination
	AdventureGameCreaturePlacementResponseData
}
