package api

import (
	"time"
)

type GameInstanceConfiguration struct {
	ID             string     `json:"id"`
	GameInstanceID string     `json:"game_instance_id"`
	ConfigKey      string     `json:"config_key"`
	ValueType      string     `json:"value_type"`
	StringValue    *string    `json:"string_value,omitempty"`
	IntegerValue   *int       `json:"integer_value,omitempty"`
	BooleanValue   *bool      `json:"boolean_value,omitempty"`
	JSONValue      *string    `json:"json_value,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

type GameInstanceConfigurationResponse struct {
	Data       *GameInstanceConfiguration `json:"data"`
	Error      *ResponseError             `json:"error,omitempty"`
	Pagination *ResponsePagination        `json:"pagination,omitempty"`
}

type GameInstanceConfigurationCollectionResponse struct {
	Data       []*GameInstanceConfiguration `json:"data"`
	Error      *ResponseError               `json:"error,omitempty"`
	Pagination *ResponsePagination          `json:"pagination,omitempty"`
}

type GameInstanceConfigurationRequest struct {
	Request
	GameInstanceID string  `json:"game_instance_id"`
	ConfigKey      string  `json:"config_key"`
	ValueType      string  `json:"value_type"`
	StringValue    *string `json:"string_value,omitempty"`
	IntegerValue   *int    `json:"integer_value,omitempty"`
	BooleanValue   *bool   `json:"boolean_value,omitempty"`
	JSONValue      *string `json:"json_value,omitempty"`
}
