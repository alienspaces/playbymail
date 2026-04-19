package mecha_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type WeaponConfigEntry struct {
	WeaponID     string `json:"weapon_id"`
	SlotLocation string `json:"slot_location"`
}

type EquipmentConfigEntry struct {
	EquipmentID  string `json:"equipment_id"`
	SlotLocation string `json:"slot_location"`
}

type MechaGameSquadMechResponseData struct {
	ID             string              `json:"id"`
	GameID         string              `json:"game_id"`
	MechaGameSquadID   string              `json:"mecha_game_squad_id"`
	MechaGameChassisID string              `json:"mecha_game_chassis_id"`
	Callsign        string                 `json:"callsign"`
	WeaponConfig    []WeaponConfigEntry    `json:"weapon_config"`
	EquipmentConfig []EquipmentConfigEntry `json:"equipment_config"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt      *time.Time          `json:"updated_at,omitempty"`
	DeletedAt      *time.Time          `json:"deleted_at,omitempty"`
}

type MechaGameSquadMechResponse struct {
	Data       *MechaGameSquadMechResponseData       `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameSquadMechCollectionResponse struct {
	Data       []*MechaGameSquadMechResponseData     `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameSquadMechRequest struct {
	common_schema.Request
	MechaGameChassisID string                 `json:"mecha_game_chassis_id"`
	Callsign           string                 `json:"callsign"`
	WeaponConfig       []WeaponConfigEntry    `json:"weapon_config,omitempty"`
	EquipmentConfig    []EquipmentConfigEntry `json:"equipment_config,omitempty"`
}

type MechaGameSquadMechQueryParams struct {
	common_schema.QueryParamsPagination
	MechaGameSquadMechResponseData
}
