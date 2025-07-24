package schema

import "time"

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
	Error      *ResponseError                          `json:"error,omitempty"`
	Pagination *ResponsePagination                     `json:"pagination,omitempty"`
}

type AdventureGameItemPlacementCollectionResponse struct {
	Data       []*AdventureGameItemPlacementResponseData `json:"data"`
	Error      *ResponseError                            `json:"error,omitempty"`
	Pagination *ResponsePagination                       `json:"pagination,omitempty"`
}

type AdventureGameItemPlacementRequest struct {
	Request
	AdventureGameItemID     string `json:"adventure_game_item_id"`
	AdventureGameLocationID string `json:"adventure_game_location_id"`
	InitialCount            int    `json:"initial_count"`
}

type AdventureGameItemPlacementQueryParams struct {
	QueryParamsPagination
	AdventureGameItemPlacementResponseData
}
