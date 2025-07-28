package api

import (
	"time"
)

type GameConfiguration struct {
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

type GameConfigurationResponse struct {
	Data       *GameConfiguration  `json:"data"`
	Error      *ResponseError      `json:"error,omitempty"`
	Pagination *ResponsePagination `json:"pagination,omitempty"`
}

type GameConfigurationCollectionResponse struct {
	Data       []*GameConfiguration `json:"data"`
	Error      *ResponseError       `json:"error,omitempty"`
	Pagination *ResponsePagination  `json:"pagination,omitempty"`
}

type GameConfigurationRequest struct {
	Request
	GameType        string  `json:"game_type"`
	ConfigKey       string  `json:"config_key"`
	ValueType       string  `json:"value_type"`
	DefaultValue    *string `json:"default_value,omitempty"`
	IsRequired      bool    `json:"is_required"`
	Description     *string `json:"description,omitempty"`
	UIHint          *string `json:"ui_hint,omitempty"`
	ValidationRules *string `json:"validation_rules,omitempty"`
}
