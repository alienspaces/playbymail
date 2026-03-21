package game_schema

import (
	"time"

	"gitlab.com/alienspaces/playbymail/schema/api/common_schema"
)

type CatalogGameInstanceResponseData struct {
	GameInstanceID        string    `json:"game_instance_id"`
	GameID                string    `json:"game_id"`
	GameName              string    `json:"game_name"`
	GameType              string    `json:"game_type"`
	GameDescription       string    `json:"game_description"`
	TurnDurationHours     int       `json:"turn_duration_hours"`
	GameSubscriptionID    string    `json:"game_subscription_id"`
	AccountName           string    `json:"account_name"`
	RequiredPlayerCount   int       `json:"required_player_count"`
	PlayerCount           int       `json:"player_count"`
	RemainingCapacity     int       `json:"remaining_capacity"`
	DeliveryEmail         bool      `json:"delivery_email"`
	DeliveryPhysicalPost  bool      `json:"delivery_physical_post"`
	DeliveryPhysicalLocal bool      `json:"delivery_physical_local"`
	IsClosedTesting       bool      `json:"is_closed_testing"`
	CreatedAt             time.Time `json:"created_at"`
}

type CatalogGameInstanceCollectionResponse struct {
	Data       []*CatalogGameInstanceResponseData `json:"data"`
	Error      *common_schema.ResponseError       `json:"error,omitempty"`
	Pagination *common_schema.ResponsePagination  `json:"pagination,omitempty"`
}
