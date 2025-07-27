package schema

import (
	"encoding/json"
	"time"
)

type AdventureGameInstance struct {
	ID                  string          `json:"id"`
	GameID              string          `json:"game_id"`
	Status              string          `json:"status"`
	CurrentTurn         int             `json:"current_turn"`
	MaxTurns            *int            `json:"max_turns,omitempty"`
	TurnDeadlineHours   int             `json:"turn_deadline_hours"`
	LastTurnProcessedAt *time.Time      `json:"last_turn_processed_at,omitempty"`
	NextTurnDeadline    *time.Time      `json:"next_turn_deadline,omitempty"`
	StartedAt           *time.Time      `json:"started_at,omitempty"`
	CompletedAt         *time.Time      `json:"completed_at,omitempty"`
	GameConfig          json.RawMessage `json:"game_config,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           *time.Time      `json:"updated_at,omitempty"`
	DeletedAt           *time.Time      `json:"deleted_at,omitempty"`
}

type AdventureGameInstanceResponse struct {
	Data       *AdventureGameInstance `json:"data"`
	Error      *ResponseError         `json:"error,omitempty"`
	Pagination *ResponsePagination    `json:"pagination,omitempty"`
}

type AdventureGameInstanceCollectionResponse struct {
	Data       []*AdventureGameInstance `json:"data"`
	Error      *ResponseError           `json:"error,omitempty"`
	Pagination *ResponsePagination      `json:"pagination,omitempty"`
}

type AdventureGameInstanceRequest struct {
	Request
	GameID              string          `json:"game_id"`
	Status              string          `json:"status,omitempty"`
	CurrentTurn         int             `json:"current_turn,omitempty"`
	MaxTurns            *int            `json:"max_turns,omitempty"`
	TurnDeadlineHours   int             `json:"turn_deadline_hours,omitempty"`
	LastTurnProcessedAt *time.Time      `json:"last_turn_processed_at,omitempty"`
	NextTurnDeadline    *time.Time      `json:"next_turn_deadline,omitempty"`
	StartedAt           *time.Time      `json:"started_at,omitempty"`
	CompletedAt         *time.Time      `json:"completed_at,omitempty"`
	GameConfig          json.RawMessage `json:"game_config,omitempty"`
}
