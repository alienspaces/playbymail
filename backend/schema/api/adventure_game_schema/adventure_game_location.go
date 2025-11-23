package adventure_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// AdventureGameLocationResponseData -
type AdventureGameLocationResponseData struct {
	ID                 string     `json:"id"`
	GameID             string     `json:"game_id"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	IsStartingLocation bool       `json:"is_starting_location"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type AdventureGameLocationResponse struct {
	Data       *AdventureGameLocationResponseData `json:"data"`
	Error      *common_schema.ResponseError       `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination  `json:"pagination,omitempty"`
}

type AdventureGameLocationCollectionResponse struct {
	Data       []*AdventureGameLocationResponseData `json:"data"`
	Error      *common_schema.ResponseError         `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination    `json:"pagination,omitempty"`
}

type AdventureGameLocationRequest struct {
	common_schema.Request
	Name               string `json:"name"`
	Description        string `json:"description"`
	IsStartingLocation bool   `json:"is_starting_location,omitempty"`
}

type AdventureGameLocationQueryParams struct {
	common_schema.QueryParamsPagination
	AdventureGameLocationResponseData
}
