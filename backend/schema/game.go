package schema

import (
	"time"
)

// GameResponseData -
type GameResponseData struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	ProcessedMessageAt *time.Time `json:"processed_message_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

type GameResponse struct {
	Response
	*GameResponseData
}

type GameCollectionResponse = []*GameResponseData

type GameRequest struct {
	Request
	Name string `json:"name"`
}

type GameQueryParams struct {
	QueryParamsPagination
	GameResponseData
}
