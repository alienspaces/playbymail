package game_schema

import (
	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type GameParameter struct {
	GameType     string  `json:"game_type"`
	ConfigKey    string  `json:"config_key"`
	Description  *string `json:"description,omitempty"`
	ValueType    string  `json:"value_type"`
	DefaultValue *string `json:"default_value,omitempty"`
}

type GameParameterCollectionResponse struct {
	Data       []*GameParameter                  `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}
