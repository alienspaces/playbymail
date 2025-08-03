package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameItemResponseData -
type AdventureGameItemResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameItemResponse struct {
	Data       *AdventureGameItemResponseData    `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type AdventureGameItemCollectionResponse struct {
	Data       []*AdventureGameItemResponseData  `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type AdventureGameItemRequest struct {
	common_schema.Request
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AdventureGameItemQueryParams struct {
	common_schema.QueryParamsPagination
	AdventureGameItemResponseData
}
