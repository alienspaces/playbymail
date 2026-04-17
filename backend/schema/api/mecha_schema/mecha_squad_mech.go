package mecha_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type WeaponConfigEntry struct {
	WeaponID     string `json:"weapon_id"`
	SlotLocation string `json:"slot_location"`
}

type MechaSquadMechResponseData struct {
	ID             string              `json:"id"`
	GameID         string              `json:"game_id"`
	MechaSquadID   string              `json:"mecha_squad_id"`
	MechaChassisID string              `json:"mecha_chassis_id"`
	Callsign       string              `json:"callsign"`
	WeaponConfig   []WeaponConfigEntry `json:"weapon_config"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      *time.Time          `json:"updated_at,omitempty"`
	DeletedAt      *time.Time          `json:"deleted_at,omitempty"`
}

type MechaSquadMechResponse struct {
	Data       *MechaSquadMechResponseData       `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaSquadMechCollectionResponse struct {
	Data       []*MechaSquadMechResponseData     `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaSquadMechRequest struct {
	common_schema.Request
	MechaChassisID string              `json:"mecha_chassis_id"`
	Callsign       string              `json:"callsign"`
	WeaponConfig   []WeaponConfigEntry `json:"weapon_config,omitempty"`
}

type MechaSquadMechQueryParams struct {
	common_schema.QueryParamsPagination
	MechaSquadMechResponseData
}
