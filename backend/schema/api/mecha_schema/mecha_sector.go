package mecha_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaSectorResponseData struct {
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

type MechaSectorResponse struct {
	Data       *MechaSectorResponseData          `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaSectorCollectionResponse struct {
	Data       []*MechaSectorResponseData        `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaSectorRequest struct {
	common_schema.Request
	Name             string `json:"name"`
	Description      string `json:"description"`
	TerrainType      string `json:"terrain_type,omitempty"`
	Elevation        int    `json:"elevation,omitempty"`
	CoverModifier    int    `json:"cover_modifier,omitempty"`
	IsStartingSector bool   `json:"is_starting_sector,omitempty"`
}

type MechaSectorQueryParams struct {
	common_schema.QueryParamsPagination
	MechaSectorResponseData
}
