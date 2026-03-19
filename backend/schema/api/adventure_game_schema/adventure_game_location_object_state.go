package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameLocationObjectStateResponseData -
type AdventureGameLocationObjectStateResponseData struct {
	ID                            string     `json:"id"`
	GameID                        string     `json:"game_id"`
	AdventureGameLocationObjectID string     `json:"adventure_game_location_object_id"`
	Name                          string     `json:"name"`
	Description                   string     `json:"description"`
	SortOrder                     int        `json:"sort_order"`
	CreatedAt                     time.Time  `json:"created_at"`
	UpdatedAt                     *time.Time `json:"updated_at,omitempty"`
	DeletedAt                     *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationObjectStateResponse struct {
	Data       *AdventureGameLocationObjectStateResponseData `json:"data"`
	Error      *common_schema.ResponseError                 `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination            `json:"pagination,omitempty"`
}

type AdventureGameLocationObjectStateCollectionResponse struct {
	Data       []*AdventureGameLocationObjectStateResponseData `json:"data"`
	Error      *common_schema.ResponseError                   `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination              `json:"pagination,omitempty"`
}

type AdventureGameLocationObjectStateRequest struct {
	common_schema.Request
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
}
