package game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type ManagerGameInstanceResponseData struct {
	GameID               string     `json:"game_id"`
	GameName             string     `json:"game_name"`
	GameType             string     `json:"game_type"`
	GameDescription      string     `json:"game_description"`
	GameSubscriptionID   string     `json:"game_subscription_id"`
	GameInstanceID       *string    `json:"game_instance_id,omitempty"`
	InstanceStatus       *string    `json:"instance_status,omitempty"`
	CurrentTurn          *int       `json:"current_turn,omitempty"`
	RequiredPlayerCount  *int       `json:"required_player_count,omitempty"`
	DeliveryEmail        *bool      `json:"delivery_email,omitempty"`
	DeliveryPhysicalPost *bool      `json:"delivery_physical_post,omitempty"`
	DeliveryPhysicalLocal *bool     `json:"delivery_physical_local,omitempty"`
	IsClosedTesting      *bool      `json:"is_closed_testing,omitempty"`
	StartedAt            *time.Time `json:"started_at,omitempty"`
	NextTurnDueAt        *time.Time `json:"next_turn_due_at,omitempty"`
	InstanceCreatedAt    *time.Time `json:"instance_created_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
}

type ManagerGameInstanceCollectionResponse struct {
	Data       []*ManagerGameInstanceResponseData `json:"data"`
	Error      *common_schema.ResponseError       `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination  `json:"pagination,omitempty"`
}
