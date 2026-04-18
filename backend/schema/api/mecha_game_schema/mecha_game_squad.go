package mecha_game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaGameSquadResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	SquadType   string     `json:"squad_type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MechaGameSquadResponse struct {
	Data       *MechaGameSquadResponseData           `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameSquadCollectionResponse struct {
	Data       []*MechaGameSquadResponseData         `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaGameSquadRequest struct {
	common_schema.Request
	SquadType   string `json:"squad_type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MechaGameSquadQueryParams struct {
	common_schema.QueryParamsPagination
	MechaGameSquadResponseData
}
