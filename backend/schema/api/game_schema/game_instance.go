package game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type GameInstanceResponseData struct {
	ID                  string     `json:"id"`
	GameID              string     `json:"game_id"`
	Status              string     `json:"status"`
	CurrentTurn         int        `json:"current_turn"`
	LastTurnProcessedAt *time.Time `json:"last_turn_processed_at,omitempty"`
	NextTurnDueAt       *time.Time `json:"next_turn_due_at,omitempty"`
	StartedAt           *time.Time `json:"started_at,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

type GameInstanceResponse struct {
	Data       *GameInstanceResponseData         `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameInstanceCollectionResponse struct {
	Data       []*GameInstanceResponseData       `json:"data"`
	Error      *common_schema.ResponseError      `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination `json:"pagination,omitempty"`
}

type GameInstanceRequest struct {
	common_schema.Request
	GameID              string     `json:"game_id"`
	Status              string     `json:"status,omitempty"`
	CurrentTurn         int        `json:"current_turn,omitempty"`
	LastTurnProcessedAt *time.Time `json:"last_turn_processed_at,omitempty"`
	NextTurnDueAt       *time.Time `json:"next_turn_due_at,omitempty"`
	StartedAt           *time.Time `json:"started_at,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
}
