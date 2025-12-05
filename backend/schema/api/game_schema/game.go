package game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// GameResponseData -
type GameResponseData struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	GameType          string     `json:"game_type"`
	TurnDurationHours int        `json:"turn_duration_hours"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

type GameResponse struct {
	Data       *GameResponseData                 `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameCollectionResponse struct {
	Data       []*GameResponseData               `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameRequest struct {
	common_schema.Request
	Name              string `json:"name"`
	Description       string `json:"description"`
	GameType          string `json:"game_type"`
	TurnDurationHours int    `json:"turn_duration_hours"`
}

type GameQueryParams struct {
	common_schema.QueryParamsPagination
	GameResponseData
}
