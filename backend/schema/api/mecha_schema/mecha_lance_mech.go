package mecha_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type WeaponConfigEntry struct {
	WeaponID     string `json:"weapon_id"`
	SlotLocation string `json:"slot_location"`
}

type MechaLanceMechResponseData struct {
	ID                   string              `json:"id"`
	GameID               string              `json:"game_id"`
	MechaLanceID   string              `json:"mecha_lance_id"`
	MechaChassisID string              `json:"mecha_chassis_id"`
	Callsign             string              `json:"callsign"`
	WeaponConfig         []WeaponConfigEntry `json:"weapon_config"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            *time.Time          `json:"updated_at,omitempty"`
	DeletedAt            *time.Time          `json:"deleted_at,omitempty"`
}

type MechaLanceMechResponse struct {
	Data       *MechaLanceMechResponseData `json:"data"`
	Error      *common_schema.ResponseError       `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination  `json:"pagination,omitempty"`
}

type MechaLanceMechCollectionResponse struct {
	Data       []*MechaLanceMechResponseData `json:"data"`
	Error      *common_schema.ResponseError         `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination    `json:"pagination,omitempty"`
}

type MechaLanceMechRequest struct {
	common_schema.Request
	MechaChassisID string              `json:"mecha_chassis_id"`
	Callsign             string              `json:"callsign"`
	WeaponConfig         []WeaponConfigEntry `json:"weapon_config,omitempty"`
}

type MechaLanceMechQueryParams struct {
	common_schema.QueryParamsPagination
	MechaLanceMechResponseData
}
