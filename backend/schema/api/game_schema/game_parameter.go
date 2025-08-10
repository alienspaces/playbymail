package game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type GameParameter struct {
	ID              string     `json:"id"`
	GameType        string     `json:"game_type"`
	ConfigKey       string     `json:"config_key"`
	ValueType       string     `json:"value_type"`
	DefaultValue    *string    `json:"default_value,omitempty"`
	IsRequired      bool       `json:"is_required"`
	Description     *string    `json:"description,omitempty"`
	UIHint          *string    `json:"ui_hint,omitempty"`
	ValidationRules *string    `json:"validation_rules,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type GameParameterResponse struct {
	Data       *GameParameter                    `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameParameterCollectionResponse struct {
	Data       []*GameParameter                  `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameParameterRequest struct {
	common_schema.Request
	GameType        string  `json:"game_type"`
	ConfigKey       string  `json:"config_key"`
	ValueType       string  `json:"value_type"`
	DefaultValue    *string `json:"default_value,omitempty"`
	IsRequired      bool    `json:"is_required"`
	Description     *string `json:"description,omitempty"`
	UIHint          *string `json:"ui_hint,omitempty"`
	ValidationRules *string `json:"validation_rules,omitempty"`
}
