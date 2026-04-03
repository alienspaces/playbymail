package mecha_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaLanceResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	LanceType   string     `json:"lance_type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MechaLanceResponse struct {
	Data       *MechaLanceResponseData           `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaLanceCollectionResponse struct {
	Data       []*MechaLanceResponseData         `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechaLanceRequest struct {
	common_schema.Request
	LanceType   string `json:"lance_type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MechaLanceQueryParams struct {
	common_schema.QueryParamsPagination
	MechaLanceResponseData
}
