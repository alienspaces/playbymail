package game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type GameInstanceParameter struct {
	ID             string     `json:"id"`
	GameInstanceID string     `json:"game_instance_id"`
	ConfigKey      string     `json:"config_key"`
	ValueType      string     `json:"value_type"`
	StringValue    *string    `json:"string_value,omitempty"`
	IntegerValue   *int32     `json:"integer_value,omitempty"`
	BooleanValue   *bool      `json:"boolean_value,omitempty"`
	JSONValue      *string    `json:"json_value,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type GameInstanceParameterResponse struct {
	Data       *GameInstanceParameter            `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameInstanceParameterCollectionResponse struct {
	Data       []*GameInstanceParameter          `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameInstanceParameterRequest struct {
	common_schema.Request
	GameInstanceID string  `json:"game_instance_id"`
	ConfigKey      string  `json:"config_key"`
	ValueType      string  `json:"value_type"`
	StringValue    *string `json:"string_value,omitempty"`
	IntegerValue   *int32  `json:"integer_value,omitempty"`
	BooleanValue   *bool   `json:"boolean_value,omitempty"`
	JSONValue      *string `json:"json_value,omitempty"`
}
