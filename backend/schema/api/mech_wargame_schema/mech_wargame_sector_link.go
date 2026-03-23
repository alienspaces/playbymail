package mech_wargame_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechWargameSectorLinkResponseData struct {
	ID                       string     `json:"id"`
	GameID                   string     `json:"game_id"`
	FromMechWargameSectorID  string     `json:"from_mech_wargame_sector_id"`
	ToMechWargameSectorID    string     `json:"to_mech_wargame_sector_id"`
	CoverModifier            int        `json:"cover_modifier"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty"`
	DeletedAt                *time.Time `json:"deleted_at,omitempty"`
}

type MechWargameSectorLinkResponse struct {
	Data       *MechWargameSectorLinkResponseData `json:"data"`
	Error      *common_schema.ResponseError        `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination   `json:"pagination,omitempty"`
}

type MechWargameSectorLinkCollectionResponse struct {
	Data       []*MechWargameSectorLinkResponseData `json:"data"`
	Error      *common_schema.ResponseError          `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination     `json:"pagination,omitempty"`
}

type MechWargameSectorLinkRequest struct {
	common_schema.Request
	FromMechWargameSectorID string `json:"from_mech_wargame_sector_id"`
	ToMechWargameSectorID   string `json:"to_mech_wargame_sector_id"`
	CoverModifier           int    `json:"cover_modifier,omitempty"`
}

type MechWargameSectorLinkQueryParams struct {
	common_schema.QueryParamsPagination
	MechWargameSectorLinkResponseData
}
