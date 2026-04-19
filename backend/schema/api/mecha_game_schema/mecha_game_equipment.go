package mecha_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaGameEquipmentResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	MountSize   string     `json:"mount_size"`
	EffectKind  string     `json:"effect_kind"`
	Magnitude   int        `json:"magnitude"`
	HeatCost    int        `json:"heat_cost"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MechaGameEquipmentResponse struct {
	Data       *MechaGameEquipmentResponseData   `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameEquipmentCollectionResponse struct {
	Data       []*MechaGameEquipmentResponseData `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameEquipmentRequest struct {
	common_schema.Request
	Name        string `json:"name"`
	Description string `json:"description"`
	MountSize   string `json:"mount_size,omitempty"`
	EffectKind  string `json:"effect_kind,omitempty"`
	Magnitude   int    `json:"magnitude,omitempty"`
	HeatCost    int    `json:"heat_cost,omitempty"`
}

type MechaGameEquipmentQueryParams struct {
	common_schema.QueryParamsPagination
	MechaGameEquipmentResponseData
}
