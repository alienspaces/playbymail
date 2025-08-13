package game_schema

import (
	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

// GameParameterConfiguration represents the configuration template for game parameters
type GameParameterConfiguration struct {
	GameType     string  `json:"game_type"`
	ConfigKey    string  `json:"config_key"`
	Description  *string `json:"description,omitempty"`
	ValueType    string  `json:"value_type"`
	DefaultValue *string `json:"default_value,omitempty"`
	IsRequired   bool    `json:"is_required"`
	IsGlobal     bool    `json:"is_global"`
}

type GameParameterConfigurationResponse struct {
	Data       *GameParameterConfiguration       `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameParameterConfigurationCollectionResponse struct {
	Data       []*GameParameterConfiguration     `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameParameterConfigurationQueryParams struct {
	common_schema.QueryParamsPagination
	GameType  string `json:"game_type,omitempty"`
	ConfigKey string `json:"config_key,omitempty"`
	ValueType string `json:"value_type,omitempty"`
	IsGlobal  *bool  `json:"is_global,omitempty"`
}
