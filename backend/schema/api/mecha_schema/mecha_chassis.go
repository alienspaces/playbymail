package mecha_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaChassisResponseData struct {
	ID              string     `json:"id"`
	GameID          string     `json:"game_id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	ChassisClass    string     `json:"chassis_class"`
	ArmorPoints     int        `json:"armor_points"`
	StructurePoints int        `json:"structure_points"`
	HeatCapacity    int        `json:"heat_capacity"`
	Speed           int        `json:"speed"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type MechaChassisResponse struct {
	Data       *MechaChassisResponseData   `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaChassisCollectionResponse struct {
	Data       []*MechaChassisResponseData `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaChassisRequest struct {
	common_schema.Request
	Name            string `json:"name"`
	Description     string `json:"description"`
	ChassisClass    string `json:"chassis_class,omitempty"`
	ArmorPoints     int    `json:"armor_points,omitempty"`
	StructurePoints int    `json:"structure_points,omitempty"`
	HeatCapacity    int    `json:"heat_capacity,omitempty"`
	Speed           int    `json:"speed,omitempty"`
}

type MechaChassisQueryParams struct {
	common_schema.QueryParamsPagination
	MechaChassisResponseData
}
