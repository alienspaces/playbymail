package mecha_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechaComputerOpponentResponseData struct {
	ID          string     `json:"id"`
	GameID      string     `json:"game_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Aggression  int        `json:"aggression"`
	IQ          int        `json:"iq"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MechaComputerOpponentResponse struct {
	Data       *MechaComputerOpponentResponseData `json:"data"`
	Error      *common_schema.ResponseError       `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination  `json:"pagination,omitempty"`
}

type MechaComputerOpponentCollectionResponse struct {
	Data       []*MechaComputerOpponentResponseData `json:"data"`
	Error      *common_schema.ResponseError         `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination    `json:"pagination,omitempty"`
}

type MechaComputerOpponentRequest struct {
	common_schema.Request
	Name        string `json:"name"`
	Description string `json:"description"`
	Aggression  int    `json:"aggression"`
	IQ          int    `json:"iq"`
}

type MechaComputerOpponentQueryParams struct {
	common_schema.QueryParamsPagination
	MechaComputerOpponentResponseData
}
