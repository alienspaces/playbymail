package schema

import (
	"time"
)

// GameResponseData -
type GameResponseData struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	GameType           string     `json:"game_type"`
	ProcessedMessageAt *time.Time `json:"processed_message_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

type GameResponse struct {
	Data       *GameResponseData   `json:"data"`
	Error      *ResponseError      `json:"error,omitempty"`
	Pagination *ResponsePagination `json:"pagination,omitempty"`
}

type GameCollectionResponse struct {
	Data       []*GameResponseData `json:"data"`
	Error      *ResponseError      `json:"error,omitempty"`
	Pagination *ResponsePagination `json:"pagination,omitempty"`
}

type GameRequest struct {
	Request
	Name     string `json:"name"`
	GameType string `json:"game_type"`
}

type GameQueryParams struct {
	QueryParamsPagination
	GameResponseData
}
