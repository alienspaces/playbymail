package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameLocationLinkResponseData -
type AdventureGameLocationLinkResponseData struct {
	ID                          string     `json:"id"`
	GameID                      string     `json:"game_id"`
	Name                        string     `json:"name"`
	Description                 string     `json:"description"`
	FromAdventureGameLocationID string     `json:"from_adventure_game_location_id"`
	ToAdventureGameLocationID   string     `json:"to_adventure_game_location_id"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at,omitempty"`
	DeletedAt                   *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationLinkResponse struct {
	Data       *AdventureGameLocationLinkResponseData `json:"data"`
	Error      *common_schema.ResponseError           `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination      `json:"pagination,omitempty"`
}

type AdventureGameLocationLinkCollectionResponse struct {
	Data       []*AdventureGameLocationLinkResponseData `json:"data"`
	Error      *common_schema.ResponseError             `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination        `json:"pagination,omitempty"`
}

type AdventureGameLocationLinkRequest struct {
	common_schema.Request
	Name                        string `json:"name"`
	Description                 string `json:"description"`
	FromAdventureGameLocationID string `json:"from_adventure_game_location_id"`
	ToAdventureGameLocationID   string `json:"to_adventure_game_location_id"`
}

type AdventureGameLocationLinkQueryParams struct {
	common_schema.QueryParamsPagination
	AdventureGameLocationLinkResponseData
}
