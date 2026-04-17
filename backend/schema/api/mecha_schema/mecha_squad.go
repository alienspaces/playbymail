package mecha_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaSquadResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	SquadType   string     `json:"squad_type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MechaSquadResponse struct {
	Data       *MechaSquadResponseData           `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaSquadCollectionResponse struct {
	Data       []*MechaSquadResponseData         `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaSquadRequest struct {
	common_schema.Request
	SquadType   string `json:"squad_type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MechaSquadQueryParams struct {
	common_schema.QueryParamsPagination
	MechaSquadResponseData
}
