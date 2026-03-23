package mecha_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaWeaponResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Damage      int        `json:"damage"`
	HeatCost    int        `json:"heat_cost"`
	RangeBand   string     `json:"range_band"`
	MountSize   string     `json:"mount_size"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MechaWeaponResponse struct {
	Data       *MechaWeaponResponseData   `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaWeaponCollectionResponse struct {
	Data       []*MechaWeaponResponseData `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaWeaponRequest struct {
	common_schema.Request
	Name        string `json:"name"`
	Description string `json:"description"`
	Damage      int    `json:"damage,omitempty"`
	HeatCost    int    `json:"heat_cost,omitempty"`
	RangeBand   string `json:"range_band,omitempty"`
	MountSize   string `json:"mount_size,omitempty"`
}

type MechaWeaponQueryParams struct {
	common_schema.QueryParamsPagination
	MechaWeaponResponseData
}
