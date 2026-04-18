package mecha_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaGameSectorLinkResponseData struct {
	ID                string     `json:"id"`
	GameID            string     `json:"game_id"`
	FromMechaGameSectorID string     `json:"from_mecha_game_sector_id"`
	ToMechaGameSectorID   string     `json:"to_mecha_game_sector_id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

type MechaGameSectorLinkResponse struct {
	Data       *MechaGameSectorLinkResponseData      `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameSectorLinkCollectionResponse struct {
	Data       []*MechaGameSectorLinkResponseData    `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameSectorLinkRequest struct {
	common_schema.Request
	FromMechaGameSectorID string `json:"from_mecha_game_sector_id"`
	ToMechaGameSectorID   string `json:"to_mecha_game_sector_id"`
}

type MechaGameSectorLinkQueryParams struct {
	common_schema.QueryParamsPagination
	MechaGameSectorLinkResponseData
}
