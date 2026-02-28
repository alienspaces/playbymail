package game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type GameInstanceResponseData struct {
	ID                                string     `json:"id"`
	GameID                            string     `json:"game_id"`
	Status                            string     `json:"status"`
	CurrentTurn                       int        `json:"current_turn"`
	LastTurnProcessedAt               *time.Time `json:"last_turn_processed_at,omitempty"`
	NextTurnDueAt                     *time.Time `json:"next_turn_due_at,omitempty"`
	StartedAt                         *time.Time `json:"started_at,omitempty"`
	CompletedAt                       *time.Time `json:"completed_at,omitempty"`
	DeliveryPhysicalPost              bool       `json:"delivery_physical_post"`
	DeliveryPhysicalLocal             bool       `json:"delivery_physical_local"`
	DeliveryEmail                     bool       `json:"delivery_email"`
	RequiredPlayerCount               int        `json:"required_player_count"`
	PlayerCount                       int        `json:"player_count"`
	IsClosedTesting                   bool       `json:"is_closed_testing"`
	ClosedTestingJoinGameKey          *string    `json:"join_game_key,omitempty"`
	ClosedTestingJoinGameKeyExpiresAt *time.Time `json:"join_game_key_expires_at,omitempty"`
	CreatedAt                         time.Time  `json:"created_at"`
	UpdatedAt                         *time.Time `json:"updated_at,omitempty"`
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
	GameID                string     `json:"game_id"`
	Status                string     `json:"status,omitempty"`
	CurrentTurn           int        `json:"current_turn,omitempty"`
	LastTurnProcessedAt   *time.Time `json:"last_turn_processed_at,omitempty"`
	NextTurnDueAt         *time.Time `json:"next_turn_due_at,omitempty"`
	StartedAt             *time.Time `json:"started_at,omitempty"`
	CompletedAt           *time.Time `json:"completed_at,omitempty"`
	DeliveryPhysicalPost  bool       `json:"delivery_physical_post,omitempty"`
	DeliveryPhysicalLocal bool       `json:"delivery_physical_local,omitempty"`
	DeliveryEmail         bool       `json:"delivery_email,omitempty"`
	RequiredPlayerCount   int        `json:"required_player_count,omitempty"`
	IsClosedTesting       bool       `json:"is_closed_testing,omitempty"`
}

type JoinGameLinkResponseData struct {
	JoinGameURL string `json:"join_game_url"`
	JoinGameKey string `json:"join_game_key"`
}

type JoinGameLinkResponse struct {
	Data  *JoinGameLinkResponseData    `json:"data"`
	Error *common_schema.ResponseError `json:"error,omitempty"`
}

type InviteTesterRequest struct {
	common_schema.Request
	Email string `json:"email"`
}

type InviteTesterResponseData struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

type InviteTesterResponse struct {
	Data  *InviteTesterResponseData    `json:"data"`
	Error *common_schema.ResponseError `json:"error,omitempty"`
}
