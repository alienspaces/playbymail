package mech_wargame_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type WeaponConfigEntry struct {
	WeaponID     string `json:"weapon_id"`
	SlotLocation string `json:"slot_location"`
}

type MechWargameLanceMechResponseData struct {
	ID                   string              `json:"id"`
	GameID               string              `json:"game_id"`
	MechWargameLanceID   string              `json:"mech_wargame_lance_id"`
	MechWargameChassisID string              `json:"mech_wargame_chassis_id"`
	Callsign             string              `json:"callsign"`
	WeaponConfig         []WeaponConfigEntry `json:"weapon_config"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            *time.Time          `json:"updated_at,omitempty"`
	DeletedAt            *time.Time          `json:"deleted_at,omitempty"`
}

type MechWargameLanceMechResponse struct {
	Data       *MechWargameLanceMechResponseData `json:"data"`
	Error      *common_schema.ResponseError       `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination  `json:"pagination,omitempty"`
}

type MechWargameLanceMechCollectionResponse struct {
	Data       []*MechWargameLanceMechResponseData `json:"data"`
	Error      *common_schema.ResponseError         `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination    `json:"pagination,omitempty"`
}

type MechWargameLanceMechRequest struct {
	common_schema.Request
	MechWargameChassisID string              `json:"mech_wargame_chassis_id"`
	Callsign             string              `json:"callsign"`
	WeaponConfig         []WeaponConfigEntry `json:"weapon_config,omitempty"`
}

type MechWargameLanceMechQueryParams struct {
	common_schema.QueryParamsPagination
	MechWargameLanceMechResponseData
}
