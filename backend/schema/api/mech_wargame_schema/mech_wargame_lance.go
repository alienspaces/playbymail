package mech_wargame_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type MechWargameLanceResponseData struct {
	ID            string     `json:"id"`
	GameID        string     `json:"game_id"`
	AccountID     string     `json:"account_id"`
	AccountUserID string     `json:"account_user_id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type MechWargameLanceResponse struct {
	Data       *MechWargameLanceResponseData    `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechWargameLanceCollectionResponse struct {
	Data       []*MechWargameLanceResponseData  `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type MechWargameLanceRequest struct {
	common_schema.Request
	AccountUserID string `json:"account_user_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
}

type MechWargameLanceQueryParams struct {
	common_schema.QueryParamsPagination
	MechWargameLanceResponseData
}
