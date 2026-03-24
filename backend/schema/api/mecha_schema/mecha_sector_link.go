package mecha_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaSectorLinkResponseData struct {
	ID                       string     `json:"id"`
	GameID                   string     `json:"game_id"`
	FromMechaSectorID  string     `json:"from_mecha_sector_id"`
	ToMechaSectorID    string     `json:"to_mecha_sector_id"`
	CoverModifier            int        `json:"cover_modifier"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty"`
	DeletedAt                *time.Time `json:"deleted_at,omitempty"`
}

type MechaSectorLinkResponse struct {
	Data       *MechaSectorLinkResponseData `json:"data"`
	Error      *common_schema.ResponseError        `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination   `json:"pagination,omitempty"`
}

type MechaSectorLinkCollectionResponse struct {
	Data       []*MechaSectorLinkResponseData `json:"data"`
	Error      *common_schema.ResponseError          `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination     `json:"pagination,omitempty"`
}

type MechaSectorLinkRequest struct {
	common_schema.Request
	FromMechaSectorID string `json:"from_mecha_sector_id"`
	ToMechaSectorID   string `json:"to_mecha_sector_id"`
	CoverModifier           int    `json:"cover_modifier,omitempty"`
}

type MechaSectorLinkQueryParams struct {
	common_schema.QueryParamsPagination
	MechaSectorLinkResponseData
}
