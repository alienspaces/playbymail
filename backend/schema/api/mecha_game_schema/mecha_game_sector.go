package mecha_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaGameSectorResponseData struct {
	ID               string     `json:"id"`
	GameID           string     `json:"game_id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	TerrainType      string     `json:"terrain_type"`
	Elevation        int        `json:"elevation"`
	CoverModifier    int        `json:"cover_modifier"`
	IsStartingSector bool       `json:"is_starting_sector"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type MechaGameSectorResponse struct {
	Data       *MechaGameSectorResponseData          `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameSectorCollectionResponse struct {
	Data       []*MechaGameSectorResponseData        `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameSectorRequest struct {
	common_schema.Request
	Name             string `json:"name"`
	Description      string `json:"description"`
	TerrainType      string `json:"terrain_type,omitempty"`
	Elevation        int    `json:"elevation,omitempty"`
	CoverModifier    int    `json:"cover_modifier,omitempty"`
	IsStartingSector bool   `json:"is_starting_sector,omitempty"`
}

type MechaGameSectorQueryParams struct {
	common_schema.QueryParamsPagination
	MechaGameSectorResponseData
}
