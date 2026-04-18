package mecha_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaGameChassisResponseData struct {
	ID              string     `json:"id"`
	GameID          string     `json:"game_id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	ChassisClass    string     `json:"chassis_class"`
	ArmorPoints     int        `json:"armor_points"`
	StructurePoints int        `json:"structure_points"`
	HeatCapacity    int        `json:"heat_capacity"`
	Speed           int        `json:"speed"`
	SmallSlots      int        `json:"small_slots"`
	MediumSlots     int        `json:"medium_slots"`
	LargeSlots      int        `json:"large_slots"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type MechaGameChassisResponse struct {
	Data       *MechaGameChassisResponseData   `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameChassisCollectionResponse struct {
	Data       []*MechaGameChassisResponseData `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameChassisRequest struct {
	common_schema.Request
	Name            string `json:"name"`
	Description     string `json:"description"`
	ChassisClass    string `json:"chassis_class,omitempty"`
	ArmorPoints     int    `json:"armor_points,omitempty"`
	StructurePoints int    `json:"structure_points,omitempty"`
	HeatCapacity    int    `json:"heat_capacity,omitempty"`
	Speed           int    `json:"speed,omitempty"`
	SmallSlots      int    `json:"small_slots,omitempty"`
	MediumSlots     int    `json:"medium_slots,omitempty"`
	LargeSlots      int    `json:"large_slots,omitempty"`
}

type MechaGameChassisQueryParams struct {
	common_schema.QueryParamsPagination
	MechaGameChassisResponseData
}
